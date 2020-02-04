#/bin/bash

# This script is used to create a Kubernetes secret with name 'aws-secrets' based on existing
# environment variables: AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, AWS_REGION, and AWS_ACCOUNT.
# This Secret will be used by the cloudsim server.

kubectl create secret generic aws-secrets --from-literal=aws-access-key-id=${AWS_ACCESS_KEY_ID} --from-literal=aws-secret-access-key=${AWS_SECRET_ACCESS_KEY} --from-literal=aws-account=${AWS_ACCOUNT} --from-literal=aws-region=${AWS_REGION}
