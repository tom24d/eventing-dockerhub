package source

import (
	"context"
	"fmt"

	//k8s.io imports
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	//knative.dev/serving imports
	v1 "knative.dev/serving/pkg/apis/serving/v1"
	servingclientset "knative.dev/serving/pkg/client/clientset/versioned"
	servinglisters "knative.dev/serving/pkg/client/listers/serving/v1"

	//knative/eventing imports
	eventingclient "knative.dev/eventing/pkg/client/clientset/versioned"
	reconcilersource "knative.dev/eventing/pkg/reconciler/source"

	// github.com/tom24d/eventing-dockerhub imports
	"github.com/tom24d/eventing-dockerhub/pkg/apis/sources/v1alpha1"
	dhreconciler "github.com/tom24d/eventing-dockerhub/pkg/client/injection/reconciler/sources/v1alpha1/dockerhubsource"
	"github.com/tom24d/eventing-dockerhub/pkg/reconciler/source/resources"

	"knative.dev/pkg/apis"
	"knative.dev/pkg/controller"
	"knative.dev/pkg/logging"
	pkgreconciler "knative.dev/pkg/reconciler"
	"knative.dev/pkg/resolver"
	"knative.dev/pkg/tracker"
)

const (
	raImageEnvVar = "DH_RA_IMAGE"
)

// Reconciler reconciles a DockerHubSource object
type Reconciler struct {
	kubeClientSet kubernetes.Interface

	servingClientSet servingclientset.Interface
	servingLister    servinglisters.ServiceLister

	eventingClientSet eventingclient.Interface

	receiveAdapterImage string

	sinkResolver *resolver.URIResolver

	configAccessor reconcilersource.ConfigAccessor
}

// // Check that our Reconciler implements Interface
var _ dhreconciler.Interface = (*Reconciler)(nil)

// // ReconcileKind implements Interface.ReconcileKind.
func (r *Reconciler) ReconcileKind(ctx context.Context, src *v1alpha1.DockerHubSource) pkgreconciler.Event {
	ksvc, err := r.getOwnedService(ctx, src)
	if apierrors.IsNotFound(err) {
		ksvc = r.getExpectedService(ctx, src)
		ksvc, err = r.servingClientSet.ServingV1().Services(src.Namespace).Create(ctx, ksvc, metav1.CreateOptions{})
		if err != nil {
			return err
		}
		src.Status.AutoCallbackDisabled = src.Spec.DisableAutoCallback
		controller.GetEventRecorder(ctx).Eventf(src, corev1.EventTypeNormal, "ServiceCreated", "Created Service %q", ksvc.Name)
	} else if err != nil {
		src.Status.MarkNoEndpoint("ServiceUnavailable", "%v", err)
		return err
	} else if !metav1.IsControlledBy(ksvc, src) {
		src.Status.MarkNoEndpoint("ServiceNotOwned", "Service %q is not owned by DockerHubSource %q", ksvc.Name, src.Name)
		return fmt.Errorf("service %q is not owned by DockerHubSource %q", ksvc.Name, src.Name)
	} else if ksvc == nil {
		logging.FromContext(ctx).Fatalf("knative service is not set without error. src: %v", src)
		return fmt.Errorf("knative service is not set without error")
	}
	src.Status.ReceiveAdapterServiceName = ksvc.Name

	// reconcile SinkBinding for the kservice.
	// remove this logic once https://github.com/knative/eventing/issues/3422 lands.
	logging.FromContext(ctx).Info("going to ReconcileSinkBinding")
	sb, errEvent := r.ReconcileSinkBinding(ctx, src, src.Spec.SourceSpec, tracker.Reference{
		APIVersion: v1.SchemeGroupVersion.String(),
		Kind:       "Service",
		Namespace:  ksvc.Namespace,
		Name:       ksvc.Name,
	})
	if sb != nil {
		s := sb.Status.GetCondition(apis.ConditionReady)
		if s.IsTrue() {
			src.Status.MarkSink(sb.Status.SinkURI)
		} else if s.IsFalse() {
			// SinkBindingReconcileFailed is a propagated status from SinkBinding controller.
			src.Status.MarkNoSink("SinkBindingReconcileFailed", "%s", s.GetMessage())
		}
	}
	if errEvent != nil {
		// FailedReconcileSinkBinding represents this controller itself failed to reconcile SinkBinding resource.
		src.Status.MarkNoSink("FailedReconcileSinkBinding", "%s", errEvent)
		return errEvent
	}

	// if user modifies DisableAutoCallback field
	if src.Status.AutoCallbackDisabled != src.Spec.DisableAutoCallback {
		if len(ksvc.Spec.Template.Spec.Containers) >= 1 {
			err = pkgreconciler.RetryUpdateConflicts(func(_ int) error {
				// Fetch ksvc in case ReconcileSinkBinding might update ksvc above.
				k, err := r.getOwnedService(ctx, src)
				if err != nil {
					return err
				}
				ksvc = k.DeepCopy()
				// override env
				ksvc.Spec.Template.Spec.Containers[0].Env = r.getServiceArgs(ctx, src).GetEnv()

				ksvc, err = r.servingClientSet.ServingV1().Services(src.Namespace).Update(ctx, ksvc, metav1.UpdateOptions{})
				return err
			})
			if err != nil {
				src.Status.MarkNoEndpoint("ServiceUpdateFailed", "failed to update service: %v", err)
				return err
			}
		}

		controller.GetEventRecorder(ctx).
			Eventf(src, corev1.EventTypeNormal,
				"ServiceUpdated", "Updated disableAutoCallback: %t", src.Spec.DisableAutoCallback)
		src.Status.AutoCallbackDisabled = src.Spec.DisableAutoCallback
	}

	if ksvc.IsReady() && ksvc.Status.URL != nil {
		src.Status.MarkEndpoint(ksvc.Status.URL)
	}

	return nil
}

func (r *Reconciler) getOwnedService(ctx context.Context, src *v1alpha1.DockerHubSource) (*v1.Service, error) {
	serviceList, err := r.servingClientSet.ServingV1().Services(src.Namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for i := range serviceList.Items {
		if metav1.IsControlledBy(&serviceList.Items[i], src) {
			return &serviceList.Items[i], nil
		}
	}
	return nil, apierrors.NewNotFound(v1.Resource("services"), "")
}

func (r *Reconciler) getExpectedService(ctx context.Context, src *v1alpha1.DockerHubSource) *v1.Service {
	ksvc := resources.MakeService(r.getServiceArgs(ctx, src))
	if firstName := src.Status.ReceiveAdapterServiceName; firstName != "" {
		ksvc.ObjectMeta.SetGenerateName("")
		ksvc.ObjectMeta.SetName(firstName)
	}
	return ksvc
}

func (r *Reconciler) getServiceArgs(ctx context.Context, src *v1alpha1.DockerHubSource) *resources.ServiceArgs {
	return &resources.ServiceArgs{
		Source:              src,
		ReceiveAdapterImage: r.receiveAdapterImage,
		EventSource:         src.Namespace + "/" + src.Name,
		Context:             ctx,
		AdditionalEnvs:      r.configAccessor.ToEnvVars(), // Grab config envs for tracing/logging/metrics
	}
}
