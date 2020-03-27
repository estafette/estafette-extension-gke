package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTemplates(t *testing.T) {

	t.Run("IncludesIngressIfVisibilityIsPrivateAndKindIsDeployment", func(t *testing.T) {

		params := Params{
			Visibility: "private",
			Kind:       "deployment",
		}

		// act
		templates := getTemplates(params)

		assert.True(t, stringArrayContains(templates, "/templates/ingress.yaml"))
	})

	t.Run("IncludesIngressIfVisibilityIsIapAndKindIsDeployment", func(t *testing.T) {

		params := Params{
			Visibility: "iap",
			Kind:       "deployment",
		}

		// act
		templates := getTemplates(params)

		assert.True(t, stringArrayContains(templates, "/templates/ingress.yaml"))
	})

	t.Run("IncludesIngressIfVisibilityIsPublicWhitelistAndKindIsDeployment", func(t *testing.T) {

		params := Params{
			Visibility: "public-whitelist",
			Kind:       "deployment",
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

	t.Run("IncludesInternalIngressIfOneOrMoreInternalHostsAreSetAndKindIsDeployment", func(t *testing.T) {

		params := Params{
			Kind:          "deployment",
			InternalHosts: []string{"ci.estafette.internal"},
		}

		// act
		templates := getTemplates(params)

		assert.True(t, stringArrayContains(templates, "/templates/ingress-internal.yaml"))
	})

	t.Run("DoesNotIncludeInternalIngressIfNoInternalHostsAreSet", func(t *testing.T) {

		params := Params{
			InternalHosts: []string{},
		}

		// act
		templates := getTemplates(params)

		assert.False(t, stringArrayContains(templates, "/templates/ingress-internal.yaml"))
	})

	t.Run("IncludesApplicationSecretsIfLengthOfSecretsIsMoreThanZero", func(t *testing.T) {

		params := Params{
			Secrets: SecretsParams{
				Keys: map[string]interface{}{
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

	t.Run("ReturnsEmptyListIfActionIsRollbackCanaray", func(t *testing.T) {

		params := Params{
			Action: "rollback-canary",
		}

		// act
		templates := getTemplates(params)

		assert.Equal(t, 0, len(templates))
	})

	t.Run("ReturnsOnlyHorizontalPodAutoscalerAndPodDisruptionBudgetIfActionIsDeployCanary", func(t *testing.T) {

		params := Params{
			Action: "deploy-canary",
		}

		// act
		templates := getTemplates(params)

		assert.False(t, stringArrayContains(templates, "/templates/horizontalpodautoscaler.yaml"))
		assert.False(t, stringArrayContains(templates, "/templates/poddisruptionbudget.yaml"))
	})

	t.Run("DoesNotIncludeCertificateSecretIfCertificateSecretIsSet", func(t *testing.T) {

		params := Params{
			Kind:              "deployment",
			CertificateSecret: "shared-wildcard-letsencrypt-certificate",
		}

		// act
		templates := getTemplates(params)

		assert.False(t, stringArrayContains(templates, "/templates/certificate-secret.yaml"))
	})

	t.Run("IncludesCertificateSecretIfCertificateSecretIsNotSet", func(t *testing.T) {

		params := Params{
			Kind: "deployment",
		}

		// act
		templates := getTemplates(params)

		assert.True(t, stringArrayContains(templates, "/templates/certificate-secret.yaml"))
	})

	t.Run("IncludesApigeeIngressIfVisibilityIsApigeeAndKindIsDeployment", func(t *testing.T) {

		params := Params{
			Visibility: "apigee",
			Kind:       "deployment",
		}

		// act
		templates := getTemplates(params)

		assert.True(t, stringArrayContains(templates, "/templates/ingress-apigee.yaml"))
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
