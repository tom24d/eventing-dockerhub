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
	knative.dev/eventing v0.26.1-0.20211014072442-a6a819dc71cf
	knative.dev/hack v0.0.0-20211015200324-86876688e735
	knative.dev/pkg v0.0.0-20211015194524-a5bb75923981
	knative.dev/serving v0.26.1-0.20211015224023-c9b2e25cb553
)
