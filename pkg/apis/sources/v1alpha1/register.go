package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime"
)

var (
	SchemeGroupVersion = schema.GroupVersion{Group: "dockerhubsource.example.com", Version: "v1alpha1"}

	SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)
	AddToScheme = SchemeBuilder.AddToScheme
)


func Kind(kind string) schema.GroupKind{
	return SchemeGroupVersion.WithKind(kind).GroupKind()
}

func Resource(resource string) schema.GroupResource {
	return SchemeGroupVersion.withResource(resource).GroupResource()
}


func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion, 
		&DockerHubSource{}, 
		&DockerHubSourceList{},
	)
	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}