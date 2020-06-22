set -o errexit
set -o nounset
set -o pipefail

readonly ROOT_DIR=$(dirname $0)/..
source ${ROOT_DIR}/vendor/knative.dev/test-infra/scripts/presubmit-tests.sh


export GO111MODULE=on

kubectl version

VERSION="0.15.0"

# Knative Serving
kubectl apply --filename https://github.com/knative/serving/releases/download/v${VERSION}/serving-crds.yaml
kubectl apply --filename https://github.com/knative/serving/releases/download/v${VERSION}/serving-core.yaml

# Install Ambassador
#kubectl create namespace ambassador
#kubectl apply --namespace ambassador \
#  --filename https://getambassador.io/yaml/ambassador/ambassador-rbac.yaml \
#  --filename https://getambassador.io/yaml/ambassador/ambassador-service.yaml

kubectl apply --namespace ambassador \
  --filename https://github.com/datawire/ambassador-operator/releases/latest/download/ambassador-operator-crds.yaml
kubectl apply -n ambassador -f https://github.com/datawire/ambassador-operator/releases/latest/download/ambassador-operator-kind.yaml
kubectl wait --timeout=180s -n ambassador --for=condition=deployed ambassadorinstallations/ambassador

kubectl patch clusterrolebinding ambassador -p '{"subjects":[{"kind": "ServiceAccount", "name": "ambassador", "namespace": "ambassador"}]}'
kubectl set env --namespace ambassador  deployments/ambassador AMBASSADOR_KNATIVE_SUPPORT=true
kubectl patch configmap/config-network \
  --namespace knative-serving \
  --type merge \
  --patch '{"data":{"ingress.class":"ambassador.ingress.networking.knative.dev"}}'


cat <<EOF | kubectl apply -n ambassador -f -
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: local-ingress
  annotations:
    kubernetes.io/ingress.class: ambassador
spec:
  rules:
  - http:
      paths:
      - backend:
          serviceName: ambassador
          servicePort: 80
  - http:
      paths:
      - path: /foo
        backend:
          serviceName: foo-service
          servicePort: 5678
EOF

cat <<EOF | kubectl apply -n ambassador -f -
kind: Pod
apiVersion: v1
metadata:
  name: foo-app
  labels:
    app: foo
spec:
  containers:
  - name: foo-app
    image: hashicorp/http-echo:0.2.3
    args:
    - "-text=foo"
---
kind: Service
apiVersion: v1
metadata:
  name: foo-service
spec:
  selector:
    app: foo
  ports:
  # Default port used by the image
  - port: 5678
EOF

# Knative Eventing
kubectl apply  --selector knative.dev/crd-install=true \
--filename https://github.com/knative/eventing/releases/download/v${VERSION}/eventing.yaml
kubectl apply --filename https://github.com/knative/eventing/releases/download/v${VERSION}/eventing.yaml

#wait_until_pods_running knative-serving|| fail_test "Knative Serving did not come up"
#wait_until_pods_running knative-eventing || fail_test "Knative Eventing did not come up"
sleep 120
kubectl --namespace ambassador get service ambassador
#kubectl --namespace contour-external get service envoy


kubectl apply -f ${ROOT_DIR}/test/e2e/normal-display.yaml
sleep 25
kubectl describe ksvc normal-display
kubectl describe route normal-display

echo "curl to the pod"
curl -d "Hello World KinD! foo" http://127.0.0.1:80/foo

echo "curl to the ksvc"
curl -H "Host: normal-display.default.example.com" -d "Hello World KinD!"  http://127.0.0.1:80
sleep 5
kubectl delete -f ${ROOT_DIR}/test/e2e/normal-display.yaml
