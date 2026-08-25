package main

type healthResponse struct {
	Status string `json:"status"`
}

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
}
