package resources

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"knative.dev/pkg/kmeta"

	v1 "knative.dev/serving/pkg/apis/serving/v1"

	sourcesv1alpha1 "github.com/tom24d/eventing-dockerhub/pkg/apis/sources/v1alpha1"
)

// ServiceArgs contains what the kservice needs.
type ServiceArgs struct {
	ReceiveAdapterImage string
	Source              *sourcesv1alpha1.DockerHubSource
	EventSource         string
	AdditionalEnvs      []corev1.EnvVar
}

// MakeService generates, but does not create, a Service for the given
// DockerHubSource.
func MakeService(args *ServiceArgs) *v1.Service {
	labels := map[string]string{
		"receive-adapter": "dockerhub",
	}
	sinkURI := args.Source.Status.SinkURI
	containerArgs := []string{fmt.Sprintf("--sink=%s", sinkURI.String())}
	return &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: fmt.Sprintf("%s-", args.Source.Name),
			Namespace:    args.Source.Namespace,
			Labels:       labels,
			OwnerReferences: []metav1.OwnerReference{
				*kmeta.NewControllerRef(args.Source),
			},
		},
		Spec: v1.ServiceSpec{
			ConfigurationSpec: v1.ConfigurationSpec{
				Template: v1.RevisionTemplateSpec{
					Spec: v1.RevisionSpec{
						PodSpec: corev1.PodSpec{
							Containers: []corev1.Container{{
								Image: args.ReceiveAdapterImage,
								Env: append(
									makeEnv(args.EventSource),
									args.AdditionalEnvs...,
								),
								Args: containerArgs,
							}},
						},
					},
				},
			},
		},
	}
}

func makeEnv(eventSource string) []corev1.EnvVar {
	return []corev1.EnvVar{{
		Name:  "EVENT_SOURCE",
		Value: eventSource,
	}, {
		Name: "NAMESPACE",
		ValueFrom: &corev1.EnvVarSource{
			FieldRef: &corev1.ObjectFieldSelector{
				FieldPath: "metadata.namespace",
			},
		},
	}, {
		Name: "NAME",
		ValueFrom: &corev1.EnvVarSource{
			FieldRef: &corev1.ObjectFieldSelector{
				FieldPath: "metadata.name",
			},
		},
	}, {
		Name:  "METRICS_DOMAIN",
		Value: "knative.dev/eventing",
	}}
}
