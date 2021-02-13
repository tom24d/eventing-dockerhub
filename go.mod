module github.com/tom24d/eventing-dockerhub

go 1.15

require (
	github.com/cloudevents/sdk-go/v2 v2.2.0
	github.com/google/go-cmp v0.5.4
	github.com/google/uuid v1.2.0
	github.com/hashicorp/go-retryablehttp v0.6.7
	go.uber.org/zap v1.16.0
	gopkg.in/go-playground/webhooks.v5 v5.15.0
	k8s.io/api v0.19.7
	k8s.io/apimachinery v0.19.7
	k8s.io/client-go v0.19.7
	knative.dev/eventing v0.20.1-0.20210211204829-1c771ce4058e
	knative.dev/hack v0.0.0-20210203173706-8368e1f6eacf
	knative.dev/pkg v0.0.0-20210212203835-448ae657fb5f
	knative.dev/serving v0.20.1-0.20210212211735-fabcc2afd40d
)

replace gopkg.in/go-playground/webhooks.v5 => ./third_party/webhooks
