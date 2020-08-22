package helpers

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	eventingtestlib "knative.dev/eventing/test/lib"

	sourcev1alpha1 "github.com/tom24d/eventing-dockerhub/pkg/apis/sources/v1alpha1"
)

func CreateDockerHubSourceOrFail(c *eventingtestlib.Client, dockerHubSource *sourcev1alpha1.DockerHubSource) {
	createdDockerHubSource, err := GetSourceClient(c).SourcesV1alpha1().
		DockerHubSources(dockerHubSource.GetNamespace()).Create(dockerHubSource)
	if err != nil {
		c.T.Fatalf("Failed to create DockerHubSource %q: %v", dockerHubSource.Name, err)
	}

	c.Tracker.AddObj(createdDockerHubSource)
	dockerHubSource = createdDockerHubSource
}

func GetSourceOrFail(c *eventingtestlib.Client, namespace, name string) *sourcev1alpha1.DockerHubSource {
	gotDockerHubSource, err := GetSourceClient(c).SourcesV1alpha1().
		DockerHubSources(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		c.T.Fatalf("Failed to create DockerHubSource %q: %v", name, err)
	}

	return gotDockerHubSource
}

func DeleteKServiceOrFail(c *eventingtestlib.Client, name, namespace string) {
	err := GetServiceClient(c).ServingV1().Services(namespace).Delete(name, metav1.NewDeleteOptions(0))
	if err != nil {
		c.T.Fatalf("Failed to delete backed knative service %q: %c", name, err)
	}
}
