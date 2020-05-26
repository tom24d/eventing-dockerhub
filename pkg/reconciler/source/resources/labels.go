package resources

const (
	// controllerAgentName is the string used by this controller to identify
	// itself when creating events.
	controllerAgentName = "dockerhubsource-controller"
)

// Labels returns map which holds the label of "knative-eventing-source"
// and "knative-eventing-source-name".
func Labels(name string) map[string]string {
	return map[string]string{
		"knative-eventing-source":      controllerAgentName,
		"knative-eventing-source-name": name,
	}
}
