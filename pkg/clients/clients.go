package clients

import (
	"os"

	"github.com/go-logr/logr"
	"github.com/openshift-psap/special-resource-operator/pkg/color"
	clientconfigv1 "github.com/openshift/client-go/config/clientset/versioned/typed/config/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

var (
	log        logr.Logger
	Interface  *ClientsInterface
	RestConfig *rest.Config
	Namespace  string
)

func init() {

}

type ClientsInterface struct {
	client.Client
	kubernetes.Clientset
	clientconfigv1.ConfigV1Client
}

func init() {
	log = zap.New(zap.UseDevMode(true)).WithName(color.Print("exit", color.Red))
}

// GetKubeClientSetOrDie Add a native non-caching client for advanced CRUD operations
func GetKubeClientSetOrDie() kubernetes.Clientset {

	clientSet, err := kubernetes.NewForConfig(RestConfig)
	if err != nil {
		log.Info(color.Print("GetConfigClientOrDie: "+err.Error(), color.Red))
		os.Exit(1)
	}
	return *clientSet
}

// GetConfigClientOrDie Add a configv1 client to the reconciler
func GetConfigClientOrDie() clientconfigv1.ConfigV1Client {

	client, err := clientconfigv1.NewForConfig(RestConfig)
	if err != nil {
		log.Info(color.Print("GetConfigClientOrDie: "+err.Error(), color.Red))
		os.Exit(1)
	}
	return *client
}
