package simulations

import (
	"gitlab.com/ignitionrobotics/web/ign-go"
	"bytes"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"io"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
	"k8s.io/client-go/util/exec"
	podutil "k8s.io/kubernetes/pkg/api/v1/pod"
	"k8s.io/kubernetes/pkg/client/conditions"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Wait ideas taken from : https://github.com/kubernetes/kubernetes/blob/master/test/e2e/framework/util.go

// How often to Poll pods, nodes and claims.
const (
	pollFrequency = 2 * time.Second
)

// Deprecated: GetKubernetesConfig returns the kubernetes config file in the specified path.
// If no path is provided (i.e. nil), then the configuration in ~/.kube/config
// is returned.
func GetKubernetesConfig(kubeconfig *string) (*restclient.Config, error) {
	if kubeconfig == nil {
		kubeconfig = sptr(filepath.Join(os.Getenv("HOME"), ".kube", "config"))
	}
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		return nil, err
	}

	return config, nil
}

// Deprecated: GetKubernetesClient returns a client object to access a kubernetes master.
// Note that this kube client assumes there is a kubernetes configuration in the
// server's ~/.kube/config file. That config is used to connect to the kubernetes
// master.
func GetKubernetesClient() (*kubernetes.Clientset, error) {
	config, err := GetKubernetesConfig(nil)
	if err != nil {
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}

// Deprecated: MakeListOptions returns a ListOptions object for an array of labels.
func MakeListOptions(labels ...string) metav1.ListOptions {
	return metav1.ListOptions{
		LabelSelector: strings.Join(labels, ","),
	}
}

// Deprecated: WaitForPodsReady function waits (blocks) until the identified pods are running and ready, or until there
// is a timeout or error.
func WaitForPodsReady(ctx context.Context, c kubernetes.Interface, namespace string, groupIDLabel string, timeout time.Duration) error {
	opts := metav1.ListOptions{LabelSelector: groupIDLabel}
	return WaitForMatchPodsCondition(ctx, c, namespace, opts, "Ready", timeout, podRunningAndReady)
}

// Deprecated: PodCondition is a function type that returns the pod condition or error by the given Kubernetes Pod.
type PodCondition func(ctx context.Context, pod *apiv1.Pod) (bool, error)

// Deprecated: WaitForMatchPodsCondition finds match pods based on the input ListOptions.
// Waits and checks if all matched pods are in the given PodCondition
var WaitForMatchPodsCondition = func(ctx context.Context, c kubernetes.Interface, namespace string,
	opts metav1.ListOptions, condStr string, timeout time.Duration, condition PodCondition) error {
	logger(ctx).Info(fmt.Sprintf("Waiting up to %v for matching pods' status to be %s", timeout, condStr))
	for start := time.Now(); time.Since(start) < timeout; Sleep(pollFrequency) {
		pods, err := c.CoreV1().Pods(namespace).List(opts)
		if err != nil {
			return err
		}
		var conditionNotMatch []string
		for _, pod := range pods.Items {
			done, err := condition(ctx, &pod)
			if done && err != nil {
				return fmt.Errorf("unexpected error: %v", err)
			}
			if !done {
				conditionNotMatch = append(conditionNotMatch, pod.Name)
			}
		}
		if len(conditionNotMatch) <= 0 {
			logger(ctx).Info(fmt.Sprintf("Success. Pods match condition '%s'", condStr))
			return err
		}
		logger(ctx).Debug(fmt.Sprintf("%d pods are not %s: %v", len(conditionNotMatch), condStr, conditionNotMatch))
	}
	return errors.Errorf("gave up waiting for matching pods to be '%s' after %v", condStr, timeout)
}

// Deprecated: podRunningAndReady checks if a pod by name is running. This function is used
// for Wait polls.
func podRunningAndReady(ctx context.Context, pod *apiv1.Pod) (bool, error) {
	switch pod.Status.Phase {
	case apiv1.PodFailed, apiv1.PodSucceeded:
		return false, conditions.ErrPodCompleted
	case apiv1.PodRunning:
		return podutil.IsPodReady(pod), nil
	}
	return false, nil
}

// Deprecated: WaitForNodesReady function waits (blocks) until the identified nodes are ready, or until there
// is a timeout or error.
func WaitForNodesReady(ctx context.Context, c kubernetes.Interface, namespace string, groupIDLabel string, timeout time.Duration) error {
	opts := metav1.ListOptions{LabelSelector: groupIDLabel}
	return WaitForMatchNodesCondition(ctx, c, namespace, opts, timeout)
}

// Deprecated: WaitForMatchNodesCondition finds match Nodes based on the input ListOptions.
// Waits and checks if all matched nodes are in the given PodCondition
func WaitForMatchNodesCondition(ctx context.Context, c kubernetes.Interface, namespace string, opts metav1.ListOptions, timeout time.Duration) error {
	logger(ctx).Info(fmt.Sprintf("Waiting up to %v for match nodes to be ready", timeout))

	maxAllowedNotReadyNodes := 0
	var notReady []*apiv1.Node
	err := wait.PollImmediate(pollFrequency, timeout, func() (bool, error) {
		notReady = nil
		// It should be OK to list unschedulable Nodes here.
		nodes, err := c.CoreV1().Nodes().List(opts)
		if err != nil {
			return false, err
		}
		logger(ctx).Debug(fmt.Sprintf("Found nodes %v", nodes.Items))
		for i := range nodes.Items {
			node := &nodes.Items[i]
			if !isNodeConditionSetAsExpected(ctx, node, apiv1.NodeReady, true) {
				notReady = append(notReady, node)
			}
		}
		return len(notReady) <= maxAllowedNotReadyNodes, nil
	})

	if err != nil && err != wait.ErrWaitTimeout {
		return err
	}

	if len(notReady) > maxAllowedNotReadyNodes {
		msg := ""
		for _, node := range notReady {
			msg = fmt.Sprintf("%s, %s", msg, node.Name)
		}
		return errors.Errorf("Nodes not ready: %#v", msg)
	}
	return nil
}

// Deprecated: isNodeConditionSetAsExpected checks if the given condition is met by a node.
func isNodeConditionSetAsExpected(ctx context.Context, node *apiv1.Node, conditionType apiv1.NodeConditionType, wantTrue bool) bool {
	// Check the node readiness condition (logging all).
	for _, cond := range node.Status.Conditions {
		// Ensure that the condition type and the status matches as desired.
		if cond.Type == conditionType {
			if cond.Type == apiv1.NodeReady {
				// For NodeReady we should check if Taints are gone as well
				// See taint.MatchTaint
				if wantTrue {
					if cond.Status == apiv1.ConditionTrue {
						return true
					}
					msg := fmt.Sprintf("Condition %s of node %s is %v instead of %t. Reason: %v, message: %v",
						conditionType, node.Name, cond.Status == apiv1.ConditionTrue, wantTrue, cond.Reason, cond.Message)
					logger(ctx).Debug(msg)
					return false

				}
				// TODO: check if the Node is tainted once we enable NC notReady/unreachable taints by default
				if cond.Status != apiv1.ConditionTrue {
					return true
				}
				logger(ctx).Debug(fmt.Sprintf("Condition %s of node %s is %v instead of %t. Reason: %v, message: %v",
					conditionType, node.Name, cond.Status == apiv1.ConditionTrue, wantTrue, cond.Reason, cond.Message))
				return false

			}
			if (wantTrue && (cond.Status == apiv1.ConditionTrue)) || (!wantTrue && (cond.Status != apiv1.ConditionTrue)) {
				return true
			}
			logger(ctx).Debug(fmt.Sprintf("Condition %s of node %s is %v instead of %t. Reason: %v, message: %v",
				conditionType, node.Name, cond.Status == apiv1.ConditionTrue, wantTrue, cond.Reason, cond.Message))
			return false
		}
	}

	logger(ctx).Debug(fmt.Sprintf("Couldn't find condition %v on node %v", conditionType, node.Name))
	return false
}

// Deprecated: createKubernetesPodExecErrorMsg creates and returns an error message that includes
// the standard output and standard error of a command executed with KubernetesPodExec
func createKubernetesPodExecErrorMsg(errorMsg string, options *remotecommand.StreamOptions) string {
	return fmt.Sprintf("%s\n%s\n%s",
		errorMsg,
		fmt.Sprintf("STDOUT dump:\n%s", options.Stdout.(*bytes.Buffer).String()),
		fmt.Sprintf("STDERR dump:\n%s", options.Stderr.(*bytes.Buffer).String()),
	)
}

// Deprecated: KubernetesPodExec creates a command for a specific kubernetes pod. stdin, stdout and stderr io can be defined
// through the options parameter.
func KubernetesPodExec(ctx context.Context, kc kubernetes.Interface, namespace string, podName string, container string,
	command []string, options *remotecommand.StreamOptions) (opts *remotecommand.StreamOptions, err error) {

	// Handle panics if pod exec cannot be run
	defer func() {
		if r := recover(); r != nil {
			opts = nil
			err = errors.New("could not run exec command on pod")
		}
	}()

	logger(ctx).Info(fmt.Sprintf("Executing %v in pod [%v]", command, podName))

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

	config, err := GetKubernetesConfig(nil)
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

// Deprecated: KubernetesPodReadFile reads and returns the contents of a file inside a pod.
func KubernetesPodReadFile(ctx context.Context, kc kubernetes.Interface, namespace string, podName string,
	container string, paths ...string) (*bytes.Buffer, error) {

	command := append([]string{"cat"}, paths...)

	options, err := KubernetesPodExec(ctx, kc, namespace, podName, container, command, nil)
	if err != nil {
		return nil, err
	}

	return options.Stdout.(*bytes.Buffer), nil
}

// Deprecated: KubernetesPodGetLog returns the log of a pod
func KubernetesPodGetLog(ctx context.Context, kc kubernetes.Interface, namespace string, podName string,
	container string, lines int64) (log *string, err error) {

	// Handle panics if get pod logs cannot be run
	defer func() {
		if r := recover(); r != nil {
			log = nil
			err = errors.New("could not run get logs on pod")
		}
	}()

	logger(ctx).Info(fmt.Sprintf("Getting logs from pod [%v]", podName))

	// Set default pod log options and get logs from that pod.
	podLogOpts := apiv1.PodLogOptions{
		Container: container,
		TailLines: int64ptr(lines),
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
	log = sptr(buffer.String())

	return log, nil
}

// Deprecated: KubernetesPodSendS3CopyCommand sends a command to a pod to upload a file to S3.
// The pod that receives the command must have `aws` and `tar` installed for this to work,
// and the container running the commands must have AWS env vars configured.
// If `target` is a directory, its contents are `tar`'d and `gzipped` before being uploaded.
// The `target` is sent to the specified S3 `bucket` with name `filename`.
func KubernetesPodSendS3CopyCommand(ctx context.Context, kc kubernetes.Interface, namespace string, podName string,
	container string, bucket string, target string, filename string) error {

	// Prepare the script
	scriptParams := struct {
		Target   string
		Filename string
		Bucket   string
	}{
		Target:   target,
		Filename: filename,
		Bucket:   PrepareS3Address(bucket, filename),
	}
	script, err := ign.ParseTemplate("simulations/scripts/copy_to_s3.sh", scriptParams)
	if err != nil {
		logger(ctx).Error("Could not prepare S3 copy script.", err)
		return err
	}

	command := []string{"sh", "-c", script}

	var options *remotecommand.StreamOptions
	if options, err = KubernetesPodExec(ctx, kc, namespace, podName, container, command, nil); err != nil {
		if execErr, ok := err.(exec.CodeExitError); ok && options != nil {
			msg := createKubernetesPodExecErrorMsg(
				fmt.Sprintf("Pod exec failed with code %d", execErr.Code),
				options,
			)
			logger(ctx).Error(msg, err)
		}
		return err
	}

	return nil
}
