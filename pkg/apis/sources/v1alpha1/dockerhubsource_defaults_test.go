package v1alpha1

import (
	"context"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	duckv1 "knative.dev/pkg/apis/duck/v1"

	"github.com/google/go-cmp/cmp"
)

func TestDockerhubSourceDefaults(t *testing.T) {
	tests := []struct {
		name string
		in   *DockerHubSource
		want *DockerHubSource
	}{{
		name: "namespace is defaulted",
		in: &DockerHubSource{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "name",
				Namespace: "source-namespace",
			},
			Spec: DockerHubSourceSpec{
				SourceSpec: duckv1.SourceSpec{
					Sink: duckv1.Destination{
						Ref: &duckv1.KReference{
							Kind:       "Pod",
							Name:       "podname",
							APIVersion: "apps/v1",
						},
					},
				},
			},
		},
		want: &DockerHubSource{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "name",
				Namespace: "source-namespace",
			},
			Spec: DockerHubSourceSpec{
				SourceSpec: duckv1.SourceSpec{
					Sink: duckv1.Destination{
						Ref: &duckv1.KReference{
							Kind:       "Pod",
							Namespace:  "source-namespace",
							Name:       "podname",
							APIVersion: "apps/v1",
						},
					},
				},
			},
		},
	}, {
		name: "no sink, given namespace",
		in: &DockerHubSource{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "name",
				Namespace: "source-namespace",
			},
			Spec: DockerHubSourceSpec{
				SourceSpec: duckv1.SourceSpec{
					Sink: duckv1.Destination{
						Ref: &duckv1.KReference{
							Kind:       "Pod",
							Namespace:  "ref-namespace",
							Name:       "podname",
							APIVersion: "apps/v1",
						},
					},
				},
			},
		},
		want: &DockerHubSource{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "name",
				Namespace: "source-namespace",
			},
			Spec: DockerHubSourceSpec{
				SourceSpec: duckv1.SourceSpec{
					Sink: duckv1.Destination{
						Ref: &duckv1.KReference{
							Kind:       "Pod",
							Namespace:  "ref-namespace",
							Name:       "podname",
							APIVersion: "apps/v1",
						},
					},
				},
			},
		},
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.in
			got.SetDefaults(context.Background())
			if !cmp.Equal(test.want, got) {
				t.Errorf("SetDefaults (-want, +got) = %v", cmp.Diff(test.want, got))
			}
		})
	}
}
