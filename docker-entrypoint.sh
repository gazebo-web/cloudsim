#!/bin/bash
set -e
set -x

# Get the EKS cluster kubectl config
if [ -n "$AWS_CLUSTER_NAME" ]; then
    aws eks update-kubeconfig --name $AWS_CLUSTER_NAME --kubeconfig /root/.kube/config
fi

exec "$@"
