package resources

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

// parse errors
var (
	ErrInvalidHTTPMethod = errors.New("invalid HTTP Method")
	ErrParsingPayload    = errors.New("error parsing payload")
)

type CallbackPayload struct {
	State       Status `json:"state"`
	Description string `json:"description"`
	Context     string `json:"context"`
	TargetURL   string `json:"target_url"`
}

type Status string

const (
	StatusSuccess Status = "success"
	StatusFailure Status = "failure"
	StatusError   Status = "error"
)

// EmitValidationCallback send callback to given callbackURL
func (callback *CallbackPayload) EmitValidationCallback(callbackURL string) error {
	if callbackURL == "" {
		return fmt.Errorf("callbackURL is not set")
	}

	payload, err := json.Marshal(callback)
	if err != nil {
		return err
	}

	resp, err := http.Post(callbackURL, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= 300 {
		return fmt.Errorf("sending callback failed")
	}
	return nil
}

// Parse verifies and parses the events specified and returns the payload object or an error
func Parse(r *http.Request) (interface{}, error) {
	defer func() {
		_, _ = io.Copy(ioutil.Discard, r.Body)
		_ = r.Body.Close()
	}()

	if r.Method != http.MethodPost {
		return nil, ErrInvalidHTTPMethod
	}

	payload, err := ioutil.ReadAll(r.Body)
	if err != nil || len(payload) == 0 {
		return nil, ErrParsingPayload
	}

	var pl CallbackPayload
	err = json.Unmarshal([]byte(payload), &pl)
	if err != nil {
		return nil, ErrParsingPayload
	}
	return pl, err
}
