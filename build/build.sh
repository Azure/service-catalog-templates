#!/usr/bin/env bash

set -xeuo pipefail

go build -i -o bin/service-catalog-templates ./cmd/service-catalog-templates
