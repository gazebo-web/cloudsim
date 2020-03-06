package main

import (
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	k8testing "k8s.io/client-go/testing"
	"log"
	"time"
)

type K8Mock struct {
	kubernetes.Interface
}

// ChainedMockFunction is a generic function used to create Kubernetes Reactor
// chains. Reactor chains allows modifying resources in response to events.
// ChainedMockFunction differs from MockFunction in that it returns a `handled`
// value, which lets the reactor chain short circuit if necessary.
type ChainedMockFunction func(args ...interface{}) (handled bool, res interface{})

// PodMutator mutates a pod as a Reaction.
type PodMutator func(pod *corev1.Pod) (handled bool, err error)

// NewK8Mock creates a new K8Mock.
func NewK8Mock(ctx context.Context) *K8Mock {
	kcli := fake.NewSimpleClientset()
	f := &K8Mock{kcli}
	return f.AddLocalNodes(3)
}

// Reset (re)initializes the mock, clearing all reactors and mutators.
func (f *K8Mock) Reset() *K8Mock {
	f.Interface = fake.NewSimpleClientset()
	return f
}

func (f *K8Mock) getSimpleClientset() *fake.Clientset {
	return f.Interface.(*fake.Clientset)
}

// AddNode adds a given Node to the mock
func (f *K8Mock) AddNode(node *corev1.Node) *K8Mock {
	node, err := f.getSimpleClientset().CoreV1().Nodes().Create(node)
	if err != nil {
		log.Fatalf("Error creating a test Node. %+v\n", err)
	}
	return f
}

// AddLocalNodes is an example helper function that adds 3 Nodes to be used with
// the local_machines.go.
func (f *K8Mock) AddLocalNodes(n int) *K8Mock {
	// Create Nodes with Ready Status, and with the labels expected
	// by the LocalNodes Client (eg. minikube/dind).
	for i := 1; i <= n; i++ {
		node := f.MakeNode(
			fmt.Sprintf("fakeNode-%d", i),
			map[string]string{
				"cloudsim_free_node": "true",
			},
		)
		f.AddNode(node)
	}
	return f
}

// MakeNode is a helper function that creates a new K8 node with the given name
// and labels. The returned node will have 'Ready' status. The returned node
// must be registered by calling `AddNode` with it.
func (f *K8Mock) MakeNode(name string, labels map[string]string) *corev1.Node {
	fakeNow := metav1.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC)
	readyCondition := corev1.NodeCondition{
		Type:               corev1.NodeReady,
		Status:             corev1.ConditionTrue,
		LastHeartbeatTime:  fakeNow,
		LastTransitionTime: fakeNow,
	}

	return &corev1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name:   name,
			Labels: labels,
		},
		Status: corev1.NodeStatus{
			Conditions: []corev1.NodeCondition{
				readyCondition,
			},
		},
	}
}

// MakePod is a helper function that creates a new K8 Pod with the given name,
// image and labels.
func (f *K8Mock) MakePod(name, image string, labels map[string]string) *corev1.Pod {
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:   name,
			Labels: labels,
		},
		Spec: corev1.PodSpec{
			RestartPolicy: corev1.RestartPolicyNever,
			Containers: []corev1.Container{
				{Image: image},
			},
		},
	}
}

// AddReactor adds a mock function that will be executed each time a given
// Verb (create, delete, update, list) is called on a given Resource
// (pods, nodes, etc). "*" can be passed to verbs and resources to be called in
// all cases.
// For a complete list of available verbs and resources, run `kubectl api-resources -o wide`
func (f *K8Mock) AddReactor(verb, resource string, fn ChainedMockFunction) *K8Mock {
	reactor := func(action k8testing.Action) (handled bool, ret runtime.Object, err error) {
		handled, res := fn(action)
		// If the mock result is an error, return that error
		if err, ok := res.(error); ok {
			return handled, nil, err
		}
		ret = res.(runtime.Object)
		return handled, ret, nil
	}
	f.getSimpleClientset().Fake.PrependReactor(verb, resource, reactor)
	return f
}

// AddPodCreationMutator adds a mock function that will be executed each time a
// given Verb (create, delete, update, list) is called on a pod.
// "*" can be passed to `verb` to be called in all cases.
// For a complete list of available verbs and resources, run `kubectl api-resources -o wide`
func (f *K8Mock) AddPodCreationMutator(podFn PodMutator) *K8Mock {
	reactor := func(action k8testing.Action) (handled bool, ret runtime.Object, err error) {
		pod := action.(k8testing.CreateActionImpl).GetObject().(*corev1.Pod)
		handled, err = podFn(pod)
		return
	}
	f.getSimpleClientset().Fake.PrependReactor("create", "pods", reactor)
	return f
}

// SetFixedMutatorsForPodCreation configures the mock to run a single specific
// pod mutator for each new Create Pod invocation. Mutators are provided in an
// array that is executed sequentially; once per Create Pod invocation. This
// will panic if more pods are created than total number of mutators provided.
func (f *K8Mock) SetFixedMutatorsForPodCreation(podFn ...PodMutator) *K8Mock {
	count := 0
	podUpdaterReactor := func(pod *corev1.Pod) (handled bool, err error) {
		fn := podFn[count]
		count++
		handled, err = fn(pod)
		return
	}
	return f.AddPodCreationMutator(podUpdaterReactor)
}

