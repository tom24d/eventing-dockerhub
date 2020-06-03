package resources

import (
	"bytes"
	"encoding/json"
	"github.com/google/go-cmp/cmp"
	"net/http"
	"net/http/httptest"
	"testing"
)

type testCase struct {
	// name is a descriptive name for this test suitable as a first argument to t.Run()
	name string

	// payload contains the CallbackPayload event payload
	payload interface{}
}

var testCases = []testCase{
	{
		name: "valid payload",
		payload: func() interface{} {
			cb := &CallbackPayload{
				State:       StatusSuccess,
				Description: "This is description field.",
				Context:     "This is context field.",
				TargetURL:   "http://example.com",
			}
			return cb
		}(),
	},
}


func TestEmitValidationCallback(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			notify := make(chan *CallbackPayload, 1)
			server := newServer(func(writer http.ResponseWriter, r *http.Request) {
				payload, err := Parse(r)
				if err != nil {
					t.Fatal(err)
				}
				notify <- payload
			})
			defer server.Close()

			want, ok := tc.payload.(*CallbackPayload)
			if !ok {
				t.Fatal("type assertion failed")
			}

			err := want.EmitValidationCallback(server.URL)
			if err != nil {
				t.Fatal(err)
			}

			got := <- notify
			if diff := cmp.Diff(want, got); diff != "" {
				t.Fatalf("unexpected event data (-want, +got) = %v", diff)
			}
		})
	}
}

func TestParse(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var parseError error
			var got interface{}
			server := newServer(func(writer http.ResponseWriter, r *http.Request) {
				got, parseError = Parse(r)
				if parseError != nil {
					t.Fatal(parseError)
				}
			})
			defer server.Close()

			body, err := json.Marshal(tc.payload)
			if err != nil {
				t.Fatal(err)
			}
			req, err := http.NewRequest(http.MethodPost, server.URL, bytes.NewReader(body))
			if err != nil {
				t.Fatal(err)
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

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
