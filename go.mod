module github.com/tom24d/eventing-dockerhub

go 1.15

require (
	github.com/cloudevents/sdk-go/v2 v2.2.0
	github.com/google/go-cmp v0.5.5
	github.com/google/uuid v1.2.0
	github.com/hashicorp/go-retryablehttp v0.6.7
	go.uber.org/zap v1.16.0
	gopkg.in/go-playground/webhooks.v5 v5.15.0
	k8s.io/api v0.19.7
	k8s.io/apimachinery v0.19.7
	k8s.io/client-go v0.19.7
	knative.dev/eventing v0.21.1-0.20210319230253-4e1a66b9629b
	knative.dev/hack v0.0.0-20210317214554-58edbdc42966
	knative.dev/pkg v0.0.0-20210318052054-dfeeb1817679
	knative.dev/serving v0.21.1-0.20210319222153-b7ef994c0994
)

replace gopkg.in/go-playground/webhooks.v5 => ./third_party/webhooks
