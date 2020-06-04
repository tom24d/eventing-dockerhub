package resources

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
	dh "gopkg.in/go-playground/webhooks.v5/docker"
)

type testCallbackCase struct {
	// name is a descriptive name for this test suitable as a first argument to t.Run()
	name string

	// payload contains the CallbackPayload event payload
	payload interface{}

	// expectedSuccess represents any process is expected to be "success"
	expectedSuccess bool

	// wantError represents first expected error
	wantError interface{}
}

type testParseCase struct {
	// payload contains the CallbackPayload event payload
	payload interface{}

	// expectedSuccess represents any process is expected to be "success"
	expectedSuccess bool

	// wantError represents first expected error
	wantError interface{}

	// name is a descriptive name for this test suitable as a first argument to t.Run()
	name string

	// httpMethod is used method
	httpMethod string
}

var testCallbackCases = []testCallbackCase{
	{
		name: "valid payload",
		payload: func() interface{} {
			cb := CallbackPayload{
				State:       StatusSuccess,
				Description: "This is description field.",
				Context:     "This is context field.",
				TargetURL:   "http://example.com",
			}
			return cb
		}(),
		expectedSuccess: true,
		wantError:       nil,
	},
	{
		name: "nil payload",
		payload: func() interface{} {
			cb := CallbackPayload{}
			return cb
		}(),
		expectedSuccess: false,
		wantError: func() interface{} {
			return errors.New("error parsing payload")
		},
	},
}

var testParseCases = []testParseCase{
	{
		name:       "valid case",
		httpMethod: http.MethodPost,
		payload: func() interface{} {
			cb := CallbackPayload{
				State:       StatusSuccess,
				Description: "This is description field.",
				Context:     "This is context field.",
				TargetURL:   "http://example.com",
			}
			return cb
		}(),
		expectedSuccess: true,
		wantError:       nil,
	},
	{
		name:       "invalid httpMethod",
		httpMethod: http.MethodPatch,
		payload: func() interface{} {
			cb := CallbackPayload{
				State:       StatusSuccess,
				Description: "This is description field.",
				Context:     "This is context field.",
				TargetURL:   "http://example.com",
			}
			return cb
		}(),
		expectedSuccess: false,
		wantError:       dh.ErrInvalidHTTPMethod.Error(),
	},
	{
		name:            "nil payload",
		httpMethod:      http.MethodPost,
		payload:         "",
		expectedSuccess: false,
		wantError:       dh.ErrParsingPayload.Error(),
	},
}

func TestEmitValidationCallback(t *testing.T) {
	for _, tc := range testCallbackCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			notify := make(chan interface{}, 1)
			server := newServer(func(writer http.ResponseWriter, r *http.Request) {
				payload, err := Parse(r)
				if err != nil {
					t.Fatal(err)
				}
				notify <- payload
			})
			defer server.Close()

			want, ok := tc.payload.(CallbackPayload)
			if !ok {
				t.Fatal("type assertion failed")
			}

			err := want.EmitValidationCallback(server.URL)
			if err != nil {
				t.Fatal(err)
			}

			got := <-notify
			if diff := cmp.Diff(want, got); diff != "" {
				t.Fatalf("unexpected event data (-want, +got) = %v", diff)
			}

			// nil callbackURL
			err = want.EmitValidationCallback("")
			t.Log(err)
			if err == nil {
				t.Fatal("no expected error detected")
			}
		})
	}
}

func TestParse(t *testing.T) {
	for _, tc := range testParseCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			parsed := make(chan bool, 1)

			var gotError error
			var got interface{}

			server := newServer(func(writer http.ResponseWriter, r *http.Request) {
				got, gotError = Parse(r)
				if gotError != nil {
					parsed <- false
					return
				}
				parsed <- true
			})
			defer server.Close()

			body, err := json.Marshal(tc.payload)
			if err != nil {
				t.Fatal(err)
			}
			req, err := http.NewRequest(tc.httpMethod, server.URL, bytes.NewReader(body))
			if err != nil {
				t.Fatal(err)
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			t.Log("waiting for incoming payload ...")
			p := <-parsed

			if !tc.expectedSuccess && p {
				t.Fatalf("expected error, but no error detected (want) = %v", tc.wantError)
			}

			if !p {
				if diff := cmp.Diff(tc.wantError, gotError.Error()); diff != "" {
					t.Fatalf("unexpected error data (-want, +got) = %v", diff)
				}
				return
			}

			if diff := cmp.Diff(tc.payload, got); diff != "" {
				t.Fatalf("unexpected event data (-want, +got) = %v", diff)
			}
		})
	}
}

func newServer(handler http.HandlerFunc) *httptest.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)
	return httptest.NewServer(mux)
}
