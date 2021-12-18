module github.com/tom24d/eventing-dockerhub

go 1.16

require (
	github.com/cloudevents/sdk-go/v2 v2.7.0
	github.com/emicklei/go-restful v2.15.0+incompatible // indirect
	github.com/go-openapi/spec v0.20.2 // indirect
	github.com/google/go-cmp v0.5.6
	github.com/google/uuid v1.3.0
	github.com/googleapis/gnostic v0.5.3 // indirect
	github.com/hashicorp/go-retryablehttp v0.6.7
	go.uber.org/zap v1.19.1
	gopkg.in/go-playground/webhooks.v5 v5.15.0
	k8s.io/api v0.21.4
	k8s.io/apimachinery v0.21.4
	k8s.io/client-go v0.21.4
	k8s.io/utils v0.0.0-20210111153108-fddb29f9d009 // indirect
	knative.dev/eventing v0.28.1-0.20211217092418-fede720191d3
	knative.dev/hack v0.0.0-20211216134818-6fc030496333
	knative.dev/pkg v0.0.0-20211216142117-79271798f696
	knative.dev/serving v0.28.1-0.20211216134718-2df3ceda5f9b
)
