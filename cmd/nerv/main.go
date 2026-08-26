package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type clusterResponse struct {
	Status string `json:"status"`
	Nodes  int    `json:"nodes"`
	Pods   int    `json:"pods"`
}

type namespaceResponse struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

type deploymentResponse struct {
	Name    string `json:"name"`
	Ready   int32  `json:"ready"`
	Desired int32  `json:"desired"`
	Status  string `json:"status"`
}

func newClient() *http.Client {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	return &http.Client{
		Transport: transport,
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: nerv <command>")
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "cluster":
		if err := getCluster(); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	case "namespaces":
		if err := getNamespaces(); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	case "deployments":
		if len(os.Args) < 3 {
			fmt.Println("usage: nerv deployments <namespace>")
			os.Exit(1)
		}

		namespace := os.Args[2]

		if err := getDeployments(namespace); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	default:
		fmt.Printf("unknown command: %s\n", command)
		os.Exit(1)
	}
}

func getCluster() error {

	client := newClient()

	resp, err := client.Get("https://nerv-platform.nerv.local/api/cluster")
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API returned %s", resp.Status)
	}

	var cluster clusterResponse

	if err := json.NewDecoder(resp.Body).Decode(&cluster); err != nil {
		return err
	}

	fmt.Println("NERV Cluster")
	fmt.Println()
	fmt.Printf("Status: %s\n", cluster.Status)
	fmt.Printf("Nodes: %d\n", cluster.Nodes)
	fmt.Printf("Pods: %d\n", cluster.Pods)

	return nil
}

func getNamespaces() error {
	client := newClient()

	resp, err := client.Get("https://nerv-platform.nerv.local/api/namespaces")
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API returned %s", resp.Status)
	}

	var namespaces []namespaceResponse

	if err := json.NewDecoder(resp.Body).Decode(&namespaces); err != nil {
		return err
	}

	fmt.Println("NERV Namespaces")
	fmt.Println()

	for _, namespace := range namespaces {
		fmt.Printf("%-24s %s\n", namespace.Name, namespace.Status)
	}

	return nil

}

func getDeployments(namespace string) error {
	client := newClient()

	url := fmt.Sprintf(
		"https://nerv-platform.nerv.local/api/namespaces/%s/deployments",
		namespace,
	)

	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API returned %s", resp.Status)
	}

	var deployments []deploymentResponse

	if err := json.NewDecoder(resp.Body).Decode(&deployments); err != nil {
		return err
	}

	fmt.Printf("Deployments in %s\n\n", namespace)

	for _, deployment := range deployments {
		fmt.Printf(
			"%-28s %d/%d  %s\n",
			deployment.Name,
			deployment.Ready,
			deployment.Desired,
			deployment.Status,
		)
	}

	return nil
}
