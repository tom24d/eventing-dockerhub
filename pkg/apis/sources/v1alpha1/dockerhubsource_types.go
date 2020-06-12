package v1alpha1

import (
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"knative.dev/pkg/apis"

	"knative.dev/pkg/webhook/resourcesemantics"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	duckv1 "knative.dev/pkg/apis/duck/v1"
)

// +genclient
// +genreconciler
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DockerHubSource is the Schema for the dockerhubsources API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:categories=all,knative,eventing,sources
type DockerHubSource struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	//Spec
	Spec DockerHubSourceSpec `json:"spec,omitempty"`

	//Status
	Status DockerHubSourceStatus `json:"status,omitempty"`
}

func (s *DockerHubSource) GetGroupVersionKind() schema.GroupVersionKind {
	return SchemeGroupVersion.WithKind("DockerHubSource")
}

var _ runtime.Object = (*DockerHubSource)(nil)
var _ resourcesemantics.GenericCRD = (*DockerHubSource)(nil)

const (
	dockerHubEventTypePrefix   = "dev.knative.source.dockerhub"
	dockerHubEventSourcePrefix = "https://hub.docker.com"
)

func DockerHubCloudEventsEventType(dhEventType string) string {
	return fmt.Sprintf("%s.%s", dockerHubEventTypePrefix, dhEventType)
}

func DockerHubEventSource(repoName string) string {
	return fmt.Sprintf("%s/r/%s", dockerHubEventSourcePrefix, repoName)
}

type DockerHubSourceSpec struct {
	// DisableAutoCallback flag allows users to make their own validation callback.
	//If unspecified this will default to false.
	DisableAutoCallback bool `json:"disableAutoCallback,omitempty"`

	// inherits duck/v1 SourceSpec, which currently provides:
	//  Sink - a reference to an object that will resolve to a domain name or
	//   a URI directly to use as the sink.
	//  CloudEventOverrides - defines overrides to control the output format
	//   and modifications of the event sent to the sink.
	duckv1.SourceSpec `json:",inline"`
}

type DockerHubSourceStatus struct {
	// inherits duck/v1 SourceStatus, which currently provides:
	// * ObservedGeneration - the 'Generation' of the Service that was last
	//   processed by the controller.
	// * Conditions - the latest available observations of a resource's current
	//   state.
	// * SinkURI - the current active sink URI that has been configured for the
	//   Source.
	duckv1.SourceStatus `json:",inline"`

	// AutoCallbackDisabled represents the state of itself.
	AutoCallbackDisabled bool `json:"autoCallbackDisabled,omitempty"`

	// Url holds the information needed to connect this up to receive events.
	// +optional
	Url *apis.URL `json:"url,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type DockerHubSourceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DockerHubSource `json:"items"`
}
