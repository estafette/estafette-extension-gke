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
		ImageRepository: "estafette",
		ImageName:       "my-app",
		ImageTag:        "1.0.0",
	}
)

func TestSetDefaults(t *testing.T) {

	t.Run("DefaultsAppToAppLabelIfEmpty", func(t *testing.T) {

		params := Params{
			App: "",
		}
		appLabel := "myapp"

		// act
		params.SetDefaults(appLabel, "", "", map[string]string{})

		assert.Equal(t, "myapp", params.App)
	})

	t.Run("KeepsAppIfNotEmpty", func(t *testing.T) {

		params := Params{
			App: "yourapp",
		}
		appLabel := "myapp"

		// act
		params.SetDefaults(appLabel, "", "", map[string]string{})

		assert.Equal(t, "yourapp", params.App)
	})

	t.Run("DefaultsImageNameToAppLabelIfEmpty", func(t *testing.T) {

		params := Params{
			ImageName: "",
		}
		appLabel := "myapp"

		// act
		params.SetDefaults(appLabel, "", "", map[string]string{})

		assert.Equal(t, "myapp", params.ImageName)
	})

	t.Run("KeepsImageTagIfNotEmpty", func(t *testing.T) {

		params := Params{
			ImageName: "my-app",
		}
		appLabel := "myapp"

		// act
		params.SetDefaults(appLabel, "", "", map[string]string{})

		assert.Equal(t, "my-app", params.ImageName)
	})

	t.Run("DefaultsImageTagToBuildVersionIfEmpty", func(t *testing.T) {

		params := Params{
			ImageTag: "",
		}
		buildVersion := "1.0.0"

		// act
		params.SetDefaults("", buildVersion, "", map[string]string{})

		assert.Equal(t, "1.0.0", params.ImageTag)
	})

	t.Run("KeepsImageTagIfNotEmpty", func(t *testing.T) {

		params := Params{
			ImageTag: "2.1.3",
		}
		buildVersion := "1.0.0"

		// act
		params.SetDefaults("", buildVersion, "", map[string]string{})

		assert.Equal(t, "2.1.3", params.ImageTag)
	})

	t.Run("DefaultsCredentialsToReleaseNamePrefixedByGKEIfEmpty", func(t *testing.T) {

		params := Params{
			Credentials: "",
		}
		releaseName := "production"

		// act
		params.SetDefaults("", "", releaseName, map[string]string{})

		assert.Equal(t, "gke-production", params.Credentials)
	})

	t.Run("KeepsCredentialsIfNotEmpty", func(t *testing.T) {

		params := Params{
			Credentials: "staging",
		}
		releaseName := "production"

		// act
		params.SetDefaults("", "", releaseName, map[string]string{})

		assert.Equal(t, "staging", params.Credentials)
	})

	t.Run("DefaultsLabelsToEstafetteLabelsIfEmpty", func(t *testing.T) {

		params := Params{
			Labels: map[string]string{},
		}
		estafetteLabels := map[string]string{
			"app":      "myapp",
			"team":     "myteam",
			"language": "golang",
		}

		// act
		params.SetDefaults("", "", "", estafetteLabels)

		assert.Equal(t, 3, len(params.Labels))
		assert.Equal(t, "myapp", params.Labels["app"])
		assert.Equal(t, "myteam", params.Labels["team"])
		assert.Equal(t, "golang", params.Labels["language"])
	})

	t.Run("KeepsLabelsIfNotEmpty", func(t *testing.T) {

		params := Params{
			Labels: map[string]string{
				"app":  "yourapp",
				"team": "yourteam",
			},
		}
		estafetteLabels := map[string]string{
			"app":      "myapp",
			"team":     "myteam",
			"language": "golang",
		}

		// act
		params.SetDefaults("", "", "", estafetteLabels)

		assert.Equal(t, 2, len(params.Labels))
		assert.Equal(t, "yourapp", params.Labels["app"])
		assert.Equal(t, "yourteam", params.Labels["team"])
	})

	t.Run("AddsAppLabelToLabelsIfNotSet", func(t *testing.T) {

		params := Params{
			Labels: map[string]string{
				"team": "yourteam",
			},
		}
		appLabel := "myapp"
		estafetteLabels := map[string]string{
			"app":      "myapp",
			"team":     "myteam",
			"language": "golang",
		}

		// act
		params.SetDefaults(appLabel, "", "", estafetteLabels)

		assert.Equal(t, 2, len(params.Labels))
		assert.Equal(t, "myapp", params.Labels["app"])
		assert.Equal(t, "yourteam", params.Labels["team"])
	})

	t.Run("OverwritesAppLabelToAppIfSetFromEstafetteLabels", func(t *testing.T) {

		params := Params{}
		appLabel := "yourapp"
		estafetteLabels := map[string]string{
			"app":      "myapp",
			"team":     "myteam",
			"language": "golang",
		}

		// act
		params.SetDefaults(appLabel, "", "", estafetteLabels)

		assert.Equal(t, 3, len(params.Labels))
		assert.Equal(t, "yourapp", params.Labels["app"])
		assert.Equal(t, "myteam", params.Labels["team"])
		assert.Equal(t, "golang", params.Labels["language"])
	})

}

