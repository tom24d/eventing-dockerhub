package v1alpha1

import (
	"context"

	"knative.dev/pkg/apis"
)

// sinkURIKey is used as the key for associating information
// with a context.Context.
type sinkURIKey struct{}

// WithSinkURI notes on the context for binding that the resolved SinkURI
// is the provided apis.URL.
func WithSinkURI(ctx context.Context, uri *apis.URL) context.Context {
	return context.WithValue(ctx, sinkURIKey{}, uri)
}

// GetSinkURI accesses the apis.URL for the Sink URI that has been associated
// with this context.
func GetSinkURI(ctx context.Context) *apis.URL {
	value := ctx.Value(sinkURIKey{})
	if value == nil {
		return nil
	}
	return value.(*apis.URL)
}
