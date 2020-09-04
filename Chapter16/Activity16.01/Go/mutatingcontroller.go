package main

import (
	"errors"
	"k8s.io/api/admission/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"encoding/json"
)

/**
This method checks inserts a kubernetes annoation to the pod specification. The annotation name is podModified and value is true
This function generates the response for the API call made by Kubernetes Admission Controller for the Mutation webhook. This functions
updates the pod specification.


This code is not ready for production.

*/

func MutateCustomAnnotation(admissionRequest *v1beta1.AdmissionRequest ) (*v1beta1.AdmissionResponse, error){

	// Parse the Pod object.
	raw := admissionRequest.Object.Raw
	pod := corev1.Pod{}
	if _, _, err := deserializer.Decode(raw, nil, &pod); err != nil {
		return nil, errors.New("unable to parse pod")
	}

	//create annotation to add
	annotations := map[string]string{"podModified" : "true"}

	//prepare the patch to be applied to the object
	var patch []patchOperation
	patch = append(patch, patchOperation{
		Op:   "add",
		Path: "/metadata/annotations",
		Value: annotations,
	})

	//convert patch into bytes
	patchBytes, err := json.Marshal(patch)
	if err != nil {
		return nil, errors.New("unable to parse the patch")
	}

	//create the response with patch bytes
	var admissionResponse *v1beta1.AdmissionResponse
	admissionResponse = &v1beta1.AdmissionResponse {
		Allowed: true,
		Patch:   patchBytes,
		PatchType: func() *v1beta1.PatchType {
			pt := v1beta1.PatchTypeJSONPatch
			return &pt
		}(),
	}

	//return the response
	return admissionResponse, nil

}
