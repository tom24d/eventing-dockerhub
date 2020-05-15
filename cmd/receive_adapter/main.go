package main

import (
	myadapter "github.com/tom24d/eventing-dockerhub/pkg/adapter"
	"knative.dev/eventing/pkg/adapter/v2"
)

func main() {
	// TODO impl to read env.
	adapter.Main("dockerhub-source", myadapter.NewEnv, myadapter.NewAdapter)
}
