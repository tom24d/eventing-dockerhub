module github.com/tom24d/eventing-dockerhub

go 1.16

require (
	github.com/cloudevents/sdk-go/v2 v2.4.1
	github.com/google/go-cmp v0.5.6
	github.com/google/uuid v1.2.0
	github.com/hashicorp/go-retryablehttp v0.6.7
	go.uber.org/zap v1.17.0
	gopkg.in/go-playground/webhooks.v5 v5.15.0
	k8s.io/api v0.20.7
	k8s.io/apimachinery v0.20.7
	k8s.io/client-go v0.20.7
	knative.dev/eventing v0.23.1-0.20210618081650-b164d3c042d2
	knative.dev/hack v0.0.0-20210614141220-66ab1a098940
	knative.dev/pkg v0.0.0-20210618060751-f454995ff92b
	knative.dev/serving v0.23.1-0.20210618204038-1c042a1a8c97
)

replace gopkg.in/go-playground/webhooks.v5 => ./third_party/webhooks
