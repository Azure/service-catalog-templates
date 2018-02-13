#!/usr/bin/env bash

set -xeuo pipefail

kubectl apply -f artifacts/instance-template.yaml
kubectl apply -f artifacts/instance.yaml
kubectl apply -f contrib/examples/instance-template.yaml
