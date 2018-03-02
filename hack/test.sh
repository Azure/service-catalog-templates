#!/usr/bin/env bash

set -xeuo pipefail

kubectl apply -f contrib/examples/broker-instance-template.yaml
kubectl get brokerinstancetemplates
kubectl get binstt

kubectl apply -f contrib/examples/cluster-instance-template.yaml
kubectl get clusterinstancetemplates
kubectl get cinstt

kubectl apply -f contrib/examples/instance-template.yaml
kubectl get instancetemplates
kubectl get instt

kubectl apply -f contrib/examples/binding-template.yaml
kubectl get bindingtemplates
kubectl get bndt

kubectl apply -f contrib/examples/templated-instance.yaml
kubectl get templatedinstances
kubectl get tinst

kubectl apply -f contrib/examples/templated-binding.yaml
kubectl get templatedbindings
kubectl get tbnd

watch kubectl describe secret testdb-creds
