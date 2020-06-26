package helpers

import (
	eventingtestlib "knative.dev/eventing/test/lib"

	dhsource "github.com/tom24d/eventing-dockerhub/pkg/client/clientset/versioned"
)

func GetSourceClient(c *eventingtestlib.Client) *dhsource.Clientset {
	client, err := dhsource.NewForConfig(c.Config)
	if err != nil {
		c.T.Fatalf("Failed to create DockerHubSource client: %v", err)
	}
	return client
}
