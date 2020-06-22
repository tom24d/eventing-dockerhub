package main

import (
	"flag"
	"net/http"
	"os"
	"time"

	"github.com/tom24d/eventing-dockerhub/pkg/adapter/resources"
)

var (
	patient int
)

func init() {
	flag.IntVar(&patient, "patient", 30, "The seconds to wait")
}

// wait for sec PATIENT_DURATION, exit 0 or 1
func main() {
	flag.Parse()

	h := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
		_, err := resources.Parse(r)
		if err != nil {
			os.Exit(1)
		}
		os.Exit(0)
	}
	r := http.NewServeMux()
	r.HandleFunc("/", h)

	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go server.ListenAndServe()

	counter := 0

	ticker := time.NewTicker(time.Second)
	for {
		<- ticker.C
		counter += 1
		if counter > patient {
			os.Exit(1)
		}
	}
}
