package helpers

import (
	"fmt"
	"strconv"

	"github.com/google/uuid"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"knative.dev/eventing/test/lib"
	"knative.dev/pkg/test"
)

func CreateValidationReceiverOrFail(client *lib.Client) *v1.Pod {
	const receiverImageName = "validation-receiver"
	args := []string{"--patient", strconv.Itoa(60)}

	receiverPod := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: client.Namespace,
			Name:      receiverImageName,
			Labels:    map[string]string{"e2etest": uuid.New().String()},
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{{
				Name:            receiverImageName,
				Image:           test.ImagePath(receiverImageName),
				ImagePullPolicy: v1.PullAlways,
				Args:            args,
				Ports: []v1.ContainerPort{
					{
						ContainerPort: 8080,
					},
				},
			}},
			RestartPolicy: v1.RestartPolicyNever,
		},
	}

	client.CreatePodOrFail(receiverPod, lib.WithService(receiverPod.GetName()))

	err := test.WaitForPodState(client.Kube, func(pod *v1.Pod) (bool, error) {
		if pod.Status.Phase == v1.PodFailed {
			return true, fmt.Errorf("validation receiver pod failed to get up with message %s", pod.Status.Message)
		} else if pod.Status.Phase != v1.PodRunning {
			return false, nil
		}
		return true, nil
	}, receiverPod.Name, receiverPod.Namespace)

	if err != nil {
		client.T.Fatalf("Failed waiting for pod running %q: %v", receiverPod.Name, receiverPod)
	}
	return receiverPod
}

func WaitForValidationReceiverPodSuccessOrFail(client *lib.Client, receiverPod *v1.Pod, notify chan bool) {
	err := test.WaitForPodState(client.Kube, func(pod *v1.Pod) (bool, error) {
		if pod.Status.Phase == v1.PodFailed {
			return true, fmt.Errorf("validation receiver pod failed: %v", pod)
		} else if pod.Status.Phase != v1.PodSucceeded {
			return false, nil
		}
		return true, nil
	}, receiverPod.Name, receiverPod.Namespace)

	if err != nil {
		client.T.Fatalf("Failed waiting for completeness of the pod: %v", err)
	}
	notify <- true
}
