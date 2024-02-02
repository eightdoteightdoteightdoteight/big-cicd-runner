package main

import (
	"bytes"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"regexp"
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
	fmt.Println("pipeline:", pipelineId)
	tag = strings.ReplaceAll(tag, "/", "-")
	pipeline, err := readYaml(path)
	if errorAndFinish(err, pipelineId, "Setup", "Erreur lors de la récupération du fichier de CI") {
		return
	}

	var fullOutput bytes.Buffer

	for _, stageName := range pipeline.StagesList {
		fmt.Println("job:", stageName)
		stageContent := pipeline.Stages[stageName]

		os.Setenv("TAG", tag)
		os.Setenv("SONAR_PROJECT_KEY", projectName)
		os.Setenv("SONAR_PROJECT_VERSION", tag)
		os.Setenv("SONAR_PROJECT_NAME", projectName)

		for _, command := range stageContent {
			re, err := regexp.Compile("({([^{}]+)})")
			if errorAndFinish(err, pipelineId, stageName, "Erreur lors de la préparation de la commande") {
				return
			}
			matches := re.FindAllStringSubmatch(command, -1)
			for _, match := range matches {
				value := os.Getenv(match[2])
				if value == "" {
					continue
				}
				command = strings.ReplaceAll(command, match[1], value)
			}
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
		fullOutput.Reset()
	}
	cd(pipelineId, projectName, tag)
}
