package helpers

import (
	"fmt"
	"github.com/cloudevents/sdk-go/v2/test"
	testing2 "github.com/tom24d/eventing-dockerhub/pkg/reconciler/testing"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"knative.dev/eventing/test/lib/recordevents"
	"knative.dev/eventing/test/lib/resources"
	"knative.dev/pkg/apis/duck/v1"
	"testing"

	eventingtestlib "knative.dev/eventing/test/lib"
	pkgTest "knative.dev/pkg/test"

	sourcesv1alpha1 "github.com/tom24d/eventing-dockerhub/pkg/apis/sources/v1alpha1"
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
				Name:  SenderImageName,
				Image: pkgTest.ImagePath(SenderImageName),
				ImagePullPolicy: corev1.PullAlways,
				Args:  args,
			}},
			RestartPolicy: corev1.RestartPolicyNever,
		},
	}
	client.CreatePodOrFail(eventSender)

	err := pkgTest.WaitForPodState(client.Kube, func(pod *corev1.Pod) (bool, error) {
		if pod.Status.Phase == corev1.PodFailed {
			return true, fmt.Errorf("event sender pod failed with message %s", pod.Status.Message)
		} else if pod.Status.Phase != corev1.PodSucceeded {
			return false, nil
		}
		return true, nil
	}, eventSender.Name, eventSender.Namespace)

	if err != nil {
		client.T.Fatalf("Failed waiting for pod for completeness %q: %v", eventSender.Name, eventSender)
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

func DockerHubSourceV1Alpha1EnabledAutoCallback(t *testing.T, payload *dockerhub.BuildPayload, matcherGen func(namespace string) test.EventMatcher) {
	const (
		dockerHubSourceName = "e2e-dockerhub-source"
		recordEventPodName  = "e2e-dockerhub-source-logger-event-tracker"
	)

	client := eventingtestlib.Setup(t, true)
	defer eventingtestlib.TearDown(client)

	// create event logger eventSender and service
	eventTracker, _ := recordevents.StartEventRecordOrFail(client, recordEventPodName)
	defer eventTracker.Cleanup()

	dockerHubSource := testing2.NewDockerHubSourceV1Alpha1(
		dockerHubSourceName,
		client.Namespace,
		testing2.WithSinkV1A1(v1.Destination{
			Ref: resources.KnativeRefForService(recordEventPodName, client.Namespace)},
		),
	)

	t.Log("Creating DockerHubSource")
	createdDHS := CreateDockerHubSourceOrFail(client, dockerHubSource)

	// wait for DockerHubSource to be URL allocated
	dhtestresources.WaitForAllTestResourcesReadyOrFail(client)

	// set URL
	allocatedURL := GetURLOrFail(client, createdDHS)

	validationReceiverPod := CreateValidationReceiverOrFail(client)

	dhtestresources.WaitForAllTestResourcesReadyOrFail(client)

	t.Log("Setting CallbackURL to its payload")
	t.Log(validationReceiverPod.GetObjectMeta())
	// set callbackURL
	SetCallbackURLOrFail(client, payload, validationReceiverPod.GetName())

	// wait for validation webhook received
	notify := make(chan bool)
	t.Log("Waiting for validation started...")
	go WaitForValidationReceiverPodSuccessOrFail(client, validationReceiverPod, notify)

	t.Log("Send webhook to DockerHubSource")
	MustSendWebhook(client, allocatedURL, payload)

	t.Log("Waiting for validation receiver report...")
	if n := <-notify; !n {
		t.Fatal("Failed to wait for validation receiver report")
	}
	eventTracker.AssertAtLeast(1, recordevents.MatchEvent(matcherGen(client.Namespace)))
}
