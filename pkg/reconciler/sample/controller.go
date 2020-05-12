package sample

// import (
// 	"context"

// 	"knative.dev/sample-source/pkg/apis/samples/v1alpha1"

// 	"github.com/kelseyhightower/envconfig"
// 	"k8s.io/client-go/tools/cache"

// 	"knative.dev/pkg/configmap"
// 	"knative.dev/pkg/controller"
// 	"knative.dev/pkg/logging"

// 	"knative.dev/sample-source/pkg/reconciler"

// 	eventingclient "knative.dev/eventing/pkg/client/injection/client"
// 	sinkbindinginformer "knative.dev/eventing/pkg/client/injection/informers/sources/v1alpha2/sinkbinding"
// 	kubeclient "knative.dev/pkg/client/injection/kube/client"
// 	deploymentinformer "knative.dev/pkg/client/injection/kube/informers/apps/v1/deployment"
// 	samplesourceinformer "knative.dev/sample-source/pkg/client/injection/informers/samples/v1alpha1/samplesource"
// 	"knative.dev/sample-source/pkg/client/injection/reconciler/samples/v1alpha1/samplesource"
// )

// // NewController initializes the controller and is called by the generated code
// // Registers event handlers to enqueue events
// func NewController(
// 	ctx context.Context,
// 	cmw configmap.Watcher,
// ) *controller.Impl {
// 	deploymentInformer := deploymentinformer.Get(ctx)
// 	sinkBindingInformer := sinkbindinginformer.Get(ctx)
// 	sampleSourceInformer := samplesourceinformer.Get(ctx)

// 	r := &Reconciler{
// 		dr:  &reconciler.DeploymentReconciler{KubeClientSet: kubeclient.Get(ctx)},
// 		sbr: &reconciler.SinkBindingReconciler{EventingClientSet: eventingclient.Get(ctx)},
// 	}
// 	if err := envconfig.Process("", r); err != nil {
// 		logging.FromContext(ctx).Panicf("required environment variable is not defined: %v", err)
// 	}

// 	impl := samplesource.NewImpl(ctx, r)

// 	logging.FromContext(ctx).Info("Setting up event handlers")

// 	sampleSourceInformer.Informer().AddEventHandler(controller.HandleAll(impl.Enqueue))

// 	deploymentInformer.Informer().AddEventHandler(cache.FilteringResourceEventHandler{
// 		FilterFunc: controller.FilterGroupKind(v1alpha1.Kind("SampleSource")),
// 		Handler:    controller.HandleAll(impl.EnqueueControllerOf),
// 	})

// 	sinkBindingInformer.Informer().AddEventHandler(cache.FilteringResourceEventHandler{
// 		FilterFunc: controller.FilterGroupKind(v1alpha1.Kind("SampleSource")),
// 		Handler:    controller.HandleAll(impl.EnqueueControllerOf),
// 	})

// 	return impl
// }
