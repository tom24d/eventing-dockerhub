package resources

import (
	"context"
	"fmt"
	"strconv"

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
	Context             context.Context
}

// MakeService generates, but does not create, a Service for the given
// DockerHubSource.
func MakeService(args *ServiceArgs) *v1.Service {
	labels := Labels(args.Source.Name)

	ksvc := &v1.Service{
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
								Env:   args.GetEnv(),
							}},
						},
					},
				},
			},
		},
	}
	return ksvc
}

// GetEnv return EnvVar used by the ksvc.
func (args *ServiceArgs) GetEnv() []corev1.EnvVar {

	envs := []corev1.EnvVar{{
		Name:  "EVENT_SOURCE",
		Value: args.EventSource,
	}, {
		Name: "METRICS_DOMAIN",
		Value: "knative.dev/eventing",
	}, {
		Name:  "NAMESPACE",
		Value: args.Source.Namespace,
	}, {
		Name:  "DISABLE_AUTO_CALLBACK",
		Value: strconv.FormatBool(args.Source.Spec.DisableAutoCallback),
	}, {
		Name: "METRICS_PROMETHEUS_PORT",
		Value: "9089",
	}}

	return append(
		envs,
		args.AdditionalEnvs...,
	)
}
