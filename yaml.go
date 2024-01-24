package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

type JobResult struct {
	ID     string    `json:"id"`
	Logs   string    `json:"logs"`
	Status string    `json:"status"`
	Date   time.Time `json:"date"`
	Stage  string    `json:"stage"`
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

func stagesExecution(path string) {
	pipeline, err := readYaml(path)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	for _, stageName := range pipeline.StagesList {
		stageContent := pipeline.Stages[stageName]
		fmt.Printf("Execution de %s:\n", stageName)

		for _, command := range stageContent {
			toExec := strings.Fields(command)
			fmt.Println(toExec)
			cmd := exec.Command("cmd", append([]string{"/c", toExec[0]}, toExec[1:]...)...) //les ... permettent de traiter chaque élément de la liste individuellement
			output, err := cmd.CombinedOutput()

			if err != nil {
				fmt.Println("Erreur lors de l'exécution de la commande:", string(output))
				break
			}

			// Affichez la sortie de la commande
			fmt.Println("Sortie de la commande:")
			fmt.Println(string(output))
		}

		fmt.Println() // Ligne vide entre les étapes pour une meilleure lisibilité
	}
}

func sendPostRequest(jobID, logs, status string, date time.Time, stage string) {
	// Create the payload
	payload := JobResult{
		ID:     jobID,
		Logs:   logs,
		Status: status,
		Date:   date,
		Stage:  stage,
	}

	// Convert payload to JSON
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Make the POST request
	url := fmt.Sprintf("http://localhost:8080/v1/test") // Update with your actual endpoint
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		fmt.Println("Error sending POST request:", err)
		return
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Unexpected response status:", resp.Status)
	}
}
