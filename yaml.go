package main

import (
	"bytes"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
)

type Pipeline struct {
	StagesList []string            `yaml:"stages"`
	Stages     map[string][]string `yaml:",inline"`
}

func (p *Pipeline) Unmarshal(data []byte) error {
	return yaml.Unmarshal(data, &p)
}

func readYaml(path string) (Pipeline, error) {
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
	tag = strings.ReplaceAll(tag, "/", "-")
	pipeline, err := readYaml(path)
	if errorAndFinish(err, pipelineId, "Setup", "Erreur lors de la récupération du fichier de CI") {
		return
	}

	var fullOutput bytes.Buffer

	for _, stageName := range pipeline.StagesList {
		fmt.Println("pipeline:", stageName)
		stageContent := pipeline.Stages[stageName]

		for _, command := range stageContent {
			command = strings.ReplaceAll(command, "{tag}", tag)
			fmt.Println("command:", command)
			args := strings.Fields(command)
			output := execCmd(pipelineId, stageName, "Erreur lors de l'exécution de la commande", args...)
			if output == nil {
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
