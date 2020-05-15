package main

import (
	dhadapter "github.com/tom24d/eventing-dockerhub/pkg/adapter"
	"knative.dev/eventing/pkg/adapter/v2"
)

func main() {
	adapter.Main("dockerhub-source", dhadapter.NewEnv, dhadapter.NewAdapter)
}
