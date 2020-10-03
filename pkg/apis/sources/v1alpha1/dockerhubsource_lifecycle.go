package v1alpha1

import (
	"context"
	"encoding/json"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"knative.dev/pkg/apis"
	"knative.dev/pkg/apis/duck"
	duckv1 "knative.dev/pkg/apis/duck/v1"
	"knative.dev/pkg/logging"
	"knative.dev/pkg/tracker"

	"go.uber.org/zap"

	"github.com/tom24d/eventing-dockerhub/pkg/reconciler/source/resources"
)

const (
	// DockerHubSourceConditionReady has status True when the
	// DockerHubSource is ready to receive webhook and send events.
	DockerHubSourceConditionReady = apis.ConditionReady

	// DockerHubSourceConditionSinkBound has status True when the
	// DockerHubSource has been configured with a sink target.
	DockerHubSourceConditionSinkBound apis.ConditionType = "SinkBound"

	// DockerHubSourceConditionEndpointProvided has status True when the
	// backing knative service gets ready.
	DockerHubSourceConditionEndpointProvided apis.ConditionType = "EndpointProvided"
)

var dockerHubCondSet = apis.NewLivingConditionSet(
	DockerHubSourceConditionSinkBound,
	DockerHubSourceConditionEndpointProvided,
)

// GetConditionSet retrieves the condition set for this resource. Implements the KRShaped interface.
func (*DockerHubSource) GetConditionSet() apis.ConditionSet {
	return dockerHubCondSet
}

// GetCondition returns the condition currently associated with the given type, or nil.
func (s *DockerHubSourceStatus) GetCondition(t apis.ConditionType) *apis.Condition {
	return dockerHubCondSet.Manage(s).GetCondition(t)
}

// InitializeConditions sets relevant unset conditions to Unknown state.
func (s *DockerHubSourceStatus) InitializeConditions() {
	dockerHubCondSet.Manage(s).InitializeConditions()
}

// MarkEndpoint sets the URL endpoint that the source has been provided.
func (s *DockerHubSourceStatus) MarkEndpoint(uri *apis.URL) {
	s.URL = uri
	if len(uri.String()) > 0 {
		dockerHubCondSet.Manage(s).MarkTrue(DockerHubSourceConditionEndpointProvided)
	} else {
		dockerHubCondSet.Manage(s).MarkUnknown(DockerHubSourceConditionEndpointProvided,
			"EndpointEmpty", "Endpoint URL has resolved to empty.%s", "")
	}
}

// MarkNoSink sets the condition that the source does not have a sink configured.
func (s *DockerHubSourceStatus) MarkNoEndpoint(reason, messageFormat string, messageA ...interface{}) {
	dockerHubCondSet.Manage(s).MarkFalse(DockerHubSourceConditionEndpointProvided, reason, messageFormat, messageA...)
}

// IsReady returns true if the resource is ready overall.
func (s *DockerHubSourceStatus) IsReady() bool {
	return dockerHubCondSet.Manage(s).IsHappy()
}

// GetSubject implements psbinding.Bindable
func (dhs *DockerHubSource) GetSubject() tracker.Reference {
	return tracker.Reference{
		APIVersion: "serving.knative.dev/v1",
		Kind:       "Service",
		Namespace:  dhs.GetNamespace(),
		Selector: &metav1.LabelSelector{
			MatchLabels:      resources.Labels(dhs.GetName()),
		},
	}
}

// MarkBindingUnavailable marks the DockerHubSource's Binding Ready condition to False with
// the provided reason and message.
func (dhss *DockerHubSourceStatus) MarkBindingUnavailable(reason, message string) {
	dockerHubCondSet.Manage(dhss).MarkFalse(DockerHubSourceConditionSinkBound, reason, message)
}

// MarkBindingAvailable marks the DockerHubSource's Binding Ready condition to True.
func (dhss *DockerHubSourceStatus) MarkBindingAvailable() {
	dockerHubCondSet.Manage(dhss).MarkTrue(DockerHubSourceConditionSinkBound)
}

// SetObservedGeneration implements psbinding.BindableStatus
func (dhss *DockerHubSourceStatus) SetObservedGeneration(gen int64) {
	dhss.ObservedGeneration = gen
}

// GetBindingStatus implements psbinding.Bindable
func (dhs *DockerHubSource) GetBindingStatus() duck.BindableStatus {
	return &dhs.Status
}

// Do implements psbinding.Bindable
func (dhs *DockerHubSource) Do(ctx context.Context, ps *duckv1.WithPod) {
	// inject K_SINK, K_CE_OVERRIDES
	// TODO also AUTO_CALLBACK

	// First undo so that we can just unconditionally append below.
	dhs.Undo(ctx, ps)

	uri := GetSinkURI(ctx)
	if uri == nil {
		logging.FromContext(ctx).Errorf("No sink URI associated with context for %+v", dhs)
		return
	}

	var ceOverrides string
	if dhs.Spec.CloudEventOverrides != nil {
		if co, err := json.Marshal(dhs.Spec.SourceSpec.CloudEventOverrides); err != nil {
			logging.FromContext(ctx).Errorw(fmt.Sprintf("Failed to marshal CloudEventOverrides into JSON for %+v", dhs), zap.Error(err))
		} else if len(co) > 0 {
			ceOverrides = string(co)
		}
	}

	spec := ps.Spec.Template.Spec
	for i := range spec.InitContainers {
		spec.InitContainers[i].Env = append(spec.InitContainers[i].Env, corev1.EnvVar{
			Name:  "K_SINK",
			Value: uri.String(),
		})
		spec.InitContainers[i].Env = append(spec.InitContainers[i].Env, corev1.EnvVar{
			Name:  "K_CE_OVERRIDES",
			Value: ceOverrides,
		})
	}
	for i := range spec.Containers {
		spec.Containers[i].Env = append(spec.Containers[i].Env, corev1.EnvVar{
			Name:  "K_SINK",
			Value: uri.String(),
		})
		spec.Containers[i].Env = append(spec.Containers[i].Env, corev1.EnvVar{
			Name:  "K_CE_OVERRIDES",
			Value: ceOverrides,
		})
	}
}

func (dhs *DockerHubSource) Undo(ctx context.Context, ps *duckv1.WithPod) {
	// eliminate K_SINK, K_CE_OVERRIDES
	spec := ps.Spec.Template.Spec
	for i, c := range spec.InitContainers {
		if len(c.Env) == 0 {
			continue
		}
		env := make([]corev1.EnvVar, 0, len(spec.InitContainers[i].Env))
		for j, ev := range c.Env {
			switch ev.Name {
			case "K_SINK", "K_CE_OVERRIDES":
				continue
			default:
				env = append(env, spec.InitContainers[i].Env[j])
			}
		}
		spec.InitContainers[i].Env = env
	}
	for i, c := range spec.Containers {
		if len(c.Env) == 0 {
			continue
		}
		env := make([]corev1.EnvVar, 0, len(spec.Containers[i].Env))
		for j, ev := range c.Env {
			switch ev.Name {
			case "K_SINK", "K_CE_OVERRIDES":
				continue
			default:
				env = append(env, spec.Containers[i].Env[j])
			}
		}
		spec.Containers[i].Env = env
	}
}
