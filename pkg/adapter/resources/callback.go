package resources

import (
	"bytes"
	"fmt"
	"net/http"
	"encoding/json"
)

type CallbackPayload struct {
	State Status `json:"state"`
	Description string `json:"description"`
	Context string `json:"context"`
	TargetURL string `json:"target_url"`
}

type Status string

const (
	StatusSuccess Status = "success"
	StatusFailure Status = "failure"
	StatusError Status = "error"
)

// TODO vocabulary check
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
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("sending callback failed")
	}
	return nil
}