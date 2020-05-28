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
	knative.dev/pkg v0.0.0-20200528190300-08a86da47d28
	knative.dev/serving v0.15.1-0.20200528205001-f9c3c1e62a57
	knative.dev/test-infra v0.0.0-20200528222301-350178ab2a0e // indirect
)

replace k8s.io/api => k8s.io/api v0.17.5

replace k8s.io/apimachinery => k8s.io/apimachinery v0.17.5

replace k8s.io/client-go => k8s.io/client-go v0.17.5

replace k8s.io/code-generator => k8s.io/code-generator v0.17.5
