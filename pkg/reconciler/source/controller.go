package source

import (
	"context"
	"os"

	// k8s.io imports
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/cache"

	//Injection imports
	sourcescheme "github.com/tom24d/eventing-dockerhub/pkg/client/clientset/versioned/scheme"
	dockerhubinformer "github.com/tom24d/eventing-dockerhub/pkg/client/injection/informers/sources/v1alpha1/dockerhubsource"
	dhreconciler "github.com/tom24d/eventing-dockerhub/pkg/client/injection/reconciler/sources/v1alpha1/dockerhubsource"

	serviceclient "knative.dev/serving/pkg/client/injection/client"
	kserviceinformer "knative.dev/serving/pkg/client/injection/informers/serving/v1/service"

	//knative.dev/eventing imports
	reconcilersource "knative.dev/eventing/pkg/reconciler/source"

	//knative.dev/pkg import
	kubeclient "knative.dev/pkg/client/injection/kube/client"
	"knative.dev/pkg/configmap"
	"knative.dev/pkg/controller"
	"knative.dev/pkg/logging"
	"knative.dev/pkg/resolver"

	"github.com/tom24d/eventing-dockerhub/pkg/apis/sources/v1alpha1"
)

func NewController(
	ctx context.Context,
	cmw configmap.Watcher,
) *controller.Impl {

	raImage, defined := os.LookupEnv(raImageEnvVar)
	if !defined {
		logging.FromContext(ctx).Errorf("required environment variable '%s' not defined", raImageEnvVar)
		return nil
	}

	ksvcInformer := kserviceinformer.Get(ctx)
	dockerhubInformer := dockerhubinformer.Get(ctx)

	r := &Reconciler{
		kubeClientSet:       kubeclient.Get(ctx),
		servingLister:       ksvcInformer.Lister(),
		servingClientSet:    serviceclient.Get(ctx),
		receiveAdapterImage: raImage,
		configAccessor:      reconcilersource.WatchConfigurations(ctx, "dockerhub_source", cmw),
	}

	impl := dhreconciler.NewImpl(ctx, r)

	r.sinkResolver = resolver.NewURIResolver(ctx, impl.EnqueueKey)

	logging.FromContext(ctx).Info("Setting up DockerHub event handlers")

	dockerhubInformer.Informer().AddEventHandler(controller.HandleAll(impl.Enqueue))

	ksvcInformer.Informer().AddEventHandler(cache.FilteringResourceEventHandler{
		FilterFunc: controller.FilterControllerGK(v1alpha1.Kind("DockerHubSource")),
		Handler:    controller.HandleAll(impl.EnqueueControllerOf),
	})

	return impl
}

func init() {
	sourcescheme.AddToScheme(scheme.Scheme)
}
