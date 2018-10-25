package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateTemplateData(t *testing.T) {

	t.Run("SetsNameToAppParam", func(t *testing.T) {

		params := Params{
			App: "myapp",
		}

		// act
		templateData := generateTemplateData(params)

		assert.Equal(t, "myapp", templateData.Name)
	})

	t.Run("SetsNamespaceToNamespaceParam", func(t *testing.T) {

		params := Params{
			Namespace: "mynamespace",
		}

		// act
		templateData := generateTemplateData(params)

		assert.Equal(t, "mynamespace", templateData.Namespace)
	})

	t.Run("SetsLabelsToLabelsParam", func(t *testing.T) {

		params := Params{
			Labels: map[string]string{
				"app":  "myapp",
				"team": "myteam",
			},
		}

		// act
		templateData := generateTemplateData(params)

		assert.Equal(t, 2, len(templateData.Labels))
		assert.Equal(t, "myapp", templateData.Labels["app"])
		assert.Equal(t, "myteam", templateData.Labels["team"])
	})

	t.Run("SetsAppLabelSelectorToAppParam", func(t *testing.T) {

		params := Params{
			App: "myapp",
		}

		// act
		templateData := generateTemplateData(params)

		assert.Equal(t, "myapp", templateData.AppLabelSelector)
	})

	t.Run("SetsContainerRepositoryToImageRepositoryParam", func(t *testing.T) {

		params := Params{
			ImageRepository: "myproject",
		}

		// act
		templateData := generateTemplateData(params)

		assert.Equal(t, "myproject", templateData.Container.Repository)
	})

	t.Run("SetsContainerNameToImageNameParam", func(t *testing.T) {

		params := Params{
			ImageName: "my-app",
		}

		// act
		templateData := generateTemplateData(params)

		assert.Equal(t, "my-app", templateData.Container.Name)
	})

	t.Run("SetsContainerTagToImageTagParam", func(t *testing.T) {

		params := Params{
			ImageTag: "1.0.0",
		}

		// act
		templateData := generateTemplateData(params)

		assert.Equal(t, "1.0.0", templateData.Container.Tag)
	})

}
