package upgrade

import (
	"errors"

	"github.com/go-logr/logr"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/openshift-psap/special-resource-operator/pkg/cache"
	"github.com/openshift-psap/special-resource-operator/pkg/color"
	"github.com/openshift-psap/special-resource-operator/pkg/registry"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

var (
	log logr.Logger
)

func init() {
	log = zap.New(zap.UseDevMode(true)).WithName(color.Print("upgrade", color.Blue))
}

type NodeVersion struct {
	OSVersion      string
	ClusterVersion string
}

func NodeVersionInfo() (map[string]NodeVersion, error) {

	var found bool
	var info = make(map[string]NodeVersion)

	// Assuming all nodes are running the same kernel version,
	// one could easily add driver-kernel-versions for each node.
	for _, node := range cache.Node.List.Items {

		var rhelVersion string
		var kernelFullVersion string
		var clusterVersion string

		labels := node.GetLabels()
		// We only need to check for the key, the value
		// is available if the key is there
		short := "feature.node.kubernetes.io/kernel-version.full"
		if kernelFullVersion, found = labels[short]; !found {
			return nil, errors.New("Label " + short + " not found is NFD running? Check node labels")
		}

		short = "feature.node.kubernetes.io/system-os_release.RHEL_VERSION"
		if rhelVersion, found = labels[short]; !found {
			return nil, errors.New("Label " + short + " not found is NFD running? Check node labels")
		}

		short = "feature.node.kubernetes.io/system-os_release.VERSION_ID"
		if clusterVersion, found = labels[short]; !found {
			return nil, errors.New("Label " + short + " not found is NFD running? Check node labels")
		}

		info[kernelFullVersion] = NodeVersion{OSVersion: rhelVersion, ClusterVersion: clusterVersion}
	}

	return info, nil
}

func DriverToolkit(entries []string) error {

	for _, entry := range entries {

		log.Info("History", "entry", entry)
		var layer v1.Layer
		if layer = registry.LastLayer(entry); layer == nil {
			continue
		}
		version, imageURL := registry.ReleaseManifests("driver-toolkit", layer)

		if version != "" {
			log.Info("version", "V", version+" : "+imageURL)
		}
	}

	return nil
}
