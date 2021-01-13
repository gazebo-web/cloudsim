package runsim

import (
	"fmt"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	ignws "gitlab.com/ignitionrobotics/web/cloudsim/transport/ign"
)

type Manager interface {
	Add(groupID simulations.GroupID, rs RunningSimulation, t ignws.PubSubWebsocketTransporter) error
	ListExpiredSimulations() []RunningSimulation
	GetTransporter(groupID simulations.GroupID) ignws.PubSubWebsocketTransporter
	Remove(groupID simulations.GroupID) error
}

type manager struct {
	transporters       map[simulations.GroupID]ignws.PubSubWebsocketTransporter
	runningSimulations map[simulations.GroupID]RunningSimulation
}

func (m *manager) Add(groupID simulations.GroupID, rs RunningSimulation, t ignws.PubSubWebsocketTransporter) error {
	if _, exists := m.transporters[groupID]; exists {
		return fmt.Errorf("websocket transport [%s] already exists", groupID)
	}
	m.transporters[groupID] = t

	if _, exists := m.runningSimulations[groupID]; exists {
		return fmt.Errorf("running simulation [%s] already exists", groupID)
	}
	m.runningSimulations[groupID] = rs

	return nil
}

func (m *manager) ListExpiredSimulations() []RunningSimulation {
	return m.listByCriteria(func(rs RunningSimulation) bool {
		return rs.IsExpired()
	})
}

func (m *manager) ListFinishedSimulations() []RunningSimulation {
	return m.listByCriteria(func(rs RunningSimulation) bool {
		return rs.Finished
	})
}

func (m *manager) listByCriteria(criteria func(rs RunningSimulation) bool) []RunningSimulation {
	rss := make([]RunningSimulation, 0, len(m.runningSimulations))
	for _, rs := range m.runningSimulations {
		if criteria(rs) {
			rss = append(rss, rs)
		}
	}
	return rss
}

func (m *manager) GetTransporter(groupID simulations.GroupID) ignws.PubSubWebsocketTransporter {
	t, ok := m.transporters[groupID]
	if !ok {
		return nil
	}
	return t
}

func (m *manager) Remove(groupID simulations.GroupID) error {
	if _, exists := m.transporters[groupID]; !exists {
		return fmt.Errorf("websocket transport [%s] does not exist", groupID)
	}
	delete(m.transporters, groupID)

	if _, exists := m.runningSimulations[groupID]; !exists {
		return fmt.Errorf("running simulation [%s] does not exists", groupID)
	}
	delete(m.runningSimulations, groupID)

	return nil
}

func NewManager() Manager {
	return &manager{
		transporters:       make(map[simulations.GroupID]ignws.PubSubWebsocketTransporter),
		runningSimulations: make(map[simulations.GroupID]RunningSimulation),
	}
}
