package v1alpha1

import (
	"context"
	"testing"
)

func TestDockerhubSourceDefaults(t *testing.T) {
	// None yet
	d := &DockerHubSource{}
	d.SetDefaults(context.TODO())
	d.Spec.SetDefaults(context.TODO())
}
