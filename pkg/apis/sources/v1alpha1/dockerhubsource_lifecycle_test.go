package v1alpha1

import (
	"context"
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"knative.dev/pkg/apis"
	"knative.dev/pkg/apis/duck"
	duckv1 "knative.dev/pkg/apis/duck/v1"
	"knative.dev/pkg/tracker"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

var _ = duck.VerifyType(&DockerHubSource{}, &duckv1.Conditions{})

func TestDockerHubSourceGetters(t *testing.T) {
	ns := "namespace"
	name := "name"

	d := &DockerHubSource{}
	d.SetNamespace(ns)
	d.SetName(name)

	want := tracker.Reference{
		APIVersion: "serving.knative.dev/v1",
		Kind:       "Service",
		Namespace:  ns,
		Selector: &metav1.LabelSelector{
			MatchLabels: map[string]string{
				"knative-eventing-source":      "dockerhubsource-controller",
				"knative-eventing-source-name": name,
				"receive-adapter":              "dockerhub",
			},
		},
	}

	if diff := cmp.Diff(want, d.GetSubject()); diff != "" {
		t.Errorf("%s: unexpected object (-want, +got) = %v", t.Name(), diff)
	}
	if diff := cmp.Diff(&d.Status, d.GetBindingStatus()); diff != "" {
		t.Errorf("%s: unexpected object (-want, +got) = %v", t.Name(), diff)
	}
}

func TestDockerHubSourceSetObsGen(t *testing.T) {
	d := DockerHubSource{}
	want := int64(1234)
	d.GetBindingStatus().SetObservedGeneration(want)
	if diff := cmp.Diff(want, d.Status.ObservedGeneration); diff != "" {
		t.Errorf("%s: unexpected observedGeneration (-want, +got) = %v", t.Name(), diff)
	}
}

func TestDockerHubSource_Do(t *testing.T) {
	sinkURI := &apis.URL{
		Scheme: "http",
		Host:   "thing.ns.svc.cluster.local",
		Path:   "/a/path",
	}

	overrides := duckv1.CloudEventOverrides{Extensions: map[string]string{"foo": "bar"}}

	tests := []struct {
		name string
		in   *duckv1.WithPod
		want *duckv1.WithPod
	}{{
		name: "nothing to add",
		in: &duckv1.WithPod{
			Spec: duckv1.WithPodSpec{
				Template: duckv1.PodSpecable{
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{{
							Name:  "blah",
							Image: "busybox",
							Env: []corev1.EnvVar{{
								Name:  "K_SINK",
								Value: sinkURI.String(),
							}, {
								Name:  "K_CE_OVERRIDES",
								Value: `{"extensions":{"foo":"bar"}}`,
							}},
						}},
					},
				},
			},
		},
		want: &duckv1.WithPod{
			Spec: duckv1.WithPodSpec{
				Template: duckv1.PodSpecable{
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{{
							Name:  "blah",
							Image: "busybox",
							Env: []corev1.EnvVar{{
								Name:  "K_SINK",
								Value: sinkURI.String(),
							}, {
								Name:  "K_CE_OVERRIDES",
								Value: `{"extensions":{"foo":"bar"}}`,
							}},
						}},
					},
				},
			},
		},
	}, {
		name: "fix the URI",
		in: &duckv1.WithPod{
			Spec: duckv1.WithPodSpec{
				Template: duckv1.PodSpecable{
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{{
							Name:  "blah",
							Image: "busybox",
							Env: []corev1.EnvVar{{
								Name:  "K_SINK",
								Value: "the wrong value",
							}, {
								Name:  "K_CE_OVERRIDES",
								Value: `{"extensions":{"wrong":"value"}}`,
							}},
						}},
					},
				},
			},
		},
		want: &duckv1.WithPod{
			Spec: duckv1.WithPodSpec{
				Template: duckv1.PodSpecable{
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{{
							Name:  "blah",
							Image: "busybox",
							Env: []corev1.EnvVar{{
								Name:  "K_SINK",
								Value: sinkURI.String(),
							}, {
								Name:  "K_CE_OVERRIDES",
								Value: `{"extensions":{"foo":"bar"}}`,
							}},
						}},
					},
				},
			},
		},
	}, {
		name: "lots to add",
		in: &duckv1.WithPod{
			Spec: duckv1.WithPodSpec{
				Template: duckv1.PodSpecable{
					Spec: corev1.PodSpec{
						InitContainers: []corev1.Container{{
							Name:  "setup",
							Image: "busybox",
						}},
						Containers: []corev1.Container{{
							Name:  "blah",
							Image: "busybox",
							Env: []corev1.EnvVar{{
								Name:  "FOO",
								Value: "BAR",
							}, {
								Name:  "BAZ",
								Value: "INGA",
							}},
						}, {
							Name:  "sidecar",
							Image: "busybox",
							Env: []corev1.EnvVar{{
								Name:  "BAZ",
								Value: "INGA",
							}},
						}},
					},
				},
			},
		},
		want: &duckv1.WithPod{
			Spec: duckv1.WithPodSpec{
				Template: duckv1.PodSpecable{
					Spec: corev1.PodSpec{
						InitContainers: []corev1.Container{{
							Name:  "setup",
							Image: "busybox",
							Env: []corev1.EnvVar{{
								Name:  "K_SINK",
								Value: sinkURI.String(),
							}, {
								Name:  "K_CE_OVERRIDES",
								Value: `{"extensions":{"foo":"bar"}}`,
							}},
						}},
						Containers: []corev1.Container{{
							Name:  "blah",
							Image: "busybox",
							Env: []corev1.EnvVar{{
								Name:  "FOO",
								Value: "BAR",
							}, {
								Name:  "BAZ",
								Value: "INGA",
							}, {
								Name:  "K_SINK",
								Value: sinkURI.String(),
							}, {
								Name:  "K_CE_OVERRIDES",
								Value: `{"extensions":{"foo":"bar"}}`,
							}},
						}, {
							Name:  "sidecar",
							Image: "busybox",
							Env: []corev1.EnvVar{{
								Name:  "BAZ",
								Value: "INGA",
							}, {
								Name:  "K_SINK",
								Value: sinkURI.String(),
							}, {
								Name:  "K_CE_OVERRIDES",
								Value: `{"extensions":{"foo":"bar"}}`,
							}},
						}},
					},
				},
			},
		},
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.in

			ctx := WithSinkURI(context.Background(), sinkURI)

			ds := &DockerHubSource{Spec: DockerHubSourceSpec{
				SourceSpec: duckv1.SourceSpec{
					CloudEventOverrides: &overrides,
				},
			}}
			ds.Do(ctx, got)

			if !cmp.Equal(got, test.want) {
				t.Error("Undo (-want, +got):", cmp.Diff(test.want, got))
			}
		})
	}
}

