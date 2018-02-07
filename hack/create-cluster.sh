#!/usr/bin/env bash

set -xeuo pipefail

if [[ "$(minikube status)" != *"Running"* ]]; then
    minikube start --vm-driver=virtualbox \
    --kubernetes-version=v1.9.2 \
    --bootstrapper=kubeadm
fi

kubectl apply -f https://raw.githubusercontent.com/Azure/helm-charts/master/docs/prerequisities/helm-rbac-config.yaml
helm init --service-account tiller
watch kubectl get pods -n kube-system

helm repo add svc-cat https://svc-catalog-charts.storage.googleapis.com
helm upgrade --install catalog --namespace svc-cat --wait svc-cat/catalog
watch kubectl get pods -n svc-cat

# TODO: replace with minibroker
helm repo add azure https://kubernetescharts.blob.core.windows.net/azure
helm upgrade osba --namespace osba \
  --set azure.subscriptionId=$AZURE_SUBSCRIPTION_ID \
  --set azure.tenantId=$AZURE_TENANT_ID \
  --set azure.clientId=$AZURE_CLIENT_ID \
  --set azure.clientSecret=$AZURE_CLIENT_SECRET \
  --wait \
  --install azure/open-service-broker-azure
watch kubectl get pods -n osba
