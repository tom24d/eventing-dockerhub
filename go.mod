module github.com/tom24d/eventing-dockerhub

go 1.14

require (
	github.com/cloudevents/sdk-go/v2 v2.1.0
	github.com/google/go-cmp v0.5.0
	github.com/google/uuid v1.1.1
	go.uber.org/zap v1.15.0
	gopkg.in/go-playground/webhooks.v5 v5.14.0
	k8s.io/api v0.18.1
	k8s.io/apimachinery v0.18.5
	k8s.io/client-go v11.0.1-0.20190805182717-6502b5e7b1b5+incompatible
	knative.dev/eventing v0.16.1-0.20200718021233-4f19a2c35971
	knative.dev/pkg v0.0.0-20200718034033-6cd456823f13
	knative.dev/serving v0.16.1-0.20200717224133-ad1aafd07749
	knative.dev/test-infra v0.0.0-20200718013433-a19b8b74cc35
)

replace k8s.io/api => k8s.io/api v0.17.6

replace k8s.io/apimachinery => k8s.io/apimachinery v0.17.6

replace k8s.io/client-go => k8s.io/client-go v0.17.6

replace k8s.io/code-generator => k8s.io/code-generator v0.17.6
