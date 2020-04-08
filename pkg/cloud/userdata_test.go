package cloud

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewRunUserDataCommand_NotPanic(t *testing.T) {
	assert.NotPanics(t, func() {
		NewRunUserDataCommand()
	})
}

func TestNewRunUserDataCommand(t *testing.T) {
	tmpl := NewRunUserDataCommand()
	assert.NotNil(t, tmpl)
}

func TestNewRunUserDataConfig(t *testing.T) {
	join := "kubeadm --join test"
	group := "cloudsim-groupid=test-test-test"
	extra1 := "testA=ValueA"
	extra2 := "testB=ValueB"
	config := NewRunUserDataConfig(join, group, extra1, extra2)

	assert.Equal(t, join, config.JoinCmd)
	assert.Equal(t, group, config.GroupLabels)
	assert.Equal(t, fmt.Sprintf("%s,%s", extra1, extra2), config.ExtraLabels)
}

func TestFillUserDataCommand(t *testing.T) {
	const cmd = `#!/bin/bash
set -x
exec > >(tee /var/log/user-data.log|logger -t user-data ) 2>&1
echo BEGIN
date '+%Y-%m-%d %H:%M:%S'
`

	join := "kubeadm --join test"
	group := "cloudsim-groupid=test-test-test"
	extra1 := "testA=ValueA"
	extra2 := "testB=ValueB"

	// Extracted from cloudsim 1.0
	nodeGroupLabel := group
	nodeLabels := fmt.Sprintf(`cat > /etc/systemd/system/kubelet.service.d/20-labels-taints.conf <<EOF
[Service]
Environment="KUBELET_EXTRA_ARGS=--node-labels=%s,%s,%s,"
EOF
`, nodeGroupLabel, extra1, extra2)
	expectedResult := cmd + nodeLabels + join

	// Test
	tmpl := NewRunUserDataCommand()
	config := NewRunUserDataConfig(join, group, extra1, extra2)
	_, result := FillUserDataCommand(tmpl, config)
	assert.Equal(t, expectedResult, result)
}