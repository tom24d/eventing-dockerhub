module github.com/tom24d/eventing-dockerhub

go 1.13

require (
	github.com/cloudevents/sdk-go/v2 v2.0.0
	go.uber.org/zap v1.15.0
	gopkg.in/go-playground/webhooks.v5 v5.14.0
	k8s.io/api v0.17.6
	k8s.io/apimachinery v0.17.6
	k8s.io/client-go v11.0.1-0.20190805182717-6502b5e7b1b5+incompatible
	knative.dev/eventing v0.15.1-0.20200528004559-c607c293b2a5
	knative.dev/pkg v0.0.0-20200527173759-2d1a04d1ff82
	knative.dev/serving v0.15.0
	knative.dev/test-infra v0.0.0-20200527191759-5ba38edc41ae // indirect
)

replace k8s.io/api => k8s.io/api v0.17.5

replace k8s.io/apimachinery => k8s.io/apimachinery v0.17.5

replace k8s.io/client-go => k8s.io/client-go v0.17.5

replace k8s.io/code-generator => k8s.io/code-generator v0.17.5
