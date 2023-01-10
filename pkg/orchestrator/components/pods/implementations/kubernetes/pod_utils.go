package kubernetes

import (
	"errors"
	"fmt"
	"github.com/gazebo-web/cloudsim/v4/pkg/orchestrator/components/spdy"
	"io"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/remotecommand"
	"k8s.io/client-go/util/exec"
)

var (
	// ErrPodExecFailed is returned when a panic is triggered after running a command in a resource.
	ErrPodExecFailed = errors.New("could not run exec command on resource")
)

// runExecInput is the input of runExec.
type runExecInput struct {
	// kubernetes has a reference to the kubernetes client.
	kubernetes kubernetes.Interface
	// namespace is the namespace where the command should be executed.
	namespace string
	// name is the name of the pod where the command should be executed.
	name string
	// container has the name of the container inside a pod that will end up running the command.
	container string
	// command is the actual command that will be executed, it includes the command name and the arguments,
	// each in a new element of the slice.
	command []string
	// options holds information pertaining to the current streaming session:
	// input/output streams, if the client is requesting a TTY, and a terminal size queue to
	// support terminal resizing.
	options remotecommand.StreamOptions
	// spdy is used to initialize a new SPDY executor.
	spdy spdy.Initializer
}

// runExec runs an exec operation inside a resource.
func runExec(input runExecInput) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = ErrPodExecFailed
		}
	}()

	// TODO: Find a way to avoid this line panicking on tests.
	req := input.kubernetes.CoreV1().RESTClient().Post().
		Resource("pods").
		Namespace(input.namespace).
		Name(input.name).
		SubResource("exec")

	scheme := runtime.NewScheme()
	if err := apiv1.AddToScheme(scheme); err != nil {
		return err
	}

	parameterCodec := runtime.NewParameterCodec(scheme)
	req.VersionedParams(&apiv1.PodExecOptions{
		Command:   input.command,
		Container: input.container,
		Stdin:     input.options.Stdin != nil,
		Stdout:    input.options.Stdout != nil,
		Stderr:    input.options.Stderr != nil,
		TTY:       input.options.Tty,
	}, parameterCodec)

	ex, err := input.spdy.NewSPDYExecutor("POST", req.URL())
	if err != nil {
		return err
	}

	err = ex.Stream(input.options)
	if err != nil {
		return err
	}
	return nil
}

// parseExecError parses the errors returned from runExec.
func parseExecError(err error, stdout io.Writer, stderr io.Writer) error {
	execErr, ok := err.(exec.CodeExitError)
	if !ok {
		return err
	}
	msg, err := createExecErrorMessage(
		fmt.Sprintf("Pod exec failed with code %d", execErr.Code),
		stdout,
		stderr,
	)
	if err != nil {
		return err
	}
	return errors.New(msg)
}

// createExecErrorMessage is a helper function to create an error text message for parseExecError.
func createExecErrorMessage(msg string, stdout io.Writer, stderr io.Writer) (string, error) {
	var outDump []byte
	_, err := stdout.Write(outDump)
	if err != nil {
		return "", err
	}

	var errDump []byte
	_, err = stderr.Write(errDump)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(
		"Executing a command inside a resource failed. Error: %s\nSTDOUT: [%s]\nSTDERR: [%s]",
		msg, outDump, errDump,
	), nil
}
