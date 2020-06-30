//+build e2e

package e2e

import (
	"testing"
	"time"

	cetestv2 "github.com/cloudevents/sdk-go/v2/test"

	dockerhub "gopkg.in/go-playground/webhooks.v5/docker"

	adapterresource "github.com/tom24d/eventing-dockerhub/pkg/adapter/resources"
	sourcev1alpha1 "github.com/tom24d/eventing-dockerhub/pkg/apis/sources/v1alpha1"
	"github.com/tom24d/eventing-dockerhub/test/e2e/helpers"
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
			helpers.DockerHubSourceV1Alpha1(t, &testData.webhookPayload, false, testData.matcherGen)
			helpers.DockerHubSourceV1Alpha1(t, &testData.webhookPayload, true, testData.matcherGen)
		})
	}
}
