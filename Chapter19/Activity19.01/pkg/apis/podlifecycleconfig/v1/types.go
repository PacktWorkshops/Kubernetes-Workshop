package v1

import (

	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type PodLifecycleConfig struct {

	// TypeMeta is the metadata for the resource, like kind and apiversion
	meta_v1.TypeMeta `json:",inline"`

	// ObjectMeta contains the metadata for the particular object like labels
	meta_v1.ObjectMeta `json:"metadata,omitempty"`

	Spec PodLifecycleConfigSpec `json:"spec"`
}

type PodLifecycleConfigSpec struct{
	NamespaceName   string `json:"namespaceName"`
	PodLiveForMinutes int `json:"podLiveForThisMinutes"`
}



// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type PodLifecycleConfigList struct {
	meta_v1.TypeMeta `json:",inline"`
	meta_v1.ListMeta `json:"metadata"`

	Items []PodLifecycleConfig `json:"items"`
}
