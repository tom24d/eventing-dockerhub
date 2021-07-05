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
REPO_ROOT_DIR=$(git rev-parse --show-toplevel)
source ${REPO_ROOT_DIR}/vendor/knative.dev/hack/e2e-tests.sh

# Use GNU tools on macOS. Requires the 'grep' and 'gnu-sed' Homebrew formulae.
if [ "$(uname)" == "Darwin" ]; then
  sed=gsed
  grep=ggrep
fi

TMP_DIR=$(mktemp -d -t ci-$(date +%Y-%m-%d-%H-%M-%S)-XXXXXXXXXX)
readonly TMP_DIR

readonly DOCKERHUB_INSTALLATION_CONFIG="config"


# Vendored eventing test image.
readonly VENDOR_EVENTING_TEST_IMAGES="vendor/knative.dev/eventing/test/test_images/"

# The number of pods for leader-election test
readonly REPLICAS=3

readonly KNATIVE_SOURCE_DEFAULT_NAMESPACE="knative-sources"


if [[ ! -v CI ]]; then
  # GitHub Action cannot run (cat /dev/urandom).
  # This the namespace used to install and test DockerHubSource.
  export TEST_SOURCE_NAMESPACE
  TEST_SOURCE_NAMESPACE="${TEST_SOURCE_NAMESPACE:-"knative-sources-"$(cat /dev/urandom \
  | LC_CTYPE=C tr -dc 'a-z0-9' | fold -w 10 | head -n 1)}"
fi

if [[ ! -v TEST_SOURCE_NAMESPACE ]]; then
  TEST_SOURCE_NAMESPACE=${KNATIVE_SOURCE_DEFAULT_NAMESPACE}
  readonly TEST_SOURCE_NAMESPACE
  echo "using 'knative-sources' for test installation namespace"
fi

# this is needed for calling eventingtestlib.Setup
# TODO migrate off its own test Setup client.
export SYSTEM_NAMESPACE
SYSTEM_NAMESPACE="knative-eventing"

ON_KIND=0
USE_ISTIO=0
USE_KOURIER=0

function parse_flags() {
  if [[ "$1" == "--run-on-kind" ]]; then
  ON_KIND=1
  return 1
  elif [[ "$1" == "--use-istio" ]]; then
  USE_ISTIO=1
  return 1
  elif [[ "$1" == "--use-kourier" ]]; then
  USE_KOURIER=1
  return 1
  fi
  return 0
}

function configure_dns() {
  echo "Set up Magic DNS"
  kubectl apply -f "$(get_latest_knative_yaml_source "serving" "serving-default-domain")"
}

function install_net_kourier() {
  subheader "Installing net-kourier"
  kubectl apply -f "$(get_latest_knative_yaml_source "net-kourier" "kourier")"
  kubectl patch configmap/config-network --namespace knative-serving --type merge \
  --patch '{"data":{"ingress.class":"kourier.ingress.networking.knative.dev"}}'

  if [[ ${ON_KIND} ]]; then
    kubectl apply -f "${REPO_ROOT_DIR}/test/config/kourier.yaml"
#    kubectl patch configmap/config-domain --namespace knative-serving --type merge \
#      --patch '{"data":{"127.0.0.1.nip.io":""}}'
  else
    configure_dns
    wait_until_service_has_external_http_address kourier-system kourier
  fi

  wait_until_pods_running kourier-system || return 1
}

function install_net_istio() {
  # install istio ingress controller from nightly release
  subheader "Install net-istio"
  kubectl apply -f "${KNATIVE_NET_ISTIO_RELEASE}"

  configure_dns
  wait_until_pods_running istio-system || return 1
  wait_until_service_has_external_http_address istio-system istio-ingressgateway
}

function install_istio() {
  if [[ -z "${ISTIO_VERSION}" ]]; then
    readonly ISTIO_VERSION="stable"
  fi

  # TODO: Figure out solid way to install Istio
  local NET_ISTIO_COMMIT=104ee547cee5ec04752a05dae946c5852022cebf

  # And checkout the setup script based on that commit.
  local NET_ISTIO_DIR=$(mktemp -d)
  (
    cd $NET_ISTIO_DIR \
      && git init \
      && git remote add origin https://github.com/knative-sandbox/net-istio.git \
      && git fetch --depth 1 origin $NET_ISTIO_COMMIT \
      && git checkout FETCH_HEAD
  )

  if (( MESH )); then
    ISTIO_PROFILE="istio-ci-mesh.yaml"
  else
    ISTIO_PROFILE="istio-ci-no-mesh.yaml"
  fi

  echo ">> Installing Istio"
  echo "Istio version: ${ISTIO_VERSION}"
  echo "Istio profile: ${ISTIO_PROFILE}"
  ${NET_ISTIO_DIR}/third_party/istio-${ISTIO_VERSION}/install-istio.sh ${ISTIO_PROFILE}
}

function scale_control_plane() {
  for deployment in "$@"; do
    # Make sure all pods run in leader-elected mode.
    kubectl -n "${TEST_SOURCE_NAMESPACE}" scale deployment "$deployment" --replicas=0 || failed=1
    # Give it time to kill the pods.
    sleep 5
    # Scale up components for HA tests
    kubectl -n "${TEST_SOURCE_NAMESPACE}" scale deployment "$deployment" --replicas="${REPLICAS}" || failed=1
  done
}

