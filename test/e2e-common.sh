VERSION="0.15.0"
ISTIO_VERSION=""
MESH=0
INGRESS_CLASS="istio.ingress.networking.knative.dev"

# Latest release. If user does not supply this as a flag, the latest
# tagged release on the current branch will be used.
readonly LATEST_RELEASE_VERSION=$(git describe --match "v[0-9]*" --abbrev=0)

readonly ROOT_DIR=$(dirname $0)/..
export GO111MODULE=on

source $(dirname $0)/../vendor/knative.dev/test-infra/scripts/e2e-tests.sh


# Setup the Knative environment for running tests.
function knative_setup() {
  install_knative_serving
  install_istio
  install_knative_eventing
}

function install_knative_serving() {
  header ">> Installing Knative Serving latest public release"

  kubectl apply --filename https://github.com/knative/serving/releases/download/v${VERSION}/serving-crds.yaml
  kubectl apply --filename https://github.com/knative/serving/releases/download/v${VERSION}/serving-core.yaml

  wait_until_pods_running knative-serving || fail_test "Knative Serving did not come up"
}

function install_knative_eventing() {
  header ">> Installing Knative Eventing latest public release"
  kubectl apply  --selector knative.dev/crd-install=true \
--filename https://github.com/knative/eventing/releases/download/v${VERSION}/eventing.yaml
  kubectl apply --filename https://github.com/knative/eventing/releases/download/v${VERSION}/eventing.yaml

  wait_until_pods_running knative-eventing || fail_test "Knative Eventing did not come up"
}

function install_istio() {
  curl -sL https://istio.io/downloadIstioctl | sh -
  export PATH=$PATH:$HOME/.istioctl/bin

  cat << EOF > ./istio-minimal-operator.yaml
apiVersion: install.istio.io/v1alpha1
kind: IstioOperator
spec:
  values:
    global:
      proxy:
        autoInject: disabled
      useMCP: false
      # The third-party-jwt is not enabled on all k8s.
      # See: https://istio.io/docs/ops/best-practices/security/#configure-third-party-service-account-tokens
      jwtPolicy: first-party-jwt

  addonComponents:
    pilot:
      enabled: true
    prometheus:
      enabled: false

  components:
    ingressGateways:
      - name: istio-ingressgateway
        enabled: true
      - name: cluster-local-gateway
        enabled: true
        label:
          istio: cluster-local-gateway
          app: cluster-local-gateway
        k8s:
          service:
            type: ClusterIP
            ports:
            - port: 15020
              name: status-port
            - port: 80
              name: http2
            - port: 443
              name: https
EOF

  istioctl manifest apply -f istio-minimal-operator.yaml

  kubectl label namespace knative-serving istio-injection=enabled
  cat <<EOF | kubectl apply -f -
apiVersion: "security.istio.io/v1beta1"
kind: "PeerAuthentication"
metadata:
  name: "default"
  namespace: "knative-serving"
spec:
  mtls:
    mode: PERMISSIVE
EOF
kubectl get pods --namespace istio-system
}


# Teardown the Knative environment after tests finish.
function knative_teardown() {
  echo ">> Stopping Knative Eventing"
  echo "Uninstalling Knative Eventing"
  kubectl delete --filename https://github.com/knative/eventing/releases/download/v${VERSION}/eventing.yaml
  kubectl delete  --selector knative.dev/crd-install=true \
--filename https://github.com/knative/eventing/releases/download/v${VERSION}/eventing.yaml
  wait_until_object_does_not_exist namespaces knative-eventing

  echo ">> Stopping Knative Serving"
  echo "Uninstalling Knative Serving"
  kubectl apply --filename https://github.com/knative/serving/releases/download/v${VERSION}/serving-core.yaml
  kubectl apply --filename https://github.com/knative/serving/releases/download/v${VERSION}/serving-crds.yaml

  wait_until_object_does_not_exist namespaces knative-serving
}