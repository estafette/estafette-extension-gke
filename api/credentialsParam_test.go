package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	validCredentialsParam = CredentialsParam{
		Credentials: "gke-production",
	}
)

func TestCredentialsParamSetDefaults(t *testing.T) {

	t.Run("DefaultsCredentialsToReleaseNamePrefixedByGKEIfEmpty", func(t *testing.T) {

		params := CredentialsParam{
			Credentials: "",
		}
		releaseName := "production"

		// act
		params.SetDefaults(releaseName)

		assert.Equal(t, "gke-production", params.Credentials)
	})

	t.Run("KeepsCredentialsIfNotEmpty", func(t *testing.T) {

		params := CredentialsParam{
			Credentials: "staging",
		}
		releaseName := "production"

		// act
		params.SetDefaults(releaseName)

		assert.Equal(t, "staging", params.Credentials)
	})
}

func TestCredentialsParamValidateRequiredProperties(t *testing.T) {
	t.Run("ReturnsFalseIfCredentialsIsNotSet", func(t *testing.T) {

		params := validCredentialsParam
		params.Credentials = ""

		// act
		valid, errors := params.ValidateRequiredProperties()

		assert.False(t, valid)
		assert.True(t, len(errors) > 0)
	})

	t.Run("ReturnsTrueIfCredentialsIsSet", func(t *testing.T) {

		params := validCredentialsParam
		params.Credentials = "gke-production"

		// act
		valid, errors := params.ValidateRequiredProperties()

		assert.True(t, valid)
		assert.True(t, len(errors) == 0)
	})

}
