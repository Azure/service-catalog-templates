#!/bin/bash

# Copyright 2017 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -x
set -o errexit
set -o nounset
set -o pipefail

SCRIPT_ROOT=$(realpath $(dirname ${BASH_SOURCE})/..)
CODEGEN_PKG=${CODEGEN_PKG:-$(cd ${SCRIPT_ROOT}; ls -d -1 ./vendor/k8s.io/code-generator 2>/dev/null || echo ../code-generator)}

echo "Generating code in ${SCRIPT_ROOT}"
rm -fr ${SCRIPT_ROOT}/pkg/client

# generate the code with:
# --output-base    because this script should also be able to run inside the vendor dir of
#                  k8s.io/kubernetes. The output-base is needed for the generators to output into the vendor dir
#                  instead of the $GOPATH directly. For normal projects this can be dropped.
${CODEGEN_PKG}/generate-groups.sh "deepcopy,client,informer,lister" \
  github.com/Azure/service-catalog-templates/pkg/client github.com/Azure/service-catalog-templates/pkg/apis \
  templates:experimental \
  --go-header-file ${SCRIPT_ROOT}/hack/custom-boilerplate.go.txt

# Fix casing problem, where one directory is generated into a directory under lower-case "azure"
DEST_DIR=${SCRIPT_ROOT}/pkg/client/clientset/versioned/typed/templates/
SRC_DIR=${SCRIPT_ROOT/Azure/azure}/pkg/client/clientset/versioned/typed/templates/experimental
mkdir -p ${DEST_DIR}
mv ${SRC_DIR} ${DEST_DIR}
find ./pkg/client -type f -exec sed -i 's/azure/Azure/g' {} \; # rewrite bad imports to use Azure
