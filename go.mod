module github.com/tom24d/eventing-dockerhub

go 1.15

require (
	github.com/cloudevents/sdk-go/v2 v2.4.0
	github.com/google/go-cmp v0.5.5
	github.com/google/uuid v1.2.0
	github.com/hashicorp/go-retryablehttp v0.6.7
	go.uber.org/zap v1.16.0
	gopkg.in/go-playground/webhooks.v5 v5.15.0
	k8s.io/api v0.19.7
	k8s.io/apimachinery v0.19.7
	k8s.io/client-go v0.19.7
	knative.dev/eventing v0.22.1-0.20210407214954-4a3216ca221e
	knative.dev/hack v0.0.0-20210325223819-b6ab329907d3
	knative.dev/pkg v0.0.0-20210409203851-3a2ae6db7097
	knative.dev/serving v0.22.1-0.20210409173351-daab3e36761a
)

replace gopkg.in/go-playground/webhooks.v5 => ./third_party/webhooks
