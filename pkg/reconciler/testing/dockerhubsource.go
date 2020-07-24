package testing

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	duckv1 "knative.dev/pkg/apis/duck/v1"

	sourcev1alpha1 "github.com/tom24d/eventing-dockerhub/pkg/apis/sources/v1alpha1"
)

type DockerHubSourceV1Alpha1Option func(source *sourcev1alpha1.DockerHubSource)

func NewDockerHubSourceV1Alpha1(name, namespace string, o ...DockerHubSourceV1Alpha1Option) *sourcev1alpha1.DockerHubSource {
	c := &sourcev1alpha1.DockerHubSource{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}
	for _, opt := range o {
		opt(c)
	}
	return c
}

func WithSinkV1A1(sink duckv1.Destination) DockerHubSourceV1Alpha1Option {
	return func(dhs *sourcev1alpha1.DockerHubSource) {
		dhs.Spec.Sink = sink
	}
}

func DisabledAutoCallback(flag bool) DockerHubSourceV1Alpha1Option {
	return func(dhs *sourcev1alpha1.DockerHubSource) {
		dhs.Spec.DisableAutoCallback = flag
	}
}
