package main

import (
	"context"
	"log"
	"os"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	dockerhub "gopkg.in/go-playground/webhooks.v5/docker"

	"github.com/tom24d/eventing-dockerhub/pkg/adapter/resources"
)

// This image shows received payload and send validation webhook.

func display(event cloudevents.Event) {
	data := &dockerhub.BuildPayload{}
	if err := event.DataAs(data); err != nil {
		log.Printf("Got Data Error: %s\n", err.Error())
		return
	}
	log.Printf("Got Data: %+v\n", data)

	// ensure RA finished processing incoming webhook.
	time.Sleep(time.Second * 5)
	// if RA send validation webhook behalf, operation below should fail.

	if data.CallbackURL != "" {
		message := "Event has been sent successfully to the sink."
		callbackData := &resources.CallbackPayload{
			State:       resources.StatusSuccess,
			Description: message,
			Context:     "from sink display",
			TargetURL:   "",
		}

		err := callbackData.EmitValidationCallback(data.CallbackURL)
		if err != nil {
			log.Printf("failed to send validation callback: %v", err)
			os.Exit(1)
		} else {
			log.Printf("callback is sent from callback-display: %v", callbackData)
			os.Exit(0)
		}
	}
}

func main() {
	c, err := cloudevents.NewDefaultClient()
	if err != nil {
		log.Fatal("Failed to create client, ", err)
	}

	log.Fatal(c.StartReceiver(context.Background(), display))
}
