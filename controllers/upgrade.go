package controllers

import (
	"context"

	"github.com/openshift-psap/special-resource-operator/pkg/cache"
	"github.com/openshift-psap/special-resource-operator/pkg/clients"
	"github.com/openshift-psap/special-resource-operator/pkg/color"
	"github.com/openshift-psap/special-resource-operator/pkg/exit"
	"github.com/pkg/errors"
	errs "github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	ctrl "sigs.k8s.io/controller-runtime"
	client "sigs.k8s.io/controller-runtime/pkg/client"
)

type NodeUpgradeVersion struct {
	OSVersion      string
	ClusterVersion string
}

// SpecialResourceUpgrade upgrade special resources
func SpecialResourceUpgrade(r *SpecialResourceReconciler, req ctrl.Request) (ctrl.Result, error) {

	var err error

	log = r.Log.WithName(color.Print("upgrade", color.Red))

	cache.Node.List, err = cacheNodes(r, false)
	exit.OnError(errs.Wrap(err, "Failed to cache nodes"))

	runInfo.ClusterUpgradeInfo, err = getUpgradeInfo()
	exit.OnError(errs.Wrap(err, "Failed to get upgrade info"))

	log.Info("TODO: preflight checks")

	return ctrl.Result{Requeue: false}, nil
}

func cacheNodes(r *SpecialResourceReconciler, force bool) (*unstructured.UnstructuredList, error) {

	// The initial list is what we're working with
	// a SharedInformer will update the list of nodes if
	// more nodes join the cluster.
	cached := int64(len(cache.Node.List.Items))
	if cached == cache.Node.Count && !force {
		return cache.Node.List, nil
	}

	cache.Node.List.SetAPIVersion("v1")
	cache.Node.List.SetKind("NodeList")

	opts := []client.ListOption{}

	// Only filter if we have a selector set, otherwise zero nodes will be
	// returned and no labels can be extracted. Set the default worker label
	// otherwise.
	if len(r.specialresource.Spec.NodeSelector) > 0 {
		opts = append(opts, client.MatchingLabels{r.specialresource.Spec.NodeSelector: "true"})
	} else {
		opts = append(opts, client.MatchingLabels{"node-role.kubernetes.io/worker": ""})
	}

	err := clients.Interface.List(context.TODO(), cache.Node.List, opts...)
	if err != nil {
		return nil, errors.Wrap(err, "Client cannot get NodeList")
	}

	return cache.Node.List, err
}

func getUpgradeInfo() (map[string]NodeUpgradeVersion, error) {

	var found bool
	var info = make(map[string]NodeUpgradeVersion)

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
			return nil, errs.New("Label " + short + " not found is NFD running? Check node labels")
		}

		short = "feature.node.kubernetes.io/system-os_release.RHEL_VERSION"
		if rhelVersion, found = labels[short]; !found {
			return nil, errs.New("Label " + short + " not found is NFD running? Check node labels")
		}

		short = "feature.node.kubernetes.io/system-os_release.VERSION_ID"
		if clusterVersion, found = labels[short]; !found {
			return nil, errs.New("Label " + short + " not found is NFD running? Check node labels")
		}

		info[kernelFullVersion] = NodeUpgradeVersion{OSVersion: rhelVersion, ClusterVersion: clusterVersion}
	}

	return info, nil
}
