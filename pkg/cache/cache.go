package cache

import "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

var Node NodesCache

func init() {
	Node.Count = 0xDEADBEEF
	Node.List = &unstructured.UnstructuredList{
		Object: map[string]interface{}{},
		Items:  []unstructured.Unstructured{},
	}
}

type NodesCache struct {
	List  *unstructured.UnstructuredList
	Count int64
}
