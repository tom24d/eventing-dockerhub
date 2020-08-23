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

export GO111MODULE=on

source $(dirname $0)/e2e-common.sh


# Script entry point.
# TODO(tom24d) cannot install Serving and Eventing together. working on it.
#initialize $@ --skip-istio-addon --use-istio
initialize $@ --skip-istio-addon --use-kourier

go_test_e2e -timeout=5m ./test/e2e -tag=e2e || fail_test

success