func TestSinkBindingDoNoURI(t *testing.T) {
	want := &duckv1.WithPod{
		Spec: duckv1.WithPodSpec{
			Template: duckv1.PodSpecable{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Name:  "blah",
						Image: "busybox",
						Env:   []corev1.EnvVar{},
					}},
				},
			},
		},
	}
	got := &duckv1.WithPod{
		Spec: duckv1.WithPodSpec{
			Template: duckv1.PodSpecable{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Name:  "blah",
						Image: "busybox",
						Env: []corev1.EnvVar{{
							Name:  "K_SINK",
							Value: "this should be removed",
						}, {
							Name:  "K_CE_OVERRIDES",
							Value: `{"extensions":{"tobe":"removed"}}`,
						}},
					}},
				},
			},
		},
	}

	ds := &DockerHubSource{}
	ds.Do(context.Background(), got)

	if !cmp.Equal(got, want) {
		t.Error("Undo (-want, +got):", cmp.Diff(want, got))
	}
}

func TestDockerHubSource_Undo(t *testing.T) {
	tests := []struct {
		name string
		in   *duckv1.WithPod
		want *duckv1.WithPod
	}{{
		name: "nothing to remove",
		in: &duckv1.WithPod{
			Spec: duckv1.WithPodSpec{
				Template: duckv1.PodSpecable{
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{{
							Name:  "blah",
							Image: "busybox",
						}},
					},
				},
			},
		},
		want: &duckv1.WithPod{
			Spec: duckv1.WithPodSpec{
				Template: duckv1.PodSpecable{
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{{
							Name:  "blah",
							Image: "busybox",
						}},
					},
				},
			},
		},
	}, {
		name: "lots to remove",
		in: &duckv1.WithPod{
			Spec: duckv1.WithPodSpec{
				Template: duckv1.PodSpecable{
					Spec: corev1.PodSpec{
						InitContainers: []corev1.Container{{
							Name:  "setup",
							Image: "busybox",
							Env: []corev1.EnvVar{{
								Name:  "FOO",
								Value: "BAR",
							}, {
								Name:  "K_SINK",
								Value: "http://localhost:8080",
							}, {
								Name:  "BAZ",
								Value: "INGA",
							}, {
								Name:  "K_CE_OVERRIDES",
								Value: `{"extensions":{"foo":"bar"}}`,
							}},
						}},
						Containers: []corev1.Container{{
							Name:  "blah",
							Image: "busybox",
							Env: []corev1.EnvVar{{
								Name:  "FOO",
								Value: "BAR",
							}, {
								Name:  "K_SINK",
								Value: "http://localhost:8080",
							}, {
								Name:  "BAZ",
								Value: "INGA",
							}, {
								Name:  "K_CE_OVERRIDES",
								Value: `{"extensions":{"foo":"bar"}}`,
							}},
						}, {
							Name:  "sidecar",
							Image: "busybox",
							Env: []corev1.EnvVar{{
								Name:  "K_SINK",
								Value: "http://localhost:8080",
							}, {
								Name:  "BAZ",
								Value: "INGA",
							}, {
								Name:  "K_CE_OVERRIDES",
								Value: `{"extensions":{"foo":"bar"}}`,
							}},
						}},
					},
				},
			},
		},
		want: &duckv1.WithPod{
			Spec: duckv1.WithPodSpec{
				Template: duckv1.PodSpecable{
					Spec: corev1.PodSpec{
						InitContainers: []corev1.Container{{
							Name:  "setup",
							Image: "busybox",
							Env: []corev1.EnvVar{{
								Name:  "FOO",
								Value: "BAR",
							}, {
								Name:  "BAZ",
								Value: "INGA",
							}},
						}},
						Containers: []corev1.Container{{
							Name:  "blah",
							Image: "busybox",
							Env: []corev1.EnvVar{{
								Name:  "FOO",
								Value: "BAR",
							}, {
								Name:  "BAZ",
								Value: "INGA",
							}},
						}, {
							Name:  "sidecar",
							Image: "busybox",
							Env: []corev1.EnvVar{{
								Name:  "BAZ",
								Value: "INGA",
							}},
						}},
					},
				},
			},
		},
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.in
			sb := &DockerHubSource{}
			sb.Undo(context.Background(), got)

			if !cmp.Equal(got, test.want) {
				t.Error("Undo (-want, +got):", cmp.Diff(test.want, got))
			}
		})
	}
}

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
		name: "mark bound",
		s: func() *DockerHubSourceStatus {
			s := &DockerHubSourceStatus{}
			s.InitializeConditions()
			s.MarkBindingAvailable()
			return s
		}(),
		want: false,
	}, {
		name: "mark endpoint",
		s: func() *DockerHubSourceStatus {
			s := &DockerHubSourceStatus{}
			s.InitializeConditions()
			s.MarkEndpoint(apis.HTTP("example"))
			return s
		}(),
		want: false,
	}, {
		name: "mark bound, then no bound",
		s: func() *DockerHubSourceStatus {
			s := &DockerHubSourceStatus{}
			s.InitializeConditions()
			s.MarkBindingAvailable()
			s.MarkBindingUnavailable("Testing", "")
			return s
		}(),
		want: false,
	}, {
		name: "mark endpoint, then no endpoint",
		s: func() *DockerHubSourceStatus {
			s := &DockerHubSourceStatus{}
			s.InitializeConditions()
			s.MarkEndpoint(apis.HTTP("example"))
			s.MarkNoEndpoint("Testing", "")
			return s
		}(),
		want: false,
	}, {
		name: "mark bound, endpoint",
		s: func() *DockerHubSourceStatus {
			s := &DockerHubSourceStatus{}
			s.InitializeConditions()
			s.MarkBindingAvailable()
			s.MarkEndpoint(apis.HTTP("example"))
			return s
		}(),
		want: true,
	}, {
		name: "mark unbound",
		s: func() *DockerHubSourceStatus {
			s := &DockerHubSourceStatus{}
			s.InitializeConditions()
			s.MarkBindingUnavailable("Test", "")
			return s
		}(),
		want: false,
	}, {
		name: "mark endpoint nil",
		s: func() *DockerHubSourceStatus {
			s := &DockerHubSourceStatus{}
			s.InitializeConditions()
			s.MarkEndpoint(nil)
			return s
		}(),
		want: false,
	}, {
		name: "mark endpoint nil, then endpoint",
		s: func() *DockerHubSourceStatus {
			s := &DockerHubSourceStatus{}
			s.InitializeConditions()
			s.MarkEndpoint(nil)
			s.MarkEndpoint(apis.HTTP("example"))
			return s
		}(),
		want: false,
	}, {
		name: "mark endpoint, bound, then no endpoint",
		s: func() *DockerHubSourceStatus {
			s := &DockerHubSourceStatus{}
			s.InitializeConditions()
			s.MarkEndpoint(apis.HTTP("example"))
			s.MarkBindingAvailable()
			s.MarkNoEndpoint("Testing", "")
			return s
		}(),
		want: false,
	}, {
		name: "mark endpoint, bound, then no bound",
		s: func() *DockerHubSourceStatus {
			s := &DockerHubSourceStatus{}
			s.InitializeConditions()
			s.MarkEndpoint(apis.HTTP("example"))
			s.MarkBindingAvailable()
			s.MarkBindingUnavailable("Testing", "")
			return s
		}(),
		want: false,
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
		name: "mark bound",
		s: func() *DockerHubSourceStatus {
			s := &DockerHubSourceStatus{}
			s.InitializeConditions()
			s.MarkBindingAvailable()
			return s
		}(),
		condQuery: DockerHubSourceConditionReady,
		want: &apis.Condition{
			Type:   DockerHubSourceConditionReady,
			Status: corev1.ConditionUnknown,
		},
	}, {
		name: "mark sink, then no sink",
		s: func() *DockerHubSourceStatus {
			s := &DockerHubSourceStatus{}
			s.InitializeConditions()
			s.MarkBindingAvailable()
			s.MarkBindingUnavailable("Testing", "hi")
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
		name: "mark endpoint",
		s: func() *DockerHubSourceStatus {
			s := &DockerHubSourceStatus{}
			s.InitializeConditions()
			s.MarkEndpoint(apis.HTTP("example"))
			return s
		}(),
		condQuery: DockerHubSourceConditionReady,
		want: &apis.Condition{
			Type:   DockerHubSourceConditionReady,
			Status: corev1.ConditionUnknown,
		},
	}, {
		name: "mark endpoint, then no endpoint",
		s: func() *DockerHubSourceStatus {
			s := &DockerHubSourceStatus{}
			s.InitializeConditions()
			s.MarkEndpoint(apis.HTTP("example"))
			s.MarkNoEndpoint("Testing", "hi%s", "")
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
		name: "mark unbound",
		s: func() *DockerHubSourceStatus {
			s := &DockerHubSourceStatus{}
			s.InitializeConditions()
			s.MarkBindingUnavailable("Test", "hi")
			return s
		}(),
		condQuery: DockerHubSourceConditionReady,
		want: &apis.Condition{
			Type:    DockerHubSourceConditionReady,
			Status:  corev1.ConditionFalse,
			Reason:  "Test",
			Message: "hi",
		},
	}, {
		name: "mark endpoint nil",
		s: func() *DockerHubSourceStatus {
			s := &DockerHubSourceStatus{}
			s.InitializeConditions()
			s.MarkEndpoint(nil)
			return s
		}(),
		condQuery: DockerHubSourceConditionReady,
		want: &apis.Condition{
			Type:    DockerHubSourceConditionReady,
			Status:  corev1.ConditionUnknown,
			Reason:  "EndpointEmpty",
			Message: "Endpoint URL has resolved to empty.",
		},
	}, {
		name: "mark endpoint nil, then endpoint",
		s: func() *DockerHubSourceStatus {
			s := &DockerHubSourceStatus{}
			s.InitializeConditions()
			s.MarkEndpoint(nil)
			s.MarkEndpoint(apis.HTTP("example"))
			return s
		}(),
		condQuery: DockerHubSourceConditionReady,
		want: &apis.Condition{
			Type:   DockerHubSourceConditionReady,
			Status: corev1.ConditionUnknown,
		},
	}, {
		name: "mark bound, endpoint",
		s: func() *DockerHubSourceStatus {
			s := &DockerHubSourceStatus{}
			s.InitializeConditions()
			s.MarkBindingAvailable()
			s.MarkEndpoint(apis.HTTP("example"))
			return s
		}(),
		condQuery: DockerHubSourceConditionReady,
		want: &apis.Condition{
			Type:   DockerHubSourceConditionReady,
			Status: corev1.ConditionTrue,
		},
	}, {
		name: "mark bound, endpoint, then unbound",
		s: func() *DockerHubSourceStatus {
			s := &DockerHubSourceStatus{}
			s.InitializeConditions()
			s.MarkBindingAvailable()
			s.MarkEndpoint(apis.HTTP("example"))
			s.MarkBindingUnavailable("Test", "hi")
			return s
		}(),
		condQuery: DockerHubSourceConditionReady,
		want: &apis.Condition{
			Type:    DockerHubSourceConditionReady,
			Status:  corev1.ConditionFalse,
			Reason:  "Test",
			Message: "hi",
		},
	}, {
		name: "mark bound, endpoint, then no endpoint",
		s: func() *DockerHubSourceStatus {
			s := &DockerHubSourceStatus{}
			s.InitializeConditions()
			s.MarkBindingAvailable()
			s.MarkEndpoint(apis.HTTP("example"))
			s.MarkNoEndpoint("Testing", "hi%s", "")
			return s
		}(),
		condQuery: DockerHubSourceConditionReady,
		want: &apis.Condition{
			Type:    DockerHubSourceConditionReady,
			Status:  corev1.ConditionFalse,
			Reason:  "Testing",
			Message: "hi",
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

func TestDockerHubSource_GetConditionSet(t *testing.T) {
	r := &DockerHubSource{}

	if got, want := r.GetConditionSet().GetTopLevelConditionType(), apis.ConditionReady; got != want {
		t.Errorf("GetTopLevelCondition=%v, want=%v", got, want)
	}
}
