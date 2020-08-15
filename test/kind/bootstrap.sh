GO111MODULE=on go get sigs.k8s.io/kind@v0.8.1

node_image='kindest/node:v1.16.9@sha256:7175872357bc85847ec4b1aba46ed1d12fa054c83ac7a8a11f5c268957fd5765'

cat <<EOF | kind create cluster --config=-
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
- role: control-plane
  image: ${node_image}
- role: worker
  image: ${node_image}
- role: worker
  image: ${node_image}
EOF
