package helpers

import (
	"fmt"
	"testing"
	"time"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"

	eventingtestlib "knative.dev/eventing/test/lib"
	"knative.dev/eventing/test/lib/recordevents"
	"knative.dev/eventing/test/lib/resources"

	"knative.dev/pkg/apis/duck/v1"
	pkgTest "knative.dev/pkg/test"

	sourcesv1alpha1 "github.com/tom24d/eventing-dockerhub/pkg/apis/sources/v1alpha1"
	dhsOptions "github.com/tom24d/eventing-dockerhub/pkg/reconciler/testing"
	dhtestresources "github.com/tom24d/eventing-dockerhub/test/resources"

	dockerhub "gopkg.in/go-playground/webhooks.v5/docker"

	"github.com/cloudevents/sdk-go/v2/test"
	"github.com/google/go-cmp/cmp"
)

const (
	SenderImageName = "webhook-sender"
)

func MustSendWebhook(client *eventingtestlib.Client, targetURL string, data *dockerhub.BuildPayload) {
	args := []string{fmt.Sprintf("--%s=%s", dhtestresources.ArgSink, targetURL),
		fmt.Sprintf("--%s=%s", dhtestresources.ArgPayload, dhtestresources.MarshalPayload(data))}

	retryBackoff := int32(1)
	ttl := int32(30)

	// create webhook sender
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: client.Namespace,
			Name:      SenderImageName,
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Name:            SenderImageName,
						Image:           pkgTest.ImagePath(SenderImageName),
						ImagePullPolicy: corev1.PullAlways,
						Args:            args,
					}},
					RestartPolicy: corev1.RestartPolicyNever,
				},
			},
			BackoffLimit:            &retryBackoff,
			TTLSecondsAfterFinished: &ttl,
		},
	}
	CreateJobOrFail(client, job)

	err := WaitForJobState(client.Kube, func(job *batchv1.Job) (bool, error) {
		if len(job.Status.Conditions) >= 1 {
			if job.Status.Conditions[0].Type == batchv1.JobFailed {
				// JobFailed. Get log and return them with error.
				l, err := client.Kube.PodLogs(job.Name, SenderImageName, client.Namespace)
				if err != nil {
					// retry
					return false, nil
				}
				return true, fmt.Errorf("job:%s is failed. Log: %s", job.Name, l)
			} else {
				// JobSuccess
				return true, nil
			}
		}
		return false, nil
	}, job.Name, job.Namespace)

	if err != nil {
		client.T.Fatalf("Failed sending webhook: %v", err)
	}
}

func GetServiceAddressOrFail(client *eventingtestlib.Client, source *sourcesv1alpha1.DockerHubSource) string {

	dhCli := GetSourceClient(client).SourcesV1alpha1().DockerHubSources(client.Namespace)
	ksvcName := ""

	err := wait.PollImmediate(1*time.Second, 1*time.Minute, func() (bool, error) {
		dhs, err := dhCli.Get(source.Name, metav1.GetOptions{})
		if err != nil {
			return true, fmt.Errorf("failed to get DockerHubSource: %v", source.Name)
		}
		ksvcName = dhs.Status.ReceiveAdapterServiceName
		if ksvcName == "" {
			return false, nil
		}
		return true, nil
	})
	if err != nil {
		client.T.Fatalf("failed to get ReceiveAdapterServiceName: %v", err)
	}

	// TODO use lib if exists
	return fmt.Sprintf("http://%s.%s.svc.cluster.local", ksvcName, source.Namespace)
}

func MustHasSameServiceName(t *testing.T, c *eventingtestlib.Client, dockerHubSource *sourcesv1alpha1.DockerHubSource) {
	before := GetSourceOrFail(c, c.Namespace, dockerHubSource.Name).Status.ReceiveAdapterServiceName
	if before == "" {
		t.Fatalf("Failed to get DockerHubSource Service for %q", dockerHubSource.Name)
	}
	DeleteKServiceOrFail(c, before, c.Namespace)

	// wait for DockerHubSource to be URL allocated
	c.WaitForAllTestResourcesReadyOrFail()

	after := GetSourceOrFail(c, c.Namespace, dockerHubSource.Name).Status.ReceiveAdapterServiceName
	if before == "" {
		t.Fatalf("Failed to get DockerHubSource Service for %q", dockerHubSource.Name)
	}

	if diff := cmp.Diff(before, after); diff != "" {
		c.T.Fatalf("Source Service name should be same: (-want, +got) = %v", diff)
	}
}

func DockerHubSourceV1Alpha1(t *testing.T, payload *dockerhub.BuildPayload, disableAutoCallback bool, matcherGen func(namespace string) test.EventMatcher) {
	const (
		dockerHubSourceName = "e2e-dockerhub-source"
		recordEventPodName  = "e2e-dockerhub-source-logger-event-tracker"
	)

	notify := make(chan bool)
	defer close(notify)

	client := eventingtestlib.Setup(t, false)
	defer eventingtestlib.TearDown(client)

	// create event logger eventSender and service
	eventTracker, _ := recordevents.StartEventRecordOrFail(client, recordEventPodName)
	defer eventTracker.Cleanup()

	dockerHubSource := dhsOptions.NewDockerHubSourceV1Alpha1(
		dockerHubSourceName,
		client.Namespace,
		dhsOptions.DisabledAutoCallback(disableAutoCallback),
		dhsOptions.WithSinkV1A1(v1.Destination{
			Ref: resources.KnativeRefForService(recordEventPodName, client.Namespace)},
		),
	)

	t.Log("Creating DockerHubSource")
	createdDHS := CreateDockerHubSourceOrFail(client, dockerHubSource)

	// wait for DockerHubSource to be URL allocated
	client.WaitForAllTestResourcesReadyOrFail()

	// set URL
	allocatedURL := GetServiceAddressOrFail(client, createdDHS)

	if !disableAutoCallback {
		validationReceiverPod := CreateValidationReceiverOrFail(client)

		client.WaitForAllTestResourcesReadyOrFail()

		// set callbackURL
		payload.CallbackURL = fmt.Sprintf("http://%s", client.GetServiceHost(validationReceiverPod.GetName()))
		t.Logf("Webhook payload: %v", payload)

		// wait for validation webhook received
		t.Log("Waiting for validation started...")
		go WaitForValidationReceiverPodSuccessOrFail(client, validationReceiverPod, notify)
	}

	// access test from cluster inside
	t.Log("Send webhook to DockerHubSource")
	MustSendWebhook(client, allocatedURL, payload)

	if !disableAutoCallback {
		t.Log("Waiting for validation receiver report...")
		if n := <-notify; !n {
			t.Fatal("Failed to wait for validation receiver report")
		}
		t.Log("Validation receiver confirmed its callback.")
	}

	eventTracker.AssertAtLeast(1, recordevents.MatchEvent(matcherGen(client.Namespace)))

	MustHasSameServiceName(t, client, dockerHubSource)
}
