package source

import (
	"context"
	"fmt"

	//k8s.io imports
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	//knative.dev/serving imports
	servingv1 "knative.dev/serving/pkg/apis/serving/v1"
	servingclientset "knative.dev/serving/pkg/client/clientset/versioned"
	servinglisters "knative.dev/serving/pkg/client/listers/serving/v1"

	// knative.dev/eventing imports
	eventingv1 "knative.dev/eventing/pkg/apis/sources/v1"
	reconcilersource "knative.dev/eventing/pkg/reconciler/source"

	// knative.dev/pkg imports
	duckv1 "knative.dev/pkg/apis/duck/v1"
	"knative.dev/pkg/controller"
	"knative.dev/pkg/logging"
	pkgreconciler "knative.dev/pkg/reconciler"
	"knative.dev/pkg/resolver"

	// github.com/tom24d/eventing-dockerhub imports
	"github.com/tom24d/eventing-dockerhub/pkg/apis/sources/v1alpha1"
	dhreconciler "github.com/tom24d/eventing-dockerhub/pkg/client/injection/reconciler/sources/v1alpha1/dockerhubsource"
	"github.com/tom24d/eventing-dockerhub/pkg/reconciler/source/resources"
)

const (
	raImageEnvVar = "DH_RA_IMAGE"
)

// Reconciler reconciles a DockerHubSource object
type Reconciler struct {
	kubeClientSet kubernetes.Interface

	servingClientSet servingclientset.Interface
	servingLister    servinglisters.ServiceLister

	sinkResolver *resolver.URIResolver

	receiveAdapterImage string

	configAccessor reconcilersource.ConfigAccessor
}

// Check that our Reconciler implements Interface
var _ dhreconciler.Interface = (*Reconciler)(nil)

// ReconcileKind implements Interface.ReconcileKind.
func (r *Reconciler) ReconcileKind(ctx context.Context, src *v1alpha1.DockerHubSource) pkgreconciler.Event {

	ctx = eventingv1.WithURIResolver(ctx, r.sinkResolver)

	ksvc, err := r.getOwnedService(ctx, src)
	if apierrors.IsNotFound(err) {
		ksvc = r.getExpectedService(ctx, src)
		ksvc, err = r.servingClientSet.ServingV1().Services(src.Namespace).Create(ctx, ksvc, metav1.CreateOptions{})
		if err != nil {
			return err
		}
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
	} else if expected := r.getExpectedService(ctx, src); podSpecChanged(expected, ksvc) {
		err := pkgreconciler.RetryUpdateConflicts(func(int) error {
			ksvc, err := r.getOwnedService(ctx, src)
			if err != nil {
				return err
			}
			syncPodSpec(expected, ksvc)
			ksvc, err = r.servingClientSet.ServingV1().Services(src.Namespace).Update(ctx, ksvc, metav1.UpdateOptions{})
			return err
		})
		if err != nil {
			src.Status.MarkNoEndpoint("ServiceUpdateFailed", "failed to update service: %v", err)
			return err
		}

		controller.GetEventRecorder(ctx).
			Eventf(src, corev1.EventTypeNormal,
				"ServiceUpdated", "Updated disableAutoCallback: %t", src.Spec.DisableAutoCallback)
	}

	// status propagation
	src.Status.ReceiveAdapterServiceName = ksvc.Name
	src.Status.AutoCallbackDisabled = src.Spec.DisableAutoCallback
	if ksvc.IsReady() && ksvc.Status.URL != nil {
		src.Status.MarkEndpoint(ksvc.Status.URL)
	}

	return nil
}

func podSpecChanged(expected *servingv1.Service, now *servingv1.Service) bool {
	old := now.DeepCopy()
	syncPodSpec(expected, now)
	return !equality.Semantic.DeepEqual(old.Spec.Template.Spec.PodSpec, now.Spec.Template.Spec.PodSpec)
}

func syncPodSpec(expected *servingv1.Service, now *servingv1.Service) {
	now.Spec.Template.Spec.PodSpec = expected.Spec.GetTemplate().Spec.PodSpec
}

func (r *Reconciler) getOwnedService(ctx context.Context, src *v1alpha1.DockerHubSource) (*servingv1.Service, error) {
	serviceList, err := r.servingClientSet.ServingV1().Services(src.Namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for i := range serviceList.Items {
		if metav1.IsControlledBy(&serviceList.Items[i], src) {
			return &serviceList.Items[i], nil
		}
	}
	return nil, apierrors.NewNotFound(servingv1.Resource("services"), "")
}

func (r *Reconciler) getExpectedService(ctx context.Context, src *v1alpha1.DockerHubSource) *servingv1.Service {
	ksvc := resources.MakeService(r.getServiceArgs(ctx, src))
	if firstName := src.Status.ReceiveAdapterServiceName; firstName != "" {
		ksvc.ObjectMeta.SetGenerateName("")
		ksvc.ObjectMeta.SetName(firstName)
	}

	ps := &duckv1.WithPod{}
	ps.Spec.Template.Spec = ksvc.Spec.Template.Spec.PodSpec

	sb := &eventingv1.SinkBinding{}
	sb.Spec.SourceSpec = src.Spec.SourceSpec

	// hack hack hack
	// this is necessary to requeue reconciliation when the sink reference changes.
	sb.SetName(src.GetName())
	sb.SetNamespace(src.GetNamespace())

	sb.Do(ctx, ps)

	if c := sb.Status.GetCondition(eventingv1.SinkBindingConditionSinkProvided); c.IsTrue() && sb.Status.SinkURI != nil {
		src.Status.MarkSink(sb.Status.SinkURI)
	} else if c.IsFalse() {
		src.Status.MarkNoSink(c.GetReason(), "%s", c.GetMessage())
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
