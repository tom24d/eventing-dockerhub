package main

import (
	"context"
	"log"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	dockerhub "gopkg.in/go-playground/webhooks.v5/docker"

	"github.com/tom24d/eventing-dockerhub/pkg/adapter/resources"
)

func display(event cloudevents.Event) {
	data := &dockerhub.BuildPayload{}
	if err := event.DataAs(data); err != nil {
		log.Fatalf("Got Data Error: %s\n", err.Error())
		return
	}
	log.Printf("Got Data: %+v\n", data)

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
			log.Fatalf("failed to send validation callback: %v", err)
		} else {
			log.Printf("callback is sent from callback-display: %v", callbackData)
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
