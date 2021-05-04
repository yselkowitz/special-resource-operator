package kernel

import (
	"strings"

	"github.com/go-logr/logr"
	errs "github.com/pkg/errors"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	"github.com/openshift-psap/special-resource-operator/pkg/color"
	"github.com/openshift-psap/special-resource-operator/pkg/exit"
	"github.com/openshift-psap/special-resource-operator/pkg/hash"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

var (
	log logr.Logger
)

func init() {
	log = zap.New(zap.UseDevMode(true)).WithName(color.Print("kernel", color.Green))
}

func SetAffineAttributes(obj *unstructured.Unstructured,
	kernelFullVersion string,
	operatingSystemMajorMinor string) error {

	kernelVersion := strings.ReplaceAll(kernelFullVersion, "_", "-")
	hash64 := hash.FNV64a(operatingSystemMajorMinor + "-" + kernelVersion)
	name := obj.GetName() + "-" + hash64
	obj.SetName(name)

	if obj.GetKind() == "DaemonSet" {
		err := unstructured.SetNestedField(obj.Object, name, "metadata", "labels", "app")
		exit.OnError(err)
		err = unstructured.SetNestedField(obj.Object, name, "spec", "selector", "matchLabels", "app")
		exit.OnError(err)
		err = unstructured.SetNestedField(obj.Object, name, "spec", "template", "metadata", "labels", "app")
		exit.OnError(err)
		err = unstructured.SetNestedField(obj.Object, name, "spec", "template", "metadata", "labels", "app")
		exit.OnError(err)
	}

	if err := SetVersionNodeAffinity(obj, kernelFullVersion); err != nil {
		return errs.Wrap(err, "Cannot set kernel version node affinity for obj: "+obj.GetKind())
	}
	return nil
}

func SetVersionNodeAffinity(obj *unstructured.Unstructured, kernelFullVersion string) error {

	if strings.Compare(obj.GetKind(), "DaemonSet") == 0 {
		if err := versionNodeAffinity(kernelFullVersion, obj, "spec", "template", "spec", "nodeSelector"); err != nil {
			return errs.Wrap(err, "Cannot setup DaemonSet kernel version affinity")
		}
	}
	if strings.Compare(obj.GetKind(), "Pod") == 0 {
		if err := versionNodeAffinity(kernelFullVersion, obj, "spec", "nodeSelector"); err != nil {
			return errs.Wrap(err, "Cannot setup Pod kernel version affinity")
		}
	}
	if strings.Compare(obj.GetKind(), "BuildConfig") == 0 {
		if err := versionNodeAffinity(kernelFullVersion, obj, "spec", "nodeSelector"); err != nil {
			return errs.Wrap(err, "Cannot setup BuildConfig kernel version affinity")
		}
	}

	return nil
}

func versionNodeAffinity(kernelFullVersion string, obj *unstructured.Unstructured, fields ...string) error {

	nodeSelector, found, err := unstructured.NestedMap(obj.Object, fields...)
	exit.OnError(err)

	if !found {
		nodeSelector = make(map[string]interface{})
	}

	nodeSelector["feature.node.kubernetes.io/kernel-version.full"] = kernelFullVersion

	if err := unstructured.SetNestedMap(obj.Object, nodeSelector, fields...); err != nil {
		return errs.Wrap(err, "Cannot update nodeSelector")
	}

	return nil
}

func IsObjectAffine(obj *unstructured.Unstructured) bool {

	annotations := obj.GetAnnotations()

	if affine, found := annotations["specialresource.openshift.io/kernel-affine"]; found && affine == "true" {
		log.Info("Object is Kernel Affine")
		return true
	}

	return false
}
