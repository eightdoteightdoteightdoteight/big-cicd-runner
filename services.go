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

	url := fmt.Sprintf("https://cicd-back.nathanaudvard.fr/v1/jobs/%s", jobID)
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

	fmt.Println("Response status:", resp.Status)
}

func finishPipeline(jobID string) {
	url := fmt.Sprintf("https://cicd-back.nathanaudvard.fr/v1/pipelines/%s/finish", jobID)
	_, err := http.Post(url, "application/json", nil)
	if err != nil {
		fmt.Println("Error sending POST request:", err)
		return
	}
}
