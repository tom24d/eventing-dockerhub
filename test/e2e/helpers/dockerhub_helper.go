package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	// k8s.io imports
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"

	// eventing imports
	eventingtestlib "knative.dev/eventing/test/lib"
	"knative.dev/eventing/test/lib/recordevents"
	"knative.dev/eventing/test/lib/resources"

	// pkg imports
	"knative.dev/pkg/apis/duck/v1"
	"knative.dev/pkg/test"

	sourcesv1alpha1 "github.com/tom24d/eventing-dockerhub/pkg/apis/sources/v1alpha1"
	dhsOptions "github.com/tom24d/eventing-dockerhub/pkg/reconciler/testing"

	dockerhub "gopkg.in/go-playground/webhooks.v5/docker"

	cetestv2 "github.com/cloudevents/sdk-go/v2/test"
	"github.com/google/go-cmp/cmp"
)

// MustSendWebhook sends data to the given targetURL.
func MustSendWebhook(client *eventingtestlib.Client, targetURL string, data *dockerhub.BuildPayload) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		client.T.Fatalf("failed to marshal payload: %v", err)
	}

	client.T.Logf("Sending %+v to %s", data, targetURL)

	res, err := http.Post(targetURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		client.T.Fatalf("failed to send payload: %v", err)
	}
	if code := res.StatusCode; code < http.StatusOK || http.StatusBadRequest <= code {
		client.T.Fatalf("status code got: %d", res.StatusCode)
	}
}

// GetSourceEndpointOrFail gets source's endpoint or fail.
func GetSourceEndpointOrFail(client *eventingtestlib.Client, source *sourcesv1alpha1.DockerHubSource) string {
	dhCli := GetSourceClient(client).SourcesV1alpha1().DockerHubSources(client.Namespace)
	url := ""

	err := wait.PollImmediate(1*time.Second, 1*time.Minute, func() (bool, error) {
		dhs, err := dhCli.Get(source.Name, metav1.GetOptions{})
		if err != nil {
			return true, fmt.Errorf("failed to get DockerHubSource: %v", source.Name)
		}
		url = dhs.Status.URL.String()
		if url == "" {
			return false, nil
		}
		return true, nil
	})
	if err != nil {
		client.T.Fatalf("failed to source endpoint: %v", err)
	}
	return url
}

// MustHasSameServiceName ensures the source keeps same ReceiveAdapterServiceName even if ksvc gets accidentally deleted.
func MustHasSameServiceName(c *eventingtestlib.Client, dockerHubSource *sourcesv1alpha1.DockerHubSource) {
	before := GetSourceOrFail(c, c.Namespace, dockerHubSource.Name).Status.ReceiveAdapterServiceName
	if before == "" {
		c.T.Fatalf("Failed to get DockerHubSource Service for %q", dockerHubSource.Name)
	}
	DeleteKServiceOrFail(c, before, c.Namespace)

	// wait for DockerHubSource to re-make ksvc
	c.WaitForAllTestResourcesReadyOrFail()

	after := GetSourceOrFail(c, c.Namespace, dockerHubSource.Name).Status.ReceiveAdapterServiceName
	if after == "" {
		c.T.Fatalf("Failed to get DockerHubSource Service for %q", dockerHubSource.Name)
	}

	if diff := cmp.Diff(before, after); diff != "" {
		c.T.Fatalf("Source Service name should be same: (-want, +got) = %v", diff)
	}
}

func DockerHubSourceV1Alpha1(t *testing.T, data dockerhub.BuildPayload, disableAutoCallback bool, matcherGen func(namespace string) cetestv2.EventMatcher) {
	const (
		dockerHubSourceName = "e2e-dockerhub-source"
		recordEventPodName  = "e2e-dockerhub-source-logger-event-tracker"
	)

	payload := &data

	client := eventingtestlib.Setup(t, true)
	defer eventingtestlib.TearDown(client)

	var eventTracker *recordevents.EventInfoStore
	var ref *v1.KReference
	var callbackPod *corev1.Pod

	if !disableAutoCallback {
		// create event logger eventSender and service
		eventTracker, _ = recordevents.StartEventRecordOrFail(client, recordEventPodName)
		ref = resources.KnativeRefForService(recordEventPodName, client.Namespace)
	} else {
		// create callback-display pod and service
		callbackPod = CreateCallbackDisplayOrFail(client)
		ref = resources.KnativeRefForService(callbackPod.Name, client.Namespace)
	}

	dockerHubSource := dhsOptions.NewDockerHubSourceV1Alpha1(
		dockerHubSourceName,
		client.Namespace,
		dhsOptions.DisabledAutoCallback(disableAutoCallback),
		dhsOptions.WithSinkV1A1(v1.Destination{Ref: ref}),
	)

	t.Log("Creating DockerHubSource")
	CreateDockerHubSourceOrFail(client, dockerHubSource)

	validationReceiverPod := CreateValidationReceiverOrFail(client)


	// wait for DockerHubSource to be URL allocated
	client.WaitForAllTestResourcesReadyOrFail()

	// set URL
	allocatedURL := GetSourceEndpointOrFail(client, dockerHubSource)


	// set callbackURL
	payload.CallbackURL = fmt.Sprintf("http://%s", client.GetServiceHost(validationReceiverPod.GetName()))

	// access test
	t.Log("Send webhook to DockerHubSource")
	MustSendWebhook(client, allocatedURL, payload)

	if eventTracker != nil { // == !disableAutoCallback
		t.Log("Asserting CloudEvents...")
		eventTracker.AssertExact(1, recordevents.MatchEvent(matcherGen(client.Namespace)))
	} else if callbackPod != nil {
		t.Log("Confirming callback-display reports succeeded...")
		waitForPodSuccessOrFail(client, callbackPod)
	}

	t.Log("Waiting for validation receiver report...")
	waitForPodSuccessOrFail(client, validationReceiverPod)

	MustHasSameServiceName(client, dockerHubSource)
}

// waitForPodSuccessOrFail waits for v1.PodSucceeded or fail.
func waitForPodSuccessOrFail(client *eventingtestlib.Client, pod *corev1.Pod) {
	err := test.WaitForPodState(client.Kube, func(p *corev1.Pod) (bool, error) {
		if p.Status.Phase == corev1.PodFailed {
			log, e := client.Kube.PodLogs(p.Name, p.Spec.Containers[0].Name, p.Namespace)
			return true, fmt.Errorf("pod %s failed. (log, err)=: (\n%v,\n%v)", p.Name, string(log), e)
		} else if p.Status.Phase != corev1.PodSucceeded {
			return false, nil
		}
		return true, nil
	}, pod.Name, pod.Namespace)

	if err != nil {
		client.T.Fatalf("Failed waiting for completeness of the pod: %v", err)
	}
}
