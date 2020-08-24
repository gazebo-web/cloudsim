package application

import "gitlab.com/ignitionrobotics/web/cloudsim/pkg/application/requester"

type requesters struct {
	start    requester.Requester
	launch   requester.Requester
	restart  requester.Requester
	shutdown requester.Requester
}

func (r *requesters) Start() requester.Requester {
	return r.start
}

func (r *requesters) Shutdown() requester.Requester {
	return r.shutdown
}

func (r *requesters) Launch() requester.Requester {
	return r.launch
}

func (r *requesters) Restart() requester.Requester {
	return r.restart
}
