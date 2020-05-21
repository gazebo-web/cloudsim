#!/bin/bash
set -e
set -x

# Extract ign-transport IP
if [ -z "$IGN_TRANSPORT_IP_INTERFACE" ]; then
    ip=$(ifconfig $IGN_TRANSPORT_IP_INTERFACE | sed -En 's/127.0.0.1//;s/.*inet (addr:)?(([0-9]*\.){3}[0-9]*).*/\2/p')
    export IGN_IP="$ip"
fi

# Get the EKS cluster kubectl config
if [ -z "$AWS_CLUSTER_NAME" ]; then
    aws eks update-kubeconfig --name $AWS_CLUSTER_NAME
fi

exec "$@"