function unleash_duck() {
  subheader "unleash the duck"
  cat test/config/chaosduck.yaml | \
    sed "s/namespace: ${KNATIVE_SOURCE_DEFAULT_NAMESPACE}/namespace: ${TEST_SOURCE_NAMESPACE}/g" | \
    ko apply -f - || return $?
}

function knative_setup() {
  # Install serving
  header "Starting Knative Serving"
  subheader "Installing Knative Serving"
  echo "Installing Serving CRDs from ${KNATIVE_SERVING_RELEASE_CRDS}"
  kubectl apply -f "${KNATIVE_SERVING_RELEASE_CRDS}"
  echo "Installing Serving core components from ${KNATIVE_SERVING_RELEASE_CORE}"
  kubectl apply -f "${KNATIVE_SERVING_RELEASE_CORE}"
  # install ingress controller
  if (( USE_KOURIER )); then
    install_net_kourier || fail_test "net-kourier not up"
  elif (( USE_ISTIO )); then
    install_istio || fail_test "Istio not up"
    install_net_istio || fail_test "istio not up"
  fi

  wait_until_pods_running knative-serving || fail_test "Knative Serving not up"

  # Install eventing
  start_latest_knative_eventing
  wait_until_pods_running knative-eventing || fail_test "Knative Eventing not up"
}

function smoke_test() {
    header "Smoke Test for example resources"
    ko apply -f ${REPO_ROOT_DIR}/example/
    wait_until_pods_running default || fail_test "example resource not up"
    ko delete -f ${REPO_ROOT_DIR}/example/
}

function test_setup() {
  echo ">> Setting up logging..."
  # Install kail if needed.
  if ! which kail > /dev/null; then
    bash <( curl -sfL https://raw.githubusercontent.com/boz/kail/master/godownloader.sh) -b "$GOPATH/bin"
  fi
  # Capture all logs.
  kail > ${ARTIFACTS}/k8s.log.txt &
  local kail_pid=$!
  # Clean up kail so it doesn't interfere with job shutting down
  add_trap "kill $kail_pid || true" EXIT

  echo ">> Creating ${TEST_SOURCE_NAMESPACE} namespace if it does not exist"
  kubectl get ns ${TEST_SOURCE_NAMESPACE} || kubectl create namespace ${TEST_SOURCE_NAMESPACE}

  dockerhub_setup || return 1

  unleash_duck || fail_test "Could not unleash the chaos duck"

  smoke_test || fail_test


  # Publish test images.
  echo ">> Publishing test images from eventing-dockerhub"
  ${REPO_ROOT_DIR}/test/upload-test-images.sh "test/test_images" e2e || fail_test "Error uploading test images"
  echo ">> Publishing test images from eventing"
  # We vendor test image code from eventing, in order to use ko to resolve them into Docker images, the
  # path has to be a GOPATH.
  local knative="knative.dev/eventing/test/test_images"
  local repo="github.com/tom24d/eventing-dockerhub/vendor/knative.dev/eventing/test/test_images"
  sed -i "s@${knative}@${repo}@g" "${VENDOR_EVENTING_TEST_IMAGES}"*/*.yaml
  ${REPO_ROOT_DIR}/test/upload-test-images.sh ${VENDOR_EVENTING_TEST_IMAGES} e2e || fail_test "Error uploading eventing test images"
  # rollback
  sed -i "s@${repo}@${knative}@g" "${VENDOR_EVENTING_TEST_IMAGES}"*/*.yaml
}

function test_teardown() {
  dockerhub_teardown

  if [[ ${TEST_SOURCE_NAMESPACE} != ${KNATIVE_SOURCE_DEFAULT_NAMESPACE} ]]; then
    kubectl delete ns ${TEST_SOURCE_NAMESPACE}
  fi
}

function dockerhub_setup() {
  header "Installing DockerHubSource"

  local TMP_SOURCE_CONTROLLER_CONFIG_DIR=${TMP_DIR}/${DOCKERHUB_INSTALLATION_CONFIG}
  mkdir -p ${TMP_SOURCE_CONTROLLER_CONFIG_DIR}
  cp -r ${DOCKERHUB_INSTALLATION_CONFIG}/* ${TMP_SOURCE_CONTROLLER_CONFIG_DIR}
  find ${TMP_SOURCE_CONTROLLER_CONFIG_DIR} -type f -name "*.yaml" -exec sed -i "s/namespace: ${KNATIVE_SOURCE_DEFAULT_NAMESPACE}/namespace: ${TEST_SOURCE_NAMESPACE}/g" {} +
  ko apply -f ${TMP_SOURCE_CONTROLLER_CONFIG_DIR} || return 1

  scale_control_plane dockerhub-source-controller dockerhub-source-webhook

  wait_until_pods_running ${TEST_SOURCE_NAMESPACE} || fail_test "DockerHubSource controller not up"
}

function dockerhub_teardown() {
  header "Uninstalling DockerHubSource"
  local TMP_SOURCE_CONTROLLER_CONFIG_DIR=${TMP_DIR}/${DOCKERHUB_INSTALLATION_CONFIG}
  ko delete --ignore-not-found=true --now --timeout 60s -f ${TMP_SOURCE_CONTROLLER_CONFIG_DIR}
}
