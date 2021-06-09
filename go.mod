module gitlab.com/ignitionrobotics/web/cloudsim

go 1.15

require (
	cloud.google.com/go v0.46.3 // indirect
	github.com/auth0/go-jwt-middleware v0.0.0-20200810150920-a32d7af194d1 // indirect
	github.com/avast/retry-go v2.4.3+incompatible // indirect
	github.com/aws/aws-sdk-go v1.34.25
	github.com/caarlos0/env v3.5.0+incompatible
	github.com/casbin/casbin/v2 v2.6.11
	github.com/casbin/gorm-adapter/v2 v2.1.0
	github.com/creasty/defaults v1.5.1
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/elastic/go-elasticsearch/v8 v8.0.0-20200916133832-96c1726edbaa // indirect
	github.com/elazarl/goproxy v0.0.0-20190421051319-9d40249d3c2f // indirect
	github.com/elazarl/goproxy/ext v0.0.0-20190421051319-9d40249d3c2f // indirect
	github.com/envoyproxy/go-control-plane v0.9.6-0.20200529035633-fc42e08917e9 // indirect
	github.com/envoyproxy/protoc-gen-validate v0.4.0 // indirect
	github.com/fatih/color v1.10.0 // indirect
	github.com/fatih/structtag v1.2.0
	github.com/go-playground/form v3.1.4+incompatible
	github.com/golang/mock v1.4.3 // indirect
	github.com/golang/protobuf v1.4.2
	github.com/google/go-cmp v0.5.4 // indirect
	github.com/gophercloud/gophercloud v0.6.0 // indirect
	github.com/gorilla/mux v1.8.0
	github.com/gorilla/websocket v1.4.2
	github.com/imdario/mergo v0.3.8
	github.com/jinzhu/gorm v1.9.16
	github.com/johannesboyne/gofakes3 v0.0.0-20210116212202-e8b5dbd08102
	github.com/joho/godotenv v1.3.0
	github.com/lib/pq v1.8.0 // indirect
	github.com/mattn/go-pointer v0.0.0-20180825124634-49522c3f3791
	github.com/mitchellh/copystructure v1.0.0
	github.com/mitchellh/mapstructure v1.1.2
	github.com/mitchellh/reflectwalk v1.0.1 // indirect
	github.com/onsi/ginkgo v1.14.1 // indirect
	github.com/onsi/gomega v1.10.2 // indirect
	github.com/panjf2000/ants v0.0.0-20190122063359-2ba69cd1384d
	github.com/pkg/errors v0.9.1
	github.com/prometheus/common v0.11.1 // indirect
	github.com/radovskyb/watcher v1.0.7 // indirect
	github.com/satori/go.uuid v1.2.0
	github.com/solo-io/gloo v1.1.0
	github.com/solo-io/go-utils v0.16.2 // indirect
	github.com/solo-io/solo-kit v0.13.8-patch4
	github.com/stretchr/testify v1.6.1
	gitlab.com/ignitionrobotics/web/fuelserver v0.0.0-20200916210816-e30ab5ed9d47
	gitlab.com/ignitionrobotics/web/ign-go v0.0.0-20201013152111-8655ead5c276
	gitlab.com/ignitionrobotics/web/scheduler v0.5.1-0.20200114185916-4bd85f4ff2d6 // indirect
	golang.org/x/crypto v0.0.0-20200820211705-5c72a883971a // indirect
	golang.org/x/mod v0.4.0 // indirect
	golang.org/x/net v0.0.0-20201021035429-f5854403a974
	golang.org/x/sys v0.0.0-20201009025420-dfb3f7c4e634 // indirect
	golang.org/x/text v0.3.4 // indirect
	golang.org/x/tools v0.0.0-20210105210202-9ed45478a130 // indirect
	google.golang.org/protobuf v1.25.0 // indirect
	gopkg.in/go-playground/assert.v1 v1.2.1 // indirect
	gopkg.in/go-playground/validator.v9 v9.31.0
	gopkg.in/yaml.v2 v2.4.0
	honnef.co/go/tools v0.0.1-2020.1.6 // indirect
	k8s.io/api v0.20.4
	k8s.io/apiextensions-apiserver v0.17.2 // indirect
	k8s.io/apimachinery v0.20.4
	k8s.io/client-go v0.20.4
	k8s.io/kubectl v0.20.4
	k8s.io/utils v0.20.4 // indirect
)

replace (
	github.com/docker/docker => github.com/moby/moby v0.7.3-0.20190826074503-38ab9da00309
	github.com/golang/mock => github.com/golang/mock v1.4.3
	// kube 1.17
	k8s.io/api => k8s.io/api v0.17.1
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.17.1
	k8s.io/apimachinery => k8s.io/apimachinery v0.17.1
	k8s.io/apiserver => k8s.io/apiserver v0.17.1
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.17.1
	k8s.io/client-go => k8s.io/client-go v0.17.1
	k8s.io/code-generator => k8s.io/code-generator v0.17.1
	k8s.io/component-base => k8s.io/component-base v0.17.1
	k8s.io/gengo => k8s.io/gengo v0.0.0-20190822140433-26a664648505
	k8s.io/klog => github.com/stefanprodan/klog v0.0.0-20190418165334-9cbb78b20423
	k8s.io/kube-openapi => k8s.io/kube-openapi v0.0.0-20190816220812-743ec37842bf
	k8s.io/kubectl => k8s.io/kubectl v0.17.1
	k8s.io/metrics => k8s.io/metrics v0.17.1
	k8s.io/utils => k8s.io/utils v0.0.0-20190801114015-581e00157fb1
)
