package v1alpha1


import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"knative.dev/pkg/webhook/resourcesemantics"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"knative.dev/pkg/apis"
	duckv1 "knative.dev/pkg/apis/duck/v1"
)

// +genclient
// +genreconciler
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DockerHubSource is the Schema for the dockerhubsources API
// +k8s:openapi-gen=true
type DockerHubSource struct {

	//Metadata
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	//Spec
	Spec DockerHubSourceSpec `json:"spec,omitempty"`

	//Status
	// +optional
	Status DockerHubSourceStatus `json:"status,omitempty"`

}

func (s *DockerHubSource) GetGroupVersionKind() schema.GroupVersionKind {
	return SchemeGroupVersion.WithKind("DockerHubSource")
}

var _ runtime.Object = (*DockerHubSource)(nil)
var _ resourcesemantics.GenericCRD = (*DockerHubSource)(nil)

const (
	dockerHubEventTypePrefix = "dev.knative.source.dockerhub"
	dockerHubSourcePrefix = "https://hub.docker.com"
	//owner and repo?
)

func DockerHubCloudEventsEventType () string {
	return dockerHubEventTypePrefix + ".webhook"
}

func DockerHubCloudEventsSource() string {
	return dockerHubSourcePrefix //+ owner and repo?
}

type DockerHubSourceSpec struct {
	// OwnerAndRepository contains DockerHub owner/org and repository to
	// receive events from. The repository may be left off to receive
	// events from an entire organization.
	// Examples:
	//  myuser/project
	// +kubebuilder:validation:MinLength=1
	OwnerAndRepository string `json:"ownerAndRepository"`

	// +optional
	ServiceAccountName string `json:"serviceAccountName,omitempty"`

	// inherits duck/v1 SourceSpec, which currently provides:
	//  Sink - a reference to an object that will resolve to a domain name or
	//   a URI directly to use as the sink.
	//  CloudEventOverrides - defines overrides to control the output format
	//   and modifications of the event sent to the sink.
	duckv1.SourceSpec `json:",inline"`
}

type DockerHubSourceStatus struct {
	duckv1.Status `json:",inline"`
	// +optional
	SinkURI *apis.URL `json:"sinkUri,omitempty"`
}


// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type DockerHubSourceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`

	Items []DockerHubSource `json:"items"`
}


func init() {
	SchemeBuilder.Register(&DockerHubSource{}, &DockerHubSourceList{})
}