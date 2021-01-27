module github.com/tom24d/eventing-dockerhub

go 1.14

require (
	github.com/cloudevents/sdk-go/v2 v2.2.0
	github.com/google/go-cmp v0.5.4
	github.com/google/uuid v1.1.2
	github.com/hashicorp/go-retryablehttp v0.6.7
	go.uber.org/zap v1.16.0
	gopkg.in/go-playground/webhooks.v5 v5.15.0
	k8s.io/api v0.19.7
	k8s.io/apimachinery v0.19.7
	k8s.io/client-go v0.19.7
	knative.dev/eventing v0.20.1-0.20210127002130-ce099cf65fec
	knative.dev/hack v0.0.0-20210120165453-8d623a0af457
	knative.dev/pkg v0.0.0-20210125222030-6040b3af4803
	knative.dev/serving v0.20.1-0.20210127022430-b083383b340f
)

replace gopkg.in/go-playground/webhooks.v5 => ./third_party/webhooks
