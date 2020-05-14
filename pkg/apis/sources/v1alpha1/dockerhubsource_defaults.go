package v1alpha1

import (
	"context"
)

func (d *DockerHubSource) SetDefaults(ctx context.Context){
	d.Spec.SetDefaults(ctx)
}

func (ds *DockerHubSourceSpec) SetDefaults(ctx context.Context){
	//initialize here if needed
}