func TestSetDefaultsFromCredentials(t *testing.T) {

	t.Run("DefaultsNamespaceToCredentialDefaultNamespaceIfEmpty", func(t *testing.T) {

		params := Params{
			Namespace: "",
		}
		credentials := GKECredentials{
			Name: "gke-1",
			Type: "kubernetes-engine",
			AdditionalProperties: GKECredentialAdditionalProperties{
				DefaultNamespace: "mynamespace",
			},
		}

		// act
		params.SetDefaultsFromCredentials(credentials)

		assert.Equal(t, "mynamespace", params.Namespace)
	})

	t.Run("KeepsNamespaceIfNotEmpty", func(t *testing.T) {

		params := Params{
			Namespace: "yournamespace",
		}
		credentials := GKECredentials{
			Name: "gke-1",
			Type: "kubernetes-engine",
			AdditionalProperties: GKECredentialAdditionalProperties{
				DefaultNamespace: "mynamespace",
			},
		}

		// act
		params.SetDefaultsFromCredentials(credentials)

		assert.Equal(t, "yournamespace", params.Namespace)
	})

	t.Run("DefaultsImageRepositoryToCredentialProjectIfEmpty", func(t *testing.T) {

		params := Params{
			ImageRepository: "",
		}
		credentials := GKECredentials{
			Name: "gke-1",
			Type: "kubernetes-engine",
			AdditionalProperties: GKECredentialAdditionalProperties{
				Project: "myproject",
			},
		}

		// act
		params.SetDefaultsFromCredentials(credentials)

		assert.Equal(t, "myproject", params.ImageRepository)
	})

	t.Run("KeepsImageRepositoryIfNotEmpty", func(t *testing.T) {

		params := Params{
			ImageRepository: "extensions",
		}
		credentials := GKECredentials{
			Name: "gke-1",
			Type: "kubernetes-engine",
			AdditionalProperties: GKECredentialAdditionalProperties{
				Project: "myproject",
			},
		}

		// act
		params.SetDefaultsFromCredentials(credentials)

		assert.Equal(t, "extensions", params.ImageRepository)
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

	t.Run("ReturnsFalseIfImageTagIsNotSet", func(t *testing.T) {

		params := validParams
		params.ImageTag = ""

		// act
		valid, errors := params.ValidateRequiredProperties()

		assert.False(t, valid)
		assert.True(t, len(errors) > 0)
	})

	t.Run("ReturnsTrueIfImageTagIsSet", func(t *testing.T) {

		params := validParams
		params.ImageTag = "1.0.0"

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
