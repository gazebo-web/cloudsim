package transporter

import igntransport "gitlab.com/ignitionrobotics/web/cloudsim/third_party/ign-transport"

type config struct {

}

type Transporter struct {
	Transport igntransport.GoIgnTransportNode
	Topic string
}

func New()  {

}
