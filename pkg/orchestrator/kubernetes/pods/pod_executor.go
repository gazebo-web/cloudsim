package pods

import (
	"bytes"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/kubernetes/spdy"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/remotecommand"
)

// executor is a orchestrator.Executor implementation.
type executor struct {
	API      kubernetes.Interface
	pod      orchestrator.Resource
	spdyInit spdy.Initializer
}

// Cmd is used to run a command in a container inside a pod.
func (e executor) Cmd(command []string) error {
	var stdout, stderr bytes.Buffer
	options := remotecommand.StreamOptions{
		Stdin:  nil,
		Stdout: &stdout,
		Stderr: &stderr,
		Tty:    false,
	}
	err := runExec(runExecInput{
		kubernetes: e.API,
		namespace:  e.pod.Namespace(),
		name:       e.pod.Name(),
		command:    command,
		options:    options,
		spdy:       e.spdyInit,
	})
	if err == nil {
		return nil
	}
	return parseExecError(err, &stdout, &stderr)
}

// Script is used to run a bash script inside a container.
func (e executor) Script(script string) error {
	return e.Cmd([]string{"sh", "-c", script})
}

// newExecutor initializes a new executor.
func newExecutor(api kubernetes.Interface, pod orchestrator.Resource, spdyInit spdy.Initializer) orchestrator.Executor {
	return &executor{
		API:      api,
		pod:      pod,
		spdyInit: spdyInit,
	}
}
