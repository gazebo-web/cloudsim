#!/bin/bash

export KUBECONFIG=/home/nkoenig/.kube/config
export KUBERNETES_CONFIG_MAP="cloudsim-config-nps"
export CI_COMMIT_SHA=`git rev-parse HEAD`
export CI_REGISTRY_IMAGE=registry.gitlab.com/ignitionrobotics/web/cloudsim
export APPLICATION_NAME=web-cloudsim-nps
export CONTAINER_IMAGE="$CI_REGISTRY_IMAGE:$CI_COMMIT_SHA"
export APPLICATION_ENVIRONMENT=web-cloudsim-nps-staging
export HEY_CMD="hey -z 15s -q 5 -c 2 http://$APPLICATION_NAME-canary.$APPLICATION_ENVIRONMENT.svc.cluster.local/healthz"
export AWS_ACCESS_KEY_ID=AKIAS5OHIRKDE2L344FO
export AWS_SECRET_ACCESS_KEY=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

docker pull $CI_REGISTRY_IMAGE || true
 
docker build --no-cache --tag $CI_REGISTRY_IMAGE:$CI_COMMIT_SHA --tag $CI_REGISTRY_IMAGE:next -f Dockerfile.nps .
docker login -u natekoenig@gmail.com -p beUuW-Nnz3pcBZcxoTwL registry.gitlab.com/ignitionrobotics/web/cloudsim
docker push $CI_REGISTRY_IMAGE:$CI_COMMIT_SHA
  
envsubst < "./deployments/00-namespace.yaml"
envsubst < "./deployments/00-namespace.yaml"  | kubectl apply -f -
envsubst < "./deployments/01-deployment.yaml"
envsubst < "./deployments/01-deployment.yaml" | kubectl apply -f -
envsubst < "./deployments/02-blue-green.yaml"
envsubst < "./deployments/02-blue-green.yaml" | kubectl apply -f -
