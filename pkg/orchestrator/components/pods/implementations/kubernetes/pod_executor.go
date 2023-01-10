package kubernetes

import (
	"bytes"
	"fmt"
	"github.com/gazebo-web/cloudsim/v4/pkg/orchestrator/components/pods"
	"github.com/gazebo-web/cloudsim/v4/pkg/orchestrator/components/spdy"
	"github.com/gazebo-web/cloudsim/v4/pkg/orchestrator/resource"
	"github.com/gazebo-web/gz-go/v7"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/remotecommand"
)

// executor is a pods.Executor implementation.
type executor struct {
	API      kubernetes.Interface
	pod      resource.Resource
	spdyInit spdy.Initializer
	logger   gz.Logger
}

// Cmd is used to run a command in a container inside a resource.
func (e *executor) Cmd(container string, command []string) error {
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
		container:  container,
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
func (e *executor) Script(container, script string) error {
	e.logger.Debug(fmt.Sprintf("Running script [%s] on pod [%s]", script, e.pod.Name()))
	return e.Cmd(container, []string{"sh", "-c", script})
}

// newExecutor initializes a new executor.
func newExecutor(api kubernetes.Interface, pod resource.Resource, spdyInit spdy.Initializer, logger gz.Logger) pods.Executor {
	return &executor{
		API:      api,
		pod:      pod,
		spdyInit: spdyInit,
		logger:   logger,
	}
}
