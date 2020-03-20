package orchestrator

import (
	"bytes"
	"fmt"
	"k8s.io/client-go/tools/remotecommand"
)

// PodCreateExecErrorMessage creates and returns an error message that includes
// the standard output and standard error of a command executed with KubernetesPodExec
func (kc Kubernetes) PodCreateExecErrorMessage(errorMsg string, options *remotecommand.StreamOptions) string {
	return fmt.Sprintf("%s\n%s\n%s",
		errorMsg,
		fmt.Sprintf("STDOUT dump:\n%s", options.Stdout.(*bytes.Buffer).String()),
		fmt.Sprintf("STDERR dump:\n%s", options.Stderr.(*bytes.Buffer).String()),
	)
}
