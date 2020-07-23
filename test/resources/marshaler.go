package resources

import (
	"encoding/json"
)

func MarshalPayload(payload interface{}) string {
	byteData, _ := json.Marshal(payload)
	return string(byteData)
}
