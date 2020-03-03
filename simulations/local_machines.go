package simulations

import (
	"gitlab.com/ignitionrobotics/web/ign-go"
	"context"
	"fmt"
	"github.com/caarlos0/env"
	"github.com/jinzhu/gorm"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// This module works with a local kubernetes (eg. minikube or dind-cluster) with a predefined set
// of machines / nodes.
// It is local implementation of the NodeManager interface. And an alternative to ec2_machines.go.

// LocalNodes is a client to interact with a local k8 cluster
type LocalNodes struct {
	// A reference to the kubernetes client
	clientset kubernetes.Interface
}

const (
	nodeLabelKey     string = "cloudsim_groupid"
	freeNodeLabelKey string = "cloudsim_free_node"
)

// NewLocalNodesClient creates a client to interact with a local k8 cluster and
// set of machines.
func NewLocalNodesClient(ctx context.Context, kcli kubernetes.Interface) (*LocalNodes, error) {

	logger(ctx).Info("Creating local Nodes Manager")
	l := LocalNodes{}
	if err := env.Parse(&l); err != nil {
		return nil, err
	}

	l.clientset = kcli

	return &l, nil
}

// CloudMachinesList returns a paginated list with all cloud machines. In the local impl
// we just return an empty set (for now).
// @public
func (s *LocalNodes) CloudMachinesList(ctx context.Context, p *ign.PaginationRequest,
	tx *gorm.DB, byStatus *MachineStatus, invertStatus bool, groupID *string, application *string) (*MachineInstances, *ign.PaginationResult, *ign.ErrMsg) {
	// Create the DB query
	var machines MachineInstances

	q := tx.Model(&MachineInstance{})
	// Force an empty result
	q = q.Where("group_id = ?", "-1")

	pagination, err := ign.PaginateQuery(q, &machines, *p)
	if err != nil {
		return nil, nil, ign.NewErrorMessageWithBase(ign.ErrorInvalidPaginationRequest, err)
	}
	if !pagination.PageFound {
		return nil, nil, ign.NewErrorMessage(ign.ErrorPaginationPageNotFound)
	}

	return &machines, pagination, nil
}

// launchNodes will try to re-use existing nodes (if not already used) for the given groupID.
// It will return an error if all nodes are currently being used.
// Returns the node labels that can be used to identify the chosen nodes.
// @public
func (s *LocalNodes) launchNodes(ctx context.Context, tx *gorm.DB, dep *SimulationDeployment) (*string, *ign.ErrMsg) {
	ignlog := logger(ctx)
	nodesInterface := s.clientset.CoreV1().Nodes()
	var err error

	groupID := *dep.GroupId

	// Find a free node or the one already used by same groupID
	// First, try to use same node
	sameNodeLabel := nodeLabelKey + "=" + groupID
	sameNode := true
	nodes, _ := nodesInterface.List(metav1.ListOptions{LabelSelector: sameNodeLabel})
	if len(nodes.Items) == 0 {
		sameNode = false
		// Now try using any free node
		freeNodeLabel := freeNodeLabelKey + "!=false"
		nodes, err = nodesInterface.List(metav1.ListOptions{LabelSelector: freeNodeLabel})
		if err != nil || len(nodes.Items) == 0 {
			return nil, NewErrorMessageWithBase(ErrorFreeNodeNotFound, err)
		}
	}

	// Update its labels (mark it as being used)
	if !sameNode {
		node := nodes.Items[0]
		node.ObjectMeta.Labels[freeNodeLabelKey] = "false"
		node.ObjectMeta.Labels[nodeLabelKey] = groupID
		_, err = nodesInterface.Update(&node)
		if err != nil {
			return nil, NewErrorMessageWithBase(ErrorMarkingLocalNodeAsUsed, err)
		}
	}

	ignlog.Info(fmt.Sprintf("Configured local node for Cloudsim GroupId: %s\n", groupID))
	nodeLabel := nodeLabelKey + "=" + groupID
	return &nodeLabel, nil
}

// deleteK8Nodes deletes the kubernetes nodes used to run containers.
// It is expected that if the labeled Node cannot be found, this func should return
// an ErrorLabeledNodeNotFound.
// @public
func (s *LocalNodes) deleteK8Nodes(ctx context.Context, tx *gorm.DB, groupID string) (interface{}, *ign.ErrMsg) {

	ignlog := logger(ctx)
	nodesInterface := s.clientset.CoreV1().Nodes()
	nodeLabel := nodeLabelKey + "=" + groupID

	// Find the nodes
	nodes, err := nodesInterface.List(metav1.ListOptions{LabelSelector: nodeLabel})
	if err != nil || len(nodes.Items) == 0 {
		logger(ctx).Info("Nodes not found for the groupID: " + groupID)
		return nil, NewErrorMessageWithBase(ErrorLabeledNodeNotFound, err)
	}

	// Update its labels (mark nodes as free)
	node := nodes.Items[0]
	node.ObjectMeta.Labels[freeNodeLabelKey] = "true"
	node.ObjectMeta.Labels[nodeLabelKey] = ""
	_, err = nodesInterface.Update(&node)
	if err != nil {
		return nil, NewErrorMessageWithBase(ErrorMarkingLocalNodeAsFree, err)
	}

	ignlog.Info(fmt.Sprintf("Stopped using local node for Cloudsim GroupId: %s\n", groupID))
	return &nodeLabel, nil
}

// deleteHosts is a helper function to remove used Hosts/Instances.
// It also updates the MachineInstance DB records with the status of the terminated instances.
// @public
func (s *LocalNodes) deleteHosts(ctx context.Context, tx *gorm.DB,
	dep *SimulationDeployment) (interface{}, *ign.ErrMsg) {
	return nil, nil
}
