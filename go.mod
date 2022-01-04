module gitlab.com/ignitionrobotics/web/cloudsim

go 1.15

require (
	cloud.google.com/go v0.46.3 // indirect
	github.com/auth0/go-jwt-middleware v0.0.0-20200810150920-a32d7af194d1 // indirect
	github.com/aws/aws-sdk-go v1.34.25
	github.com/caarlos0/env v3.5.0+incompatible
	github.com/casbin/casbin/v2 v2.6.11
	github.com/casbin/gorm-adapter/v2 v2.1.0
	github.com/creasty/defaults v1.5.1
	github.com/elastic/go-elasticsearch/v8 v8.0.0-20200916133832-96c1726edbaa // indirect
	github.com/fatih/color v1.10.0 // indirect
	github.com/fatih/structtag v1.2.0
	github.com/golang/protobuf v1.4.2
	github.com/gorilla/mux v1.8.0 // indirect
	github.com/gorilla/websocket v1.4.2
	github.com/imdario/mergo v0.3.8
	github.com/itchyny/gojq v0.12.5
	github.com/jinzhu/gorm v1.9.16
	github.com/johannesboyne/gofakes3 v0.0.0-20210116212202-e8b5dbd08102
	github.com/lib/pq v1.8.0 // indirect
	github.com/mitchellh/copystructure v1.0.0
	github.com/mitchellh/mapstructure v1.1.2
	github.com/onsi/ginkgo v1.14.1 // indirect
	github.com/onsi/gomega v1.10.2 // indirect
	github.com/pkg/errors v0.9.1
	github.com/prometheus/common v0.11.1 // indirect
	github.com/satori/go.uuid v1.2.0
	github.com/solo-io/gloo v1.4.12
	github.com/solo-io/solo-kit v0.13.8-patch4
	github.com/stretchr/testify v1.7.0
	gitlab.com/ignitionrobotics/billing/credits v0.0.0-20211110124124-7dbf77a10783
	gitlab.com/ignitionrobotics/billing/payments v0.0.0-20211110174459-d690b1b701e8
	gitlab.com/ignitionrobotics/web/fuelserver v0.0.0-20200916210816-e30ab5ed9d47
	gitlab.com/ignitionrobotics/web/ign-go v0.0.0-20211109185912-2322af1b387c
	golang.org/x/crypto v0.0.0-20200820211705-5c72a883971a // indirect
	golang.org/x/mod v0.4.0 // indirect
	golang.org/x/net v0.0.0-20201021035429-f5854403a974
	golang.org/x/text v0.3.4 // indirect
	golang.org/x/tools v0.0.0-20210105210202-9ed45478a130 // indirect
	google.golang.org/protobuf v1.25.0
	gopkg.in/go-playground/assert.v1 v1.2.1 // indirect
	gopkg.in/go-playground/validator.v9 v9.31.0
	gopkg.in/yaml.v2 v2.4.0
	honnef.co/go/tools v0.0.1-2020.1.6 // indirect
	k8s.io/api v0.17.2
	k8s.io/apimachinery v0.17.2
	k8s.io/client-go v11.0.0+incompatible
	k8s.io/kubernetes v1.17.1
)

replace (
	github.com/Azure/go-autorest => github.com/Azure/go-autorest v13.0.0+incompatible
	github.com/Sirupsen/logrus => github.com/sirupsen/logrus v1.4.2
	github.com/docker/docker => github.com/moby/moby v0.7.3-0.20190826074503-38ab9da00309
	github.com/golang/mock => github.com/golang/mock v1.4.3
	// kube 1.17
	k8s.io/api => k8s.io/api v0.17.1
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.17.1
	k8s.io/apimachinery => k8s.io/apimachinery v0.17.1
	k8s.io/apiserver => k8s.io/apiserver v0.17.1
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.17.1
	k8s.io/client-go => k8s.io/client-go v0.17.1
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.17.1
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.17.1
	k8s.io/code-generator => k8s.io/code-generator v0.17.1
	k8s.io/component-base => k8s.io/component-base v0.17.1
	k8s.io/cri-api => k8s.io/cri-api v0.17.1
	k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.17.1
	k8s.io/gengo => k8s.io/gengo v0.0.0-20190822140433-26a664648505
	k8s.io/heapster => k8s.io/heapster v1.2.0-beta.1
	k8s.io/klog => github.com/stefanprodan/klog v0.0.0-20190418165334-9cbb78b20423
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.17.1
	k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.17.1
	k8s.io/kube-openapi => k8s.io/kube-openapi v0.0.0-20190816220812-743ec37842bf
	k8s.io/kube-proxy => k8s.io/kube-proxy v0.17.1
	k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.17.1
	k8s.io/kubectl => k8s.io/kubectl v0.17.1
	k8s.io/kubelet => k8s.io/kubelet v0.17.1
	k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.17.1
	k8s.io/metrics => k8s.io/metrics v0.17.1
	k8s.io/node-api => k8s.io/node-api v0.17.1
	k8s.io/repo-infra => k8s.io/repo-infra v0.0.0-20181204233714-00fe14e3d1a3
	k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.17.1
	k8s.io/sample-cli-plugin => k8s.io/sample-cli-plugin v0.17.1
	k8s.io/sample-controller => k8s.io/sample-controller v0.17.1
	k8s.io/utils => k8s.io/utils v0.0.0-20190801114015-581e00157fb1
)
