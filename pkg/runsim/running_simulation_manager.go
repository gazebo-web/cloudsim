package runsim

import (
	"fmt"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	ignws "gitlab.com/ignitionrobotics/web/cloudsim/pkg/transport/ign"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"sync"
)

// Manager describes a set of methods to handle a set of RunningSimulation and their connections to different websocket servers.
type Manager interface {
	Add(groupID simulations.GroupID, rs *RunningSimulation, t ignws.PubSubWebsocketTransporter) error
	ListExpiredSimulations() []*RunningSimulation
	ListFinishedSimulations() []*RunningSimulation
	GetTransporter(groupID simulations.GroupID) ignws.PubSubWebsocketTransporter
	Free(groupID simulations.GroupID)
	Remove(groupID simulations.GroupID) error
	Exists(groupID simulations.GroupID) bool
	Debug(gid simulations.GroupID) (interface{}, *ign.ErrMsg)
}

// manager is a Manager implementation.
type manager struct {
	runningSimulations map[simulations.GroupID]*RunningSimulation
	lock               sync.RWMutex
}

type debugResponse struct {
	Exists      bool   `json:"exists"`
	IsConnected bool   `json:"is_connected"`
	Message     string `json:"message"`
}

// Debug returns information about a certain websocket connection.
func (m *manager) Debug(gid simulations.GroupID) (interface{}, *ign.ErrMsg) {
	exists := m.Exists(gid)

	var isConnected bool
	var msg string
	if exists {
		t := m.GetTransporter(gid)
		isConnected = t.IsConnected()
		_, bmsg, err := t.Connection().ReadMessage()
		if err != nil {
			msg = "error"
		} else {
			msg = string(bmsg)
		}
	}

	res := debugResponse{
		Exists:      exists,
		IsConnected: isConnected,
		Message:     msg,
	}

	return res, nil
}

// Exists checks if the given group id is registered as a running simulation.
func (m *manager) Exists(groupID simulations.GroupID) bool {
	_, ok := m.runningSimulations[groupID]
	return ok
}

// Free disconnects the websocket client for the given GroupID.
func (m *manager) Free(groupID simulations.GroupID) {
	t := m.GetTransporter(groupID)

	m.lock.Lock()
	defer m.lock.Unlock()

	rs, ok := m.runningSimulations[groupID]
	if !ok {
		return
	}

	rs.publishing = false
	m.runningSimulations[groupID] = rs

	if t != nil && t.IsConnected() {
		t.Disconnect()
	}
}

// Add adds a running simulation and a websocket transport to the given groupID.
func (m *manager) Add(groupID simulations.GroupID, rs *RunningSimulation, t ignws.PubSubWebsocketTransporter) error {
	if rs == nil {
		return fmt.Errorf("invalid running simulation for [%s], it should be not nil", groupID)
	}
	if t == nil {
		return fmt.Errorf("invalid websocket transport for [%s], it should not be nil", groupID)
	}

	m.lock.Lock()
	defer m.lock.Unlock()

	if _, exists := m.runningSimulations[groupID]; exists {
		return fmt.Errorf("running simulation [%s] already exists", groupID)
	}

	rs.Transport = t

	m.runningSimulations[groupID] = rs

	return nil
}

// ListExpiredSimulations list all expired simulations from the list of running simulations.
func (m *manager) ListExpiredSimulations() []*RunningSimulation {
	return m.listByCriteria(func(rs *RunningSimulation) bool {
		return rs.IsExpired()
	})
}

// ListFinishedSimulations list all finished simulations from the list of running simulations.
func (m *manager) ListFinishedSimulations() []*RunningSimulation {
	return m.listByCriteria(func(rs *RunningSimulation) bool {
		return rs.Finished
	})
}

// listByCriteria allows you to list running simulations by a given criteria.
func (m *manager) listByCriteria(criteria func(rs *RunningSimulation) bool) []*RunningSimulation {
	m.lock.RLock()
	defer m.lock.RUnlock()
	rss := make([]*RunningSimulation, 0, len(m.runningSimulations))
	for _, rs := range m.runningSimulations {
		if criteria(rs) {
			rss = append(rss, rs)
		}
	}
	return rss
}

// GetTransporter returns a websocket transporter for the given groupID.
// It returns nil if there's no connection available for the given groupID.
func (m *manager) GetTransporter(groupID simulations.GroupID) ignws.PubSubWebsocketTransporter {
	m.lock.RLock()
	defer m.lock.RUnlock()
	r, ok := m.runningSimulations[groupID]
	if !ok {
		return nil
	}
	return r.Transport
}

// Remove removes a running simulation and its websocket connection.
// If the websocket connection is still active, it will return an error.
func (m *manager) Remove(groupID simulations.GroupID) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	if r, exists := m.runningSimulations[groupID]; !exists || r.Transport.IsConnected() {
		return fmt.Errorf("websocket transport [%s] does not exist or it's still connected to the websocket server", groupID)
	}
	delete(m.runningSimulations, groupID)

	return nil
}

// NewManager initializes a running simulation's manager in charge of communicating to websocket servers.
func NewManager() Manager {
	return &manager{
		runningSimulations: make(map[simulations.GroupID]*RunningSimulation),
	}
}
