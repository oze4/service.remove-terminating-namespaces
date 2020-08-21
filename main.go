package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	kubeconfig := kubeConfig{}

	/* OUT OF CLUSTER CONFIG */
	// clientset, initErr := kubeconfig.initOutOfClusterClientConfig()

	/* IN CLUSTER CONFIG */
	clientset, initErr := kubeconfig.initInClusterClientConfig()

	if initErr != nil {
		panic(initErr.Error())
	}

	nsList, nsListErr := clientset.CoreV1().Namespaces().List(metav1.ListOptions{})
	if nsListErr != nil {
		panic(nsListErr.Error())
	}

	ishealthy := true
	for _, ns := range nsList.Items {
		if ns.Status.Phase == "Terminating" {
			resp := clientset.CoreV1().RESTClient().Get().AbsPath("/api/v1/namespaces/" + ns.Name).Do()
			if getNsErr := resp.Error(); getNsErr != nil {
				panic(getNsErr.Error())
			}

			respRaw, rawErr := resp.Raw()
			if rawErr != nil {
				panic(rawErr.Error())
			}

			var nsjson v1.Namespace
			if umErr := json.Unmarshal(respRaw, &nsjson); umErr != nil {
				panic(umErr.Error())
			}

			// Clear out finalizers *this is what actually fixes our issue*
			nsjson.Spec.Finalizers = make([]v1.FinalizerName, 0)

			nsbody, mErr := json.Marshal(nsjson)
			if mErr != nil {
				panic(mErr.Error())
			}

			r := clientset.CoreV1().RESTClient().Put().AbsPath("/api/v1/namespaces/" + ns.Name + "/finalize").Body(nsbody).Do()
			if rErr := r.Error(); rErr != nil {
				panic(rErr.Error())
			}

			fmt.Printf("Successfully removed namespace: %s\n", ns.Name)
			ishealthy = false
		}
	}

	if ishealthy {
		fmt.Println("All namespaces healthy! (eg: we could not find a namespace in Terminating state)")
	}
}

type kubeConfig struct{}

func (kc *kubeConfig) initOutOfClusterClientConfig() (*kubernetes.Clientset, error) {
	home := os.Getenv("USERPROFILE")
	if home == "" {
		home = os.Getenv("HOME")
	} else if home != "" {
		flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}

	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", flag.Lookup("kubeconfig").Value.String())
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}

func (kc *kubeConfig) initInClusterClientConfig() (*kubernetes.Clientset, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}
