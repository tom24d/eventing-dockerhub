module github.com/tom24d/eventing-dockerhub

go 1.16

require (
	github.com/cloudevents/sdk-go/v2 v2.4.1
	github.com/google/go-cmp v0.5.6
	github.com/google/uuid v1.2.0
	github.com/hashicorp/go-retryablehttp v0.6.7
	go.uber.org/zap v1.18.1
	gopkg.in/go-playground/webhooks.v5 v5.15.0
	k8s.io/api v0.20.7
	k8s.io/apimachinery v0.20.7
	k8s.io/client-go v0.20.7
	knative.dev/eventing v0.24.1-0.20210714200632-25bd8efb7179
	knative.dev/hack v0.0.0-20210622141627-e28525d8d260
	knative.dev/pkg v0.0.0-20210715175632-d9b7180af6f2
	knative.dev/serving v0.24.1-0.20210715190932-0b09cdda41ac
)

replace gopkg.in/go-playground/webhooks.v5 => ./third_party/webhooks
