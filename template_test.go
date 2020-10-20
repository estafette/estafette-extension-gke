package main

import (
	"bytes"
	"strings"
	"testing"
	"text/template"

	"github.com/stretchr/testify/assert"
)

func TestInjectSteps(t *testing.T) {

	t.Run("RenderNamespace", func(t *testing.T) {

		data := TemplateData{
			Namespace: "mynamespace",
		}
		tmpl, err := template.ParseFiles("templates/namespace.yaml")

		// act
		var renderedTemplate bytes.Buffer
		err = tmpl.Execute(&renderedTemplate, data)

		assert.Nil(t, err)
		assert.Equal(t, "apiVersion: v1\nkind: Namespace\nmetadata:\n  name: mynamespace", renderedTemplate.String())
		assert.True(t, strings.Contains(renderedTemplate.String(), "mynamespace"))
	})

	t.Run("RenderServiceAccount", func(t *testing.T) {

		data := TemplateData{
			Name:      "myapp",
			Namespace: "mynamespace",
			Labels: map[string]string{
				"app":  "myapp",
				"team": "myteam",
			},
		}
		tmpl, err := template.ParseFiles("templates/serviceaccount.yaml")

		// act
		var renderedTemplate bytes.Buffer
		err = tmpl.Execute(&renderedTemplate, data)

		assert.Nil(t, err)
		assert.Equal(t, "apiVersion: v1\nkind: ServiceAccount\nmetadata:\n  name: myapp\n  namespace: mynamespace\n  labels:\n    app: myapp\n    team: myteam", renderedTemplate.String())
		assert.True(t, strings.Contains(renderedTemplate.String(), "mynamespace"))
	})

	t.Run("RenderHorizontalPodAutoscaler", func(t *testing.T) {

		data := TemplateData{
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
		tmpl, err := template.ParseFiles("templates/horizontalpodautoscaler.yaml")

		// act
		var renderedTemplate bytes.Buffer
		err = tmpl.Execute(&renderedTemplate, data)

		assert.Nil(t, err)
		assert.Equal(t, "apiVersion: autoscaling/v1\nkind: HorizontalPodAutoscaler\nmetadata:\n  name: myapp-canary\n  namespace: mynamespace\n  labels:\n    app: myapp\n    team: myteam\nspec:\n  scaleTargetRef:\n    apiVersion: apps/v1\n    kind: Deployment\n    name: myapp-canary\n  minReplicas: 3\n  maxReplicas: 19\n  targetCPUUtilizationPercentage: 65", renderedTemplate.String())
		assert.True(t, strings.Contains(renderedTemplate.String(), "mynamespace"))
	})
}
