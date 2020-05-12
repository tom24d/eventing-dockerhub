package v1alpha1

// import (
// 	"context"
// 	"testing"
// 	"time"

// 	"github.com/google/go-cmp/cmp"
// 	"knative.dev/pkg/webhook/resourcesemantics"

// 	"knative.dev/pkg/apis"
// )

// func TestSampleSourceValidation(t *testing.T) {
// 	testCases := map[string]struct {
// 		cr   resourcesemantics.GenericCRD
// 		want *apis.FieldError
// 	}{
// 		"nil spec": {
// 			cr: &SampleSource{
// 				Spec: SampleSourceSpec{},
// 			},
// 			want: func() *apis.FieldError {
// 				var errs *apis.FieldError

// 				feSink := apis.ErrGeneric("expected at least one, got none", "ref", "uri")
// 				feSink = feSink.ViaField("sink").ViaField("spec")
// 				errs = errs.Also(feSink)

// 				_, timeErr := time.ParseDuration("")
// 				feInterval := apis.ErrInvalidValue(timeErr, "interval")
// 				feInterval = feInterval.ViaField("spec")
// 				errs = errs.Also(feInterval)

// 				feServiceAccountName := apis.ErrMissingField("serviceAccountName")
// 				feServiceAccountName = feServiceAccountName.ViaField("spec")
// 				errs = errs.Also(feServiceAccountName)

// 				return errs
// 			}(),
// 		},
// 	}

// 	for n, test := range testCases {
// 		t.Run(n, func(t *testing.T) {
// 			got := test.cr.Validate(context.Background())
// 			if diff := cmp.Diff(test.want.Error(), got.Error()); diff != "" {
// 				t.Errorf("%s: validate (-want, +got) = %v", n, diff)
// 			}
// 		})
// 	}
// }
