module github.com/tom24d/eventing-dockerhub

go 1.14

require (
	github.com/cloudevents/sdk-go/v2 v2.1.0
	github.com/google/go-cmp v0.5.1
	github.com/google/uuid v1.1.1
	go.uber.org/zap v1.15.0
	gopkg.in/go-playground/webhooks.v5 v5.15.0
	k8s.io/api v0.18.1
	k8s.io/apimachinery v0.18.6
	k8s.io/client-go v11.0.1-0.20190805182717-6502b5e7b1b5+incompatible
	knative.dev/eventing v0.16.1-0.20200723071157-82399b82ed81
	knative.dev/pkg v0.0.0-20200723060257-ae9c3f7fa8d3
	knative.dev/serving v0.16.1-0.20200723083057-234bd6b1761c
	knative.dev/test-infra v0.0.0-20200722142057-3ca910b5a25e
)

replace k8s.io/api => k8s.io/api v0.17.6

replace k8s.io/apimachinery => k8s.io/apimachinery v0.17.6

replace k8s.io/client-go => k8s.io/client-go v0.17.6

replace k8s.io/code-generator => k8s.io/code-generator v0.17.6

replace gopkg.in/go-playground/webhooks.v5 => github.com/tom24d/webhooks v5.15.1-0.20200724062239-a4d0e87c76c3+incompatible
