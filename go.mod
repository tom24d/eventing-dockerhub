module github.com/tom24d/eventing-dockerhub

go 1.16

require (
	github.com/cloudevents/sdk-go/v2 v2.4.1
	github.com/google/go-cmp v0.5.6
	github.com/google/uuid v1.3.0
	github.com/hashicorp/go-retryablehttp v0.6.7
	go.uber.org/zap v1.19.0
	gopkg.in/go-playground/webhooks.v5 v5.15.0
	k8s.io/api v0.21.4
	k8s.io/apimachinery v0.21.4
	k8s.io/client-go v0.21.4
	knative.dev/eventing v0.25.1-0.20210827141738-ea5ed9adf51f
	knative.dev/hack v0.0.0-20210806075220-815cd312d65c
	knative.dev/pkg v0.0.0-20210827184538-2bd91f75571c
	knative.dev/serving v0.25.1-0.20210827140938-e6a7166509e6
)
