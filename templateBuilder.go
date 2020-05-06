package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/rs/zerolog/log"
)

func buildTemplates(params Params, includePodDisruptionBudget bool) (*template.Template, error) {

	// merge templates
	templatesToMerge := getTemplates(params, includePodDisruptionBudget)

	if len(templatesToMerge) == 0 {
		return nil, nil
	}

	log.Info().Msgf("Merging templates %v...", strings.Join(templatesToMerge, ", "))

	templateStrings := []string{}
	for _, t := range templatesToMerge {
		data, err := ioutil.ReadFile(t)
		if err != nil {
			log.Fatal().Err(err).Msgf("Failed reading file %v. Do you have a git-clone stage before running this extension? For releases git-clone is not automatically handled to save time in case it's not needed. ", t)
		}
		templateStrings = append(templateStrings, string(data))
	}
	templateString := strings.Join(templateStrings, "\n---\n")

	// parse templates
	log.Info().Msg("Parsing merged templates...")
	return template.New("kubernetes.yaml").Funcs(sprig.TxtFuncMap()).Parse(templateString)
}

func getTemplates(params Params, includePodDisruptionBudget bool) []string {

	if params.Action == "rollback-canary" {
		return []string{}
	}

	templatesToMerge := []string{}

	switch params.Kind {
	case "config":
		templatesToMerge = append(templatesToMerge, []string{
			"namespace.yaml",
		}...)
	case "job":
		templatesToMerge = append(templatesToMerge, []string{
			"namespace.yaml",
			"serviceaccount.yaml",
			"job.yaml",
		}...)

	case "cronjob":
		templatesToMerge = append(templatesToMerge, []string{
			"namespace.yaml",
			"serviceaccount.yaml",
			"cronjob.yaml",
		}...)

	case "statefulset":
		templatesToMerge = append(templatesToMerge, []string{
			"namespace.yaml",
			"service.yaml",
			"service-headless.yaml",
			"serviceaccount.yaml",
			"statefulset.yaml",
		}...)
		if params.CertificateSecret == "" {
			templatesToMerge = append(templatesToMerge, "certificate-secret.yaml")
		}

	case "deployment":
		templatesToMerge = append(templatesToMerge, []string{
			"namespace.yaml",
			"service.yaml",
			"serviceaccount.yaml",
			"deployment.yaml",
		}...)
		if params.CertificateSecret == "" {
			templatesToMerge = append(templatesToMerge, "certificate-secret.yaml")
		}

	case "headless-deployment":
		templatesToMerge = append(templatesToMerge, []string{
			"namespace.yaml",
			"serviceaccount.yaml",
			"deployment.yaml",
		}...)
	}

	if includePodDisruptionBudget && (params.Kind == "deployment" || params.Kind == "headless-deployment" || params.Kind == "statefulset") && (params.Action == "deploy-simple" || params.Action == "deploy-stable" || params.Action == "diff-simple" || params.Action == "diff-canary" || params.Action == "diff-stable") {
		templatesToMerge = append(templatesToMerge, "poddisruptionbudget.yaml")
	}
	if (params.Kind == "deployment" || params.Kind == "headless-deployment") && params.Autoscale.Enabled != nil && *params.Autoscale.Enabled && (params.Action == "deploy-simple" || params.Action == "deploy-stable" || params.Action == "diff-simple" || params.Action == "diff-canary" || params.Action == "diff-stable") {
		templatesToMerge = append(templatesToMerge, "horizontalpodautoscaler.yaml")
	}
	if (params.Kind == "deployment" || params.Kind == "statefulset") && (params.Visibility == "private" || params.Visibility == "iap" || params.Visibility == "public-whitelist") {
		templatesToMerge = append(templatesToMerge, "ingress.yaml")
	}

	if params.Kind == "deployment" && params.Visibility == "apigee" {
		templatesToMerge = append(templatesToMerge, "ingress-apigee.yaml")
		templatesToMerge = append(templatesToMerge, "ingress.yaml")
	}

	if (params.Kind == "deployment" || params.Kind == "statefulset") && params.Visibility == "iap" {
		templatesToMerge = append(templatesToMerge, "backend-config.yaml", "iap-oauth-credentials-secret.yaml")
	}
	if (params.Kind == "deployment" || params.Kind == "statefulset") && len(params.InternalHosts) > 0 {
		templatesToMerge = append(templatesToMerge, "ingress-internal.yaml")
	}
	if len(params.Secrets.Keys) > 0 {
		templatesToMerge = append(templatesToMerge, "application-secrets.yaml")
	}
	if params.UseGoogleCloudCredentials || params.LegacyGoogleCloudServiceAccountKeyFile != "" {
		templatesToMerge = append(templatesToMerge, "service-account-secret.yaml")
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
		log.Info().Msg("Prerendering config files...")

		// render files passed with configs.files property, replacing placeholders with values specified in configs.data property
		for _, cf := range params.Configs.Files {

			data, err := ioutil.ReadFile(cf)
			if err != nil {
				log.Fatal().Err(err).Msgf("Failed reading file %v. Do you have a git-clone stage before running this extension? For releases git-clone is not automatically handled to save time in case it's not needed. ", cf)
			}
			tmpl, err := template.New(cf).Parse(string(data))
			if err != nil {
				log.Fatal().Err(err).Msgf("Failed building template from file %v: ", cf)
			}

			var renderedTemplate bytes.Buffer
			err = tmpl.Execute(&renderedTemplate, params.Configs.Data)
			if err != nil {
				log.Fatal().Err(err).Msgf("Failed rendering template from file %v: ", cf)
			}

			renderedConfigFiles[filepath.Base(cf)] = renderedTemplate.String()
		}

		// add files passed with configs.inline property, replacing placeholders with values specified in configs.data property
		for filename, content := range params.Configs.InlineFiles {
			tmpl, err := template.New(filename).Parse(content)
			if err != nil {
				log.Fatal().Err(err).Msgf("Failed building template from inline file %v: ", filename)
			}
			var renderedTemplate bytes.Buffer
			err = tmpl.Execute(&renderedTemplate, params.Configs.Data)
			if err != nil {
				log.Fatal().Err(err).Msgf("Failed rendering template from file %v: ", filename)
			}

			renderedConfigFiles[filename] = renderedTemplate.String()
		}
	}

	return
}

func renderTemplate(tmpl *template.Template, templateData TemplateData) (bytes.Buffer, error) {

	if tmpl == nil {
		return bytes.Buffer{}, nil
	}

	// render templates
	log.Info().Msg("Rendering merged templates...")
	var renderedTemplate bytes.Buffer
	err := tmpl.Execute(&renderedTemplate, templateData)
	if err != nil {
		return renderedTemplate, err
	}

	log.Info().Msg("Template after rendering:")
	log.Info().Msg(renderedTemplate.String())
	log.Info().Msg("")

	return renderedTemplate, err
}
