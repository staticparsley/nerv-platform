package main

import (
	"encoding/json"
	"errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	"net/http"
	"net/http/httptest"
	"testing"

	"k8s.io/apimachinery/pkg/runtime"
	ktesting "k8s.io/client-go/testing"
)

func TestHealthHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rr := httptest.NewRecorder()

	healthHandler(rr, req)

	contentType := rr.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", contentType)
	}
	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	var response healthResponse

	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Status != "healthy" {
		t.Errorf("expected status healthy, got %s", response.Status)
	}

}

func TestClusterHandler(t *testing.T) {
	clientset := fake.NewSimpleClientset(
		&corev1.Node{
			ObjectMeta: metav1.ObjectMeta{
				Name: "nerv",
			},
		},
		&corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "pod-1",
				Namespace: "default",
			},
		},
		&corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "pod-2",
				Namespace: "argocd",
			},
		},
		&corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "pod-3",
				Namespace: "monitoring",
			},
		},
	)

	req := httptest.NewRequest(http.MethodGet, "/api/cluster", nil)
	rr := httptest.NewRecorder()

	clusterHandler(clientset)(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	var response clusterResponse

	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Nodes != 1 {
		t.Errorf("expected 1 node, got %d", response.Nodes)
	}

	if response.Pods != 3 {
		t.Errorf("expected 3 pods, got %d", response.Pods)
	}

}

func TestClusterHandlerNodeError(t *testing.T) {
	clientset := fake.NewSimpleClientset()

	clientset.PrependReactor(
		"list",
		"nodes",
		func(action ktesting.Action) (bool, runtime.Object, error) {
			return true, nil, errors.New("Kubernetes API unavailable")
		},
	)

	req := httptest.NewRequest(http.MethodGet, "/api/cluster", nil)
	rr := httptest.NewRecorder()

	handler := clusterHandler(clientset)
	handler(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf(
			"expected status %d, got %d",
			http.StatusInternalServerError,
			rr.Code,
		)
	}
}
