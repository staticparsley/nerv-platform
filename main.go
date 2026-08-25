package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {

	config, err := kubernetesConfig()

	if err != nil {
		log.Fatalf("failed to build Kubernetes config: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("failed to create Kubernetes client: %v", err)
	}

	http.HandleFunc("GET /healthz", healthHandler)

	http.HandleFunc("GET /api/cluster", clusterHandler(clientset))

	http.HandleFunc("GET /api/namespaces", namespaceHandler(clientset))

	http.HandleFunc("GET /api/namespaces/{namespace}/deployments", deploymentsHandler(clientset))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func kubernetesConfig() (*rest.Config, error) {
	config, err := rest.InClusterConfig()
	if err == nil {
		return config, nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	kubeconfig := filepath.Join(home, ".kube", "config")

	return clientcmd.BuildConfigFromFlags("", kubeconfig)
}
