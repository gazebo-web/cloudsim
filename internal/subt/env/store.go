package env

import "gitlab.com/ignitionrobotics/web/cloudsim/pkg/store"

type envStore struct {
	machines store.Machines
}

func (e envStore) Machines() store.Machines {
	return e.machines
}

func NewStore() store.Store {
	return &envStore{
		machines: newMachinesStore(),
	}
}
