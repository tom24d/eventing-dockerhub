package main

import (
	"context"
	"fmt"
	"log"

	cloudevents "github.com/cloudevents/sdk-go/v2"

	"github.com/tom24d/eventing-dockerhub/pkg/adapter/resources"
)

func display(event cloudevents.Event) {
	data := &resources.CallbackPayload{}
	if err := event.DataAs(data); err != nil {
		fmt.Printf("Got Data Error: %s\n", err.Error())
	}

	callbackURL := data.TargetURL
	if callbackURL != "" {
		message := "Event has been sent successfully to the sink."
		callbackData := &resources.CallbackPayload{
			State:       resources.StatusSuccess,
			Description: message,
			Context:     "from sink display",
			TargetURL:   "",
		}

		err := callbackData.EmitValidationCallback(callbackURL)
		if err != nil {
			fmt.Printf("failed to send validation callback: %v", err)
		} else {
			fmt.Printf("callback is sent from callback-display: %v", callbackData)
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
