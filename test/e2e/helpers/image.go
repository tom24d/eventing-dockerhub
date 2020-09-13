package helpers

import (
	"context"
	"strconv"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"knative.dev/eventing/test/lib"
	"knative.dev/pkg/test"

	"github.com/google/uuid"
)

// CreateValidationReceiverOrFail creates validation-receiver pod or fail.
func CreateValidationReceiverOrFail(ctx context.Context, client *lib.Client) *v1.Pod {
	const receiverImageName = "validation-receiver"
	args := []string{"--patient=" + strconv.Itoa(180)}

	receiverPod := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: client.Namespace,
			Name:      receiverImageName,
			Labels:    map[string]string{"e2e-test": uuid.New().String()},
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{{
				Name:            receiverImageName,
				Image:           test.ImagePath(receiverImageName),
				ImagePullPolicy: v1.PullAlways,
				Args:            args,
			}},
			RestartPolicy: v1.RestartPolicyNever,
		},
	}

	createPodOrFailWithMessage(ctx, client, receiverPod)

	return receiverPod
}

// CreateCallbackDisplayOrFail creates callback-display pod or fail.
func CreateCallbackDisplayOrFail(ctx context.Context, client *lib.Client) *v1.Pod {
	const receiverImageName = "callback-display"

	callbackPod := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: client.Namespace,
			Name:      receiverImageName,
			Labels:    map[string]string{"e2e-test": uuid.New().String()},
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{{
				Name:            receiverImageName,
				Image:           test.ImagePath(receiverImageName),
				ImagePullPolicy: v1.PullAlways,
			}},
			RestartPolicy: v1.RestartPolicyNever,
		},
	}

	createPodOrFailWithMessage(ctx, client, callbackPod)

	return callbackPod
}

func createPodOrFailWithMessage(ctx context.Context, client *lib.Client, pod *v1.Pod) {
	client.CreatePodOrFail(pod, lib.WithService(pod.GetName()))
	err := test.WaitForPodRunning(ctx, client.Kube, pod.GetName(), client.Namespace)
	if err != nil {
		client.T.Fatalf("Failed waiting for pod running %q: %v", pod.Name, pod)
	}
	client.WaitForServiceEndpointsOrFail(ctx, pod.GetName(), 1)
}
