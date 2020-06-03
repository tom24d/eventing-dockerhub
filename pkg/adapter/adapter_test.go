package adapter

import (
	"bytes"
	"context"
	"fmt"
	"github.com/google/go-cmp/cmp"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"testing"
	"encoding/json"
	"time"

	"knative.dev/eventing/pkg/adapter/v2"
	pkgtesting "knative.dev/pkg/reconciler/testing"
	"knative.dev/pkg/logging"
	adaptertest "knative.dev/eventing/pkg/adapter/v2/test"


	dh "gopkg.in/go-playground/webhooks.v5/docker"

	cloudevents "github.com/cloudevents/sdk-go/v2"

	"github.com/tom24d/eventing-dockerhub/pkg/adapter/resources"
)

const (
	testSubject   = "1234"
	testOwnerRepo = "test-user/test-repo"
	testCallbackPort = "4320"
	testAdapterPort = "8765"
)

type testCase struct {
	// name is a descriptive name for this test suitable as a first argument to t.Run()
	name string

	// payload contains the DockerHub event payload
	payload interface{}

	// eventType is the DockerHub event type
	eventType string

	// wantEventType is the expected CloudEvent EventType
	wantCloudEventType string

	// wantCloudEventSubject is the expected CloudEvent subject
	wantCloudEventSubject string

	//wantCallbackStatus is the expected resources.Status
	wantCallbackStatus resources.Status
}

var testCases = []testCase{
	{
		name: "valid build payload",
		payload: func() interface{} {
			bp := &dh.BuildPayload{}
			// TODO populate callback server url
			bp.CallbackURL = fmt.Sprintf("http://:%s/", testCallbackPort)
			return bp
		}(),
		eventType:             "push",
		//wantCloudEventSubject: testSubject,
		wantCallbackStatus: resources.StatusSuccess,
	},
}


func TestServer(t *testing.T) {
	for _, tc := range testCases {
		ce := adaptertest.NewTestClient()
		testAdapter := newTestAdapter(t, ce)
		hook, err := dh.New()
		if err != nil {
			t.Fatal(err)
		}
		router := testAdapter.newRouter(hook)
		server := httptest.NewServer(router)

		notify := make(chan string, 1)

		callbackServer := newCallbackServer(t, testCallbackPort, notify)
		defer server.Close()
		defer callbackServer.Close()

		t.Run(tc.name, tc.runner(t, server.URL, ce, notify))
	}
}

func TestGracefulShutdown(t *testing.T) {
	ce := adaptertest.NewTestClient()
	ra := newTestAdapter(t, ce)
	ctx, cancel := context.WithCancel(context.Background())

	go func(cancel context.CancelFunc) {
		defer cancel()
		time.Sleep(time.Second)

	}(cancel)

	t.Logf("starting webhook server")
	err := ra.Start(ctx)
	if err != nil {
		t.Error(err)
	}
}

func newTestAdapter(t *testing.T, ce cloudevents.Client) *Adapter {
	env := envConfig{
		EnvConfig: adapter.EnvConfig{
			Namespace: "default",
		},
		Port: testAdapterPort,
	}
	ctx, _ := pkgtesting.SetupFakeContext(t)
	logger := zap.NewExample().Sugar()
	ctx = logging.WithLogger(ctx, logger)

	return NewAdapter(ctx, &env, ce).(*Adapter)
}

func newCallbackServer(t *testing.T, port string, notify chan string) *httptest.Server {
	h := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Callback POSTed.")
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("accepted"))
		notify <- "posted"
	}
	r := http.NewServeMux()
	r.HandleFunc("/", h)
	server := httptest.NewServer(r)
	return server
}

// runner returns a testing func that can be passed to t.Run.
func (tc *testCase) runner(t *testing.T, url string, ceClient *adaptertest.TestCloudEventsClient, nc chan string) func(t *testing.T) {
	return func(t *testing.T) {
		if tc.eventType == "" {
			t.Fatal("eventType is required for table tests")
		}
		body, _ := json.Marshal(tc.payload)
		req, err := http.NewRequest("POST", url, bytes.NewReader(body))
		if err != nil {
			t.Fatal(err)
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Error(err)
		}
		defer resp.Body.Close()

		waitCallbackReport(t, nc)

		tc.validateAcceptedPayload(t, ceClient)
	}
}

func waitCallbackReport(t *testing.T, notify chan string) {
	ticker := time.NewTicker(time.Second)
	th := 0
	for {
		select {
		case <-notify:
			return
		case <-ticker.C:
			th += 1
			if th > 4 {
				t.Fatal("could not receive validation callback")
			}
		}
	}
}

func (tc *testCase) validateAcceptedPayload(t *testing.T, ce *adaptertest.TestCloudEventsClient) {
	t.Helper()
	if len(ce.Sent()) != 1 {
		return
	}
	eventSubject := ce.Sent()[0].Subject()
	if eventSubject != tc.wantCloudEventSubject {
		t.Fatalf("Expected %q event subject to be sent, got %q", tc.wantCloudEventSubject, eventSubject)
	}

	if tc.wantCloudEventType != "" {
		eventType := ce.Sent()[0].Type()
		if eventType != tc.wantCloudEventType {
			t.Fatalf("Expected %q event type to be sent, got %q", tc.wantCloudEventType, eventType)
		}
	}

	data := ce.Sent()[0].Data()

	var got interface{}
	var want interface{}

	err := json.Unmarshal(data, &got)
	if err != nil {
		t.Fatalf("Could not unmarshal sent data: %v", err)
	}
	payload, err := json.Marshal(tc.payload)
	if err != nil {
		t.Fatalf("Could not marshal sent payload: %v", err)
	}
	err = json.Unmarshal(payload, &want)
	if err != nil {
		t.Fatalf("Could not unmarshal sent payload: %v", err)
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("unexpected event data (-want, +got) = %v", diff)
	}
}
