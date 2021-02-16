package builder

import (
	bytes "bytes"
	"context"
	"strings"
	"testing"
	template "text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/estafette/estafette-extension-gke/api"
	"github.com/stretchr/testify/assert"
)

func TestGetTemplates(t *testing.T) {

	t.Run("IncludesIngressIfVisibilityIsPrivateAndKindIsDeployment", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Action:     api.ActionDeploySimple,
			Visibility: api.VisibilityPrivate,
			Kind:       api.KindDeployment,
		}

		// act
		templates := service.GetTemplates(params, true)

		assert.True(t, stringArrayContains(templates, "/templates/ingress.yaml"))
	})

	t.Run("IncludesIngressIfVisibilityIsIapAndKindIsDeployment", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Action:     api.ActionDeploySimple,
			Visibility: api.VisibilityIAP,
			Kind:       api.KindDeployment,
		}

		// act
		templates := service.GetTemplates(params, true)

		assert.True(t, stringArrayContains(templates, "/templates/ingress.yaml"))
	})

	t.Run("IncludesIngressIfVisibilityIsPublicWhitelistAndKindIsDeployment", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Action:     api.ActionDeploySimple,
			Visibility: api.VisibilityPublicWhitelist,
			Kind:       api.KindDeployment,
		}

		// act
		templates := service.GetTemplates(params, true)

		assert.True(t, stringArrayContains(templates, "/templates/ingress.yaml"))
	})

	t.Run("DoesNotIncludeIngressIfVisibilityIsPublic", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Action:     api.ActionDeploySimple,
			Visibility: api.VisibilityPublic,
		}

		// act
		templates := service.GetTemplates(params, true)

		assert.False(t, stringArrayContains(templates, "/templates/ingress.yaml"))
	})

	t.Run("IncludesInternalIngressIfOneOrMoreInternalHostsAreSetAndKindIsDeployment", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Action:        api.ActionDeploySimple,
			Kind:          api.KindDeployment,
			InternalHosts: []string{"ci.estafette.internal"},
		}

		// act
		templates := service.GetTemplates(params, true)

		assert.True(t, stringArrayContains(templates, "/templates/ingress-internal.yaml"))
	})

	t.Run("DoesNotIncludeInternalIngressIfNoInternalHostsAreSet", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Action:        api.ActionDeploySimple,
			InternalHosts: []string{},
		}

		// act
		templates := service.GetTemplates(params, true)

		assert.False(t, stringArrayContains(templates, "/templates/ingress-internal.yaml"))
	})

	t.Run("IncludesApplicationSecretsIfLengthOfSecretsIsMoreThanZero", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Action: api.ActionDeploySimple,
			Secrets: api.SecretsParams{
				Keys: map[string]interface{}{
					"secret-file-1.json": "c29tZSBzZWNyZXQgdmFsdWU=",
					"secret-file-2.yaml": "YW5vdGhlciBzZWNyZXQgdmFsdWU=",
				},
			},
		}

		// act
		templates := service.GetTemplates(params, true)

		assert.True(t, stringArrayContains(templates, "/templates/application-secrets.yaml"))
	})

	t.Run("DoesNotIncludeApplicationSecretsIfLengthOfSecretsZero", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Action: api.ActionDeploySimple,
		}

		// act
		templates := service.GetTemplates(params, true)

		assert.False(t, stringArrayContains(templates, "/templates/application-secrets.yaml"))
	})

	t.Run("AddLocalManifestsIfSetInLocalManifestsParam", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Action: api.ActionDeploySimple,
			Manifests: api.ManifestsParams{
				Files: []string{
					"./gke/another-ingress.yaml",
				},
			},
		}

		// act
		templates := service.GetTemplates(params, true)

		assert.True(t, stringArrayContains(templates, "./gke/another-ingress.yaml"))
	})

	t.Run("OverrideWithLocalManifestsIfSetInLocalManifestsParamWithSameFilename", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Action: api.ActionDeploySimple,
			Manifests: api.ManifestsParams{
				Files: []string{
					"./gke/service.yaml",
				},
			},
		}

		// act
		templates := service.GetTemplates(params, true)

		assert.True(t, stringArrayContains(templates, "./gke/service.yaml"))
		assert.False(t, stringArrayContains(templates, "/templates/service.yaml"))
	})

	t.Run("ReturnsEmptyListIfActionIsRollbackCanaray", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Action: api.ActionRollbackCanary,
		}

		// act
		templates := service.GetTemplates(params, true)

		assert.Equal(t, 0, len(templates))
	})

	t.Run("ReturnsOnlyHorizontalPodAutoscalerAndPodDisruptionBudgetIfActionIsDeployCanary", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Action: api.ActionDeployCanary,
		}

		// act
		templates := service.GetTemplates(params, true)

		assert.False(t, stringArrayContains(templates, "/templates/horizontalpodautoscaler.yaml"))
		assert.False(t, stringArrayContains(templates, "/templates/poddisruptionbudget.yaml"))
	})

	t.Run("DoesNotIncludeCertificateSecretIfCertificateSecretIsSet", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Action:            api.ActionDeploySimple,
			Kind:              api.KindDeployment,
			CertificateSecret: "shared-wildcard-letsencrypt-certificate",
		}

		// act
		templates := service.GetTemplates(params, true)

		assert.False(t, stringArrayContains(templates, "/templates/certificate-secret.yaml"))
	})

	t.Run("IncludesCertificateSecretIfCertificateSecretIsNotSet", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Action: api.ActionDeploySimple,
			Kind:   api.KindDeployment,
		}

		// act
		templates := service.GetTemplates(params, true)

		assert.True(t, stringArrayContains(templates, "/templates/certificate-secret.yaml"))
	})

	t.Run("IncludesApigeeIngressIfVisibilityIsApigeeAndKindIsDeployment", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Action:     api.ActionDeploySimple,
			Visibility: api.VisibilityApigee,
			Kind:       api.KindDeployment,
		}

		// act
		templates := service.GetTemplates(params, true)

		assert.True(t, stringArrayContains(templates, "/templates/ingress-apigee.yaml"))
		assert.True(t, stringArrayContains(templates, "/templates/ingress.yaml"))
	})
}

