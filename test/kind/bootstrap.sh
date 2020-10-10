#!/usr/bin/env bash
set -o errexit


if [[ ! -v KIND_VERSION ]]; then
  KIND_VERSION="v0.11.1"
  readonly KIND_VERSION
fi

node_image="kindest/node:v1.20.7@sha256:cbeaf907fc78ac97ce7b625e4bf0de16e3ea725daf6b04f930bd14c67c671ff9"

# pull in KinD bin
TEMP_DIR=$(mktemp -d)
readonly TEMP_DIR
pushd ${TEMP_DIR}
curl -Lo ./kind https://github.com/kubernetes-sigs/kind/releases/download/${KIND_VERSION}/kind-$(uname)-amd64
chmod +x ./kind
KIND_BIN=${TEMP_DIR}/kind
readonly KIND_BIN
popd


cat <<EOF | ${KIND_BIN} create cluster --config=-
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- role: control-plane
  image: ${node_image}
  extraPortMappings:
  - containerPort: 31080
    hostPort: 80
    protocol: TCP
  - containerPort: 31443
    hostPort: 443
    protocol: TCP
- role: worker
  image: ${node_image}
- role: worker
  image: ${node_image}
- role: worker
  image: ${node_image}
EOF
