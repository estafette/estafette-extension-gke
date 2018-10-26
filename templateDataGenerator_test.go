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

	t.Run("SetsServiceTypeToLoadBalancerIfVisibilityParamIsPublic", func(t *testing.T) {

		params := Params{
			Visibility: "public",
		}

		// act
		templateData := generateTemplateData(params)

		assert.Equal(t, "LoadBalancer", templateData.ServiceType)
	})

	t.Run("SetsUseDNSAnnotationsOnIngressToTrueIfVisibilityParamIsPrivate", func(t *testing.T) {

		params := Params{
			Visibility: "private",
		}

		// act
		templateData := generateTemplateData(params)

		assert.True(t, templateData.UseDNSAnnotationsOnIngress)
	})

	t.Run("SetsUseDNSAnnotationsOnIngressToFalseIfVisibilityParamIsPublic", func(t *testing.T) {

		params := Params{
			Visibility: "public",
		}

		// act
		templateData := generateTemplateData(params)

		assert.False(t, templateData.UseDNSAnnotationsOnIngress)
	})

	t.Run("SetsUseDNSAnnotationsOnServiceToTrueIfVisibilityParamIsPublic", func(t *testing.T) {

		params := Params{
			Visibility: "public",
		}

		// act
		templateData := generateTemplateData(params)

		assert.True(t, templateData.UseDNSAnnotationsOnService)
	})

	t.Run("SetsUseDNSAnnotationsOnServiceToFalseIfVisibilityParamIsPrivate", func(t *testing.T) {

		params := Params{
			Visibility: "private",
		}

		// act
		templateData := generateTemplateData(params)

		assert.False(t, templateData.UseDNSAnnotationsOnService)
	})

	t.Run("SetsContainerCPURequestToCPURequestParam", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				CPU: CPUParams{
					Request: "1200m",
				},
			},
		}

		// act
		templateData := generateTemplateData(params)

		assert.Equal(t, "1200m", templateData.Container.CPURequest)
	})

	t.Run("SetsContainerCPULimitToCPULimitParam", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				CPU: CPUParams{
					Limit: "1500m",
				},
			},
		}

		// act
		templateData := generateTemplateData(params)

		assert.Equal(t, "1500m", templateData.Container.CPULimit)
	})

	t.Run("SetsContainerMemoryRequestToMemoryRequestParam", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				Memory: MemoryParams{
					Request: "1024Mi",
				},
			},
		}

		// act
		templateData := generateTemplateData(params)

		assert.Equal(t, "1024Mi", templateData.Container.MemoryRequest)
	})

	t.Run("SetsContainerMemoryLimitToMemoryLimitParam", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				Memory: MemoryParams{
					Limit: "2048Mi",
				},
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

	t.Run("SetsUseNginxIngressToTrueIfVisibilityIsPrivate", func(t *testing.T) {

		params := Params{
			Visibility: "private",
		}

		// act
		templateData := generateTemplateData(params)

		assert.True(t, templateData.UseNginxIngress)
	})

	t.Run("SetsUseNginxIngressToFalseIfVisibilityIsPublic", func(t *testing.T) {

		params := Params{
			Visibility: "public",
		}

		// act
		templateData := generateTemplateData(params)

		assert.False(t, templateData.UseNginxIngress)
	})

	t.Run("SetsLivenessPathToLivenessProbePathParam", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				LivenessProbe: ProbeParams{
					Path: "/liveness",
				},
			},
		}

		// act
		templateData := generateTemplateData(params)

		assert.Equal(t, "/liveness", templateData.Container.Liveness.Path)
	})

	t.Run("SetsLivenessInitialDelaySecondsToLivenessProbeInitialDelaySecondsParam", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				LivenessProbe: ProbeParams{
					InitialDelaySeconds: 30,
				},
			},
		}

		// act
		templateData := generateTemplateData(params)

		assert.Equal(t, 30, templateData.Container.Liveness.InitialDelaySeconds)
	})

	t.Run("SetsLivenessTimeoutSecondsToLivenessProbeTimeoutSecondsParam", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				LivenessProbe: ProbeParams{
					TimeoutSeconds: 1,
				},
			},
		}

		// act
		templateData := generateTemplateData(params)

		assert.Equal(t, 1, templateData.Container.Liveness.TimeoutSeconds)
	})

	t.Run("SetsReadinessPathToReadinessProbePathParam", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				ReadinessProbe: ProbeParams{
					Path: "/readiness",
				},
			},
		}

		// act
		templateData := generateTemplateData(params)

		assert.Equal(t, "/readiness", templateData.Container.Readiness.Path)
	})

	t.Run("SetsReadinessInitialDelaySecondsToReadinessProbeInitialDelaySecondsParam", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				ReadinessProbe: ProbeParams{
					InitialDelaySeconds: 30,
				},
			},
		}

		// act
		templateData := generateTemplateData(params)

		assert.Equal(t, 30, templateData.Container.Readiness.InitialDelaySeconds)
	})

	t.Run("SetsReadinessTimeoutSecondsToReadinessProbeTimeoutSecondsParam", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				ReadinessProbe: ProbeParams{
					TimeoutSeconds: 1,
				},
			},
		}

		// act
		templateData := generateTemplateData(params)

		assert.Equal(t, 1, templateData.Container.Readiness.TimeoutSeconds)
	})

	t.Run("SetsEnvironmentVariablesToContainerEnvironmentVariablesParam", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				EnvironmentVariables: map[string]string{
					"MY_CUSTOM_ENV":       "value1",
					"MY_OTHER_CUSTOM_ENV": "value2",
				},
			},
		}

		// act
		templateData := generateTemplateData(params)

		assert.Equal(t, 2, len(templateData.Container.EnvironmentVariables))
		assert.Equal(t, "value1", templateData.Container.EnvironmentVariables["MY_CUSTOM_ENV"])
		assert.Equal(t, "value2", templateData.Container.EnvironmentVariables["MY_OTHER_CUSTOM_ENV"])
	})
}
