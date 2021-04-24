module github.com/tom24d/eventing-dockerhub

go 1.15

require (
	github.com/cloudevents/sdk-go/v2 v2.4.1
	github.com/google/go-cmp v0.5.5
	github.com/google/uuid v1.2.0
	github.com/hashicorp/go-retryablehttp v0.6.7
	go.uber.org/zap v1.16.0
	gopkg.in/go-playground/webhooks.v5 v5.15.0
	k8s.io/api v0.19.7
	k8s.io/apimachinery v0.19.7
	k8s.io/client-go v0.19.7
	knative.dev/eventing v0.22.1-0.20210423044837-a0a33025aee0
	knative.dev/hack v0.0.0-20210423193138-b5f6e2587f6d
	knative.dev/pkg v0.0.0-20210423162638-78b8140ed19c
	knative.dev/serving v0.22.1-0.20210423111038-353798437452
)

replace gopkg.in/go-playground/webhooks.v5 => ./third_party/webhooks
