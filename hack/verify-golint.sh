#!/bin/bash

# Copyright 2022 The Kubernetes Authors.
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

set -o errexit
set -o nounset
set -o pipefail

readonly VERSION="v1.64.7"
KUBE_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
readonly KUBE_ROOT

cd "${KUBE_ROOT}"

# See configuration file in ${KUBE_ROOT}/.golangci.yml.
echo "Installing golangci-lint..."
curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh -s "$VERSION"

echo "Running golangci-lint..."
./bin/golangci-lint run --timeout=10m
