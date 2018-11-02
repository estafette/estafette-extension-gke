package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig"
)

func buildTemplates(params Params) (*template.Template, error) {

	// merge templates
	templatesToMerge := getTemplates(params)

	log.Printf("Merging templates %v...", strings.Join(templatesToMerge, ", "))

	templateStrings := []string{}
	for _, t := range templatesToMerge {
		data, err := ioutil.ReadFile(t)
		if err != nil {
			log.Fatal(fmt.Sprintf("Failed reading file %v: ", t), err)
		}
		templateStrings = append(templateStrings, string(data))
	}
	templateString := strings.Join(templateStrings, "\n---\n")

	// parse templates
	log.Printf("Parsing merged templates...")
	return template.New("kubernetes.yaml").Funcs(sprig.TxtFuncMap()).Parse(templateString)
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

	if params.Visibility == "private" || params.Visibility == "iap" {
		templatesToMerge = append(templatesToMerge, "ingress.yaml")
	}
	if len(params.Secrets) > 0 {
		templatesToMerge = append(templatesToMerge, "application-secrets.yaml")
	}
	if len(params.ConfigFiles) > 0 {
		templatesToMerge = append(templatesToMerge, "configmap.yaml")
	}

	// prefix all filenames with templates dir
	for i, t := range templatesToMerge {
		templatesToMerge[i] = fmt.Sprintf("/templates/%v", t)
	}

	// add or override with local manifests
	for _, lm := range params.LocalManifests {
		filename := filepath.Base(lm)

		overridesExistingTemplate := false
		for i, t := range templatesToMerge {
			if filename == filepath.Base(t) {
				overridesExistingTemplate = true
				templatesToMerge[i] = lm
				break
			}
		}

		if !overridesExistingTemplate {
			templatesToMerge = append(templatesToMerge, lm)
		}
	}

	return templatesToMerge
}

func renderConfig(params Params) (renderedConfigFiles []ConfigFileParams) {

	if len(params.ConfigFiles) > 0 {
		log.Printf("Prerendering config files...")

		for _, cf := range params.ConfigFiles {

			data, err := ioutil.ReadFile(cf.File)
			if err != nil {
				log.Fatal(fmt.Sprintf("Failed reading file %v: ", cf.File), err)
			}
			tmpl, err := template.New(cf.File).Parse(string(data))
			if err != nil {
				log.Fatal(fmt.Sprintf("Failed building template from file %v: ", cf.File), err)
			}

			var renderedTemplate bytes.Buffer
			err = tmpl.Execute(&renderedTemplate, cf.Data)
			if err != nil {
				log.Fatal(fmt.Sprintf("Failed rendering template from file %v: ", cf.File), err)
			}

			cf.RenderedFileContent = renderedTemplate.String()

			renderedConfigFiles = append(renderedConfigFiles, cf)
		}
	}

	return
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
