package resources

import (
	"encoding/json"

	dockerhub "gopkg.in/go-playground/webhooks.v5/docker"
)

func MarshalPayload(payload interface{}) string {
	byteData, _ := json.Marshal(payload)
	return string(byteData)
}

func UnmarshalPayload(str string) dockerhub.BuildPayload {
	var p dockerhub.BuildPayload
	_ = json.Unmarshal([]byte(str), &p)
	return p
}
