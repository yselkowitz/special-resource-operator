package e2e

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"strings"
	"time"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/openshift-psap/special-resource-operator/pkg/warn"
	"github.com/openshift-psap/special-resource-operator/test/framework"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = ginkgo.Describe("[basic][ping-pong] create and deploy ping-poing", func() {

	cs := framework.NewClientSet()
	cl := framework.NewControllerRuntimeClient()
	ginkgo.It("Can create and deploy ping-pong", func() {
		ginkgo.By("Creating ping-pong #1")
		specialResourceCreate(cs, cl, "../../../charts/example/ping-pong-0.0.1/ping-pong.yaml")
		checkPingPong(cs, cl)
		specialResourceDelete(cs, cl, "../../../charts/example/ping-pong-0.0.1/ping-pong.yaml")
	})
})

func checkPingPong(cs *framework.ClientSet, cl client.Client) {

	for {
		time.Sleep(60 * time.Second)

		ginkgo.By("Waiting for ping-poing Pods to be ready")
		opts := metav1.ListOptions{}
		pods, err := cs.Pods("ping-pong").List(context.TODO(), opts)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		for _, pod := range pods.Items {
			//run command in pod
			ginkgo.By("Ensuring that ping-pong is working")
			log := getPodLogs(pod)
			if !strings.Contains(log, "Ping") || !strings.Contains(log, "Pong") {
				warn.OnError(errors.New("Did not see Ping or either Pong, waiting"))
			}

			if strings.Contains(log, "Ping") && strings.Contains(log, "Pong") {
				ginkgo.By("Found Ping, Pong in logs, done")
				return
			}

		}
	}

}

func specialResourceDelete(cs *framework.ClientSet, cl client.Client, path string) {

	ginkgo.By("deleting ping-pong")
	sr, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	framework.DeleteFromYAMLWithCR(sr, cl)
}

func specialResourceCreate(cs *framework.ClientSet, cl client.Client, path string) {

	ginkgo.By("creating ping-pong")
	sr, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	framework.CreateFromYAML(sr, cl)
}

func getPodLogs(pod corev1.Pod) string {
	podLogOpts := corev1.PodLogOptions{}
	config, err := rest.InClusterConfig()
	if err != nil {
		return "error in getting config"
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return "error in getting access to K8S"
	}
	req := clientset.CoreV1().Pods(pod.Namespace).GetLogs(pod.Name, &podLogOpts)
	podLogs, err := req.Stream(context.TODO())
	if err != nil {
		return "error in opening stream"
	}
	defer podLogs.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, podLogs)
	if err != nil {
		return "error in copy information from podLogs to buf"
	}
	str := buf.String()

	return str
}
