package main

import (
	"io/ioutil"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

func generateValues(templateData TemplateData, dir string) error {
	bytes, err := yaml.Marshal(templateData)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filepath.Join(dir, "values.yaml"), bytes, 0600)
	if err != nil {
		return err
	}

	return nil
}