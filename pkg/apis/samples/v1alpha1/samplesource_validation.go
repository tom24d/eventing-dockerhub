package v1alpha1

// import (
// 	"context"
// 	"time"

// 	"knative.dev/pkg/apis"
// )

// // Validate validates SampleSource.
// func (s *SampleSource) Validate(ctx context.Context) *apis.FieldError {
// 	var errs *apis.FieldError

// 	//example: validation for "spec" field.
// 	errs = errs.Also(s.Spec.Validate(ctx).ViaField("spec"))

// 	//errs is nil if everything is fine.
// 	return errs
// }

// // Validate validates SampleSourceSpec.
// func (sspec *SampleSourceSpec) Validate(ctx context.Context) *apis.FieldError {
// 	//Add code for validation webhook for SampleSourceSpec.
// 	var errs *apis.FieldError

// 	//example: validation for sink field.
// 	if fe := sspec.Sink.Validate(ctx); fe != nil {
// 		errs = errs.Also(fe.ViaField("sink"))
// 	}

// 	//example: validation for interval field.
// 	if _, fe := time.ParseDuration(sspec.Interval); fe != nil {
// 		errs = errs.Also(apis.ErrInvalidValue(fe, "interval"))
// 	}

// 	//example: validation for serviceAccountName field.
// 	if sspec.ServiceAccountName == "" {
// 		errs = errs.Also(apis.ErrMissingField("serviceAccountName"))
// 	}

// 	return errs
// }
