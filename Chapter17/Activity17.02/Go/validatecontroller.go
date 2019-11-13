package main

import (
	"errors"
	"k8s.io/api/admission/v1beta1"
	corev1 "k8s.io/api/core/v1"
)

/**
	This method checks if a particular teamName label is available in the pod spec.
	This function generates the response for the API call made by Kubernetes ADmission Controller for the Validation webhook.


	This code is not ready for production.

 */
func ValidateTeamAnnotation(admissionRequest *v1beta1.AdmissionRequest ) (*v1beta1.AdmissionResponse, error){

	// Parse the Pod object.
	raw := admissionRequest.Object.Raw
	pod := corev1.Pod{}
	if _, _, err := deserializer.Decode(raw, nil, &pod); err != nil {
		return nil, errors.New("unable to parse pod")
	}

	//Get all the Labels of the Pod
	podLabels := pod.ObjectMeta.GetLabels()

	//check if the teamName label is available, if not generate an error.
	if podLabels == nil || podLabels[teamNameLabel] == "" {
		return nil, errors.New("teamName label not found")
	}

	//if the teamName label exists, return the response with Allowed flag set to true.
	var admissionResponse *v1beta1.AdmissionResponse
	admissionResponse = &v1beta1.AdmissionResponse {
		Allowed: true,
	}

	//return the response with Allowed set to true
	return admissionResponse, nil

}



const (
	//this is the name of the label that is expected to be part of the pods to allow them to be created.
	teamNameLabel = `teamName`

)
