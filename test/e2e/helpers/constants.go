package helpers

const (
	Pusher   = "webhook-sender-pod"
	Tag      = "latest"
	RepoName = "e2e/sender"

	ValidationReceivePort = 8080 // same number as knative/eventing/test/lib/resources/ServiceDefaultHTTP()
)
