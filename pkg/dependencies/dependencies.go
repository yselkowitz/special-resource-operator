package dependencies

import (
	"context"
	"os"

	"github.com/openshift-psap/special-resource-operator/pkg/clients"
	"github.com/openshift-psap/special-resource-operator/pkg/exit"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/types"
)

func UpdateConfigMap(parent string, child string) {

	cm := &unstructured.Unstructured{}
	cm.SetAPIVersion("v1")
	cm.SetKind("ConfigMap")

	dep := types.NamespacedName{Namespace: os.Getenv("OPERATOR_NAMESPACE"), Name: "special-resource-depedencies"}

	err := clients.Interface.Get(context.TODO(), dep, cm)

	if apierrors.IsNotFound(err) {
		exit.OnError(err)
	}

	data, found, err := unstructured.NestedMap(cm.Object, "data")
	exit.OnError(err)

	dependencies := make(map[string]interface{})
	dependencies[child] = parent

	if !found {
		data = make(map[string]interface{})
		data["data"] = dependencies
		err := unstructured.SetNestedMap(cm.Object, dependencies, "data")
		exit.OnError(err)
	}

	err = unstructured.SetNestedMap(cm.Object, dependencies, "data")
	exit.OnError(err)

	err = clients.Interface.Update(context.TODO(), cm)
	exit.OnError(err)
}
