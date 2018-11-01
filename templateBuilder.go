package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"text/template"
)

func buildTemplates(params Params) (*template.Template, error) {

	// merge templates
	templatesToMerge := getTemplates(params)

	log.Printf("Merging templates %v...", strings.Join(templatesToMerge, ", "))

	templateStrings := []string{}
	for _, t := range templatesToMerge {
		filePath := fmt.Sprintf("/templates/%v", t)
		data, err := ioutil.ReadFile(filePath)
		if err != nil {
			log.Fatal(fmt.Sprintf("Failed reading file %v: ", filePath), err)
		}
		templateStrings = append(templateStrings, string(data))
	}
	templateString := strings.Join(templateStrings, "\n---\n")

	// parse templates
	log.Printf("Parsing merged templates...")
	return template.New("kubernetes.yaml").Parse(templateString)
}

func getTemplates(params Params) []string {

	templatesToMerge := []string{
		"namespace.yaml",
		"service.yaml",
		"serviceaccount.yaml",
		"certificate-secret.yaml",
		"poddisruptionbudget.yaml",
		"horizontalpodautoscaler.yaml",
		"deployment.yaml",
	}

	if params.Visibility == "private" {
		templatesToMerge = append(templatesToMerge, "ingress.yaml")
	}
	if len(params.Secrets) > 0 {
		templatesToMerge = append(templatesToMerge, "application-secrets.yaml")
	}

	return templatesToMerge
}

func renderTemplate(tmpl *template.Template, templateData TemplateData) (bytes.Buffer, error) {

	// render templates
	log.Printf("Rendering merged templates...")
	var renderedTemplate bytes.Buffer
	err := tmpl.Execute(&renderedTemplate, templateData)
	if err != nil {
		return renderedTemplate, err
	}

	log.Printf("Template after rendering:\n\n")
	log.Println(renderedTemplate.String())
	log.Println("")

	return renderedTemplate, err
}
