#!/usr/bin/env bash

set -xeuo pipefail

BINDIR=${BINDIR:-}
if [[ "$BINDIR" != "" ]]; then
    OUTPUT="-o ${BINDIR}/service-catalog-templates"
else
    OUTPUT=""
fi

CGO_ENABLED=0 go build -tags netgo ${OUTPUT} ./cmd/service-catalog-templates
