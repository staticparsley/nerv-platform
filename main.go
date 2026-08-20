package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type healthResponse struct {
	Status string `json:"status"`
}

type clusterResponse struct {
	Status string `json:"status"`
	Nodes  int    `json:"nodes"`
	Pods   int    `json:"pods"`
}

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

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := healthResponse{
		Status: "healthy",
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("failed to encode health response: %v", err)
	}
}

func clusterHandler(clientset kubernetes.Interface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		nodes, err := clientset.CoreV1().Nodes().List(
			r.Context(),
			metav1.ListOptions{},
		)
		if err != nil {
			http.Error(w, "failed to query Kubernetes nodes", http.StatusInternalServerError)
			return
		}

		pods, err := clientset.CoreV1().Pods(metav1.NamespaceAll).List(
			r.Context(),
			metav1.ListOptions{},
		)
		if err != nil {
			http.Error(w, "failed to query Kubernetes pods", http.StatusInternalServerError)
			return
		}

		response := clusterResponse{
			Status: "healthy",
			Nodes:  len(nodes.Items),
			Pods:   len(pods.Items),
		}

		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf("failed to encode cluster response: %v", err)
		}
	}
}