func TestInjectSteps(t *testing.T) {

	t.Run("RenderNamespace", func(t *testing.T) {

		data := api.TemplateData{
			Namespace: "mynamespace",
		}
		tmpl, err := template.ParseFiles("../../templates/namespace.yaml")

		// act
		var renderedTemplate bytes.Buffer
		err = tmpl.Execute(&renderedTemplate, data)

		assert.Nil(t, err)
		assert.Equal(t, "apiVersion: v1\nkind: Namespace\nmetadata:\n  name: mynamespace", renderedTemplate.String())
		assert.True(t, strings.Contains(renderedTemplate.String(), "mynamespace"))
	})

	t.Run("RenderServiceAccount", func(t *testing.T) {

		data := api.TemplateData{
			Name:      "myapp",
			Namespace: "mynamespace",
			Labels: map[string]string{
				"app":  "myapp",
				"team": "myteam",
			},
		}
		tmpl, err := template.New("serviceaccount.yaml").Funcs(sprig.TxtFuncMap()).ParseFiles("../../templates/serviceaccount.yaml")
		assert.Nil(t, err)

		// act
		var renderedTemplate bytes.Buffer
		err = tmpl.Execute(&renderedTemplate, data)

		assert.Nil(t, err)
		assert.Equal(t, "apiVersion: v1\nkind: ServiceAccount\nmetadata:\n  name: myapp\n  namespace: mynamespace\n  labels:\n    \"app\": \"myapp\"\n    \"team\": \"myteam\"", renderedTemplate.String())
		assert.True(t, strings.Contains(renderedTemplate.String(), "mynamespace"))
	})

	t.Run("RenderHorizontalPodAutoscaler", func(t *testing.T) {

		data := api.TemplateData{
			Name:          "myapp",
			NameWithTrack: "myapp-canary",
			Namespace:     "mynamespace",
			Labels: map[string]string{
				"app":  "myapp",
				"team": "myteam",
			},
			MinReplicas:         3,
			MaxReplicas:         19,
			TargetCPUPercentage: 65,
		}
		tmpl, err := template.New("horizontalpodautoscaler.yaml").Funcs(sprig.TxtFuncMap()).ParseFiles("../../templates/horizontalpodautoscaler.yaml")
		assert.Nil(t, err)

		// act
		var renderedTemplate bytes.Buffer
		err = tmpl.Execute(&renderedTemplate, data)

		assert.Nil(t, err)
		assert.Equal(t, "apiVersion: autoscaling/v1\nkind: HorizontalPodAutoscaler\nmetadata:\n  name: myapp-canary\n  namespace: mynamespace\n  labels:\n    \"app\": \"myapp\"\n    \"team\": \"myteam\"\nspec:\n  scaleTargetRef:\n    apiVersion: apps/v1\n    kind: Deployment\n    name: myapp-canary\n  minReplicas: 3\n  maxReplicas: 19\n  targetCPUUtilizationPercentage: 65", renderedTemplate.String())
		assert.True(t, strings.Contains(renderedTemplate.String(), "mynamespace"))
	})
}

func stringArrayContains(array []string, search string) bool {
	for _, v := range array {
		if v == search {
			return true
		}
	}
	return false
}
