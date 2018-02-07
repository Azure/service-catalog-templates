#!/usr/bin/env bash

set -xeuo pipefail

go build -o bin/service-catalog-templates ./cmd/service-catalog-templates
