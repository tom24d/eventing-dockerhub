#!/usr/bin/env bash
set -o errexit

# almost copied from knative.dev/discovery
# license: Apache-2.0 License

if [[ ! -v KIND_VERSION ]]; then
  KIND_VERSION="v0.12.0"
  readonly KIND_VERSION
fi


# Supported flags:
# Set an alternate name on the cluster, defaults to kink
#   --name "alternate-cluster-name"
# Verify the installation:
#   --check
# Set registry port, defaults to 5000
#   --reg-port 1337
#

# create registry container unless it already exists
cluster_name=knik
reg_name='kind-registry'
reg_port='5000'

node_image="kindest/node:v1.23.4@sha256:0e34f0d0fd448aa2f2819cfd74e99fe5793a6e4938b328f657c8e3f81ee0dfb9"


# Parse flags to determine any we should pass to dep.
check=0
shutdown=0
while [[ $# -ne 0 ]]; do
  parameter=$1
  case ${parameter} in
    --check) check=1 ;;
    --shutdown) shutdown=1 ;;
    *)
      [[ $# -ge 2 ]] || abort "missing parameter after $1"
      shift
      case ${parameter} in
        --name) cluster_name=$1 ;;
        --reg-port) reg_port=$1 ;;
        *) abort "unknown option ${parameter}" ;;
      esac
  esac
  shift
done
readonly check
readonly shutdown
readonly cluster_name
readonly reg_name
readonly reg_port

# pull in KinD bin
TEMP_DIR=$(mktemp -d)
readonly TEMP_DIR
pushd ${TEMP_DIR}
curl -Lo ./kind https://github.com/kubernetes-sigs/kind/releases/download/${KIND_VERSION}/kind-$(uname)-amd64
chmod +x ./kind
KIND_BIN=${TEMP_DIR}/kind
readonly KIND_BIN
popd

reg_running="$(docker inspect -f '{{.State.Running}}' "${reg_name}" 2>/dev/null || true)"
kind_running="$(${KIND_BIN} get clusters -q | grep ${cluster_name} || true)"


if (( check )); then
  if [[ "${kind_running}" ]]; then
    echo "KinD cluster ${cluster_name} is running"
  else
    echo "KinD cluster ${cluster_name} is NOT running"
  fi
  if [ "${reg_running}" == 'true' ]; then
    echo "Docker hosted registry ${reg_name} is running"
  else
    echo "Docker hosted registry ${reg_name} is NOT running"
  fi
  ${KIND_BIN} --version
  docker --version
  exit 0
fi

if (( shutdown )); then
  if [[ "${kind_running}" == "${cluster_name}" ]]; then
    echo "Deleting KinD cluster ${cluster_name}"
    ${KIND_BIN} delete cluster -q --name ${cluster_name}
  fi
  if [ "${reg_running}" == 'true' ]; then
    echo "Stopping docker hosted registry ${reg_name}"
    docker stop "${reg_name}"
    docker rm "${reg_name}"
  fi
  echo "Shutdown."
  exit 0
fi

if [ "${reg_running}" != 'true' ]; then
  docker run \
    -d --restart=always -p "${reg_port}:5000" --name "${reg_name}" \
    registry:2
fi



cat <<EOF | ${KIND_BIN} create cluster --name ${cluster_name} --config=-
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
containerdConfigPatches:
- |-
  [plugins."io.containerd.grpc.v1.cri".registry.mirrors."localhost:${reg_port}"]
    endpoint = ["http://${reg_name}:${reg_port}"]
EOF

# connect the registry to the cluster network
if [ "${running}" != 'true' ]; then
  docker network connect "kind" "${reg_name}"
fi

# Document the local registry
# https://github.com/kubernetes/enhancements/tree/master/keps/sig-cluster-lifecycle/generic/1755-communicating-a-local-registry
cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: ConfigMap
metadata:
  name: local-registry-hosting
  namespace: kube-public
data:
  localRegistryHosting.v1: |
    host: "localhost:${reg_port}"
    help: "https://kind.sigs.k8s.io/docs/user/local-registry/"
EOF

echo "To use local registry:"
echo "export KO_DOCKER_REPO=localhost:${reg_port}"
