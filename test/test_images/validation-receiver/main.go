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
		fmt.Fprintln(w, "payload received.")
		reqDump, _ := httputil.DumpRequest(r, true)
		log.Printf("incoming request: %s", string(reqDump))

		exitCode := 0
		_, err := resources.Parse(r)
		if err != nil {
			log.Println(err.Error())
			exitCode = 1
		}

		ch := make(chan struct{})
		go func() {
			ch <- struct{}{}

			// ensure the goroutine of the handler has been finished.
			time.Sleep(time.Second)

			received <- exitCode
		}()

		<-ch
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
