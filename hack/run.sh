#!/usr/bin/env bash

set -xeuo pipefail

go run ./cmd/service-catalog-templates/*.go --kubeconfig $HOME/.kube/config \
    --logtostderr=1 -v=10
