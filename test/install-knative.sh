set -o errexit
set -o nounset
set -o pipefail

export GO111MODULE=on

VERSION="0.15.0"

# Knative Serving
kubectl apply --filename https://github.com/knative/serving/releases/download/v${VERSION}/serving-crds.yaml
kubectl apply --filename https://github.com/knative/serving/releases/download/v${VERSION}/serving-core.yaml

# Install Ambassador
kubectl create namespace ambassador
kubectl apply --namespace ambassador \
  --filename https://getambassador.io/yaml/ambassador/ambassador-rbac.yaml \
  --filename https://getambassador.io/yaml/ambassador/ambassador-service.yaml
kubectl patch clusterrolebinding ambassador -p '{"subjects":[{"kind": "ServiceAccount", "name": "ambassador", "namespace": "ambassador"}]}'
kubectl set env --namespace ambassador  deployments/ambassador AMBASSADOR_KNATIVE_SUPPORT=true
kubectl patch configmap/config-network \
  --namespace knative-serving \
  --type merge \
  --patch '{"data":{"ingress.class":"ambassador.ingress.networking.knative.dev"}}'
kubectl --namespace ambassador get service ambassador

# Knative Eventing
kubectl apply  --selector knative.dev/crd-install=true \
--filename https://github.com/knative/eventing/releases/download/v${VERSION}/eventing.yaml
kubectl apply --filename https://github.com/knative/eventing/releases/download/v${VERSION}/eventing.yaml
kubectl apply --filename https://github.com/knative/eventing/releases/download/v${VERSION}/in-memory-channel.yaml
kubectl apply --filename https://github.com/knative/eventing/releases/download/v${VERSION}/mt-channel-broker.yaml