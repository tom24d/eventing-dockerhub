package v1alpha1

import (
	"context"

	"knative.dev/pkg/apis"
)

// Validate validates SampleSource.
func (s *DockerHubSource) Validate(ctx context.Context) *apis.FieldError {
	var errs *apis.FieldError

	//example: validation for "spec" field.
	//errs = errs.Also(s.Spec.Validate(ctx).ViaField("spec"))

	//errs is nil if everything is fine.
	return errs
}

// Validate validates SampleSourceSpec.
func (sspec *DockerHubSourceSpec) Validate(ctx context.Context) *apis.FieldError {
	//Add code for validation webhook for SampleSourceSpec.
	var errs *apis.FieldError


	return errs
}
