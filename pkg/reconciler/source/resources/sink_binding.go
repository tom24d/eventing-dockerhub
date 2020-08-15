package resources

import (
	"context"
	"fmt"

	// k8s.io imports
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	// knative.dev/eventing imports
	"knative.dev/eventing/pkg/apis/sources/v1beta1"

	// knative.dev/pkg imports
	duckv1 "knative.dev/pkg/apis/duck/v1"
	duckv1beta1 "knative.dev/pkg/apis/duck/v1beta1"
	"knative.dev/pkg/kmeta"
	"knative.dev/pkg/tracker"
)

func SinkBindingName(source, subject string) string {
	return kmeta.ChildName(fmt.Sprintf("%s-%s", source, subject), "-sinkbinding")
}

func MakeSinkBinding(owner kmeta.OwnerRefable, source duckv1.SourceSpec, subject tracker.Reference) *v1beta1.SinkBinding {
	sb := &v1beta1.SinkBinding{
		ObjectMeta: metav1.ObjectMeta{
			OwnerReferences: []metav1.OwnerReference{
				*kmeta.NewControllerRef(owner),
			},
			Name:      SinkBindingName(owner.GetObjectMeta().GetName(), subject.Name),
			Namespace: owner.GetObjectMeta().GetNamespace(),
		},
		Spec: v1beta1.SinkBindingSpec{
			SourceSpec: source,
			BindingSpec: duckv1beta1.BindingSpec{
				Subject: subject,
			},
		},
	}

	sb.SetDefaults(context.Background())
	return sb
}
