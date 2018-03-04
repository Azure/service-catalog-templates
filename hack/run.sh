#!/usr/bin/env bash

set -xeuo pipefail

helm upgrade --install svcatt-crd charts/svcatt-crd
helm upgrade --install svcatt-osba charts/svcatt-osba

go run ./cmd/service-catalog-templates/*.go --logtostderr=1 -v=10
