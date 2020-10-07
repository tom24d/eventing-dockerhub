package v1alpha1

import (
	"context"
	"testing"

	"knative.dev/pkg/apis"
)

func TestGetSinkURI(t *testing.T) {
	ctx := context.Background()

	if uri := GetSinkURI(ctx); uri != nil {
		t.Errorf("GetSinkURI() = %v, wanted nil", uri)
	}

	want := &apis.URL{
		Scheme: "https",
		Host:   "host.host",
		Path:   "/path/path/path",
	}
	ctx = WithSinkURI(ctx, want)

	if got := GetSinkURI(ctx); got != want {
		t.Errorf("GetSinkURI() = %v, wanted %v", got, want)
	}
}
