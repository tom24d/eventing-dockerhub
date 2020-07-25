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

# Vendored eventing test image.
readonly VENDOR_EVENTING_TEST_IMAGES="vendor/knative.dev/eventing/test/test_images/"
# HEAD eventing test images.
readonly HEAD_EVENTING_TEST_IMAGES="${GOPATH}/src/knative.dev/eventing/test/test_images/"

# Istio version for ci
readonly ISTIO_VERSION="1.5.7"
# Istio crd yaml
readonly SERVING_ISTIO_CRD="https://raw.githubusercontent.com/knative/serving/master/third_party/istio-${ISTIO_VERSION}/istio-crds.yaml"
# istio-ci-no-mesh
readonly SERVING_ISTIO_CI_NO_MESH="https://raw.githubusercontent.com/knative/serving/master/third_party/istio-${ISTIO_VERSION}/istio-ci-no-mesh.yaml"

# Configure DNS
readonly SERVING_DNS_SETUP="$(get_latest_knative_yaml_source "serving" "serving-default-domain")"

function start_istio() {
  header "Starting Istio-${ISTIO_VERSION}"

  subheader "Installing Istio-${ISTIO_VERSION} CRD"
  echo "Installing Istio CRD from ${SERVING_VENDORED_ISTIO_CRD}"
  kubectl apply -f ${SERVING_ISTIO_CRD}
  while [[ $(kubectl get crd gateways.networking.istio.io -o jsonpath='{.status.conditions[?(@.type=="Established")].status}') != 'True' ]]; do
    echo "Waiting on Istio CRDs"; sleep 1
  done

  subheader "Installing Istio-${ISTIO_VERSION} YAML"
  echo "Installing Istio from ${SERVING_VENDORED_ISTIO_CI_NO_MESH}"
  kubectl apply -f ${SERVING_ISTIO_CI_NO_MESH}
}

function configure_dns() {
  subheader "Configuring DNS"
  echo "Configuring DNS from ${SERVING_DNS_SETUP}"
  kubectl apply -f ${SERVING_DNS_SETUP}
}

function knative_setup() {
  start_istio
  wait_until_pods_running istio-system || fail_test "Istio not up"

  start_latest_knative_serving
  configure_dns
  wait_until_pods_running knative-serving || fail_test "Knative Serving not up"

  start_latest_knative_eventing
  wait_until_pods_running knative-eventing || fail_test "Knative Eventing not up"
}

function test_setup() {
  dockerhub_setup || return 1

  # Publish test images.
  echo ">> Publishing test images from eventing-dockerhub"
  $(dirname $0)/upload-test-images.sh "test/test_images" e2e || fail_test "Error uploading test images"
  echo ">> Publishing test images from eventing"
  # We vendor test image code from eventing, in order to use ko to resolve them into Docker images, the
  # path has to be a GOPATH.
  sed -i 's@knative.dev/eventing/test/test_images@github.com/tom24d/eventing-dockerhub/vendor/knative.dev/eventing/test/test_images@g' "${VENDOR_EVENTING_TEST_IMAGES}"*/*.yaml
  $(dirname $0)/upload-test-images.sh ${VENDOR_EVENTING_TEST_IMAGES} e2e || fail_test "Error uploading eventing test images"
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
initialize $@ --skip-istio-addon

go_test_e2e -timeout=5m ./test/e2e -tag e2e || fail_test

success
