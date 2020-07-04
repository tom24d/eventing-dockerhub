package adapter

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	//knative.dev imports
	"knative.dev/eventing/pkg/adapter/v2"
	"knative.dev/pkg/logging"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
	dockerhub "gopkg.in/go-playground/webhooks.v5/docker"

	"github.com/tom24d/eventing-dockerhub/pkg/adapter/resources"
	"github.com/tom24d/eventing-dockerhub/pkg/apis/sources/v1alpha1"
)

type envConfig struct {
	// Include the standard adapter.EnvConfig used by all adapters.
	adapter.EnvConfig

	// Port to listen incoming connections
	Port string `envconfig:"PORT" default:"8080"`

	// DisableAutoCallback represents whether Receive Adapter always report its result to given callbackUrl.
	DisableAutoCallback bool `envconfig:"DISABLE_AUTO_CALLBACK" default:"false"`

	// EventSource is namespace/name of backing ksvc resource.
	EventSource string `envconfig:"EVENT_SOURCE"`
}

func NewEnv() adapter.EnvConfigAccessor { return &envConfig{} }

// Adapter converts incoming GitHub webhook events to CloudEvents
type Adapter struct {
	client         cloudevents.Client
	logger         *zap.SugaredLogger
	port           string
	autoValidation bool
	eventSource    string
}

// NewAdapter creates an adapter to convert incoming DockerHub webhook events to CloudEvents and
// then sends them to the specified Sink
func NewAdapter(ctx context.Context, aEnv adapter.EnvConfigAccessor, ceClient cloudevents.Client) adapter.Adapter {
	env := aEnv.(*envConfig) // Will always be our own envConfig type
	logger := logging.FromContext(ctx)
	return &Adapter{
		client:         ceClient,
		logger:         logger,
		port:           env.Port,
		autoValidation: !env.DisableAutoCallback,
		eventSource:    env.EventSource,
	}
}

// Start runs the adapter.
// Returns if stopCh is closed or Send() returns an error.
func (a *Adapter) Start(ctx context.Context) error {
	return a.start(ctx.Done())
}

func (a *Adapter) start(stopCh <-chan struct{}) error {
	done := make(chan bool, 1)
	hook, _ := dockerhub.New()

	server := &http.Server{
		Addr:    ":" + a.port,
		Handler: a.newRouter(hook),
	}

	go gracefulShutdown(server, a.logger, stopCh, done)

	a.logger.Infof("Server is ready to handle requests at %s", server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("could not listen on %s: %v", server.Addr, err)
	}

	<-done
	a.logger.Infof("Server stopped")
	return nil
}

func gracefulShutdown(server *http.Server, logger *zap.SugaredLogger, stopCh <-chan struct{}, done chan<- bool) {
	<-stopCh
	logger.Info("Server is shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	server.SetKeepAlivesEnabled(false)
	if err := server.Shutdown(ctx); err != nil {
		logger.Fatalf("Could not gracefully shutdown the server: %v", err)
	}
	close(done)
}

func (a *Adapter) newRouter(hook *dockerhub.Webhook) *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		payload, err := hook.Parse(r, dockerhub.BuildEvent)

		if err != nil {
			if err == dockerhub.ErrInvalidHTTPMethod {
				a.logger.Error("event not send to sink as invalid http method")
			} else if err == dockerhub.ErrParsingPayload {
				a.logger.Error("event not send to sink as parsing buildPayload err")
			}
			a.logger.Errorf("Error processing request: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		bp, _ := payload.(dockerhub.BuildPayload)

		go a.processPayload(bp)

		a.logger.Infof("event accepted: %v", bp)
		w.WriteHeader(http.StatusOK)
	})
	return router
}

func (a *Adapter) processPayload(payload dockerhub.BuildPayload) {

	a.logger.Info("processing event ...")

	err := a.sendEventToSink(payload)
	if err != nil {
		a.logger.Errorf("failed to send event to sink: %v", err)
	}

	if a.autoValidation {
		message := "Event has been sent successfully."
		if err != nil {
			message = fmt.Sprintf("failed to send event to sink: %v", err)
		}
		callbackData := &resources.CallbackPayload{
			// always StatusSuccess to continue to receive webhook
			State:       resources.StatusSuccess,
			Description: message,
			Context:     a.eventSource,
			TargetURL:   "",
		}

		err := callbackData.EmitValidationCallback(payload.CallbackURL)
		if err != nil {
			a.logger.Errorf("failed to send validation callback: %v", err)
		} else {
			a.logger.Info("http request for callbackUrl done.")
		}
	} else {
		a.logger.Info("Auto Callback does not work as specified in the spec.")
	}
}

// sendEventToSink transforms buildPayload to CloudEvent, then try to send to sink.
func (a *Adapter) sendEventToSink(payload dockerhub.BuildPayload) error {
	cloudEventType := v1alpha1.DockerHubCloudEventsEventType(resources.DockerHubEventType)
	cloudEventSource := v1alpha1.DockerHubEventSource(payload.Repository.RepoName)
	cloudEventTime, err := getTime(payload.PushData.PushedAt)
	if err != nil {
		a.logger.Warnf("failed to parse pushedAt field. Use time.Now(): %v", err)
	}
	uid, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	event := cloudevents.NewEvent()
	event.SetID(uid.String())
	event.SetType(cloudEventType)
	event.SetSource(cloudEventSource)
	event.SetTime(cloudEventTime)
	event.SetSubject(payload.PushData.Pusher)
	event.SetExtension("tag", payload.PushData.Tag)
	err = event.SetData(cloudevents.ApplicationJSON, payload)
	if err != nil {
		return fmt.Errorf("failed to marshal buildPayload :%v", err)
	}

	a.logger.Infof("Sending event: %v", event)

	result := a.client.Send(context.Background(), event)
	if !cloudevents.IsACK(result) {
		return fmt.Errorf("send() could not get ACK: %v", result)
	}
	return nil
}

// This is precise, but definitely there should be a better way :(
// If the operation failed, use time.Now()
func getTime(pushedAt float32) (time.Time, error) {
	pt := pushedAt
	if pt < 0 {
		return time.Now(), fmt.Errorf("pushedAt should not be negative: %f", pt)
	}
	ft := strconv.FormatFloat(float64(pt), 'f', -1, 64)
	u, err := strconv.ParseInt(ft, 10, 64)
	if err != nil {
		return time.Now(), fmt.Errorf("failed to process pushedAt field: %v", err)
	}
	ans := time.Unix(u, 0)
	return ans, nil
}
