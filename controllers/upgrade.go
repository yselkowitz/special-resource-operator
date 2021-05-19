package controllers

import (
	"github.com/openshift-psap/special-resource-operator/pkg/cache"
	"github.com/openshift-psap/special-resource-operator/pkg/cluster"
	"github.com/openshift-psap/special-resource-operator/pkg/color"
	"github.com/openshift-psap/special-resource-operator/pkg/exit"
	"github.com/openshift-psap/special-resource-operator/pkg/upgrade"
	"github.com/pkg/errors"
	ctrl "sigs.k8s.io/controller-runtime"
)

// SpecialResourceUpgrade upgrade special resources
func SpecialResourceUpgrade(r *SpecialResourceReconciler, req ctrl.Request) (ctrl.Result, error) {

	var err error

	log = r.Log.WithName(color.Print("upgrade", color.Blue))

	err = cache.Nodes(r.specialresource.Spec.NodeSelector, false)
	exit.OnError(errors.Wrap(err, "Failed to cache nodes"))

	RunInfo.ClusterUpgradeInfo, err = upgrade.NodeVersionInfo()
	exit.OnError(errors.Wrap(err, "Failed to get upgrade info"))

	log.Info("TODO: preflight checks")

	history, err := cluster.VersionHistory()
	exit.OnError(errors.Wrap(err, "Could not get version history"))

	upgrade.DriverToolkit(history)

	return ctrl.Result{Requeue: false}, nil
}
