package storage

import (
	"context"

	"github.com/openshift-psap/special-resource-operator/pkg/clients"
	"github.com/openshift-psap/special-resource-operator/pkg/exit"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/types"
)

var Driver string

func init() {
	Driver = "ConfigMap"
}

func GetConfigMap(namespace string, name string) *unstructured.Unstructured {

	cm := &unstructured.Unstructured{}
	cm.SetAPIVersion("v1")
	cm.SetKind("ConfigMap")

	dep := types.NamespacedName{Namespace: namespace, Name: name}

	err := clients.Interface.Get(context.TODO(), dep, cm)

	if apierrors.IsNotFound(err) {
		exit.OnError(err)
	}

	return cm
}

func CheckConfigMapEntry(key string, ins types.NamespacedName) string {

	cm := GetConfigMap(ins.Namespace, ins.Name)

	data, found, err := unstructured.NestedMap(cm.Object, "data")
	exit.OnError(err)
	// No parent found for depedency just return
	if !found {
		return ""
	}
	// We have a dependency return the value
	if value, found := data[key]; found {
		return value.(string)
	}

	return ""
}

func UpdateConfigMapEntry(key string, value string, ins types.NamespacedName) {

	cm := GetConfigMap(ins.Namespace, ins.Name)

	// Just looking if exists so we can create or update
	entries, found, err := unstructured.NestedMap(cm.Object, "data")
	exit.OnError(err)

	if !found {
		entries = make(map[string]interface{})
	}

	entries[key] = value

	err = unstructured.SetNestedMap(cm.Object, entries, "data")
	exit.OnError(err)

	err = clients.Interface.Update(context.TODO(), cm)
	exit.OnError(err)
}

func DeleteConfigMapEntry(delete string, ins types.NamespacedName) {

	cm := GetConfigMap(ins.Namespace, ins.Name)

	// Just looking if exists so we can create or update
	old, found, err := unstructured.NestedMap(cm.Object, "data")
	exit.OnError(err)

	if !found {
		return
	}

	new := make(map[string]interface{})

	for k, v := range old {
		if delete != k {
			new[k] = v
		}
	}

	err = unstructured.SetNestedMap(cm.Object, new, "data")
	exit.OnError(err)

	err = clients.Interface.Update(context.TODO(), cm)
	exit.OnError(err)

}
