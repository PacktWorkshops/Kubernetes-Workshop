package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"k8s.io/api/admission/v1beta1"
	"log"
	"net/http"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"

)

type WebHookServer struct {
	server *http.Server
}

//a function pointer to the Mutation or te Validation WebHook core logic
type controllerFunc func(admissionRequest *v1beta1.AdmissionRequest) (*v1beta1.AdmissionResponse, error)

var (
	runtimeScheme = runtime.NewScheme()
	codecs        = serializer.NewCodecFactory(runtimeScheme)
	deserializer  = codecs.UniversalDeserializer()
	podResource   = metav1.GroupVersionResource{Version: "v1", Resource: "pods"}

	// (https://github.com/kubernetes/kubernetes/issues/57982)
	defaulter = runtime.ObjectDefaulter(runtimeScheme)
)

func main() {

	log.Print("Loading the certificates ....")
	certificateFile := filepath.Join(tlsFolder, tlsCertificate)
	privateKeyFile := filepath.Join(tlsFolder, tlsPrivateKey)

	pair, err := tls.LoadX509KeyPair(certificateFile, privateKeyFile)
	if err != nil {
		log.Fatal("failed to load key pair: %v", err)
	}

	whsvr := &WebHookServer{
		server: &http.Server{
			Addr:      fmt.Sprintf(":%v", serverPort),
			TLSConfig: &tls.Config{Certificates: []tls.Certificate{pair}},
		},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/mutate", whsvr.serveMutatingRequest)
	mux.HandleFunc("/validate", whsvr.serveValidationRequest)
	server := &http.Server{
		Addr:    serverPort,
		Handler: mux,
	}
	log.Fatal(server.ListenAndServeTLS(certificateFile, privateKeyFile))
}

func (whsvr *WebHookServer) serveMutatingRequest(w http.ResponseWriter, r *http.Request) {
	serve(w, r, MutateCustomAnnotation)
}

func (whsvr *WebHookServer) serveValidationRequest(w http.ResponseWriter, r *http.Request) {
	serve(w, r, ValidateTeamAnnotation)
}

//This is not a production ready code.
func serve(w http.ResponseWriter, r *http.Request, f controllerFunc) {

	var body []byte
	if r.Body != nil {
		if data, err := ioutil.ReadAll(r.Body); err == nil {
			body = data
		}
	}
	if len(body) == 0 {
		log.Println("empty body")
		http.Error(w, "empty body", http.StatusBadRequest)
		return
	}


	var admissionReviewReq v1beta1.AdmissionReview

	if _, _, err := deserializer.Decode(body, nil, &admissionReviewReq); err != nil {
		http.Error(w, "cannot deserialize bod payload", http.StatusBadRequest)
		return
	} else if admissionReviewReq.Request == nil {
		http.Error(w, "request is null", http.StatusBadRequest)
		return
	}

	//check if this is a POD
	if admissionReviewReq.Request.Resource != podResource {
		log.Printf("expect resource to be %s", podResource)
		http.Error(w, "request is NOT pod", http.StatusBadRequest)
	}

	var admissionResponse *v1beta1.AdmissionResponse
	admissionResponse, err := f(admissionReviewReq.Request)
	log.Printf("admission Respnse created %v", admissionResponse)

	if err != nil {
		log.Printf("err from admission response %v", err)
		admissionResponse = &v1beta1.AdmissionResponse{
			Allowed: false,
			Result: &metav1.Status{
				Message: err.Error(),
			},
		}
	}

	admissionReview := v1beta1.AdmissionReview{}
	if admissionResponse != nil {
		admissionReview.Response = admissionResponse
		if admissionReviewReq.Request != nil {
			admissionReview.Response.UID = admissionReviewReq.Request.UID
			admissionReview.Response.Allowed = admissionResponse.Allowed
		}
	}

	resp, err := json.Marshal(admissionReview)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not encode response: %v", err), http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(resp); err != nil {
		http.Error(w, fmt.Sprintf("could not write response: %v", err), http.StatusInternalServerError)
	}

}

type patchOperation struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value,omitempty"`
}

//patch = append(patch, updateAnnotation(pod.Annotations, annotations)...)
/*
func updateAnnotation(target map[string]string, added map[string]string) (patch []patchOperation) {
	for key, value := range added {
		if target == nil || target[key] == "" {
			target = map[string]string{}
			patch = append(patch, patchOperation{
				Op:   "add",
				Path: "/metadata/annotations",
				Value: map[string]string{
					key: value,
				},
			})
		} else {
			patch = append(patch, patchOperation{
				Op:    "replace",
				Path:  "/metadata/annotations/" + key,
				Value: value,
			})
		}
	}
	return patch
}
*/
const (
	tlsFolder       = `/etc/secrets/tls`
	tlsCertificate  = `tls.crt`
	tlsPrivateKey   = `tls.key`
	serverPort      = `:8443`

)
