package controllers

import (
	"context"
	"fmt"

	"github.com/openshift-psap/special-resource-operator/pkg/clients"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const specialresourceFinalizer = "finalizer.sro.openshift.io"

func reconcileFinalizers(r *SpecialResourceReconciler) error {
	if contains(r.specialresource.GetFinalizers(), specialresourceFinalizer) {
		// Run finalization logic for specialresource
		if err := finalizeSpecialResource(r); err != nil {
			log.Info("Finalization logic failed.", "error", fmt.Sprintf("%v", err))
			return err
		}

		controllerutil.RemoveFinalizer(&r.specialresource, specialresourceFinalizer)
		err := clients.Interface.Update(context.TODO(), &r.specialresource)
		if err != nil {
			log.Info("Could not remove finalizer after running finalization logic", "error", fmt.Sprintf("%v", err))
			return err
		}
	}
	return nil
}

func finalizeSpecialResource(r *SpecialResourceReconciler) error {
	// TODO(user): Add the cleanup steps that the operator
	// needs to do before the CR can be deleted. Examples
	// of finalizers include performing backups and deleting
	// resources that are not owned by this CR, like a PVC.

	log.Info("Successfully finalized", "SpecialResource:", r.specialresource.Name)
	return nil
}

func addFinalizer(r *SpecialResourceReconciler) error {
	log.Info("Adding finalizer to special resource")
	controllerutil.AddFinalizer(&r.specialresource, specialresourceFinalizer)

	// Update CR
	err := clients.Interface.Update(context.TODO(), &r.specialresource)
	if err != nil {
		log.Info("Adding finalizer failed", "error", fmt.Sprintf("%v", err))
		return err
	}
	return nil
}

func contains(list []string, s string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}
