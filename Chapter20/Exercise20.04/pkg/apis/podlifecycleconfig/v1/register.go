package v1

import (
	"github.com/example-inc/pod-normaliser-controller/pkg/apis/podlifecycleconfig"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// GroupVersion is the identifier for the API which includes
// the name of the group and the version of the API
var SchemeGroupVersion = schema.GroupVersion{
	Group:   podlifecycleconfig.GroupName,
	Version: "v1",

}

// Kind takes an unqualified kind and returns back a Group qualified GroupKind
func Kind(kind string) schema.GroupKind {

	return SchemeGroupVersion.WithKind(kind).GroupKind()
}

// create a SchemeBuilder which uses functions to add types to
// the scheme
//var AddToScheme = runtime.NewSchemeBuilder(addKnownTypes).AddToScheme
var (
		SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)
		AddToScheme = SchemeBuilder.AddToScheme
	)



func Resource(resource string) schema.GroupResource {
	return SchemeGroupVersion.WithResource(resource).GroupResource()
}

// addKnownTypes adds our types to the API scheme by registering
// PodLifecycleConfig and PodLifecycleConfigList
func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(
		SchemeGroupVersion,
		&PodLifecycleConfig{},
		&PodLifecycleConfigList{},
	)

	// register the type in the scheme
	meta_v1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}
