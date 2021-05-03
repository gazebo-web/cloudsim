#!/bin/bash
set -e
set -x

# EKS cluster kubectl configs
#
# The AWS_EKS_CLUSTER_CONFIG environment variable is used to configure target EKS clusters. The information must be
# provided as an array of fields.
#
# AWS_EKS_CLUSTER_CONFIG fields (in order):
# * AWS Region - Target AWS Region.
# * Cluster Name - EKS Cluster name.
# * Kubeconfig - Kubeconfig filename.
#
# Notes:
# * All fields are required.
# * Fields are separated by commas.
# * Kubeconfig files will be stored in the running user's home .kube directory.
#
#                                        Cluster name
#                                             ↓
# Example: AWS_EKS_CLUSTER_CONFIG="us-east-1,test-cluster,config_test_cluster"
#                                  ↑                           ↑
#                               Region                Cluster config name
#
# To support multiple clusters, you can define multiple configurations and separate them by semicolons.
CheckOK() {
    if [ $1 -ne 0 ]; then
        OK=0
    fi
}

if [ -n "$AWS_EKS_CLUSTER_CONFIG" ]; then
    echo "Parsing EKS cluster configurations."

    # Get the set of cluster configurations
    OK=1
    IFS=';' read -ra CLUSTERS <<< "$AWS_EKS_CLUSTER_CONFIG"
    for i in "${!CLUSTERS[@]}"; do
        echo -e "\nParsing configuration $((i+1))"
        # Parse config values
        IFS=',' read -ra CONFIG <<< "${CLUSTERS[i]}"
        if [ ${#CONFIG[@]} -ne 3 ]; then
            echo "  Invalid cluster config values:" "${CONFIG[@]}"
            CheckOK 1
            continue
        fi
        read -r REGION CLUSTER_NAME CONFIG_NAME <<< "${CONFIG[@]}"
        echo "  Region: ${REGION}"
        echo "  Cluster Name: ${CLUSTER_NAME}"
        echo "  Config Name: ${CONFIG_NAME}"

        # Get the EKS kubeconfig
        aws eks update-kubeconfig \
            --region ${REGION} \
            --name ${CLUSTER_NAME} \
            --kubeconfig $HOME/.kube/${CONFIG_NAME} \
            > /dev/null
        CheckOK $?
    done

    if [ $OK -ne 1 ]; then
        echo -e "\nCluster configuration failed. Exiting"
        exit 1
    fi
else
    echo "No EKS cluster configurations provided."
    exit 1
fi

exec "$@"
