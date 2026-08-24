package main

type healthResponse struct {
	Status string `json:"status"`
}

type clusterResponse struct {
	Status string `json:"status"`
	Nodes  int    `json:"nodes"`
	Pods   int    `json:"pods"`
}
