package helmer

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-logr/logr"
	"github.com/openshift-psap/special-resource-operator/pkg/color"
	"github.com/openshift-psap/special-resource-operator/pkg/exit"
	"github.com/openshift-psap/special-resource-operator/pkg/slice"
	errs "github.com/pkg/errors"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/releaseutil"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

var (
	log logr.Logger
)

func init() {
	log = zap.New(zap.UseDevMode(true)).WithName(color.Print("helm", color.Blue))
	err := OpenShiftInstallOrder()
	exit.OnError(err)
}

type HelmChart struct {
	Name       string `json:"name"`
	Version    string `json:"version"`
	Repository string `json:"repository"`
}

type HelmDependency struct {
	Name       string   `json:"name"`
	Version    string   `json:"version"`
	Repository string   `json:"repository"`
	Tags       []string `json:"tags"`
}

func (in *HelmDependency) DeepCopyInto(out *HelmDependency) {
	out.Name = in.Name
	out.Version = in.Version
	out.Repository = in.Repository
	out.Tags = make([]string, len(in.Tags))
	copy(out.Tags, in.Tags)
}

func Load(ch interface{}) (*chart.Chart, error) {

	var curr HelmChart

	switch v := ch.(type) {
	case *chart.Dependency:
		curr.Name = v.Name
		curr.Version = v.Version
		curr.Repository = v.Repository
	case chart.Dependency:
		curr.Name = v.Name
		curr.Version = v.Version
		curr.Repository = v.Repository
	case HelmChart:
		curr = v
	default:
		exit.OnError(errs.New("Unknown Type:" + reflect.TypeOf(v).String()))

	}

	var repo string
	if strings.Contains(curr.Repository, "file:///") {
		repo = strings.Replace(curr.Repository, "file://", "", -1)
		log.Info("DEBUG", "repo", repo)
	} else {
		exit.OnError(errs.New("Only file:/// currently supported"))
	}
	loaded, err := loader.Load(repo + "/" + curr.Name + "-" + curr.Version)

	return loaded, err

}

func OpenShiftInstallOrder() error {

	idx := slice.Find(releaseutil.InstallOrder, "Service")
	releaseutil.InstallOrder = slice.Insert(releaseutil.InstallOrder, idx, "BuildConfig")
	releaseutil.InstallOrder = slice.Insert(releaseutil.InstallOrder, idx, "ImageStream")
	releaseutil.InstallOrder = slice.Insert(releaseutil.InstallOrder, idx, "SecurityContextConstraints")

	return nil
}

func TemplateChart(ch chart.Chart, vals map[string]interface{}) ([]byte, error) {

	actionConfig := action.Configuration{}

	client := action.NewInstall(&actionConfig)

	client.DryRun = true
	client.ReleaseName = ch.Metadata.Name
	client.Replace = true
	client.ClientOnly = true
	client.APIVersions = []string{}
	client.IncludeCRDs = true

	if client.Version == "" {
		client.Version = ">0.0.0-0"
	}

	if ch.Metadata.Type != "" && ch.Metadata.Type != "application" {
		return nil, errs.New("Chart has an unsupported type and is not installable:" + ch.Metadata.Type)
	}

	out := new(bytes.Buffer)

	rel, err := client.Run(&ch, vals)

	if rel != nil {
		var manifests bytes.Buffer
		fmt.Fprintln(&manifests, strings.TrimSpace(rel.Manifest))
		if !client.DisableHooks {
			for _, m := range rel.Hooks {
				fmt.Fprintf(&manifests, "---\n# Source: %s\n%s\n", m.Path, m.Manifest)
			}
		}
		fmt.Fprintf(out, "%s", manifests.String())
	}
	return out.Bytes(), err
}
