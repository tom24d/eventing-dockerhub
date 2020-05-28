package main

import (
	"fmt"

	dhadapter "github.com/tom24d/eventing-dockerhub/pkg/adapter"
	"knative.dev/eventing/pkg/adapter/v2"
)

func main() {
	fmt.Println("Hello Receive Adapter v2")
	adapter.Main("dockerhub-source", dhadapter.NewEnv, dhadapter.NewAdapter)
}
