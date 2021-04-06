#!/bin/bash
set -x
exec > >(tee /var/log/user-data.log|logger -t user-data ) 2>&1
echo BEGIN
date '+%Y-%m-%d %H:%M:%S'
cat > /etc/systemd/system/kubelet.service.d/20-labels-taints.conf <<EOF
[Service]
Environment="KUBELET_EXTRA_ARGS=--node-labels={{ .Labels }}"
EOF
set -o xtrace
/etc/eks/bootstrap.sh {{ .ClusterName }} {{ .Args }}
