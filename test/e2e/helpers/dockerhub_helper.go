package helpers

import (
	"fmt"
	"testing"

	"github.com/cloudevents/sdk-go/v2/test"
	"github.com/google/go-cmp/cmp"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	eventingtestlib "knative.dev/eventing/test/lib"
	"knative.dev/eventing/test/lib/recordevents"
	"knative.dev/eventing/test/lib/resources"

	"knative.dev/pkg/apis/duck/v1"
	pkgTest "knative.dev/pkg/test"

	sourcesv1alpha1 "github.com/tom24d/eventing-dockerhub/pkg/apis/sources/v1alpha1"
	dhsOptions "github.com/tom24d/eventing-dockerhub/pkg/reconciler/testing"
	dhtestresources "github.com/tom24d/eventing-dockerhub/test/resources"


	dockerhub "gopkg.in/go-playground/webhooks.v5/docker"
)

const (
	SenderImageName = "webhook-sender"
)

func MustSendWebhook(client *eventingtestlib.Client, targetURL string, data *dockerhub.BuildPayload) {
	args := []string{fmt.Sprintf("--%s=%s", dhtestresources.ArgSink, targetURL),
		fmt.Sprintf("--%s=%s", dhtestresources.ArgPayload, dhtestresources.MarshalPayload(data))}

	// create webhook sender
	eventSender := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: client.Namespace,
			Name:      SenderImageName,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{{
				Name:            SenderImageName,
				Image:           pkgTest.ImagePath(SenderImageName),
				ImagePullPolicy: corev1.PullAlways,
				Args:            args,
			}},
			RestartPolicy: corev1.RestartPolicyNever,
		},
	}
	client.CreatePodOrFail(eventSender)

	err := pkgTest.WaitForPodState(client.Kube, func(pod *corev1.Pod) (bool, error) {
		if pod.Status.Phase == corev1.PodFailed {
			log, _ := client.Kube.PodLogs(pod.Name, SenderImageName, client.Namespace)
			return true, fmt.Errorf("event sender pod failed with log %s", log)
		} else if pod.Status.Phase != corev1.PodSucceeded {
			return false, nil
		}
		return true, nil
	}, eventSender.Name, eventSender.Namespace)

	if err != nil {
		client.T.Fatalf("Failed sending webhook %q: %v", eventSender.Name, err)
	}
}

func GetURLOrFail(client *eventingtestlib.Client, source *sourcesv1alpha1.DockerHubSource) string {
	dhs, err := GetSourceClient(client).SourcesV1alpha1().
		DockerHubSources(client.Namespace).Get(source.Name, metav1.GetOptions{})
	if err != nil {
		client.T.Fatalf("failed to get DockerHubSource: %v", source.Name)
	}

	allocatedURL := dhs.Status.URL.String()
	if allocatedURL == "" {
		client.T.Fatalf("DockerHubSource URL is nil: %v", source.GetName())
	}
	return allocatedURL
}

// TODO Get and return URL might be better
func SetCallbackURLOrFail(c *eventingtestlib.Client, data *dockerhub.BuildPayload, svcName string) {
	// TODO use lib if exists
	url := fmt.Sprintf("http://%s.%s.svc.cluster.local", svcName, c.Namespace)
	data.CallbackURL = url
}

func MustHasSameServiceName(t *testing.T, c *eventingtestlib.Client, dockerHubSource *sourcesv1alpha1.DockerHubSource) {
	before := GetSourceOrFail(c, c.Namespace, dockerHubSource.Name).Status.ReceiveAdapterServiceName
	if before == "" {
		t.Fatalf("Failed to get DockerHubSource Service for %q", dockerHubSource.Name)
	}
	DeleteKServiceOrFail(c, before, c.Namespace)

	// wait for DockerHubSource to be URL allocated
	dhtestresources.WaitForAllTestResourcesReadyOrFail(c)

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
	dhtestresources.WaitForAllTestResourcesReadyOrFail(client)

	// set URL
	allocatedURL := GetURLOrFail(client, createdDHS)

	if !disableAutoCallback {
		validationReceiverPod := CreateValidationReceiverOrFail(client)

		dhtestresources.WaitForAllTestResourcesReadyOrFail(client)

		t.Log("Setting CallbackURL to its payload")
		t.Log(validationReceiverPod.GetObjectMeta())
		// set callbackURL
		SetCallbackURLOrFail(client, payload, validationReceiverPod.GetName())

		// wait for validation webhook received
		t.Log("Waiting for validation started...")
		go WaitForValidationReceiverPodSuccessOrFail(client, validationReceiverPod, notify)
	}

	t.Log("Send webhook to DockerHubSource")
	MustSendWebhook(client, allocatedURL, payload)

	if !disableAutoCallback {
		t.Log("Waiting for validation receiver report...")
		if n := <-notify; !n {
			t.Fatal("Failed to wait for validation receiver report")
		}
	}

	eventTracker.AssertAtLeast(1, recordevents.MatchEvent(matcherGen(client.Namespace)))

	MustHasSameServiceName(t, client, dockerHubSource)
}
