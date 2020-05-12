package v1alpha1

// import (
// 	"context"
// 	"testing"

// 	"github.com/google/go-cmp/cmp"
// )

// func TestSampleSourceDefaults(t *testing.T) {
// 	testCases := map[string]struct {
// 		initial  SampleSource
// 		expected SampleSource
// 	}{
// 		"nil spec": {
// 			initial: SampleSource{},
// 			expected: SampleSource{
// 				Spec: SampleSourceSpec{
// 					ServiceAccountName: "default",
// 					Interval:           "10s",
// 				},
// 			},
// 		},
// 	}
// 	for n, tc := range testCases {
// 		t.Run(n, func(t *testing.T) {
// 			tc.initial.SetDefaults(context.TODO())
// 			if diff := cmp.Diff(tc.expected, tc.initial); diff != "" {
// 				t.Fatalf("Unexpected defaults (-want, +got): %s", diff)
// 			}
// 		})
// 	}
// }
