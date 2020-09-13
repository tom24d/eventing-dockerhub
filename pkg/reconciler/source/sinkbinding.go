package source

import (
	"context"
	"fmt"

	// k8s.io imports
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	// knative.dev/eventing imports
	"knative.dev/eventing/pkg/apis/sources/v1alpha2"

	// knative.dev/pkg imports
	duckv1 "knative.dev/pkg/apis/duck/v1"
	"knative.dev/pkg/kmeta"
	"knative.dev/pkg/logging"
	pkgreconciler "knative.dev/pkg/reconciler"
	"knative.dev/pkg/tracker"

	"github.com/tom24d/eventing-dockerhub/pkg/reconciler/source/resources"
	"go.uber.org/zap"
)

// newSinkBindingFailed makes a new reconciler event with event type Warning, and
// reason SinkBindingFailed.
func newSinkBindingFailed(namespace, name string, err error) pkgreconciler.Event {
	return pkgreconciler.NewEvent(corev1.EventTypeWarning, "SinkBindingFailed", "failed to create SinkBinding: \"%s/%s\", %w", namespace, name, err)
}

func (r *Reconciler) ReconcileSinkBinding(ctx context.Context, owner kmeta.OwnerRefable, source duckv1.SourceSpec, subject tracker.Reference) (*v1alpha2.SinkBinding, pkgreconciler.Event) {
	expected := resources.MakeSinkBinding(owner, source, subject)

	namespace := owner.GetObjectMeta().GetNamespace()
	sb, err := r.eventingClientSet.SourcesV1alpha2().SinkBindings(namespace).Get(ctx, expected.Name, metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		sb, err = r.eventingClientSet.SourcesV1alpha2().SinkBindings(namespace).Create(ctx, expected, metav1.CreateOptions{})
		if err != nil {
			return nil, newSinkBindingFailed(expected.Namespace, expected.Name, err)
		}
		return sb, nil
	} else if err != nil {
		return nil, fmt.Errorf("error getting SinkBinding %q: %v", expected.Name, err)
	} else if !metav1.IsControlledBy(sb, owner.GetObjectMeta()) {
		return nil, fmt.Errorf("SinkBinding %q is not owned by %s %q",
			sb.Name, owner.GetGroupVersionKind().Kind, owner.GetObjectMeta().GetName())
	} else if r.specChanged(sb.Spec, expected.Spec) {
		sb.Spec = expected.Spec
		if sb, err = r.eventingClientSet.SourcesV1alpha2().SinkBindings(namespace).Update(ctx, sb, metav1.UpdateOptions{}); err != nil {
			return sb, err
		}
		return sb, nil
	} else {
		logging.FromContext(ctx).Debugw("Reusing existing sink binding", zap.Any("sinkBinding", sb))
	}
	return sb, nil
}

func (r *Reconciler) specChanged(oldSpec v1alpha2.SinkBindingSpec, newSpec v1alpha2.SinkBindingSpec) bool {
	if !equality.Semantic.DeepDerivative(newSpec, oldSpec) {
		return true
	}
	return false
}
