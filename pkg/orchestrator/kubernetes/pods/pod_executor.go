package pods

import (
	"bytes"
	"fmt"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/kubernetes/spdy"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/remotecommand"
)

// executor is a orchestrator.Executor implementation.
type executor struct {
	API      kubernetes.Interface
	pod      orchestrator.Resource
	spdyInit spdy.Initializer
	logger   ign.Logger
}

// Cmd is used to run a command in a container inside a resource.
func (e *executor) Cmd(command []string) error {
	e.logger.Debug(fmt.Sprintf("Running command [%s] on pod [%s]", command, e.pod.Name()))

	// Prepare buffers
	var stdout, stderr bytes.Buffer

	// Prepare options for SPDY
	options := remotecommand.StreamOptions{
		Stdin:  nil,
		Stdout: &stdout,
		Stderr: &stderr,
		Tty:    false,
	}

	// Run command
	err := runExec(runExecInput{
		kubernetes: e.API,
		namespace:  e.pod.Namespace(),
		name:       e.pod.Name(),
		command:    command,
		options:    options,
		spdy:       e.spdyInit,
	})

	if err == nil {
		e.logger.Debug(fmt.Sprintf("Command [%s] on pod [%s] sucessfully run.", command, e.pod.Name()))
		return nil
	}
	err = parseExecError(err, &stdout, &stderr)

	e.logger.Debug(fmt.Sprintf("Running ommand [%s] on pod [%s] failed. Error: %s", command, e.pod.Name(), err.Error()))
	return err
}

// Script is used to run a bash script inside a container.
func (e *executor) Script(script string) error {
	e.logger.Debug(fmt.Sprintf("Running script [%s] on pod [%s]", script, e.pod.Name()))
	return e.Cmd([]string{"sh", "-c", script})
}

// newExecutor initializes a new executor.
func newExecutor(api kubernetes.Interface, pod orchestrator.Resource, spdyInit spdy.Initializer, logger ign.Logger) orchestrator.Executor {
	return &executor{
		API:      api,
		pod:      pod,
		spdyInit: spdyInit,
		logger:   logger,
	}
}
