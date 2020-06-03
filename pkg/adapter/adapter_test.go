package adapter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/go-cmp/cmp"
	"go.uber.org/zap"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"knative.dev/eventing/pkg/adapter/v2"
	adaptertest "knative.dev/eventing/pkg/adapter/v2/test"
	"knative.dev/pkg/logging"
	pkgtesting "knative.dev/pkg/reconciler/testing"

	dh "gopkg.in/go-playground/webhooks.v5/docker"

	cloudevents "github.com/cloudevents/sdk-go/v2"

	"github.com/tom24d/eventing-dockerhub/pkg/adapter/resources"
)

const (
	testSubject   = "1234"
	testOwnerRepo = "test-user/test-repo"
	testCallbackPort = "4320"
	testAdapterPort = "8765"
	callbackServerWaitThreshold = 4
)

type testCase struct {
	// name is a descriptive name for this test suitable as a first argument to t.Run()
	name string

	// buildPayload contains the DockerHub event buildPayload
	buildPayload interface{}

	// httpMethod is the method to emit http request
	httpMethod string

	// eventType is the DockerHub event type
	eventType string

	// cloudEventSendExpected is whether event is transferred
	cloudEventSendExpected bool

	// wantEventType is the expected CloudEvent EventType
	wantCloudEventType string

	// wantCloudEventSubject is the expected CloudEvent subject
	//wantCloudEventSubject string

	//wantCallbackExpected is whether callback is expected
	wantCallbackExpected bool

	//wantCallbackStatus is the expected resources.Status
	wantCallbackStatus resources.Status
}

var testCases = []testCase{
	{
		name: "valid build buildPayload",
		buildPayload: func() interface{} {
			bp := &dh.BuildPayload{}
			bp.CallbackURL = fmt.Sprintf("http://127.0.0.1:%s/", testCallbackPort)
			return bp
		}(),
		httpMethod: http.MethodPost,
		eventType:             DockerHubEventType,
		cloudEventSendExpected: true,
		wantCloudEventType: "dev.knative.source.dockerhub.push",
		wantCallbackExpected: true,
		wantCallbackStatus:   resources.StatusSuccess,
	},
	{
		name: "invalid callback url",
		buildPayload: func() interface{} {
			bp := &dh.BuildPayload{}
			bp.CallbackURL = fmt.Sprintf("http://127.0.0.1:%s/", "10000")
			return bp
		}(),
		httpMethod: http.MethodPost,
		eventType:            DockerHubEventType,
		cloudEventSendExpected: false,
		wantCloudEventType: "dev.knative.source.dockerhub.push",//TODO
		wantCallbackExpected: false,
	},
	{
		name: "nil buildPayload",
		buildPayload: func() interface{} {
			bp := ""
			return bp
		}(),
		httpMethod: http.MethodPost,
		eventType:            DockerHubEventType,
		cloudEventSendExpected: false,
		wantCloudEventType: "dev.knative.source.dockerhub.push",
		wantCallbackExpected: false,
	},
	{
		name: "funny payload",
		buildPayload: func() interface{} {
			bp := "bazinga"
			return bp
		}(),
		httpMethod: http.MethodPost,
		eventType:            DockerHubEventType,
		cloudEventSendExpected: false,
		wantCloudEventType: "dev.knative.source.dockerhub.push",
		wantCallbackExpected: false,
	},
	{
		name: "not buildPayload",
		buildPayload: func() interface{} {
			bp := &resources.CallbackPayload{
				State:       resources.StatusError,
				Description: "This is attack webhook.",
				Context:     "",
				TargetURL:   "",
			}
			return bp
		}(),
		httpMethod: http.MethodPost,
		eventType:            DockerHubEventType,
		cloudEventSendExpected: false,
		wantCloudEventType: "dev.knative.source.dockerhub.push",
		wantCallbackExpected: false,
	},
	{
		name: "httpPatch",
		buildPayload: func() interface{} {
			bp := &dh.BuildPayload{}
			bp.CallbackURL = fmt.Sprintf("http://127.0.0.1:%s/", testCallbackPort)
			return bp
		}(),
		httpMethod: http.MethodPatch,
		eventType:            DockerHubEventType,
		cloudEventSendExpected: false,
		wantCloudEventType: "dev.knative.source.dockerhub.push",
		wantCallbackExpected: false,
	},
}

