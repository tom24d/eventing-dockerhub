package helpers

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	batchv1 "k8s.io/api/batch/v1"

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

// TODO consider move this to eventing test lib
func CreateJobOrFail(c *eventingtestlib.Client, job *batchv1.Job, options ...func(*batchv1.Job, *eventingtestlib.Client) error) {
	// set namespace for the job in case it's empty
	namespace := c.Namespace
	job.Namespace = namespace

	// apply options on the cronjob before creation
	for _, option := range options {
		if err := option(job, c); err != nil {
			c.T.Fatalf("Failed to configure job %q: %v", job.Name, err)
		}
	}

	// c.applyTracingEnv(&job.Spec.Template.Spec)

	c.T.Logf("Creating job %+v", job)
	if _, err := c.Kube.Kube.BatchV1().Jobs(job.Namespace).Create(job); err != nil {
		c.T.Fatalf("Failed to create job %q: %v", job.Name, err)
	}
	c.Tracker.Add("batch", "v1", "jobs", namespace, job.Name)
}
