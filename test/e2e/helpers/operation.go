package helpers

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	eventingtestlib "knative.dev/eventing/test/lib"

	sourcev1alpha1 "github.com/tom24d/eventing-dockerhub/pkg/apis/sources/v1alpha1"
)

func CreateDockerHubSourceOrFail(ctx context.Context, c *eventingtestlib.Client, dockerHubSource *sourcev1alpha1.DockerHubSource) {
	createdDockerHubSource, err := GetSourceClient(c).SourcesV1alpha1().
		DockerHubSources(dockerHubSource.GetNamespace()).Create(ctx, dockerHubSource, metav1.CreateOptions{})
	if err != nil {
		c.T.Fatalf("Failed to create DockerHubSource %q: %v", dockerHubSource.Name, err)
	}

	c.Tracker.AddObj(createdDockerHubSource)
	dockerHubSource = createdDockerHubSource
}

func GetSourceOrFail(ctx context.Context, c *eventingtestlib.Client, namespace, name string) *sourcev1alpha1.DockerHubSource {
	gotDockerHubSource, err := GetSourceClient(c).SourcesV1alpha1().
		DockerHubSources(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		c.T.Fatalf("Failed to create DockerHubSource %q: %v", name, err)
	}

	return gotDockerHubSource
}

func DeleteKServiceOrFail(ctx context.Context, c *eventingtestlib.Client, name, namespace string) {
	err := GetServiceClient(c).ServingV1().Services(namespace).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		c.T.Fatalf("Failed to delete backed knative service %q: %c", name, err)
	}
}
