//+build e2e

package e2e

import (
	"testing"
	"time"

	"knative.dev/eventing/test/lib"
	"knative.dev/eventing/test/lib/recordevents"
	"knative.dev/eventing/test/lib/resources"

	duckv1 "knative.dev/pkg/apis/duck/v1"

	cetestv2 "github.com/cloudevents/sdk-go/v2/test"

	dockerhub "gopkg.in/go-playground/webhooks.v5/docker"

	adapterresource "github.com/tom24d/eventing-dockerhub/pkg/adapter/resources"
	sourcev1alpha1 "github.com/tom24d/eventing-dockerhub/pkg/apis/sources/v1alpha1"
	sourcetesting "github.com/tom24d/eventing-dockerhub/pkg/reconciler/testing"

	"github.com/tom24d/eventing-dockerhub/test/e2e/helpers"
	dhtestresources "github.com/tom24d/eventing-dockerhub/test/resources"
)

func TestDockerHubSource(t *testing.T) {
	timeNow := time.Now()

	tests := map[string]struct {
		webhookPayload dockerhub.BuildPayload
		matcherGen     func(namespace string) cetestv2.EventMatcher
	}{
		"valid_payload": {
			webhookPayload: func() dockerhub.BuildPayload {
				p := dockerhub.BuildPayload{}
				p.PushData.PushedAt = float32(timeNow.Unix())
				p.PushData.Pusher = helpers.Pusher
				p.PushData.Tag = helpers.Tag
				p.Repository.RepoName = helpers.RepoName
				return p
			}(),
			matcherGen: func(namespace string) cetestv2.EventMatcher {
				return cetestv2.AllOf(
					cetestv2.HasSource(sourcev1alpha1.DockerHubEventSource(helpers.RepoName)), //TODO add more
					cetestv2.HasType(sourcev1alpha1.DockerHubCloudEventsEventType(adapterresource.DockerHubEventType)),
					//cetestv2.HasTime(timeNow),
					cetestv2.HasSubject(helpers.Pusher),
					cetestv2.HasExtension("tag", helpers.Tag),
				)
			},
		},
	}

	for name, test := range tests {
		testData := test
		t.Run(name, func(t *testing.T) {
			testDockerHubSourceV1Alpha1(t, &testData.webhookPayload, testData.matcherGen)
		})
	}
}

func testDockerHubSourceV1Alpha1(t *testing.T, payload *dockerhub.BuildPayload, matcherGen func(namespace string) cetestv2.EventMatcher) {
	const (
		dockerHubSourceName = "e2e-dockerhub-source"
		recordEventPodName = "e2e-dockerhub-source-logger-event-tracker"
	)

	client := lib.Setup(t, true)
	defer lib.TearDown(client)

	// create event logger eventSender and service
	eventTracker, _ := recordevents.StartEventRecordOrFail(client, recordEventPodName)
	defer eventTracker.Cleanup()

	dockerHubSource := sourcetesting.NewDockerHubSourceV1Alpha1(
		dockerHubSourceName,
		client.Namespace,
		sourcetesting.WithSinkV1A1(duckv1.Destination{
			Ref: resources.KnativeRefForService(recordEventPodName, client.Namespace)},
		),
	)

	t.Log("Creating DockerHubSource")
	createdDHS := helpers.CreateDockerHubSourceOrFail(client, dockerHubSource)

	// wait for DockerHubSource to be URL allocated
	dhtestresources.WaitForAllTestResourcesReadyOrFail(client)

	// set URL
	allocatedURL := helpers.GetURLOrFail(client, createdDHS)

	validationReceiverPod := helpers.CreateValidationReceiverOrFail(client)

	dhtestresources.WaitForAllTestResourcesReadyOrFail(client)

	t.Log("Setting CallbackURL to its payload")
	t.Log(validationReceiverPod.GetObjectMeta())
	// set callbackURL
	helpers.SetCallbackURLOrFail(client, payload, validationReceiverPod.GetName())



	// wait for validation webhook received
	notify := make(chan bool)
	t.Log("Waiting for validation started...")
	go helpers.WaitForValidationReceiverPodSuccessOrFail(client, validationReceiverPod, notify)

	t.Log("Send webhook to DockerHubSource")
	helpers.MustSendWebhook(client, allocatedURL, payload)

	t.Log("Waiting for validation receiver report...")
	if n := <- notify; !n {
		t.Fatal("Failed to wait for validation receiver report")
	}
	eventTracker.AssertAtLeast(1, recordevents.MatchEvent(matcherGen(client.Namespace)))
}
