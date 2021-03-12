package scripts

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
	"text/template"
)

func TestUserDataScript(t *testing.T) {
	tmpl, err := template.ParseFiles("ec2_user_data.sh")
	require.NoError(t, err)

	var b []byte
	buffer := bytes.NewBuffer(b)

	err = tmpl.Execute(buffer, map[string]interface{}{
		"Labels":      "app=test,example=test",
		"ClusterName": "testing-cluster-name",
		"Args":        "--use-max-pods false",
	})
	require.NoError(t, err)

	fmt.Println("Script:", buffer.String())
}
