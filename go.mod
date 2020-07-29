module github.com/tom24d/eventing-dockerhub

go 1.14

require (
	github.com/cloudevents/sdk-go/v2 v2.2.0
	github.com/google/go-cmp v0.5.1
	github.com/google/uuid v1.1.1
	go.uber.org/zap v1.15.0
	gopkg.in/go-playground/webhooks.v5 v5.15.0
	k8s.io/api v0.18.1
	k8s.io/apimachinery v0.18.6
	k8s.io/client-go v11.0.1-0.20190805182717-6502b5e7b1b5+incompatible
	knative.dev/eventing v0.16.1-0.20200724032657-8d83431c07bd
	knative.dev/pkg v0.0.0-20200724211057-f21f66204a5c
	knative.dev/serving v0.16.1-0.20200724232757-55ebbade754c
	knative.dev/test-infra v0.0.0-20200724213858-d5ec9cdc6b33
)

replace k8s.io/api => k8s.io/api v0.17.6

replace k8s.io/apimachinery => k8s.io/apimachinery v0.17.6

replace k8s.io/client-go => k8s.io/client-go v0.17.6

replace k8s.io/code-generator => k8s.io/code-generator v0.17.6

replace gopkg.in/go-playground/webhooks.v5 => github.com/tom24d/webhooks v5.15.1-0.20200724062239-a4d0e87c76c3+incompatible
