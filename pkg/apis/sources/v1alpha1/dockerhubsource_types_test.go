package v1alpha1

import (
	"fmt"
	"knative.dev/pkg/apis/duck"
	duckv1 "knative.dev/pkg/apis/duck/v1"
	"testing"

	"github.com/google/go-cmp/cmp"
)

var _ = duck.VerifyType(&DockerHubSource{}, &duckv1.Conditions{})

func TestDockerHubSource_GetGroupVersionKind(t *testing.T) {
	src := DockerHubSource{}
	gvk := src.GetGroupVersionKind()

	if gvk.Kind != "DockerHubSource" {
		t.Errorf("Should be 'DockerHubSource'.")
	}
}

func TestDockerHubCloudEventsEventType(t *testing.T) {
	prefix := "tom24d.source.dockerhub"
	eventType := "push"
	want := fmt.Sprintf("%s.%s", prefix, eventType)

	got := DockerHubCloudEventsEventType(eventType)
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("unexpected eventType (-want, +got) = %v", diff)
	}
}

func TestDockerHubEventSource(t *testing.T) {
	prefix := "https://hub.docker.com"
	repoName := "tom24d/test"
	want := fmt.Sprintf("%s/r/%s", prefix, repoName)

	got := DockerHubEventSource(repoName)
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("unexpected source (-want, +got) = %v", diff)
	}
}
