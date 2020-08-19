#!/usr/bin/env bash

REPO_ROOT_DIR=$(git rev-parse --show-toplevel)
source ${REPO_ROOT_DIR}/vendor/knative.dev/test-infra/scripts/e2e-tests.sh

if [[ ${KNATIVE_VERSION} == "latest-release" ]]; then
  SERVING_VERSION="v0.17.0"
  EVENTING_VERSION="v0.17.0"
else
  SERVING_VERSION="nightly"
  EVENTING_VERSION="nightly"
fi


function install_knative() {
  local repo_name="$1"
  local yaml_name="$2"
  local version="$3"
  if [[ ${version} == "nightly" ]]; then
  kubectl apply --filename "https://storage.googleapis.com/knative-nightly/${repo_name}/latest/${yaml_name}.yaml"
  return
  fi
  kubectl apply --filename "https://github.com/knative/${repo_name}/releases/download/${version}/${yaml_name}.yaml"
}


header "Install Serving ${SERVING_VERSION}"
install_knative "serving" "serving-crds" ${SERVING_VERSION}
install_knative "serving" "serving-core" ${SERVING_VERSION}

header "Install Kourier ${SERVING_VERSION}"
install_knative "net-kourier" "kourier" ${SERVING_VERSION}
cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: Service
metadata:
    name: kourier
    namespace: kourier-system
    labels:
        networking.knative.dev/ingress-provider: kourier
spec:
    ports:
    - name: http2
      port: 80
      protocol: TCP
      targetPort: 8080
      nodePort: 31080
    - name: https
      port: 443
      protocol: TCP
      targetPort: 8443
      nodePort: 31443
    selector:
        app: 3scale-kourier-gateway
    type: NodePort
EOF
kubectl patch configmap/config-network --namespace knative-serving --type merge \
  --patch '{"data":{"ingress.class":"kourier.ingress.networking.knative.dev"}}'
kubectl patch configmap/config-domain --namespace knative-serving --type merge \
  --patch '{"data":{"127.0.0.1.nip.io":""}}'

kubectl --namespace kourier-system get service kourier
wait_until_pods_running kourier-system || fail_test "Kourier not up"
wait_until_pods_running knative-serving || fail_test "Knative Serving not up"

header "Install Eventing ${EVENTING_VERSION}"
install_knative "eventing" "eventing-crds" ${EVENTING_VERSION}
install_knative "eventing" "eventing-core" ${EVENTING_VERSION}
wait_until_pods_running knative-eventing || fail_test "Knative Eventing not up"
