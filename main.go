package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {

	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("failed to determine user home directory: %v", err)
	}

	kubeconfig := filepath.Join(home, ".kube", "config")

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)

	if err != nil {
		log.Fatalf("failed to build Kubernetes config: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("failed to create Kubernetes client: %v", err)
	}

	http.HandleFunc("GET /healthz", healthHandler)

	http.HandleFunc("GET /api/cluster", clusterHandler(clientset))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
