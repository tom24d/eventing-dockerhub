package main

import (
	"bytes"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/tom24d/eventing-dockerhub/test/resources"
)

var (
	sink    string
	payload string
)

func init() {
	flag.StringVar(&sink, resources.ArgSink, "", "The sink url for the message destination.")
	flag.StringVar(&payload, resources.ArgPayload, "", "Payload JSON encoded")
}

func send() int {
	flag.Parse()

	if sink == "" {
		return 1
	}

	if payload == "" {
		return 1
	}

	resp, err := http.Post(sink, "application/json", bytes.NewBufferString(payload))
	if err != nil {
		log.Fatalf("httpPost error: %v", err)
		return 1
	}
	if c := resp.StatusCode; http.StatusOK <= c && c < http.StatusBadRequest {
		return 0
	} else {
		log.Fatalf("exit with status code: %d", c)
		return 1
	}
}

func main() {
	os.Exit(send())
}
