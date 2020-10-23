package v1alpha1

import (
	"context"
)

func (d *DockerHubSource) SetDefaults(ctx context.Context) {
	if d.Spec.Sink.Ref.Namespace == "" {
		// default the sink namespaces to the namespace of DockerHubSource.
		d.Spec.Sink.Ref.Namespace = d.GetNamespace()
	}
}
