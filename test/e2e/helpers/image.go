package helpers

import (
	"fmt"
	"strconv"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"knative.dev/eventing/test/lib"
	"knative.dev/pkg/test"
)

// CreateValidationReceiverOrFail creates validation-receiver pod or fail.
func CreateValidationReceiverOrFail(client *lib.Client) *v1.Pod {
	const receiverImageName = "validation-receiver"
	args := []string{"--patient=" + strconv.Itoa(60)}

	receiverPod := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: client.Namespace,
			Name:      receiverImageName,
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

	createPodOrFailWithMessage(client, receiverPod)

	return receiverPod
}

// CreateCallbackDisplayOrFail creates callback-display pod or fail.
func CreateCallbackDisplayOrFail(client *lib.Client) *v1.Pod {
	const receiverImageName = "callback-display"

	callbackPod := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: client.Namespace,
			Name:      receiverImageName,
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{{
				Name:            receiverImageName,
				Image:           test.ImagePath(receiverImageName),
				ImagePullPolicy: v1.PullAlways,
				Ports: []v1.ContainerPort{
					{
						ContainerPort: ValidationReceivePort,
					},
				},
			}},
			RestartPolicy: v1.RestartPolicyNever,
		},
	}

	createPodOrFailWithMessage(client, callbackPod)

	return callbackPod
}

func createPodOrFailWithMessage(client *lib.Client, pod *v1.Pod) {
	client.CreatePodOrFail(pod, lib.WithService(pod.GetName()))

	err := test.WaitForPodState(client.Kube, func(p *v1.Pod) (bool, error) {
		if p.Status.Phase == v1.PodFailed {
			return true, fmt.Errorf("pod failed to get up. message: %s", p.Status.Message)
		} else if p.Status.Phase != v1.PodRunning {
			return false, nil
		}
		return true, nil
	}, pod.Name, pod.Namespace)

	if err != nil {
		client.T.Fatalf("Failed waiting for pod running %q: %v", pod.Name, pod)
	}
}
