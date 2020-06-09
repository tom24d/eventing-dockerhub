package v1alpha1

import (
	"github.com/google/go-cmp/cmp/cmpopts"
	corev1 "k8s.io/api/core/v1"
	"testing"

	"knative.dev/pkg/apis"
	"knative.dev/pkg/apis/duck"
	duckv1 "knative.dev/pkg/apis/duck/v1"

	"github.com/google/go-cmp/cmp"
)

var _ = duck.VerifyType(&DockerHubSource{}, &duckv1.Conditions{})

func TestDockerHubSourceStatusIsReady(t *testing.T) {
	tests := []struct {
		name string
		s    *DockerHubSourceStatus
		want bool
	}{{
		name: "uninitialized",
		s:    &DockerHubSourceStatus{},
		want: false,
	}, {
		name: "initialized",
		s: func() *DockerHubSourceStatus {
			s := &DockerHubSourceStatus{}
			s.InitializeConditions()
			return s
		}(),
		want: false,
	}, {
		name: "mark sink",
		s: func() *DockerHubSourceStatus {
			s := &DockerHubSourceStatus{}
			s.InitializeConditions()
			s.MarkSink(apis.HTTP("example"))
			return s
		}(),
		want: true,
	}, {
		name: "mark sink, then no sink",
		s: func() *DockerHubSourceStatus {
			s := &DockerHubSourceStatus{}
			s.InitializeConditions()
			s.MarkSink(apis.HTTP("example"))
			s.MarkNoSink("Testing", "")
			return s
		}(),
		want: false,
	}, {
		name: "mark sink nil",
		s: func() *DockerHubSourceStatus {
			s := &DockerHubSourceStatus{}
			s.InitializeConditions()
			s.MarkSink(nil)
			return s
		}(),
		want: false,
	}, {
		name: "mark sink nil, then sink",
		s: func() *DockerHubSourceStatus {
			s := &DockerHubSourceStatus{}
			s.InitializeConditions()
			s.MarkSink(nil)
			s.MarkSink(apis.HTTP("example"))
			return s
		}(),
		want: true,
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.s.IsReady()
			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("%s: unexpected condition (-want, +got) = %v", test.name, diff)
			}
		})
	}
}

func TestDockerHubSourceStatusGetCondition(t *testing.T) {
	tests := []struct {
		name      string
		s         *DockerHubSourceStatus
		condQuery apis.ConditionType
		want      *apis.Condition
	}{{
		name:      "uninitialized",
		s:         &DockerHubSourceStatus{},
		condQuery: DockerHubSourceConditionReady,
		want:      nil,
	}, {
		name: "initialized",
		s: func() *DockerHubSourceStatus {
			s := &DockerHubSourceStatus{}
			s.InitializeConditions()
			return s
		}(),
		condQuery: DockerHubSourceConditionReady,
		want: &apis.Condition{
			Type:   DockerHubSourceConditionReady,
			Status: corev1.ConditionUnknown,
		},
	}, {
		name: "mark sink",
		s: func() *DockerHubSourceStatus {
			s := &DockerHubSourceStatus{}
			s.InitializeConditions()
			s.MarkSink(apis.HTTP("example"))
			return s
		}(),
		condQuery: DockerHubSourceConditionReady,
		want: &apis.Condition{
			Type:   DockerHubSourceConditionReady,
			Status: corev1.ConditionTrue,
		},
	}, {
		name: "mark sink, then no sink",
		s: func() *DockerHubSourceStatus {
			s := &DockerHubSourceStatus{}
			s.InitializeConditions()
			s.MarkSink(apis.HTTP("example"))
			s.MarkNoSink("Testing", "hi%s", "")
			return s
		}(),
		condQuery: DockerHubSourceConditionReady,
		want: &apis.Condition{
			Type:    DockerHubSourceConditionReady,
			Status:  corev1.ConditionFalse,
			Reason:  "Testing",
			Message: "hi",
		},
	}, {
		name: "mark sink nil",
		s: func() *DockerHubSourceStatus {
			s := &DockerHubSourceStatus{}
			s.InitializeConditions()
			s.MarkSink(nil)
			return s
		}(),
		condQuery: DockerHubSourceConditionReady,
		want: &apis.Condition{
			Type:    DockerHubSourceConditionReady,
			Status:  corev1.ConditionUnknown,
			Reason:  "SinkEmpty",
			Message: "Sink has resolved to empty.",
		},
	}, {
		name: "mark sink nil, then sink",
		s: func() *DockerHubSourceStatus {
			s := &DockerHubSourceStatus{}
			s.InitializeConditions()
			s.MarkSink(nil)
			s.MarkSink(apis.HTTP("example"))
			return s
		}(),
		condQuery: DockerHubSourceConditionReady,
		want: &apis.Condition{
			Type:   DockerHubSourceConditionReady,
			Status: corev1.ConditionTrue,
		},
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.s.GetCondition(test.condQuery)
			ignoreTime := cmpopts.IgnoreFields(apis.Condition{},
				"LastTransitionTime", "Severity")
			if diff := cmp.Diff(test.want, got, ignoreTime); diff != "" {
				t.Errorf("unexpected condition (-want, +got) = %v", diff)
			}
		})
	}
}
