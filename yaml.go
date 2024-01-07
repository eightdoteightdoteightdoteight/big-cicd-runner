package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
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
