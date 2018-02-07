#!/usr/bin/env bash

set -xeuo pipefail

./hack/update-codegen.sh
go build -i -o bin/service-catalog-templates ./cmd/service-catalog-templates
