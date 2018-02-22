#!/usr/bin/env bash

set -xeuo pipefail

kubectl apply -f artifacts/instance-template.yaml
kubectl apply -f artifacts/binding-template.yaml
kubectl apply -f artifacts/instance.yaml
kubectl apply -f artifacts/binding.yaml

go run ./cmd/service-catalog-templates/*.go --kubeconfig $HOME/.kube/config \
    --logtostderr=1 -v=10
