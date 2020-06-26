package helpers

import (
	eventingtestlib "knative.dev/eventing/test/lib"

	sourcev1alpha1 "github.com/tom24d/eventing-dockerhub/pkg/apis/sources/v1alpha1"
)

func CreateDockerHubSourceOrFail(c *eventingtestlib.Client, dockerHubSource *sourcev1alpha1.DockerHubSource) *sourcev1alpha1.DockerHubSource {
	createdDockerHubSource, err := GetSourceClient(c).SourcesV1alpha1().
		DockerHubSources(dockerHubSource.GetNamespace()).Create(dockerHubSource)
	if err != nil {
		c.T.Fatalf("Failed to create DockerHubSource %q: %v", dockerHubSource.Name, err)
	}

	c.Tracker.AddObj(createdDockerHubSource)
	return createdDockerHubSource
}
