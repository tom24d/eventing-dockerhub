#!/usr/bin/env bash

source $(dirname $0)/../vendor/knative.dev/test-infra/scripts/e2e-tests.sh

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

# Vendored eventing test image.
readonly VENDOR_EVENTING_TEST_IMAGES="vendor/knative.dev/eventing/test/test_images/"
# HEAD eventing test images.
readonly HEAD_EVENTING_TEST_IMAGES="${GOPATH}/src/knative.dev/eventing/test/test_images/"

# Publish test images.
subheader ">> Publishing test images from eventing-dockerhub"
$(dirname $0)/upload-test-images.sh "test/test_images" e2e || fail_test "Error uploading test images"
subheader ">> Publishing test images from eventing"
# We vendor test image code from eventing, in order to use ko to resolve them into Docker images, the
# path has to be a GOPATH.
sed -i 's@knative.dev/eventing/test/test_images@github.com/tom24d/eventing-dockerhub/vendor/knative.dev/eventing/test/test_images@g' "${VENDOR_EVENTING_TEST_IMAGES}"*/*.yaml
$(dirname $0)/upload-test-images.sh ${VENDOR_EVENTING_TEST_IMAGES} e2e || fail_test "Error uploading eventing test images"

LOAD_IMAGE=(
  "tom24d/webhook-sender:e2e"
  "tom24d/event-sender:e2e"
  "tom24d/recordevents:e2e"
  "tom24d/callback-display:e2e"
  "tom24d/validation-receiver:e2e"
)
for item in "${LOAD_IMAGE[@]}" ; do
    docker pull ${item}
    kind load docker-image ${item} --name kind
done
