package main

import (
	"bytes"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"os/exec"
	"strings"
)

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
		sendJobResult(jobID, stageName, logs, status)
		fmt.Println("Output complet:")
		fmt.Println(finalOutput)
	}
}
