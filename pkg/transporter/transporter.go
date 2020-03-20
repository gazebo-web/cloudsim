package transporter

import igntransport "gitlab.com/ignitionrobotics/web/cloudsim/third_party/ign-transport"

type config struct {
	Topic string `env:"IGN_TRANSPORT_TEST_TOPIC" envDefault:"/foo"`
}

type Transporter struct {
	Transport igntransport.GoIgnTransportNode
	Topic string
}

func New() {

}
