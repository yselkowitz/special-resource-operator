package resource

import (
	"github.com/openshift-psap/special-resource-operator/pkg/exit"
	errs "github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func IsNamespaced(kind string) bool {
	if kind == "Namespace" ||
		kind == "ClusterRole" ||
		kind == "ClusterRoleBinding" ||
		kind == "SecurityContextConstraint" ||
		kind == "SpecialResource" {
		return false
	}
	return true
}

func IsNotUpdateable(kind string) bool {
	// ServiceAccounts cannot be updated, maybe delete and create?
	if kind == "ServiceAccount" || kind == "Pod" {
		return true
	}
	return false
}

// Some resources need an updated resourceversion, during updates
func NeedsResourceVersionUpdate(kind string) bool {
	if kind == "SecurityContextConstraints" ||
		kind == "Service" ||
		kind == "ServiceMonitor" ||
		kind == "Route" ||
		kind == "BuildConfig" ||
		kind == "ImageStream" ||
		kind == "PrometheusRule" ||
		kind == "CSIDriver" ||
		kind == "Issuer" ||
		kind == "Certificate" {
		return true
	}
	return false

}

func UpdateResourceVersion(req *unstructured.Unstructured, found *unstructured.Unstructured) error {

	kind := found.GetKind()

	if NeedsResourceVersionUpdate(kind) {
		version, fnd, err := unstructured.NestedString(found.Object, "metadata", "resourceVersion")
		exit.OnErrorOrNotFound(fnd, err)

		if err := unstructured.SetNestedField(req.Object, version, "metadata", "resourceVersion"); err != nil {
			return errs.Wrap(err, "Couldn't update ResourceVersion")
		}

	}
	if kind == "Service" {
		clusterIP, fnd, err := unstructured.NestedString(found.Object, "spec", "clusterIP")
		exit.OnErrorOrNotFound(fnd, err)

		if err := unstructured.SetNestedField(req.Object, clusterIP, "spec", "clusterIP"); err != nil {
			return errs.Wrap(err, "Couldn't update clusterIP")
		}
		return nil
	}
	return nil
}
