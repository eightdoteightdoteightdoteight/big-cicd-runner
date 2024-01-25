package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

type JobResult struct {
	Logs   string `json:"logs"`
	Status string `json:"status"`
	Name   string `json:"name"`
}

type Pipeline struct {
	StagesList []string            `yaml:"stages"`
	Stages     map[string][]string `yaml:",inline"`
}

func (p *Pipeline) Unmarshal(data []byte) error {
	return yaml.Unmarshal(data, &p)
}

func readYaml(path string) (Pipeline, error) { // path is currently test.yml
	yamlFile, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("Erreur de lecture du fichier YAML: %v\n", err)
		return Pipeline{}, err
	}

	var pipeline Pipeline
	err = pipeline.Unmarshal(yamlFile)

	if err != nil {
		fmt.Printf("Erreur de parsing YAML: %v\n", err)
		return Pipeline{}, err
	}

	return pipeline, err
}

func stagesExecution(path string, jobID string) {
	pipeline, err := readYaml(path)

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	var fullOutput bytes.Buffer

	for _, stageName := range pipeline.StagesList {
		var status string = "Success"
		stageContent := pipeline.Stages[stageName]

		for _, command := range stageContent {
			toExec := strings.Fields(command)
			cmd := exec.Command("cmd", append([]string{"/c", toExec[0]}, toExec[1:]...)...) //les ... permettent de traiter chaque élément de la liste individuellement
			output, err := cmd.CombinedOutput()

			if err != nil {
				status = "Error"
				fmt.Println("Erreur lors de l'exécution de la commande:", string(output))
				break
			}

			fullOutput.WriteString(fmt.Sprintf("%s\n %s\n", command, output))
		}

		finalOutput := fullOutput.String()
		logs := fmt.Sprintf(fullOutput.String())
		sendPostRequest(jobID, stageName, logs, status)
		fmt.Println("Output complet:")
		fmt.Println(finalOutput)
	}
}

func sendPostRequest(jobID, name string, logs string, status string) {
	// Create the payload
	payload := JobResult{
		Name:   name,
		Logs:   logs,
		Status: status,
	}

	// Convert payload to JSON
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Make the POST request
	url := fmt.Sprintf("https://cicd-back.nathanaudvard.fr/v1/jobs/%s", jobID) // Update with your actual endpoint
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
