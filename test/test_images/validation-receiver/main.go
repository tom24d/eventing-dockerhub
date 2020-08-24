package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"time"

	"github.com/tom24d/eventing-dockerhub/pkg/adapter/resources"
	"github.com/tom24d/eventing-dockerhub/test/e2e/helpers"
)

var (
	patient int
)

func init() {
	flag.IntVar(&patient, "patient", 30, "The seconds to wait")
}

// wait for sec patient, exit 0 or 1
func main() {
	flag.Parse()

	received := make(chan int, 1)
	defer close(received)

	h := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("payload received."))
		reqDump, _ := httputil.DumpRequest(r, true)
		log.Printf("incoming request: %s", string(reqDump))
		_, err := resources.Parse(r)
		if err != nil {
			log.Println(err.Error())
			received <- 1
		} else {
			received <- 0
		}
	}
	r := http.NewServeMux()
	r.HandleFunc("/", h)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", helpers.ValidationReceivePort),
		Handler: r,
	}

	go func() {
		log.Println("start listening to validation report...")
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("failed to wake server up: %v", err)
		}
	}()

	counter := 0
	ticker := time.NewTicker(time.Second)

	for {
		select {
		case <-ticker.C:
			counter += 1
			if counter > patient {
				cancel()
				log.Fatalln("exhausted to wait for validation report.")
			}
		case exitCode := <-received:
			server.Shutdown(ctx)
			cancel()
			os.Exit(exitCode)
		}
	}
}
