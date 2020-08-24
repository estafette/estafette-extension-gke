package main

import (
	"testing"
)

func TestGenerateValues(t *testing.T) {

	t.Run("WritesValuesYamlFile", func(t *testing.T) {

		templateData := TemplateData{}

		// act
		generateValues(templateData, ".")
	})
}
