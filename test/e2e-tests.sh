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
source $(dirname $0)/../vendor/knative.dev/test-infra/scripts/e2e-tests.sh

readonly DOCKERHUB_INSTALLATION_CONFIG="./config/"

function knative_setup() {
  start_latest_knative_serving
  wait_until_pods_running knative-serving || fail_test "Knative Serving not up"

  start_latest_knative_eventing
  wait_until_pods_running knative-eventing || fail_test "Knative Eventing not up"
}

function test_setup() {
  dockerhub_setup || return 1
}

function test_teardown() {
  dockerhub_teardown
}

function dockerhub_setup() {
  header "Installing DockerHubSource"
  ko apply -f "${DOCKERHUB_INSTALLATION_CONFIG}"
  wait_until_pods_running knative-sources || fail_test "DockerHubSource controller not up"
}

function dockerhub_teardown() {
  header "Uninstalling DockerHubSource"
  kubectl delete -f "${DOCKERHUB_INSTALLATION_CONFIG}"
}

# Script entry point.
initialize $@

go_test_e2e -timeout=5m ./test/e2e || fail_test

success
