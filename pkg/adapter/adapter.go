package adapter

import (
	"context"
	"fmt"
	"net/http"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"go.uber.org/zap"
	dockerhub "gopkg.in/go-playground/webhooks.v5/docker"

	//knative.dev imports
	"knative.dev/eventing/pkg/adapter/v2"
	"knative.dev/pkg/logging"

	"github.com/google/uuid"

	"github.com/tom24d/eventing-dockerhub/pkg/adapter/resources"
	"github.com/tom24d/eventing-dockerhub/pkg/apis/sources/v1alpha1"
)

const (
	DockerHubEventType = "push"
)

type envConfig struct {
	// Include the standard adapter.EnvConfig used by all adapters.
	adapter.EnvConfig

	// Port to listen incoming connections
	Port string `envconfig:"PORT" default:"8080"`

}

func NewEnv() adapter.EnvConfigAccessor { return &envConfig{} }

// Adapter converts incoming GitHub webhook events to CloudEvents
type Adapter struct {
	client         cloudevents.Client
	logger         *zap.SugaredLogger
	port           string
	autoValidation bool
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
		autoValidation: true, //TODO
	}
}

// Start runs the adapter.
// Returns if stopCh is closed or Send() returns an error.
func (a *Adapter) Start(ctx context.Context) error {
	return a.start(ctx.Done())
}

func (a *Adapter) start(stopCh <-chan struct{}) error {
	done := make(chan bool, 1)
	hook, err := dockerhub.New()
	if err != nil {
		return fmt.Errorf("cannot create gitlab hook: %v", err)
	}

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
				w.Write([]byte("event not send to sink as invalid http method"))
				return
			} else if err == dockerhub.ErrParsingPayload {
				w.Write([]byte("event not send to sink as parsing payload err"))
				return
			}
			a.logger.Errorf("Error processing request: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		bp, ok := payload.(dockerhub.BuildPayload)
		if !ok {
			a.logger.Error("type assertion failed for payload")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// TODO think what is "event processed"?
		go a.processPayload(bp)

		a.logger.Infof("event accepted")
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("accepted"))
	})
	return router
}

func (a *Adapter)processPayload(payload dockerhub.BuildPayload) {

	a.logger.Info("processing event ...")

	err := a.sendEventToSink(payload)
	if err != nil {
		a.logger.Errorf("failed to send event to sink: %v", err)
	}

	if a.autoValidation {
		if err != nil {
			a.logger.Info("going to report that sending sink has failed")
			callbackData := &resources.CallbackPayload{
				State:       resources.StatusSuccess, // always StatusSuccess to continue receiving webhook.
				Description: fmt.Sprintf("failed to send event to sink: %v", err),
				Context:     "",// TODO adapter resource name
				TargetURL:   "",
			}
			err := callbackData.EmitValidationCallback(payload.CallbackURL)
			if err != nil {
				a.logger.Errorf("failed to send validation callback: %v", err)
				return
			}
			return
		}
		a.logger.Info("going to report that sending sink has completed successfully")
		callbackData := &resources.CallbackPayload{
			State:       resources.StatusSuccess,
			Description: "Event has been sent successfully.",
			Context:     "",// TODO adapter resource name
			TargetURL:   "",
		}

		err := callbackData.EmitValidationCallback(payload.CallbackURL)
		if err != nil {
			a.logger.Errorf("failed to send validation callback: %v", err)
		}
	}
}

// sendEventToSink transforms payload to CloudEvent, then try to send to sink.
func (a *Adapter) sendEventToSink(payload dockerhub.BuildPayload) error {
	cloudEventType := v1alpha1.DockerHubCloudEventsEventType(DockerHubEventType)
	cloudEventSource := v1alpha1.DockerHubEventSource(payload.Repository.RepoName)
	uid, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	event := cloudevents.NewEvent()
	event.SetID(uid.String())
	event.SetType(cloudEventType)
	event.SetSource(cloudEventSource)
	err = event.SetData(cloudevents.ApplicationJSON, payload)
	if err != nil {
		return err
	}

	result := a.client.Send(context.Background(), event)
	if !cloudevents.IsACK(result) {
		return result
	}
	return nil
}
