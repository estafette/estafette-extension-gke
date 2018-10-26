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
			Container: ContainerParams{
				ImageRepository: "myproject",
			},
		}

		// act
		templateData := generateTemplateData(params)

		assert.Equal(t, "myproject", templateData.Container.Repository)
	})

	t.Run("SetsContainerNameToImageNameParam", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				ImageName: "my-app",
			},
		}

		// act
		templateData := generateTemplateData(params)

		assert.Equal(t, "my-app", templateData.Container.Name)
	})

	t.Run("SetsContainerTagToImageTagParam", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				ImageTag: "1.0.0",
			},
		}

		// act
		templateData := generateTemplateData(params)

		assert.Equal(t, "1.0.0", templateData.Container.Tag)
	})

	t.Run("SetsServiceTypeToClusterIPIfVisibilityParamIsPrivate", func(t *testing.T) {

		params := Params{
			Visibility: "private",
		}

		// act
		templateData := generateTemplateData(params)

		assert.Equal(t, "ClusterIP", templateData.ServiceType)
	})

	t.Run("SetsContainerCPURequestToCPURequestParam", func(t *testing.T) {

		params := Params{
			CPU: CPUParams{
				Request: "1200m",
			},
		}

		// act
		templateData := generateTemplateData(params)

		assert.Equal(t, "1200m", templateData.Container.CPURequest)
	})

	t.Run("SetsContainerCPULimitToCPULimitParam", func(t *testing.T) {

		params := Params{
			CPU: CPUParams{
				Limit: "1500m",
			},
		}

		// act
		templateData := generateTemplateData(params)

		assert.Equal(t, "1500m", templateData.Container.CPULimit)
	})

	t.Run("SetsContainerMemoryRequestToMemoryRequestParam", func(t *testing.T) {

		params := Params{
			Memory: MemoryParams{
				Request: "1024Mi",
			},
		}

		// act
		templateData := generateTemplateData(params)

		assert.Equal(t, "1024Mi", templateData.Container.MemoryRequest)
	})

	t.Run("SetsContainerMemoryLimitToMemoryLimitParam", func(t *testing.T) {

		params := Params{
			Memory: MemoryParams{
				Limit: "2048Mi",
			},
		}

		// act
		templateData := generateTemplateData(params)

		assert.Equal(t, "2048Mi", templateData.Container.MemoryLimit)
	})

	t.Run("SetsContainerPortToContainerPortParam", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				Port: 3080,
			},
		}

		// act
		templateData := generateTemplateData(params)

		assert.Equal(t, 3080, templateData.Container.Port)
	})

	t.Run("SetsHostsToHostsParam", func(t *testing.T) {

		params := Params{
			Hosts: []string{
				"gke.estafette.io",
				"gke-deploy.estafette.io",
			},
		}

		// act
		templateData := generateTemplateData(params)

		assert.Equal(t, 2, len(templateData.Hosts))
		assert.Equal(t, "gke.estafette.io", templateData.Hosts[0])
		assert.Equal(t, "gke-deploy.estafette.io", templateData.Hosts[1])
	})

	t.Run("SetsHostsJoinedToCommaSeparatedJoinOfHostsParam", func(t *testing.T) {

		params := Params{
			Hosts: []string{
				"gke.estafette.io",
				"gke-deploy.estafette.io",
			},
		}

		// act
		templateData := generateTemplateData(params)

		assert.Equal(t, "gke.estafette.io,gke-deploy.estafette.io", templateData.HostsJoined)
	})

	t.Run("SetsMinReplicasToAutoscaleMinReplicasParam", func(t *testing.T) {

		params := Params{
			Autoscale: AutoscaleParams{
				MinReplicas: 5,
			},
		}

		// act
		templateData := generateTemplateData(params)

		assert.Equal(t, 5, templateData.MinReplicas)
	})

	t.Run("SetsMaxReplicasToAutoscaleMaxReplicasParam", func(t *testing.T) {

		params := Params{
			Autoscale: AutoscaleParams{
				MaxReplicas: 16,
			},
		}

		// act
		templateData := generateTemplateData(params)

		assert.Equal(t, 16, templateData.MaxReplicas)
	})

	t.Run("SetsTargetCPUPercentageToAutoscaleCPUPercentageParam", func(t *testing.T) {

		params := Params{
			Autoscale: AutoscaleParams{
				CPUPercentage: 75,
			},
		}

		// act
		templateData := generateTemplateData(params)

		assert.Equal(t, 75, templateData.TargetCPUPercentage)
	})
}
