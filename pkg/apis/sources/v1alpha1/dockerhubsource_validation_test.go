package v1alpha1

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"knative.dev/pkg/apis"
	"knative.dev/pkg/webhook/resourcesemantics"
)

func TestDockerHubSourceValidation(t *testing.T) {
	testCases := map[string]struct {
		cr   resourcesemantics.GenericCRD
		want *apis.FieldError
	}{
		"nil spec": {
			cr: &DockerHubSource{
				Spec: DockerHubSourceSpec{},
			},
			want: func() *apis.FieldError {
				var errs *apis.FieldError

				feSink := apis.ErrGeneric("expected at least one, got none", "ref", "uri")
				feSink = feSink.ViaField("sink").ViaField("spec")
				errs = errs.Also(feSink)

				return errs
			}(),
		},
	}

	for n, test := range testCases {
		t.Run(n, func(t *testing.T) {
			got := test.cr.Validate(context.Background())
			if diff := cmp.Diff(test.want.Error(), got.Error()); diff != "" {
				t.Errorf("%s: validate (-want, +got) = %v", n, diff)
			}
		})
	}
}
