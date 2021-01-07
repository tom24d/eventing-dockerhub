module github.com/tom24d/eventing-dockerhub

go 1.14

require (
	github.com/cloudevents/sdk-go/v2 v2.2.0
	github.com/google/go-cmp v0.5.4
	github.com/google/uuid v1.1.2
	github.com/hashicorp/go-retryablehttp v0.6.7
	github.com/onsi/gomega v1.10.2 // indirect
	go.uber.org/zap v1.16.0
	gopkg.in/go-playground/webhooks.v5 v5.15.0
	k8s.io/api v0.18.12
	k8s.io/apimachinery v0.18.12
	k8s.io/client-go v11.0.1-0.20190805182717-6502b5e7b1b5+incompatible
	knative.dev/eventing v0.19.4
	knative.dev/pkg v0.0.0-20201231084629-fb3dc711206a
	knative.dev/serving v0.19.1-0.20210101232430-6bb7c1d825cc
	knative.dev/test-infra v0.0.0-20201221075203-d99502731eb6
)

replace gopkg.in/go-playground/webhooks.v5 => ./third_party/webhooks

replace k8s.io/api => k8s.io/api v0.18.12

replace k8s.io/apimachinery => k8s.io/apimachinery v0.18.12

replace k8s.io/client-go => k8s.io/client-go v0.18.12

replace k8s.io/code-generator => k8s.io/code-generator v0.18.12

replace k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.18.12

replace k8s.io/apiserver => k8s.io/apiserver v0.18.12
