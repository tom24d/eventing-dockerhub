package main

import (
	"log"

	"k8s.io/client-go/rest"

	"knative.dev/pkg/injection"
	"knative.dev/pkg/logging"
	_ "knative.dev/pkg/system/testing"

	loggerVent "knative.dev/eventing/test/lib/recordevents/logger_vent"
	"knative.dev/eventing/test/lib/recordevents/observer"
	recorderVent "knative.dev/eventing/test/lib/recordevents/recorder_vent"
	"knative.dev/eventing/test/test_images"
)

func main() {
	cfg, err := rest.InClusterConfig()
	if err != nil {
		log.Fatal("Error while reading the cfg", err)
	}

	//nolint // nil ctx is fine here, look at the code of EnableInjectionOrDie
	ctx, _ := injection.EnableInjectionOrDie(nil, cfg)
	ctx = test_images.ConfigureLogging(ctx, "recordevents")

	obs := observer.NewFromEnv(ctx,
		loggerVent.Logger(logging.FromContext(ctx).Infof),
		recorderVent.NewFromEnv(ctx),
	)

	if err := obs.Start(ctx); err != nil {
		panic(err)
	}
}
