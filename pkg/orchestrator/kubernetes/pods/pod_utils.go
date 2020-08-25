package pods

import (
	"errors"
	"fmt"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/kubernetes/spdy"
	"io"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
	"k8s.io/client-go/util/exec"
)

var (
	// ErrPodExecFailed is returned when a panic is triggered after running a command in a resource.
	ErrPodExecFailed = errors.New("could not run exec command on resource")
)

// runExecInput is the input of runExec.
type runExecInput struct {
	kubernetes kubernetes.Interface
	config     *rest.Config
	namespace  string
	name       string
	command    []string
	options    remotecommand.StreamOptions
	spdy       spdy.Initializer
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
		Container: input.name,
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
