package v1alpha1

import (
	"knative.dev/pkg/webhook/resourcesemantics"

	"k8s.io/apimachinery/pkg/runtime"
	"knative.dev/pkg/apis"
)

// Check that GitHubSource can be validated and can be defaulted.
var _ runtime.Object = (*DockerHubSource)(nil)
var _ resourcesemantics.GenericCRD = (*DockerHubSource)(nil)

const (
	// DockerHubSourceConditionReady has status True when the
	// DockerHubSource is ready to receive webhook and send events.
	DockerHubSourceConditionReady = apis.ConditionReady

	// DockerHubSourceConditionSinkProvided has status True when the
	// DockerHubSource has been configured with a sink target.
	DockerHubSourceConditionSinkProvided apis.ConditionType = "SinkProvided"

	// DockerHubSourceConditionWebhookChainActivated has status when
	// first webhook from Docker Hub is successfully received.
	DockerHubSourceConditionWebhookChainActivated apis.ConditionType = "WebhookChainActivated"

)

var dockerHubCondSet = apis.NewLivingConditionSet(
	DockerHubSourceConditionSinkProvided,
	DockerHubSourceConditionWebhookChainActivated)

// GetCondition returns the condition currently associated with the given type, or nil.
func (s *DockerHubSourceStatus) GetCondition(t apis.ConditionType) *apis.Condition {
	return dockerHubCondSet.Manage(s).GetCondition(t)
}

// InitializeConditions sets relevant unset conditions to Unknown state.
func (s *DockerHubSourceStatus) InitializeConditions() {
	dockerHubCondSet.Manage(s).InitializeConditions()
}

// MarkSink sets the condition that the source has a sink configured.
func (s *DockerHubSourceStatus) MarkSink(uri *apis.URL) {
	s.SinkURI = uri
	if len(uri.String()) > 0 {
		dockerHubCondSet.Manage(s).MarkTrue(DockerHubSourceConditionSinkProvided)
	}else {
		dockerHubCondSet.Manage(s).MarkUnknown(DockerHubSourceConditionSinkProvided, "SinkEmpty", "Sink has resolved to empty.%s", "")
	}
}

// MarkNoSink sets the condition that the source does not have a sink configured.
func (s *DockerHubSourceStatus) MarkNoSink(reason, messageFormat  string, messageA ...interface{}){
	dockerHubCondSet.Manage(s).MarkFalse(DockerHubSourceConditionSinkProvided, reason, messageFormat, messageA...)
}

// IsReady returns true if the resource is ready overall.
func (s *DockerHubSourceStatus) IsReady() bool {
	return dockerHubCondSet.Manage(s).IsHappy()
}