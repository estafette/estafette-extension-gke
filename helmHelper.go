package main

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

func generateValues(templateData TemplateData) error {
	bytes, err := yaml.Marshal(templateData)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile("/helm-chart/values.yaml", bytes, 0600)
	if err != nil {
		return err
	}

	return nil
}