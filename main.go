package main

import (
	"encoding/json"
	"fmt"
	"os"
	"flag"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func main() {
	kubeconfig := kubeConfig{};
	kubeconfig.init();

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig.path)
	if err != nil {
		panic(err.Error());
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error());
	}

	nsList, err := clientset.CoreV1().Namespaces().List(metav1.ListOptions{});
	if err != nil {
		panic(err.Error());
	}

	for _, ns := range nsList.Items {
		if ns.Status.Phase == "Terminating" {
			var resp []byte
			if resp, err = clientset.
				CoreV1().
				RESTClient().
				Get().
				AbsPath("/api/v1/namespaces/" + ns.Name).
				DoRaw(); err != nil {
				panic(err.Error());
			}
	 
			var nsjson v1.Namespace;
			if err := json.Unmarshal(resp, &nsjson); err != nil {
				panic(err.Error());
			}

			// Clear out finalizers *this is what actually fixes our issue*
			nsjson.Spec.Finalizers = make([]v1.FinalizerName, 0);

			var nsbody []byte;
			if nsbody, err = json.Marshal(nsjson); err != nil {
				panic(err.Error());
			}

			abspath := fmt.Sprintf("/api/v1/namespaces/%s/finalize", ns.Name);
			if _, err = clientset.
				CoreV1().
				RESTClient().
				Put().
				AbsPath(abspath).
				Body(nsbody).
				DoRaw(); err != nil {
				panic(err.Error());
			}

			fmt.Println("Successfully removed namespace: " + ns.Name);
		}
	}
}

type kubeConfig struct {
	path string
}

func (kc *kubeConfig) init() {
	home := os.Getenv("USERPROFILE") // Windows
	if home == "" { // Not on Windows
		home = os.Getenv("HOME");
	}

	if home != "" { // If an 'home' path was found for current OS
		flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file");
	} else {
		flag.String("kubeconfig", "", "absolute path to the kubeconfig file");
	}

	flag.Parse();
	kc.path = flag.Lookup("kubeconfig").Value.String();
}