package helpers

import (
	eventingtestlib "knative.dev/eventing/test/lib"

	kservicesource "knative.dev/serving/pkg/client/clientset/versioned"

	dhsource "github.com/tom24d/eventing-dockerhub/pkg/client/clientset/versioned"
)

func GetSourceClient(c *eventingtestlib.Client) *dhsource.Clientset {
	client, err := dhsource.NewForConfig(c.Config)
	if err != nil {
		c.T.Fatalf("Failed to create DockerHubSource client: %v", err)
	}
	return client
}

func GetServiceClient(c *eventingtestlib.Client) *kservicesource.Clientset {
	client, err := kservicesource.NewForConfig(c.Config)
	if err != nil {
		c.T.Fatalf("Failed to create DockerHubSource client: %v", err)
	}
	return client
}
