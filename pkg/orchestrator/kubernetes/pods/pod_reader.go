package pods

import (
	"bytes"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/kubernetes/spdy"
	"io"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/remotecommand"
)

// reader is a orchestrator.Reader implementation.
type reader struct {
	API      kubernetes.Interface
	pod      orchestrator.Resource
	spdyInit spdy.Initializer
}

// File is used to read a file from the given paths.
func (r *reader) File(paths ...string) (io.Reader, error) {
	var stdout, stderr bytes.Buffer
	options := remotecommand.StreamOptions{
		Stdin:  nil,
		Stdout: &stdout,
		Stderr: &stderr,
		Tty:    false,
	}
	err := runExec(runExecInput{
		kubernetes: r.API,
		namespace:  r.pod.Namespace(),
		name:       r.pod.Name(),
		command:    append([]string{"cat"}, paths...),
		options:    options,
		spdy:       r.spdyInit,
	})
	if err == nil {
		return &stdout, nil
	}
	return nil, parseExecError(err, &stdout, &stderr)
}

// Logs returns the log from the given container running inside the pod.
func (r *reader) Logs(container string, lines int64) (string, error) {
	req := r.API.CoreV1().Pods(r.pod.Namespace()).GetLogs(r.pod.Name(), &apiv1.PodLogOptions{
		Container: container,
		TailLines: &lines,
	})

	reader, err := req.Stream()
	defer reader.Close()
	if err != nil {
		return "", err
	}

	var logs []byte
	_, err = reader.Read(logs)
	if err != nil {
		return "", err
	}

	return string(logs), nil
}

// newReader initializes a new reader.
func newReader(api kubernetes.Interface, pod orchestrator.Resource, spdy spdy.Initializer) orchestrator.Reader {
	return &reader{
		API:      api,
		pod:      pod,
		spdyInit: spdy,
	}
}
