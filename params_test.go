package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	validParams = Params{
		Credentials:     "gke-production",
		App:             "myapp",
		Namespace:       "mynamespace",
		AppContainerTag: "1.0.0",
	}
)

func TestSetDefaults(t *testing.T) {

	t.Run("DefaultsAppToAppLabelIfEmpty", func(t *testing.T) {

		params := Params{
			App: "",
		}
		appLabel := "myapp"

		// act
		params.SetDefaults(appLabel, "", "")

		assert.Equal(t, "myapp", params.App)
	})

	t.Run("KeepsAppIfNotEmpty", func(t *testing.T) {

		params := Params{
			App: "yourapp",
		}
		appLabel := "myapp"

		// act
		params.SetDefaults(appLabel, "", "")

		assert.Equal(t, "yourapp", params.App)
	})

	t.Run("DefaultsAppContainerTagToBuildVersionIfEmpty", func(t *testing.T) {

		params := Params{
			AppContainerTag: "",
		}
		buildVersion := "1.0.0"

		// act
		params.SetDefaults("", buildVersion, "")

		assert.Equal(t, "1.0.0", params.AppContainerTag)
	})

	t.Run("KeepsAppContainerTagIfNotEmpty", func(t *testing.T) {

		params := Params{
			AppContainerTag: "2.1.3",
		}
		buildVersion := "1.0.0"

		// act
		params.SetDefaults("", buildVersion, "")

		assert.Equal(t, "2.1.3", params.AppContainerTag)
	})

	t.Run("DefaultsCredentialsToReleaseNamePrefixedByGKEIfEmpty", func(t *testing.T) {

		params := Params{
			Credentials: "",
		}
		releaseName := "production"

		// act
		params.SetDefaults("", "", releaseName)

		assert.Equal(t, "gke-production", params.Credentials)
	})

	t.Run("KeepsCredentialsIfNotEmpty", func(t *testing.T) {

		params := Params{
			Credentials: "staging",
		}
		releaseName := "production"

		// act
		params.SetDefaults("", "", releaseName)

		assert.Equal(t, "staging", params.Credentials)
	})
}

func TestSetDefaultNamespace(t *testing.T) {

	t.Run("DefaultsAppToAppLabelIfEmpty", func(t *testing.T) {

		params := Params{
			Namespace: "",
		}
		defaultNamespace := "mynamespace"

		// act
		params.SetDefaultNamespace(defaultNamespace)

		assert.Equal(t, "mynamespace", params.Namespace)
	})

	t.Run("KeepsAppIfNotEmpty", func(t *testing.T) {

		params := Params{
			Namespace: "yournamespace",
		}
		defaultNamespace := "mynamespace"

		// act
		params.SetDefaultNamespace(defaultNamespace)

		assert.Equal(t, "yournamespace", params.Namespace)
	})
}

func TestValidateRequiredProperties(t *testing.T) {

	t.Run("ReturnsFalseIfAppIsNotSet", func(t *testing.T) {

		params := validParams
		params.App = ""

		// act
		valid, errors := params.ValidateRequiredProperties()

		assert.False(t, valid)
		assert.True(t, len(errors) > 0)
	})

	t.Run("ReturnsTrueIfAppIsSet", func(t *testing.T) {

		params := validParams
		params.App = "myapp"

		// act
		valid, errors := params.ValidateRequiredProperties()

		assert.True(t, valid)
		assert.True(t, len(errors) == 0)
	})

	t.Run("ReturnsFalseIfNamespaceIsNotSet", func(t *testing.T) {

		params := validParams
		params.Namespace = ""

		// act
		valid, errors := params.ValidateRequiredProperties()

		assert.False(t, valid)
		assert.True(t, len(errors) > 0)
	})

	t.Run("ReturnsTrueIfNamespaceIsSet", func(t *testing.T) {

		params := validParams
		params.Namespace = "mynamespace"

		// act
		valid, errors := params.ValidateRequiredProperties()

		assert.True(t, valid)
		assert.True(t, len(errors) == 0)
	})

	t.Run("ReturnsFalseIfAppContainerTagIsNotSet", func(t *testing.T) {

		params := validParams
		params.AppContainerTag = ""

		// act
		valid, errors := params.ValidateRequiredProperties()

		assert.False(t, valid)
		assert.True(t, len(errors) > 0)
	})

	t.Run("ReturnsTrueIfAppContainerTagIsSet", func(t *testing.T) {

		params := validParams
		params.AppContainerTag = "1.0.0"

		// act
		valid, errors := params.ValidateRequiredProperties()

		assert.True(t, valid)
		assert.True(t, len(errors) == 0)
	})

	t.Run("ReturnsFalseIfCredentialsIsNotSet", func(t *testing.T) {

		params := validParams
		params.Credentials = ""

		// act
		valid, errors := params.ValidateRequiredProperties()

		assert.False(t, valid)
		assert.True(t, len(errors) > 0)
	})

	t.Run("ReturnsTrueIfCredentialsIsSet", func(t *testing.T) {

		params := validParams
		params.Credentials = "gke-production"

		// act
		valid, errors := params.ValidateRequiredProperties()

		assert.True(t, valid)
		assert.True(t, len(errors) == 0)
	})
}
