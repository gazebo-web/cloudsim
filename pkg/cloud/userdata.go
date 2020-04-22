package cloud

import (
	"bytes"
	"encoding/base64"
	"strings"
	"text/template"
)

// RunUserDataConfig is the configuration to create an User Data command for EC2 machines initialization.
type RunUserDataConfig struct {
	GroupLabels string
	ExtraLabels string
	JoinCmd     string
}

// NewTemplate creates a new Template from the parsed file called template.gotxt.
// template.gotxt includes the ec2 user data commands to run when the EC2 instance starts.
func NewRunUserDataCommand() *template.Template {
	t := template.Must(template.ParseFiles("template.gotxt"))
	return t
}

// NewRunUserDataConfig creates a RunUserDataConfig to configure a Template.
// It includes the kubeadm join command, and the node labels to set to the Kubelet.
func NewRunUserDataConfig(joinCmd, groupLabels string, extraLabels ...string) RunUserDataConfig {
	return RunUserDataConfig{
		GroupLabels: groupLabels,
		ExtraLabels: strings.Join(extraLabels, ","),
		JoinCmd:     joinCmd,
	}
}

// FillUserDataCommand takes a template and fills it with the given configuration.
func FillUserDataCommand(t *template.Template, config RunUserDataConfig) (base64Data, userData string) {
	b := bytes.NewBuffer(nil)
	t.Execute(b, config)
	userData = b.String()
	base64Data = base64.StdEncoding.EncodeToString(b.Bytes())
	return
}
