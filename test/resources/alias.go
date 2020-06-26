package resources

import (
	eventingtestlib "knative.dev/eventing/test/lib"
)

func WaitForAllTestResourcesReadyOrFail(c *eventingtestlib.Client) {
	c.WaitForAllTestResourcesReadyOrFail()
}
