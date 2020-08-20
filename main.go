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
	// clientset, err := kubeconfig.initOutOfClusterClientConfig();
	clientset, err := kubeconfig.initInClusterClientConfig()
	if err != nil {
		panic(err.Error())
	}

	nsList, err := clientset.
		CoreV1().
		Namespaces().
		List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	for _, ns := range nsList.Items {
		if ns.Status.Phase == "Terminating" {
			resp, err := clientset.
				CoreV1().
				RESTClient().
				Get().
				AbsPath("/api/v1/namespaces/" + ns.Name).
				DoRaw()
			if err != nil {
				panic(err.Error())
			}

			var nsjson v1.Namespace
			if err := json.Unmarshal(resp, &nsjson); err != nil {
				panic(err.Error())
			}

			// Clear out finalizers 
			// *THIS IS WHAT ACTUALLY FIXES OUR ISSUE*
			nsjson.Spec.Finalizers = make([]v1.FinalizerName, 0)

			nsbody, err := json.Marshal(nsjson)
			if err != nil {
				panic(err.Error())
			}

			if _, err := clientset.
				CoreV1().
				RESTClient().
				Put().
				AbsPath(fmt.Sprintf("/api/v1/namespaces/%s/finalize", ns.Name)).
				Body(nsbody).
				DoRaw(); err != nil {
				panic(err.Error())
			}

			fmt.Printf("Successfully removed namespace: %s\n", ns.Name)
		}
	}
}

type kubeConfig struct {
	path string
}

func (kc *kubeConfig) initOutOfClusterClientConfig() (*kubernetes.Clientset, error) {
	home := os.Getenv("USERPROFILE") // Windows
	if home == "" {                  // Not on Windows
		home = os.Getenv("HOME")
	}
	if home != "" { // If an 'home' path was found for current OS
		flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()
	kc.path = flag.Lookup("kubeconfig").Value.String()
	config, err := clientcmd.BuildConfigFromFlags("", kc.path)
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
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	fmt.Println(config.Username)
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return clientset, nil
}
