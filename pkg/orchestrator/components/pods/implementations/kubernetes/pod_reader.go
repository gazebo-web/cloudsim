package kubernetes

import (
	"bytes"
	"fmt"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/components/pods"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/components/spdy"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/resource"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"io"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/remotecommand"
)

// reader is a pods.Reader implementation.
type reader struct {
	API      kubernetes.Interface
	pod      resource.Resource
	spdyInit spdy.Initializer
	logger   ign.Logger
}

// File is used to read a file from the given paths.
func (r *reader) File(paths ...string) (io.Reader, error) {
	r.logger.Debug(fmt.Sprintf("Reading file from paths [%+v] on pod [%s]", paths, r.pod.Name()))

	// Prepare buffers
	var stdout, stderr bytes.Buffer

	// Prepare options
	options := remotecommand.StreamOptions{
		Stdin:  nil,
		Stdout: &stdout,
		Stderr: &stderr,
		Tty:    false,
	}

	// Run command
	err := runExec(runExecInput{
		kubernetes: r.API,
		namespace:  r.pod.Namespace(),
		name:       r.pod.Name(),
		command:    append([]string{"cat"}, paths...),
		options:    options,
		spdy:       r.spdyInit,
	})
	if err == nil {
		r.logger.Debug(fmt.Sprintf("Reading file from paths [%+v] on pod [%s] succeeded. Output: %s", paths, r.pod.Name(), stdout.String()))
		return &stdout, nil
	}

	err = parseExecError(err, &stdout, &stderr)

	r.logger.Debug(fmt.Sprintf("Reading file from paths [%+v] on pod [%s] failed. Error: %s", paths, r.pod.Name(), err.Error()))
	return nil, err
}

// Logs returns the log from the given container running inside the resource.
func (r *reader) Logs(container string, lines int64) (string, error) {
	r.logger.Debug(fmt.Sprintf("Reading logs from container [%s] in pod [%s]", container, r.pod.Name()))

	// Prepare request to get logs
	req := r.API.CoreV1().Pods(r.pod.Namespace()).GetLogs(r.pod.Name(), &apiv1.PodLogOptions{
		Container: container,
		TailLines: &lines,
	})

	// Open data stream
	re, err := req.Stream()
	if err != nil {
		r.logger.Debug(fmt.Sprintf("Reading logs from container [%s] in pod [%s] failed. Error: %s", container, r.pod.Name(), err.Error()))
		return "", err
	}
	defer re.Close()

	// Read logs
	var logs []byte
	_, err = re.Read(logs)

	if err != nil {
		r.logger.Debug(fmt.Sprintf("Reading logs from container [%s] in pod [%s] failed. Error: %s", container, r.pod.Name(), err.Error()))
		return "", err
	}

	r.logger.Debug(fmt.Sprintf("Reading logs from container [%s] in pod [%s] succeded. Output: %s", container, r.pod.Name(), string(logs)))
	return string(logs), nil
}

// newReader initializes a new reader.
func newReader(api kubernetes.Interface, pod resource.Resource, spdy spdy.Initializer, logger ign.Logger) pods.Reader {
	return &reader{
		API:      api,
		pod:      pod,
		spdyInit: spdy,
		logger:   logger,
	}
}
