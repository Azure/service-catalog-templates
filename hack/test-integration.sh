#!/usr/bin/env bash

set -xeuo pipefail

helm install --name wordpress charts/wordpress --namespace svcatt --wait

kubectl get brokerinstancetemplates
kubectl get binstt

kubectl get clusterinstancetemplates
kubectl get cinstt

kubectl get instancetemplates
kubectl get instt

kubectl get brokerbindingtemplates
kubectl get bbndt

kubectl get clusterbindingtemplates
kubectl get cbndt

kubectl get bindingtemplates
kubectl get bndt

kubectl get templatedinstances
kubectl get tinst

kubectl get templatedbindings
kubectl get tbnd

watch kubectl get secrets -n svcatt
