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

	if len(templatesToMerge) == 0 {
		return nil, nil
	}

	logInfo("Merging templates %v...", strings.Join(templatesToMerge, ", "))

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
	logInfo("Parsing merged templates...")
	return template.New("kubernetes.yaml").Funcs(sprig.TxtFuncMap()).Parse(templateString)
}

func getTemplates(params Params) []string {

	if params.Action == "rollback-canary" {
		return []string{}
	}

	templatesToMerge := []string{}

	switch params.Kind {
	case "job":
		templatesToMerge = append(templatesToMerge, []string{
			"namespace.yaml",
			"serviceaccount.yaml",
			"job.yaml",
		}...)

	default:
		templatesToMerge = append(templatesToMerge, []string{
			"namespace.yaml",
			"service.yaml",
			"serviceaccount.yaml",
			"certificate-secret.yaml",
			"deployment.yaml",
		}...)

	}

	if params.Action == "deploy-simple" || params.Action == "deploy-stable" {
		templatesToMerge = append(templatesToMerge, []string{
			"poddisruptionbudget.yaml",
			"horizontalpodautoscaler.yaml",
		}...)
	}

	if params.Kind == "deployment" && (params.Visibility == "private" || params.Visibility == "iap" || params.Visibility == "public-whitelist") {
		templatesToMerge = append(templatesToMerge, "ingress.yaml")
	}
	if params.Kind == "deployment" && len(params.InternalHosts) > 0 {
		templatesToMerge = append(templatesToMerge, "ingress-internal.yaml")
	}
	if len(params.Secrets.Keys) > 0 {
		templatesToMerge = append(templatesToMerge, "application-secrets.yaml")
	}
	if len(params.Configs.Files) > 0 || len(params.Configs.InlineFiles) > 0 {
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

	if params.Action != "rollback-canary" && (len(params.Configs.Files) > 0 || len(params.Configs.InlineFiles) > 0) {
		logInfo("Prerendering config files...")

		// render files passed with configs.files property, replacing placeholders with values specified in configs.data property
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

		// add files passed with configs.inline property as is
		for filename, content := range params.Configs.InlineFiles {
			renderedConfigFiles[filename] = content
		}
	}

	return
}

func renderTemplate(tmpl *template.Template, templateData TemplateData) (bytes.Buffer, error) {

	if tmpl == nil {
		return bytes.Buffer{}, nil
	}

	// render templates
	logInfo("Rendering merged templates...")
	var renderedTemplate bytes.Buffer
	err := tmpl.Execute(&renderedTemplate, templateData)
	if err != nil {
		return renderedTemplate, err
	}

	logInfo("Template after rendering:")
	log.Println(renderedTemplate.String())
	log.Println("")

	return renderedTemplate, err
}
