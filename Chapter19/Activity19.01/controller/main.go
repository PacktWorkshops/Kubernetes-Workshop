package main

import (
	//"github.com/example-inc/pod-normaliser-controller/pkg/client/clientset/versioned"
	//v1 "github.com/example-inc/pod-normaliser-controller/pkg/client/informers/externalversions/podlifecycleconfig/v1"

	v1 "github.com/example-inc/pod-normaliser-controller/pkg/apis/podlifecycleconfig/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"os"
	"os/signal"
	//"sync"
	"syscall"

	/*"github.com/example-inc/pod-normaliser-controller/pkg/client/clientset/versioned"
	"k8s.io/client-go/tools/cache"
	//"github.com/operator-framework/operator-sdk/pkg/k8sutil"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/kubernetes"
	"log"
	kubeinformers "k8s.io/client-go/informers"
	. "time"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/typed/apiextensions/v1"
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/typed/apiextensions/v1beta1"
	*/
	clientset "github.com/example-inc/pod-normaliser-controller/pkg/client/clientset/versioned"
	"k8s.io/client-go/tools/cache"

	informers "github.com/example-inc/pod-normaliser-controller/pkg/client/informers/externalversions"

	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	. "time"
)

func main(){

	// create the kubernetes client configuration
	config, err := clientcmd.BuildConfigFromFlags("", "")
	if err != nil {
		log.Fatal(err)
	}

	// create the kubernetes client
	podlifecyelconfgiclient, err := clientset.NewForConfig(config)


	// create the shared informer factory and use the client to connect to kubernetes
	podlifecycleconfigfactory := informers.NewSharedInformerFactoryWithOptions(podlifecyelconfgiclient, Second*30,
										informers.WithNamespace(os.Getenv(NAMESPACE_TO_WATCH)))


	m := make(map[string]chan struct{})
	podclientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Println("Unable to get pods")
		return
	}

	// fetch the informer for the PodLifecycleConfig
	podlifecycleconfiginformer := podlifecycleconfigfactory.Controllers().V1().PodLifecycleConfigs().Informer()

	// register wit the informaer for the events
	podlifecycleconfiginformer.AddEventHandler(
		cache.ResourceEventHandlerFuncs{

			//define what to do in case if a new custom resource is created
			AddFunc: func(obj interface{}) {
				log.Printf("A Custom Resource has been Added\n")

				//type case the object into an object defined by our custom resource
				x := obj.(*v1.PodLifecycleConfig)
				log.Printf("The CRD is for namespace %s with pod active time as %v", x.Spec.NamespaceName, x.Spec.PodLiveForMinutes)

				// create a new go channel. this channel would be used to stop this go routine if the CR is deleted
				signal := make(chan struct{}, 1)

				//store the channel in a local map
				m[x.Spec.NamespaceName] = signal

				// start the subroutine to check and kill the pods for this namespace
				go checkAndRemovePodsPeriodically(signal, podclientset, x)
			},

			//define what to do in case a custom resource is removed
			DeleteFunc: func(obj interface{}) {
				log.Printf("A Custom Resource has been Deleted\n")
				x := obj.(*v1.PodLifecycleConfig)
				log.Printf("The CRD is for namespace %s with pod active time as %v", x.Spec.NamespaceName, x.Spec.PodLiveForMinutes)

				// since this is a delete event, fetch the signal stored in the local map
				signal := m[x.Spec.NamespaceName]

				// close the signal so the go routine that was periodically checking the pod life allowed time can be stopped
				close(signal)
			},
		},
		)

	podlifecycleconfigfactory.Start(wait.NeverStop)
	podlifecycleconfigfactory.WaitForCacheSync(wait.NeverStop)


	sigTerm := make(chan os.Signal, 1)
	signal.Notify(sigTerm, syscall.SIGTERM)
	signal.Notify(sigTerm, syscall.SIGINT)
	<-sigTerm


}


func checkAndRemovePodsPeriodically(signal chan struct{}, podclientset *kubernetes.Clientset, podlifecycleconfig *v1.PodLifecycleConfig) {
	//this value is not production friendly
	timeChan := Tick(Second * 10)

	for {
		select {
		case <- timeChan:
			checkAndRemovePods(podlifecycleconfig, podclientset)
		case <- signal:
			return
		}
	}
}

func changeNumberOfReplicasToZero(podlifecycleconfig *v1.PodLifecycleConfig, podclientset *kubernetes.Clientset){
	depls, _ := podclientset.AppsV1().Deployments(podlifecycleconfig.Spec.NamespaceName).List(meta_v1.ListOptions{})

	for _, deployment := range depls.Items {
		var zero int32 = 0
		deployment.Spec.Replicas = &zero
		podclientset.AppsV1().Deployments("").Update(&deployment)
	}
}

func checkAndRemovePods(podlifecycleconfig *v1.PodLifecycleConfig, podclientset *kubernetes.Clientset) {
	pods := podclientset.CoreV1().Pods(podlifecycleconfig.Spec.NamespaceName)
	podList, _ := pods.List(meta_v1.ListOptions{})
	log.Printf("Total pods found %v", len(podList.Items))
	for _, pod := range podList.Items {
		    if pod.Name == os.Getenv(SELF_POD_NAME) {
		    	continue
			}
			podStartTime := pod.Status.StartTime.Time
			now := Now()
			diff := int(now.Sub(podStartTime).Minutes())

			if diff >= podlifecycleconfig.Spec.PodLiveForMinutes {
				log.Printf("Pod %s is running for more than the allocated time of %v. The anamoly is %v with pod start time at %v", pod.Name, podlifecycleconfig.Spec.PodLiveForMinutes, diff, pod.Status.StartTime)
				var zero int64 = 0
				err := pods.Delete(pod.Name, &meta_v1.DeleteOptions{GracePeriodSeconds: &zero})
				if err != nil {
					log.Printf(err.Error())
				}

			}


	}
}

const NAMESPACE_TO_WATCH  = "namespaceToWatch"
const SELF_POD_NAME = "podName"
