package adapter

//import (
//	"knative.dev/eventing/pkg/adapter/v2"
//)
//
//type envConfig struct {
//	// Include the standard adapter.EnvConfig used by all adapters.
//	adapter.EnvConfig
//}
//
//func NewEnv() adapter.EnvConfigAccessor { return &envConfig{} }

// // Adapter generates events at a regular interval.
// type Adapter struct {
// 	client   cloudevents.Client
// 	interval time.Duration
// 	logger   *zap.SugaredLogger

// 	nextID int
// }

// type dataExample struct {
// 	Sequence  int    `json:"sequence"`
// 	Heartbeat string `json:"heartbeat"`
// }

// func (a *Adapter) newEvent() cloudevents.Event {
// 	event := cloudevents.NewEvent()
// 	event.SetType("dev.knative.sample")
// 	event.SetSource("sample.knative.dev/heartbeat-source")

// 	if err := event.SetData(cloudevents.ApplicationJSON, &dataExample{
// 		Sequence:  a.nextID,
// 		Heartbeat: a.interval.String(),
// 	}); err != nil {
// 		a.logger.Errorw("failed to set data")
// 	}
// 	a.nextID++
// 	return event
// }

// // Start runs the adapter.
// // Returns if stopCh is closed or Send() returns an error.
// func (a *Adapter) Start(stopCh <-chan struct{}) error {
// 	a.logger.Infow("Starting heartbeat", zap.String("interval", a.interval.String()))
// 	for {
// 		select {
// 		case <-time.After(a.interval):
// 			event := a.newEvent()
// 			a.logger.Infow("Sending new event", zap.String("event", event.String()))
// 			if result := a.client.Send(context.Background(), event); !cloudevents.IsACK(result) {
// 				a.logger.Infow("failed to send event", zap.String("event", event.String()), zap.Error(result))
// 				// We got an error but it could be transient, try again next interval.
// 				continue
// 			}
// 		case <-stopCh:
// 			a.logger.Info("Shutting down...")
// 			return nil
// 		}
// 	}
// }

// func NewAdapter(ctx context.Context, aEnv adapter.EnvConfigAccessor, ceClient cloudevents.Client) adapter.Adapter {
// 	env := aEnv.(*envConfig) // Will always be our own envConfig type
// 	logger := logging.FromContext(ctx)
// 	logger.Infow("Heartbeat example", zap.Duration("interval", env.Interval))
// 	return &Adapter{
// 		interval: env.Interval,
// 		client:   ceClient,
// 		logger:   logger,
// 	}
// }
