module github.com/tom24d/eventing-dockerhub

go 1.16

require (
	github.com/cloudevents/sdk-go/v2 v2.4.1
	github.com/google/go-cmp v0.5.6
	github.com/google/uuid v1.3.0
	github.com/hashicorp/go-retryablehttp v0.6.7
	go.uber.org/zap v1.19.1
	gopkg.in/go-playground/webhooks.v5 v5.15.0
	k8s.io/api v0.21.4
	k8s.io/apimachinery v0.21.4
	k8s.io/client-go v0.21.4
	knative.dev/eventing v0.26.1-0.20211029100351-4de0da062efa
	knative.dev/hack v0.0.0-20211029071251-a42c72a8fc00
	knative.dev/pkg v0.0.0-20211029145451-6ff7fb81707b
	knative.dev/serving v0.26.1-0.20211029170552-e9b8ec46bc03
)
