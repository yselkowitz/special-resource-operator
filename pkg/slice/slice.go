package slice

import (
	"github.com/go-logr/logr"
	srov1beta1 "github.com/openshift-psap/special-resource-operator/api/v1beta1"
	"github.com/openshift-psap/special-resource-operator/pkg/color"
	"helm.sh/helm/v3/pkg/chart"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

var (
	log logr.Logger
)

func init() {
	log = zap.New(zap.UseDevMode(true)).WithName(color.Print("slice", color.Blue))
}

// Find returns the smallest index i at which x == a[i],
// or len(a) if there is no such index.
func Find(a []string, x string) int {
	for i, n := range a {
		if x == n {
			return i
		}
	}
	return len(a)
}

// Contains tells whether a contains x.
func Contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

func FindSR(a []srov1beta1.SpecialResource, x string, by string) int {
	for i, n := range a {
		if by == "Name" {
			if x == n.GetName() {
				return i
			}
		}
	}
	return -1
}

func FindCRFile(a []*chart.File, x string) int {
	for i, n := range a {
		if n.Name == x+".yaml" {
			return i
		}
	}
	return -1
}