func TestNewEnv(t *testing.T) {
	want := &envConfig{}

	got := NewEnv()
	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("unexpected event data (-want, +got) = %v", diff)
	}
}

func TestServer(t *testing.T) {
	ce := adaptertest.NewTestClient()
	testAdapter := newTestAdapter(t, ce)
	hook, _ := dh.New()
	router := testAdapter.newRouter(hook)
	server := httptest.NewServer(router)

	notify := make(chan resources.Status, 1)

	callbackServer := newCallbackServer(t, testCallbackPort, notify)
	defer server.Close()
	defer callbackServer.Close()

	for _, tc := range testCases {
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

func newCallbackServer(t *testing.T, port string, notify chan resources.Status) *httptest.Server {
	h := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("accepted"))
		t.Log("callback posted")
		cp, err := resources.Parse(r)
		if err != nil {
			t.Fatalf("failed to parse callback buildPayload: %v", err)
		}
		notify <- cp.State
	}
	r := http.NewServeMux()
	r.HandleFunc("/", h)
	server := httptest.NewUnstartedServer(r)

	addr, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%s", testCallbackPort))
	if err != nil {
		t.Fatalf("cannot listen on custom address: %v", err)
	}
	server.Listener.Close()
	server.Listener = addr

	server.Start()

	return server
}

// runner returns a testing func that can be passed to t.Run.
func (tc *testCase) runner(t *testing.T, url string, ceClient *adaptertest.TestCloudEventsClient, nc chan resources.Status) func(t *testing.T) {
	return func(t *testing.T) {
		if tc.eventType == "" {
			t.Fatal("eventType is required for table tests")
		}
		body, _ := json.Marshal(tc.buildPayload)
		req, err := http.NewRequest(tc.httpMethod, url, bytes.NewReader(body))
		if err != nil {
			t.Fatal(err)
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Error(err)
		}
		defer resp.Body.Close()

		if tc.cloudEventSendExpected && tc.wantCallbackExpected {
			tc.waitCallbackReport(t, nc)
			tc.validateCESentPayload(t, ceClient)
		} else if tc.cloudEventSendExpected {
			tc.validateCESentPayload(t, ceClient)
		}
	}
}

func (tc *testCase)waitCallbackReport(t *testing.T, notify chan resources.Status) {
	ticker := time.NewTicker(time.Second)
	th := 0
	for {
		select {
		case gotStatus := <-notify:
			if diff := cmp.Diff(tc.wantCallbackStatus, gotStatus); diff != "" {
				t.Fatalf("unexpected event data (-want, +got) = %v", diff)
			}
			return
		case <-ticker.C:
			th += 1
			if th > callbackServerWaitThreshold {
				t.Fatalf("could not receive validation callback in %d seconds.", callbackServerWaitThreshold)
			}
		}
	}
}

func (tc *testCase) validateCESentPayload(t *testing.T, ce *adaptertest.TestCloudEventsClient) {
	t.Helper()
	if len(ce.Sent()) != 1 {
		return
	}
	//TODO add subject test if needed
	//eventSubject := ce.Sent()[0].Subject()
	//if eventSubject != tc.wantCloudEventSubject {
	//	t.Fatalf("Expected %q event subject to be sent, got %q", tc.wantCloudEventSubject, eventSubject)
	//}

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
	payload, err := json.Marshal(tc.buildPayload)
	if err != nil {
		t.Fatalf("Could not marshal sent buildPayload: %v", err)
	}
	err = json.Unmarshal(payload, &want)
	if err != nil {
		t.Fatalf("Could not unmarshal sent buildPayload: %v", err)
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("unexpected event data (-want, +got) = %v", diff)
	}
}
