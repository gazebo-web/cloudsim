package orchestrator

import (
	"bytes"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/logger"
	"gitlab.com/ignitionrobotics/web/cloudsim/tools"
	"io"
	apiv1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/remotecommand"
	"time"
)

const (
	pollFrequency = 2 * time.Second
)

const (
	podLabelPodGroup       = "pod-group"
	podLabelKeyGroupID     = "cloudsim-group-id"
	cloudsimTagLabelKey    = "cloudsim"
	cloudsimTagLabel       = "cloudsim=true"
	launcherRelaunchNeeded = "relaunch"
)

// Pod wraps the k8s Pod to aggregate two field related to simulations.
type Pod struct {
	apiv1.Pod
	IsRunning bool
	GroupID   string
}

type Pods []Pod

// PodExec creates a command for a specific kubernetes pod. stdin, stdout and stderr io can be defined
// through the options parameter.
func (kc *k8s) PodExec(ctx context.Context, namespace string, podName string, container string, command []string, options *remotecommand.StreamOptions) (opts *remotecommand.StreamOptions, err error) {
	// Handle panics if pod exec cannot be run
	defer func() {
		if r := recover(); r != nil {
			opts = nil
			err = errors.New("could not run exec command on pod")
		}
	}()

	logger.Logger(ctx).Info(fmt.Sprintf("Executing %v in pod [%v]", command, podName))

	// Set default stream options
	if options == nil {
		var stdout, stderr bytes.Buffer
		options = &remotecommand.StreamOptions{
			Stdin:  nil,
			Stdout: &stdout,
			Stderr: &stderr,
			Tty:    false,
		}
	}

	config, err := NewConfig(nil)
	if err != nil {
		return options, err
	}
	req := kc.CoreV1().RESTClient().Post().
		Resource("pods").
		Namespace(namespace).
		Name(podName).
		SubResource("exec")
	scheme := runtime.NewScheme()
	if err := apiv1.AddToScheme(scheme); err != nil {
		return options, err
	}

	parameterCodec := runtime.NewParameterCodec(scheme)
	req.VersionedParams(&apiv1.PodExecOptions{
		Command:   command,
		Container: container,
		Stdin:     options.Stdin != nil,
		Stdout:    options.Stdout != nil,
		Stderr:    options.Stderr != nil,
		TTY:       options.Tty,
	}, parameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
	if err != nil {
		return options, err
	}

	err = exec.Stream(*options)
	if err != nil {
		return options, err
	}

	return options, nil
}

// PodGetLog returns the log from a pod
func (kc *k8s) PodGetLog(ctx context.Context, namespace string, podName string, container string, lines int64) (log *string, err error) {
	// Handle panics if get pod logs cannot be run
	defer func() {
		if r := recover(); r != nil {
			log = nil
			err = errors.New("could not run get logs on pod")
		}
	}()

	logger.Logger(ctx).Info(fmt.Sprintf("Getting logs from pod [%v]", podName))

	// Set default pod log options and get logs from that pod.
	podLogOpts := apiv1.PodLogOptions{
		Container: container,
		TailLines: tools.Int64ptr(lines),
	}
	req := kc.CoreV1().Pods(namespace).GetLogs(podName, &podLogOpts)

	// Reading from logs pod
	var reader io.ReadCloser
	reader, err = req.Stream()
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	// Copy content from reader to buffer
	buffer := new(bytes.Buffer)
	_, err = io.Copy(buffer, reader)

	if err != nil {
		return nil, err
	}

	// Convert buffer to string and get pointer
	log = tools.Sptr(buffer.String())

	return log, nil
}

// GetAllPods returns a set of Pods that match the given label
func (kc *k8s) GetAllPods(label *string) (Pods, error) {
	if label == nil {
		label = tools.Sptr(cloudsimTagLabel)
	}
	list, err := kc.CoreV1().Pods(kc.Namespace()).List(v1.ListOptions{LabelSelector: *label})
	if err != nil {
		return nil, err
	}
	var pods Pods
	for _, p := range list.Items {
		var pod Pod
		pod.GroupID = p.Labels[podLabelKeyGroupID]
		if p.ObjectMeta.DeletionTimestamp != nil {
			pod.IsRunning = false
			continue
		}
		pod.IsRunning = p.Status.Phase == apiv1.PodRunning
		pods = append(pods, pod)
	}
	return pods, nil
}
