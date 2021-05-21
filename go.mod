module github.com/openshift-psap/special-resource-operator

go 1.16

require (
	cloud.google.com/go v0.58.0 // indirect
	github.com/go-logr/logr v0.2.1
	github.com/google/go-containerregistry v0.5.1
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/mitchellh/hashstructure/v2 v2.0.1
	github.com/onsi/ginkgo v1.14.1
	github.com/onsi/gomega v1.10.2
	github.com/openshift/api v0.0.0-20201005153912-821561a7f2a2
	github.com/openshift/client-go v0.0.0-20200827190008-3062137373b5
	github.com/openshift/library-go v0.0.0-20200911100307-610c6e9e90b8
	github.com/pkg/errors v0.9.1
	github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring v0.42.1
	github.com/prometheus/client_golang v1.7.1
	github.com/prometheus/common v0.14.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	golang.org/x/tools v0.0.0-20201013201025-64a9e34f3752 // indirect
	helm.sh/helm/v3 v3.5.1
	k8s.io/api v0.20.1
	k8s.io/apimachinery v0.20.1
	k8s.io/client-go v0.20.1
	rsc.io/letsencrypt v0.0.3 // indirect
	sigs.k8s.io/controller-runtime v0.6.3
	sigs.k8s.io/yaml v1.2.0
)
