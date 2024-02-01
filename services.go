package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type JobResult struct {
	Logs   string `json:"logs"`
	Status string `json:"status"`
	Name   string `json:"name"`
}

func sendJobResult(jobID string, name string, logs string, status string) {
	payload := JobResult{
		Name:   name,
		Logs:   logs,
		Status: status,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	url := fmt.Sprintf("http://cicd-back-service:8080/v1/jobs/%s", jobID)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		fmt.Println("Error sending POST request:", err)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)
}

func finishPipeline(jobID string, status string) {
	url := fmt.Sprintf("http://cicd-back-service:8080/v1/pipelines/%s/finish?status=%s", jobID, status)
	_, err := http.Post(url, "application/json", nil)
	if err != nil {
		fmt.Println("Error sending POST request:", err)
		return
	}
}

func isHealthy(namespace, name string) bool {
	url := fmt.Sprintf("http://%s-service.%s.svc.cluster.local:8080/actuator/health", name, namespace)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error sending GET request:", err)
		return false
	}
	return resp.StatusCode == http.StatusOK
}
