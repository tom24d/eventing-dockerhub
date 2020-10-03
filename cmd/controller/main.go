package main

import (
	// The set of controllers this process runs.
	"github.com/tom24d/eventing-dockerhub/pkg/reconciler/binding"
	"github.com/tom24d/eventing-dockerhub/pkg/reconciler/source"

	// This defines the shared main for injected controllers.
	"knative.dev/pkg/injection/sharedmain"
)

const (
	component = "dockerhub-controller"
)

func main() {
	sharedmain.Main(component, source.NewController, binding.NewController)
}
