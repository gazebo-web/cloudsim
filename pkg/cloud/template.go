package cloud

import (
	"io"
	"strings"
	"text/template"
)

type TemplateConfig struct {
	GroupLabels string
	ExtraLabels string
	JoinCmd string
}

// NewTemplate creates a new Template from the parsed file called template.gotxt.
// template.gotxt includes the ec2 user data commands to run when the EC2 instance starts.
func NewTemplate() *template.Template {
	t := template.Must(template.ParseFiles("template.gotxt"))
	return t
}

// NewTemplateConfig creates a TemplateConfig to configure a Template.
// It includes the kubeadm join command, and the node labels to set to the Kubelet.
func NewTemplateConfig(joinCmd, groupLabels string, extraLabels []string) *TemplateConfig {
	return &TemplateConfig{
		GroupLabels: groupLabels,
		ExtraLabels: strings.Join(extraLabels, ","),
		JoinCmd:     joinCmd,
	}
}

// FillTemplate takes the template and fills it with the given configuration.
func FillTemplate(t *template.Template, config TemplateConfig) string {
	var w io.Writer
	var result []byte
	t.Execute(w, config)
	w.Write(result)
	return string(result)
}