// AttachStatusToNewPods instructs the mock to attach a copy of the given PodStatus
// to any newly created Pod. It is typically used to allow the code to create any
// Pod and have them set to Ready automatically.
func (f *K8Mock) AttachStatusToNewPods(status corev1.PodStatus) *K8Mock {
	podFn := func(pod *corev1.Pod) (bool, error) {
		pod.Status = status
		// Set 'handled' to false to allow other reactors to also act on the action.
		return false, nil
	}
	return f.AddPodCreationMutator(podFn)
}

// MakePodStatusMutator is a helper function that will return a mutator function
// that sets a pod status.
func MakePodStatusMutator(status corev1.PodStatus) PodMutator {
	podFn := func(pod *corev1.Pod) (bool, error) {
		pod.Status = status
		// Set 'handled' to false to allow other reactors to also act on the action.
		return false, nil
	}
	return podFn
}

// MakePodStatus creates a PodStatus object. This function can be used to
// configure pod statuses for testing purposes. You'll typically want to call
// one of the functions that return a predefined status (e.g. MakePodReadyStatus,
// MakePodFailedStatus, etc.) instead of creating it from scratch using this
// function.
func MakePodStatus(phase corev1.PodPhase, reason string, conditions ...corev1.PodCondition) corev1.PodStatus {
	return corev1.PodStatus{
		Phase:      phase,
		Reason:     reason,
		PodIP:      "1.1.1.1",
		Conditions: conditions,
	}
}

// MakePodReadyStatus returns a PodStatus for a pod that is in a Running phase.
// Pods in this state were launched with no problems.
func MakePodReadyStatus() corev1.PodStatus {
	return MakePodStatus(
		corev1.PodRunning,
		"Running",
		corev1.PodCondition{
			Type:   corev1.PodInitialized,
			Status: corev1.ConditionTrue,
		},
		corev1.PodCondition{
			Type:   corev1.PodScheduled,
			Status: corev1.ConditionTrue,
		},
		corev1.PodCondition{
			Type:   corev1.PodReady,
			Status: corev1.ConditionTrue,
		},
	)
}

// MakePodFailedStatus returns a PodStatus for a pod that is in a Failed phase.
// Pods can get to this state for a number of reasons, but will typically reach
// this state due to containers with errors.
func MakePodFailedStatus() corev1.PodStatus {
	return MakePodStatus(
		corev1.PodFailed,
		"Failed",
		corev1.PodCondition{
			Type:   corev1.PodInitialized,
			Status: corev1.ConditionTrue,
		},
		corev1.PodCondition{
			Type:   corev1.PodScheduled,
			Status: corev1.ConditionTrue,
		},
		corev1.PodCondition{
			Type:   corev1.PodReady,
			Status: corev1.ConditionFalse,
		},
	)
}

// MakePodUnschedulableStatus returns a PodStatus for a pod that is stuck in a
// Pending phase due to having no nodes available.
func MakePodUnschedulableStatus() corev1.PodStatus {
	return MakePodStatus(
		corev1.PodPending,
		"Pending",
		corev1.PodCondition{
			Type:   corev1.PodInitialized,
			Status: corev1.ConditionTrue,
		},
		corev1.PodCondition{
			Type:   corev1.PodScheduled,
			Status: corev1.ConditionFalse,
		},
		corev1.PodCondition{
			Type:   corev1.PodReady,
			Status: corev1.ConditionFalse,
		},
	)
}

// MakePodImagePullBackoffStatus returns a PodStatus for a pod that failed to
// pull an image. Pods in this status are usually retried automatically by
// Kubernetes, but there are some cases where pods can get stuck in this state
// indefinitely.
func MakePodImagePullBackoffStatus() corev1.PodStatus {
	return MakePodStatus(
		corev1.PodPending,
		"ImagePullBackOff",
		corev1.PodCondition{
			Type:   corev1.PodInitialized,
			Status: corev1.ConditionTrue,
		},
		corev1.PodCondition{
			Type:   corev1.PodScheduled,
			Status: corev1.ConditionTrue,
		},
		corev1.PodCondition{
			Type:   corev1.PodReady,
			Status: corev1.ConditionFalse,
		},
	)
}

// MakePodErrImagePullStatus returns a PodStatus for a pod that is in a
// Pending phase due to failing to pull a container image. Pods in this status
// are usually retried automatically by Kubernetes, but there are some cases
// where pods can get stuck in this state indefinitely.
func MakePodErrImagePullStatus() corev1.PodStatus {
	return MakePodStatus(
		corev1.PodPending,
		"ErrImagePull",
		corev1.PodCondition{
			Type:   corev1.PodInitialized,
			Status: corev1.ConditionTrue,
		},
		corev1.PodCondition{
			Type:   corev1.PodScheduled,
			Status: corev1.ConditionTrue,
		},
		corev1.PodCondition{
			Type:   corev1.PodReady,
			Status: corev1.ConditionFalse,
		},
	)
}
