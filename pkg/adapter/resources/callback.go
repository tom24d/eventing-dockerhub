package resources

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

// TODO make PR to gopkg.in webhook
