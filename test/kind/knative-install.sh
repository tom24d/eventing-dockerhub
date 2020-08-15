#!/usr/bin/env bash

source $(dirname $0)/../../vendor/knative.dev/test-infra/scripts/e2e-tests.sh

SERVING_VERSION=0.16.0
EVENTING_VERSION=0.16.1


header "Install Serving"
kubectl apply --filename https://github.com/knative/serving/releases/download/v${SERVING_VERSION}/serving-crds.yaml
kubectl apply --filename https://github.com/knative/serving/releases/download/v${SERVING_VERSION}/serving-core.yaml

header "Install Kourier"
kubectl apply --filename https://github.com/knative/net-kourier/releases/download/v${SERVING_VERSION}/kourier.yaml
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

header "Install Eventing"
kubectl apply --filename https://github.com/knative/eventing/releases/download/v${EVENTING_VERSION}/eventing-crds.yaml
kubectl apply --filename https://github.com/knative/eventing/releases/download/v${EVENTING_VERSION}/eventing-core.yaml
wait_until_pods_running knative-eventing || fail_test "Knative Eventing not up"
