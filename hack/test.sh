#!/usr/bin/env bash

set -xeuo pipefail

kubectl apply -f contrib/examples/instance-template.yaml
kubectl get instancetemplates
kubectl get instt

kubectl apply -f contrib/examples/binding-template.yaml
kubectl get bindingtemplates
kubectl get bndt

kubectl apply -f contrib/examples/instance.yaml
kubectl get templatedinstances
kubectl get tinst

kubectl apply -f contrib/examples/binding.yaml
kubectl get templatedbindings
kubectl get tbnd

watch kubectl describe secret testdb-creds
