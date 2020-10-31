module github.com/tom24d/eventing-dockerhub

go 1.14

require (
	github.com/cloudevents/sdk-go/v2 v2.2.0
	github.com/google/go-cmp v0.5.2
	github.com/google/uuid v1.1.2
	github.com/hashicorp/go-retryablehttp v0.6.7
	go.uber.org/zap v1.16.0
	gopkg.in/go-playground/webhooks.v5 v5.15.0
	k8s.io/api v0.18.8
	k8s.io/apimachinery v0.18.8
	k8s.io/client-go v11.0.1-0.20190805182717-6502b5e7b1b5+incompatible
	knative.dev/eventing v0.18.1-0.20201030152235-29ff9510e6f6
	knative.dev/pkg v0.0.0-20201030200734-72451e0d5973
	knative.dev/serving v0.18.1-0.20201030225334-a45beb3ac892
	knative.dev/test-infra v0.0.0-20201030214934-cb2c0a727744
)

replace gopkg.in/go-playground/webhooks.v5 => ./third_party/webhooks

replace k8s.io/api => k8s.io/api v0.18.8

replace k8s.io/apimachinery => k8s.io/apimachinery v0.18.8

replace k8s.io/client-go => k8s.io/client-go v0.18.8

replace k8s.io/code-generator => k8s.io/code-generator v0.18.8

replace k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.18.8

replace k8s.io/apiserver => k8s.io/apiserver v0.18.8
