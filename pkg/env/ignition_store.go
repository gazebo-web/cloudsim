package env

import "gitlab.com/ignitionrobotics/web/cloudsim/pkg/store"

type ignition struct {
}

func (i ignition) IP() string {
	panic("implement me")
}

func newIgnitionStore() store.Ignition {
	return &ignition{}
}
