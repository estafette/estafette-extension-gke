package builder

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/estafette/estafette-extension-gke/api"
	"github.com/rs/zerolog/log"
)

//go:generate mockgen -package=builder -destination ./mock.go -source=service.go
type Service interface {
	BuildTemplates(params api.Params, includePodDisruptionBudget bool) (*template.Template, error)
	GetTemplates(params api.Params, includePodDisruptionBudget bool) []string
	GetAtomicUpdateServiceTemplate() (*template.Template, error)
	RenderConfig(params api.Params) (renderedConfigFiles map[string]string)
	RenderTemplate(tmpl *template.Template, templateData api.TemplateData, logTemplate bool) (bytes.Buffer, error)
}

// NewService returns a new extension.Service
func NewService(ctx context.Context) (Service, error) {
	return &service{}, nil
}

type service struct {
}

func (s *service) BuildTemplates(params api.Params, includePodDisruptionBudget bool) (*template.Template, error) {

	// merge templates
	templatesToMerge := s.GetTemplates(params, includePodDisruptionBudget)

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

func (s *service) GetTemplates(params api.Params, includePodDisruptionBudget bool) []string {

	if params.Action == api.ActionRollbackCanary || params.Action == api.ActionUnknown || params.Action == api.ActionRestartCanary || params.Action == api.ActionRestartStable || params.Action == api.ActionRestartSimple {
		return []string{}
	}

	templatesToMerge := []string{}

	switch params.Kind {
	case api.KindConfig:
		templatesToMerge = append(templatesToMerge, []string{
			"namespace.yaml",
		}...)
	case api.KindJob:
		templatesToMerge = append(templatesToMerge, []string{
			"namespace.yaml",
			"serviceaccount.yaml",
			"job.yaml",
		}...)

	case api.KindCronJob:
		templatesToMerge = append(templatesToMerge, []string{
			"namespace.yaml",
			"serviceaccount.yaml",
			"cronjob.yaml",
		}...)

	case api.KindStatefulset:
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

	case api.KindDeployment:
		templatesToMerge = append(templatesToMerge, []string{
			"namespace.yaml",
			"serviceaccount.yaml",
			"deployment.yaml",
		}...)

		if params.StrategyType != api.StrategyTypeAtomicUpdate {
			templatesToMerge = append(templatesToMerge, "service.yaml")
		}

		if params.CertificateSecret == "" {
			templatesToMerge = append(templatesToMerge, "certificate-secret.yaml")
		}

	case api.KindHeadlessDeployment:
		templatesToMerge = append(templatesToMerge, []string{
			"namespace.yaml",
			"serviceaccount.yaml",
			"deployment.yaml",
		}...)
	}

	hasImagePullSecret := params.ImagePullSecretUser != "" && params.ImagePullSecretPassword != ""

	if hasImagePullSecret && params.Kind != api.KindConfig && params.Kind != api.KindConfigToFile {
		templatesToMerge = append(templatesToMerge, []string{
			"image-pull-secret.yaml",
		}...)
	}

	if includePodDisruptionBudget && (params.Kind == api.KindDeployment || params.Kind == api.KindHeadlessDeployment || params.Kind == api.KindStatefulset) && (params.Action == api.ActionDeploySimple || params.Action == api.ActionDeployStable || params.Action == api.ActionDiffSimple || params.Action == api.ActionDiffStable) {
		templatesToMerge = append(templatesToMerge, "poddisruptionbudget.yaml")
	}
	if (params.Kind == api.KindDeployment || params.Kind == api.KindHeadlessDeployment) && params.Autoscale.Enabled != nil && *params.Autoscale.Enabled && params.StrategyType != "Recreate" && (params.Action == api.ActionDeploySimple || params.Action == api.ActionDeployStable || params.Action == api.ActionDiffSimple || params.Action == api.ActionDiffStable) {
		templatesToMerge = append(templatesToMerge, "horizontalpodautoscaler.yaml")
	}
	if (params.Kind == api.KindDeployment || params.Kind == api.KindHeadlessDeployment) && params.VerticalPodAutoscaler.Enabled != nil && *params.VerticalPodAutoscaler.Enabled && (params.Action == api.ActionDeploySimple || params.Action == api.ActionDeployStable || params.Action == api.ActionDiffSimple || params.Action == api.ActionDiffStable) {
		templatesToMerge = append(templatesToMerge, "verticalpodautoscaler.yaml")
	}
	if (params.Kind == api.KindDeployment || params.Kind == api.KindStatefulset) && (params.Visibility == api.VisibilityPrivate || params.Visibility == api.VisibilityIAP || params.Visibility == api.VisibilityPublicWhitelist) {
		templatesToMerge = append(templatesToMerge, "ingress.yaml")
	}

	if params.Kind == api.KindDeployment && params.Visibility == api.VisibilityApigee {
		templatesToMerge = append(templatesToMerge, "ingress-apigee.yaml")
		templatesToMerge = append(templatesToMerge, "ingress.yaml")
	}
	if params.Kind == api.KindDeployment && (params.Visibility == api.VisibilityESP || params.Visibility == api.VisibilityESPv2) {
		templatesToMerge = append(templatesToMerge, "ingress-esp.yaml")
		//templatesToMerge = append(templatesToMerge, "service-clusterip.yaml")
	}

	if (params.Kind == api.KindDeployment || params.Kind == api.KindStatefulset) && params.Visibility == api.VisibilityIAP {
		templatesToMerge = append(templatesToMerge, "backend-config.yaml", "iap-oauth-credentials-secret.yaml")
	}
	if (params.Kind == api.KindDeployment || params.Kind == api.KindStatefulset) && len(params.InternalHosts) > 0 {
		templatesToMerge = append(templatesToMerge, "ingress-internal.yaml")
	}
	if params.HasSecrets() {
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

func (s *service) GetAtomicUpdateServiceTemplate() (*template.Template, error) {

	// parse service template
	return template.New("service.yaml").Funcs(sprig.TxtFuncMap()).ParseFiles("/templates/service.yaml")
}

func (s *service) RenderConfig(params api.Params) (renderedConfigFiles map[string]string) {

	renderedConfigFiles = map[string]string{}

	if params.Action != api.ActionRollbackCanary && (len(params.Configs.Files) > 0 || len(params.Configs.InlineFiles) > 0) {
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

func (s *service) RenderTemplate(tmpl *template.Template, templateData api.TemplateData, logTemplate bool) (bytes.Buffer, error) {

	if tmpl == nil {
		return bytes.Buffer{}, nil
	}

	// render templates
	if logTemplate {
		log.Info().Msg("Rendering merged templates...")
	}
	var renderedTemplate bytes.Buffer
	err := tmpl.Execute(&renderedTemplate, templateData)
	if err != nil {
		return renderedTemplate, err
	}

	if logTemplate {
		log.Info().Msg("Template after rendering:")
		log.Info().Msg(renderedTemplate.String())
		log.Info().Msg("")
	}

	return renderedTemplate, err
}
