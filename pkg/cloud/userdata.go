package cloud

import (
	"io"
	"strings"
	"text/template"
)

type RunUserDataConfig struct {
	GroupLabels string
	ExtraLabels string
	JoinCmd string
}

// NewTemplate creates a new Template from the parsed file called template.gotxt.
// template.gotxt includes the ec2 user data commands to run when the EC2 instance starts.
func NewRunUserDataCommand() *template.Template {
	t := template.Must(template.ParseFiles("template.gotxt"))
	return t
}

// NewRunUserDataConfig creates a RunUserDataConfig to configure a Template.
// It includes the kubeadm join command, and the node labels to set to the Kubelet.
func NewRunUserDataConfig(joinCmd, groupLabels string, extraLabels []string) *RunUserDataConfig {
	return &RunUserDataConfig{
		GroupLabels: groupLabels,
		ExtraLabels: strings.Join(extraLabels, ","),
		JoinCmd:     joinCmd,
	}
}

// FillUserDataCommand takes a template and fills it with the given configuration.
func FillUserDataCommand(t *template.Template, config RunUserDataConfig) string {
	var w io.Writer
	var result []byte
	t.Execute(w, config)
	w.Write(result)
	return string(result)
}