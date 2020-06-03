module github.com/tom24d/eventing-dockerhub

go 1.14

require (
	github.com/cloudevents/sdk-go/v2 v2.0.0
	github.com/google/uuid v1.1.1
	go.uber.org/zap v1.15.0
	gopkg.in/go-playground/webhooks.v5 v5.14.0
	k8s.io/api v0.17.6
	k8s.io/apimachinery v0.17.6
	k8s.io/client-go v11.0.1-0.20190805182717-6502b5e7b1b5+incompatible
	knative.dev/eventing v0.15.1-0.20200601170604-bdc48890f0e4
	knative.dev/pkg v0.0.0-20200601184204-18c577c87d4f
	knative.dev/serving v0.15.1-0.20200602034816-2d0091efb2b8
	knative.dev/test-infra v0.0.0-20200602024303-e406de9f65c4 // indirect
)

replace k8s.io/api => k8s.io/api v0.17.5

replace k8s.io/apimachinery => k8s.io/apimachinery v0.17.5

replace k8s.io/client-go => k8s.io/client-go v0.17.5

replace k8s.io/code-generator => k8s.io/code-generator v0.17.5
