package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

// You are developing an automated internal SRE tool called Namespace Sentinel.
// In multi-tenant clusters, a few misconfigured workloads can quickly exhaust
// cluster resources or trigger alert fatigue.

// Your task is to write a script that scans cluster workloads,
// aggregates health metrics per namespace, and flags namespaces
// where the ratio of unhealthy Pods exceeds a configurable threshold.

type Pod struct {
	PodMetadata PodMetadata `json:"metadata"`
	Status      Status      `json:"status"`
}

type PodMetadata struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

type Status struct {
	Phase      string      `json:"phase"`
	Conditions []Condition `json:"conditions"`
}

type Condition struct {
	Type   string `json:"type"`
	Status string `json:"status"`
}

type HealthRatio struct {
	Healthy   int
	UnHealthy int
}

func main() {

	var pods []Pod

	jsonFile, err := os.ReadFile("./pods.json")
	if err != nil {
		log.Fatal(err)
	}

	// Convert File to Byte Array
	json.Unmarshal(jsonFile, &pods)
	namespaceHR := make(map[string]HealthRatio)

	for _, pod := range pods {
		namespace := pod.PodMetadata.Namespace
		hr := namespaceHR[namespace]
		fmt.Println(hr)
		isHealthy := pod.Status.Phase == "Running" && pod.Status.Conditions[0].Type == "Ready" && pod.Status.Conditions[0].Status == "True"
		if isHealthy {
			hr.Healthy++
		} else {
			hr.UnHealthy++
		}
		namespaceHR[namespace] = hr
	}

	for k, v := range namespaceHR {
		total := v.Healthy + v.UnHealthy
		fmt.Println(total)
		ratio := float64(v.UnHealthy) / float64(total)
		if ratio > 0.25 {
			fmt.Printf("Namespace %s has %d unhealthy pods and its health ratio is %f", k, v.UnHealthy, ratio)
		} else {
			fmt.Printf("Namespace %s is healthy and its health ratio is %f", k, ratio)
		}
	}
}
