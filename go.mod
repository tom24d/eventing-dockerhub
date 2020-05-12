module github.com/tom24d/eventing-dockerhub

go 1.13

require (
	github.com/cloudevents/sdk-go/v2 v2.0.0-preview8 // indirect
	github.com/robfig/cron v1.2.0 // indirect
	go.uber.org/zap v1.15.0
	k8s.io/api v0.18.2
	k8s.io/apimachinery v0.18.2
	k8s.io/client-go v11.0.1-0.20190805182717-6502b5e7b1b5+incompatible
	knative.dev/eventing v0.14.2
	knative.dev/pkg v0.0.0-20200511223446-de5c590700ff
)
