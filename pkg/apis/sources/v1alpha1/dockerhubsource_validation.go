package v1alpha1

import (
	"context"

	"knative.dev/pkg/apis"
)

// Validate validates DockerHubSource.
func (s *DockerHubSource) Validate(ctx context.Context) *apis.FieldError {
	var errs *apis.FieldError

	//validation for "spec" field.
	errs = errs.Also(s.Spec.Validate(ctx).ViaField("spec"))

	//errs is nil if everything is fine.
	return errs
}

// Validate validates SampleSourceSpec.
func (sspec *DockerHubSourceSpec) Validate(ctx context.Context) *apis.FieldError {
	//Add code for validation webhook for DockerHubSourceSpec.
	var errs *apis.FieldError

	//validation for sink field.
	if fe := sspec.Sink.Validate(ctx); fe != nil {
		errs = errs.Also(fe.ViaField("sink"))
	}

	return errs
}
