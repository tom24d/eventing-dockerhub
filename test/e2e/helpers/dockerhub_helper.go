package helpers

import (
	"fmt"
	"github.com/google/uuid"
	"strconv"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	eventingtestlib "knative.dev/eventing/test/lib"
	pkgTest "knative.dev/pkg/test"

	sourcesv1alpha1 "github.com/tom24d/eventing-dockerhub/pkg/apis/sources/v1alpha1"
	dhtestresources "github.com/tom24d/eventing-dockerhub/test/resources"

	dockerhub "gopkg.in/go-playground/webhooks.v5/docker"
)

const (
	SenderImageName = "webhook-sender"
)

func MustSendWebhook(client *eventingtestlib.Client, targetURL string, data *dockerhub.BuildPayload) {
	args := []string{fmt.Sprintf("--%s=%s", dhtestresources.ArgSink, targetURL),
		fmt.Sprintf("--%s=%s", dhtestresources.ArgPayload, dhtestresources.MarshalPayload(data))}

	// create webhook sender
	eventSender := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: client.Namespace,
			Name:      SenderImageName,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{{
				Name:  SenderImageName,
				Image: pkgTest.ImagePath(SenderImageName),
				ImagePullPolicy: corev1.PullAlways,
				Args:  args,
			}},
			RestartPolicy: corev1.RestartPolicyNever,
		},
	}
	client.CreatePodOrFail(eventSender)

	err := pkgTest.WaitForPodState(client.Kube, func(pod *corev1.Pod) (bool, error) {
		if pod.Status.Phase == corev1.PodFailed {
			return true, fmt.Errorf("event sender pod failed with message %s", pod.Status.Message)
		} else if pod.Status.Phase != corev1.PodSucceeded {
			return false, nil
		}
		return true, nil
	}, eventSender.Name, eventSender.Namespace)

	if err != nil {
		client.T.Fatalf("Failed waiting for pod for completeness %q: %v", eventSender.Name, eventSender)
	}
}

func CreateValidationReceiverOrFail(client *eventingtestlib.Client) *corev1.Pod {
	const receiverImageName = "validation-receiver"
	args := []string{"--patient", strconv.Itoa(60)}

	receiverPod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: client.Namespace,
			Name:      receiverImageName,
			Labels: map[string]string{"e2etest": uuid.New().String()},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{{
				Name:  receiverImageName,
				Image: pkgTest.ImagePath(receiverImageName),
				ImagePullPolicy: corev1.PullAlways,
				Args:  args,
				Ports: []corev1.ContainerPort{
					{
						ContainerPort: 8080,
					},
				},
			}},
			RestartPolicy: corev1.RestartPolicyNever,
		},
	}

	client.CreatePodOrFail(receiverPod, eventingtestlib.WithService(receiverPod.GetName()))

	err := pkgTest.WaitForPodState(client.Kube, func(pod *corev1.Pod) (bool, error) {
		if pod.Status.Phase == corev1.PodFailed {
			return true, fmt.Errorf("validation receiver pod failed to get up with message %s", pod.Status.Message)
		} else if pod.Status.Phase != corev1.PodRunning {
			return false, nil
		}
		return true, nil
	}, receiverPod.Name, receiverPod.Namespace)

	if err != nil {
		client.T.Fatalf("Failed waiting for pod running %q: %v", receiverPod.Name, receiverPod)
	}
	return receiverPod
}

func WaitForValidationReceiverPodSuccessOrFail(client *eventingtestlib.Client, receiverPod *corev1.Pod, notify chan bool) {
	err := pkgTest.WaitForPodState(client.Kube, func(pod *corev1.Pod) (bool, error) {
		if pod.Status.Phase == corev1.PodFailed {
			return true, fmt.Errorf("validation receiver pod failed with message %s", pod.Status.Message)
		} else if pod.Status.Phase != corev1.PodSucceeded {
			return false, nil
		}
		return true, nil
	}, receiverPod.Name, receiverPod.Namespace)

	if err != nil {
		client.T.Fatalf("Failed waiting for pod for completeness %q: %v", receiverPod.Name, receiverPod)
	}
	notify <- true
}

func GetURLOrFail(client *eventingtestlib.Client, source *sourcesv1alpha1.DockerHubSource) string {
	dhs, err := GetSourceClient(client).SourcesV1alpha1().
		DockerHubSources(client.Namespace).Get(source.Name, metav1.GetOptions{})
	if err != nil {
		client.T.Fatalf("failed to get DockerHubSource: %v", source.Name)
	}

	allocatedURL := dhs.Status.URL.String()
	if allocatedURL == "" {
		client.T.Fatalf("DockerHubSource URL is nil: %v", source.GetName())
	}
	return allocatedURL
}

// TODO Get and return URL might be better
func SetCallbackURLOrFail(c *eventingtestlib.Client, data *dockerhub.BuildPayload, svcName string) {
	// TODO use lib if exists
	url := fmt.Sprintf("http://%s.%s.svc.cluster.local", svcName, c.Namespace)
	data.CallbackURL = url
}
