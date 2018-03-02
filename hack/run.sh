#!/usr/bin/env bash

set -xeuo pipefail

kubectl apply -f artifacts/broker-instance-template.yaml
kubectl apply -f artifacts/cluster-instance-template.yaml
kubectl apply -f artifacts/instance-template.yaml
kubectl apply -f artifacts/broker-binding-template.yaml
kubectl apply -f artifacts/cluster-binding-template.yaml
kubectl apply -f artifacts/binding-template.yaml
kubectl apply -f artifacts/templated-instance.yaml
kubectl apply -f artifacts/templated-binding.yaml

go run ./cmd/service-catalog-templates/*.go --kubeconfig $HOME/.kube/config --logtostderr=1 -v=10
