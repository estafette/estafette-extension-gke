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

		assert.True(t, stringArrayContains(templates, "ingress.yaml"))
	})

	t.Run("IncludesIngressIfVisibilityIsIap", func(t *testing.T) {

		params := Params{
			Visibility: "iap",
		}

		// act
		templates := getTemplates(params)

		assert.True(t, stringArrayContains(templates, "ingress.yaml"))
	})

	t.Run("DoesNotIncludeIngressIfVisibilityIsPublic", func(t *testing.T) {

		params := Params{
			Visibility: "public",
		}

		// act
		templates := getTemplates(params)

		assert.False(t, stringArrayContains(templates, "ingress.yaml"))
	})

	t.Run("IncludesApplicationSecretsIfLengthOfSecretsIsMoreThanZero", func(t *testing.T) {

		params := Params{
			Secrets: map[string]string{
				"secret-file-1.json": "c29tZSBzZWNyZXQgdmFsdWU=",
				"secret-file-2.yaml": "YW5vdGhlciBzZWNyZXQgdmFsdWU=",
			},
		}

		// act
		templates := getTemplates(params)

		assert.True(t, stringArrayContains(templates, "application-secrets.yaml"))
	})

	t.Run("DoesNotIncludeApplicationSecretsIfLengthOfSecretsZero", func(t *testing.T) {

		params := Params{}

		// act
		templates := getTemplates(params)

		assert.False(t, stringArrayContains(templates, "application-secrets.yaml"))
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
