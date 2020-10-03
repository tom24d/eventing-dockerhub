package binding

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"

	"knative.dev/pkg/apis/duck"
	"knative.dev/pkg/client/injection/ducks/duck/v1/podspecable"
	"knative.dev/pkg/client/injection/kube/informers/core/v1/namespace"
	"knative.dev/pkg/configmap"
	"knative.dev/pkg/controller"
	"knative.dev/pkg/injection/clients/dynamicclient"
	"knative.dev/pkg/logging"
	"knative.dev/pkg/reconciler"
	"knative.dev/pkg/resolver"
	"knative.dev/pkg/tracker"
	"knative.dev/pkg/webhook/psbinding"

	"github.com/tom24d/eventing-dockerhub/pkg/apis/sources/v1alpha1"
	"github.com/tom24d/eventing-dockerhub/pkg/client/clientset/versioned/scheme"
	dhbinformer "github.com/tom24d/eventing-dockerhub/pkg/client/injection/informers/sources/v1alpha1/dockerhubsource"
)

const (
	controllerAgentName = "dockerhubbinding-controller"
)

// NewController returns a new DockerHubSource destination reconciler.
func NewController(
	ctx context.Context,
	cmw configmap.Watcher,
) *controller.Impl {
	logger := logging.FromContext(ctx)

	dhbInformer := dhbinformer.Get(ctx)
	dc := dynamicclient.Get(ctx)
	psInformerFactory := podspecable.Get(ctx)
	namespaceInformer := namespace.Get(ctx)

	c := &psbinding.BaseReconciler{
		LeaderAwareFuncs: reconciler.LeaderAwareFuncs{
			PromoteFunc: func(bkt reconciler.Bucket, enq func(reconciler.Bucket, types.NamespacedName)) error {
				all, err := dhbInformer.Lister().List(labels.Everything())
				if err != nil {
					return err
				}
				for _, elt := range all {
					enq(bkt, types.NamespacedName{
						Namespace: elt.GetNamespace(),
						Name:      elt.GetName(),
					})
				}
				return nil
			},
		},
		GVR: v1alpha1.SchemeGroupVersion.WithResource("dockerhubsources"),
		Get: func(namespace string, name string) (psbinding.Bindable, error) {
			return dhbInformer.Lister().DockerHubSources(namespace).Get(name)
		},
		DynamicClient: dc,
		Recorder: record.NewBroadcaster().NewRecorder(
			scheme.Scheme, corev1.EventSource{Component: controllerAgentName}),
		NamespaceLister: namespaceInformer.Lister(),
	}
	impl := controller.NewImpl(c, logger, "DockerHubSources")

	logger.Info("Setting up event handlers")

	dhbInformer.Informer().AddEventHandler(controller.HandleAll(impl.Enqueue))
	namespaceInformer.Informer().AddEventHandler(controller.HandleAll(impl.Enqueue))

	c.WithContext = WithContextFactory(ctx, impl.EnqueueKey)
	c.Tracker = tracker.New(impl.EnqueueKey, controller.GetTrackerLease(ctx))
	c.Factory = &duck.CachedInformerFactory{
		Delegate: &duck.EnqueueInformerFactory{
			Delegate:     psInformerFactory,
			EventHandler: controller.HandleAll(c.Tracker.OnChanged),
		},
	}

	return impl
}

func WithContextFactory(ctx context.Context, handler func(types.NamespacedName)) psbinding.BindableContext {
	r := resolver.NewURIResolver(ctx, handler)

	return func(ctx context.Context, b psbinding.Bindable) (context.Context, error) {
		sb := b.(*v1alpha1.DockerHubSource)
		if sb.Spec.Sink.Ref.Namespace == "" {
			sb.Spec.Sink.Ref.Namespace = sb.Namespace
		}
		uri, err := r.URIFromDestinationV1(ctx, sb.Spec.Sink, sb)
		if err != nil {
			return nil, err
		}
		sb.Status.SinkURI = uri
		return v1alpha1.WithSinkURI(ctx, sb.Status.SinkURI), nil
	}
}
