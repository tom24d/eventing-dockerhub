package main

import (
	// The set of controllers this controller process runs.
	"github.com/tom24d/eventing-dockerhub/pkg/reconciler/source"

	// This defines the shared main for injected controllers.
	"knative.dev/pkg/injection/sharedmain"
)

const (
	component = "dockerhub_controller"
)

func main() {
	sharedmain.Main(component, source.NewController)
}
