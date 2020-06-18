#!/usr/bin/env bash

# Copyright 2019 The Knative Authors
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

readonly ROOT_DIR=$(dirname $0)/..
export GO111MODULE=on

source ${ROOT_DIR}/vendor/knative.dev/test-infra/scripts/presubmit-tests.sh

# TODO(tom24d): integration tests

ko apply -f config

sleep 10

kubectl apply -f ${ROOT_DIR}/test/e2e/source.yaml
kubectl apply -f ${ROOT_DIR}/test/e2e/normal-display.yaml

sleep 10

kubectl get dockerhubsource dockerhub-source
