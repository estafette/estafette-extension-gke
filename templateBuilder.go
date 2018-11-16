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
			log.Fatal(fmt.Sprintf("Failed reading file %v. Do you have a git-clone stage before running this extension? For releases git-clone is not automatically handled to save time in case it's not needed. ", t), err)
		}
		templateStrings = append(templateStrings, string(data))
	}
	templateString := strings.Join(templateStrings, "\n---\n")

	// parse templates
	log.Printf("Parsing merged templates...")
	return template.New("kubernetes.yaml").Funcs(sprig.TxtFuncMap()).Parse(templateString)
}

func getTemplates(params Params) []string {

	if params.Type == "rollback" {
		return []string{
			"/templates/horizontalpodautoscaler.yaml",
		}
	}

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
	if len(params.Secrets.Keys) > 0 {
		templatesToMerge = append(templatesToMerge, "application-secrets.yaml")
	}
	if len(params.Configs.Files) > 0 {
		templatesToMerge = append(templatesToMerge, "configmap.yaml")
	}

	// prefix all filenames with templates dir
	for i, t := range templatesToMerge {
		templatesToMerge[i] = fmt.Sprintf("/templates/%v", t)
	}

	// add or override with local manifests
	for _, lm := range params.Manifests.Files {
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

func renderConfig(params Params) (renderedConfigFiles map[string]string) {

	renderedConfigFiles = map[string]string{}
	if len(params.Configs.Files) > 0 {
		log.Printf("Prerendering config files...")

		for _, cf := range params.Configs.Files {

			data, err := ioutil.ReadFile(cf)
			if err != nil {
				log.Fatal(fmt.Sprintf("Failed reading file %v. Do you have a git-clone stage before running this extension? For releases git-clone is not automatically handled to save time in case it's not needed. ", cf), err)
			}
			tmpl, err := template.New(cf).Parse(string(data))
			if err != nil {
				log.Fatal(fmt.Sprintf("Failed building template from file %v: ", cf), err)
			}

			var renderedTemplate bytes.Buffer
			err = tmpl.Execute(&renderedTemplate, params.Configs.Data)
			if err != nil {
				log.Fatal(fmt.Sprintf("Failed rendering template from file %v: ", cf), err)
			}

			renderedConfigFiles[filepath.Base(cf)] = renderedTemplate.String()
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
