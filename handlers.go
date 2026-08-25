package main

import (
	"encoding/json"
	"log"
	"net/http"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

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

func namespaceHandler(clientset kubernetes.Interface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		namespaces, err := clientset.CoreV1().Namespaces().List(
			r.Context(),
			metav1.ListOptions{},
		)

		if err != nil {
			http.Error(w, "failed to query Kubernetes namespace", http.StatusInternalServerError)
			return
		}

		response := []namespaceResponse{}

		for _, namespace := range namespaces.Items {
			item := namespaceResponse{
				Name:   namespace.Name,
				Status: string(namespace.Status.Phase),
			}

			response = append(response, item)
		}

		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf("failed to encode namespace response: %v", err)
		}

	}
}

func deploymentsHandler(clientset kubernetes.Interface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		namespace := r.PathValue("namespace")

		deployments, err := clientset.AppsV1().Deployments(namespace).List(
			r.Context(),
			metav1.ListOptions{},
		)
		if err != nil {
			http.Error(
				w,
				"failed to query Kubernetes deployments",
				http.StatusInternalServerError,
			)
		}

		response := []deploymentResponse{}

		for _, deployment := range deployments.Items {
			desired := int32(1)

			if deployment.Spec.Replicas != nil {
				desired = *deployment.Spec.Replicas
			}

			status := "healthy"

			if deployment.Status.ReadyReplicas < desired {
				status = "degraded"
			}

			item := deploymentResponse{
				Name:    deployment.Name,
				Ready:   deployment.Status.ReadyReplicas,
				Desired: desired,
				Status:  status,
			}

			response = append(response, item)
		}

		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf("failed to encode deployments response: %v", err)
		}
	}
}
