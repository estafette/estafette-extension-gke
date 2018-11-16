package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTemplates(t *testing.T) {

	t.Run("IncludesIngressIfVisibilityIsPrivate", func(t *testing.T) {

		params := Params{
			Visibility: "private",
		}

		// act
		templates := getTemplates(params)

		assert.True(t, stringArrayContains(templates, "/templates/ingress.yaml"))
	})

	t.Run("IncludesIngressIfVisibilityIsIap", func(t *testing.T) {

		params := Params{
			Visibility: "iap",
		}

		// act
		templates := getTemplates(params)

		assert.True(t, stringArrayContains(templates, "/templates/ingress.yaml"))
	})

	t.Run("DoesNotIncludeIngressIfVisibilityIsPublic", func(t *testing.T) {

		params := Params{
			Visibility: "public",
		}

		// act
		templates := getTemplates(params)

		assert.False(t, stringArrayContains(templates, "/templates/ingress.yaml"))
	})

	t.Run("IncludesApplicationSecretsIfLengthOfSecretsIsMoreThanZero", func(t *testing.T) {

		params := Params{
			Secrets: SecretsParams{
				Keys: map[string]string{
					"secret-file-1.json": "c29tZSBzZWNyZXQgdmFsdWU=",
					"secret-file-2.yaml": "YW5vdGhlciBzZWNyZXQgdmFsdWU=",
				},
			},
		}

		// act
		templates := getTemplates(params)

		assert.True(t, stringArrayContains(templates, "/templates/application-secrets.yaml"))
	})

	t.Run("DoesNotIncludeApplicationSecretsIfLengthOfSecretsZero", func(t *testing.T) {

		params := Params{}

		// act
		templates := getTemplates(params)

		assert.False(t, stringArrayContains(templates, "/templates/application-secrets.yaml"))
	})

	t.Run("AddLocalManifestsIfSetInLocalManifestsParam", func(t *testing.T) {

		params := Params{
			Manifests: ManifestsParams{
				Files: []string{
					"./gke/another-ingress.yaml",
				},
			},
		}

		// act
		templates := getTemplates(params)

		assert.True(t, stringArrayContains(templates, "./gke/another-ingress.yaml"))
	})

	t.Run("OverrideWithLocalManifestsIfSetInLocalManifestsParamWithSameFilename", func(t *testing.T) {

		params := Params{
			Manifests: ManifestsParams{
				Files: []string{
					"./gke/service.yaml",
				},
			},
		}

		// act
		templates := getTemplates(params)

		assert.True(t, stringArrayContains(templates, "./gke/service.yaml"))
		assert.False(t, stringArrayContains(templates, "/templates/service.yaml"))
	})

	t.Run("ReturnsOnlyHorizontalPodAutoscalerIfTypeIsRollback", func(t *testing.T) {

		params := Params{
			Type: "rollback",
		}

		// act
		templates := getTemplates(params)

		assert.Equal(t, 1, len(templates))
		assert.True(t, stringArrayContains(templates, "/templates/horizontalpodautoscaler.yaml"))
	})

	t.Run("ReturnsOnlyHorizontalPodAutoscalerAndPodDisruptionBudgetIfTypeIsCanary", func(t *testing.T) {

		params := Params{
			Type: "canary",
		}

		// act
		templates := getTemplates(params)

		assert.False(t, stringArrayContains(templates, "/templates/horizontalpodautoscaler.yaml"))
		assert.False(t, stringArrayContains(templates, "/templates/poddisruptionbudget.yaml"))
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
