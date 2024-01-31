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
		return Pipeline{}, err
	}

	var pipeline Pipeline
	err = pipeline.Unmarshal(yamlFile)

	if err != nil {
		return Pipeline{}, err
	}

	return pipeline, err
}

func stagesExecution(path string, pipelineId string, projectName string, tag string) {
	pipeline, err := readYaml(path)
	if errorAndFinish(err, pipelineId, "Setup", "Erreur lors de la récupération du fichier de CI") {
		return
	}

	var fullOutput bytes.Buffer
	i := 0

	for _, stageName := range pipeline.StagesList {
		stageContent := pipeline.Stages[stageName]

		for _, command := range stageContent {
			if strings.Contains(command, "docker push") || strings.Contains(command, "docker image push") {
				command = strings.ReplaceAll(command, "{tag	}", tag)
			}
			i++
			fmt.Println("pipeline:", i)
			cmd := exec.Command(command)
			output, err := cmd.CombinedOutput()
			if errorAndFinish(err, pipelineId, stageName, "Erreur lors de l'exécution du stage") {
				return
			}

			fullOutput.WriteString(fmt.Sprintf("%s\n %s\n", command, output))
		}

		finalOutput := fullOutput.String()
		logs := fmt.Sprintf(finalOutput)
		sendJobResult(pipelineId, stageName, logs, "Success")
		fmt.Println("Output complet:")
		fmt.Println(finalOutput)
		fullOutput.Reset()
	}
	cd(pipelineId, projectName, tag)
}
