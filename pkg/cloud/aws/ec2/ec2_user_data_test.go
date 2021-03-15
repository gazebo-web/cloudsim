package ec2

import (
	"bytes"
	"github.com/stretchr/testify/assert"
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

	const expected = `#!/bin/bash
set -x
exec > >(tee /var/log/user-data.log|logger -t user-data ) 2>&1
echo BEGIN
date '+%Y-%m-%d %H:%M:%S'
cat > /etc/systemd/system/kubelet.service.d/20-labels-taints.conf <<EOF
[Service]
Environment="KUBELET_EXTRA_ARGS=--node-labels=app=test,example=test"
EOF
set -o xtrace
/etc/eks/bootstrap.sh testing-cluster-name --use-max-pods false
`

	assert.Equal(t, expected, buffer.String())
}
