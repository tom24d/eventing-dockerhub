module github.com/tom24d/eventing-dockerhub

go 1.13

require (
	github.com/cloudevents/sdk-go/v2 v2.0.0-RC2
	go.uber.org/zap v1.15.0
	gopkg.in/go-playground/webhooks.v5 v5.14.0
	k8s.io/api v0.18.2
	k8s.io/apimachinery v0.18.2
	k8s.io/client-go v11.0.1-0.20190805182717-6502b5e7b1b5+incompatible
	knative.dev/eventing v0.14.1-0.20200513200558-1689191fdf85
	knative.dev/pkg v0.0.0-20200514052058-c75d324f8b8b
	knative.dev/serving v0.14.0
	knative.dev/test-infra v0.0.0-20200513224158-2b7ecf0da961 // indirect
)
