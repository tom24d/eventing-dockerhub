module github.com/tom24d/eventing-dockerhub

go 1.16

require (
	github.com/cloudevents/sdk-go/v2 v2.4.1
	github.com/google/go-cmp v0.5.6
	github.com/google/uuid v1.3.0
	github.com/hashicorp/go-retryablehttp v0.6.7
	go.uber.org/zap v1.18.1
	gopkg.in/go-playground/webhooks.v5 v5.15.0
	k8s.io/api v0.20.7
	k8s.io/apimachinery v0.20.7
	k8s.io/client-go v0.20.7
	knative.dev/eventing v0.25.1-0.20210812141535-41e36eea88aa
	knative.dev/hack v0.0.0-20210806075220-815cd312d65c
	knative.dev/pkg v0.0.0-20210812143236-2c8a7b6bd74e
	knative.dev/serving v0.25.1-0.20210813163652-db6e59d3c173
)

replace gopkg.in/go-playground/webhooks.v5 => ./third_party/webhooks
