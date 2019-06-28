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
		templateData := generateTemplateData(params, -1, "", "")

		assert.Equal(t, "myapp", templateData.Name)
	})

	t.Run("SetsNamespaceToNamespaceParam", func(t *testing.T) {

		params := Params{
			Namespace: "mynamespace",
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

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
		templateData := generateTemplateData(params, -1, "", "")

		assert.Equal(t, 2, len(templateData.Labels))
		assert.Equal(t, "myapp", templateData.Labels["app"])
		assert.Equal(t, "myteam", templateData.Labels["team"])
	})

	t.Run("SetsAppLabelSelectorToAppParam", func(t *testing.T) {

		params := Params{
			App: "myapp",
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		assert.Equal(t, "myapp", templateData.AppLabelSelector)
	})

	t.Run("ReplacesAppLabelValueWithAppParamIfAppLabelExists", func(t *testing.T) {

		params := Params{
			Labels: map[string]string{
				"app":  "myapp",
				"team": "myteam",
			},
			App: "yourapp",
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		assert.Equal(t, 2, len(templateData.Labels))
		assert.Equal(t, "yourapp", templateData.Labels["app"])
	})

	t.Run("AddsAppLabelValueWithAppParamIfAppLabelDoesNotExists", func(t *testing.T) {

		params := Params{
			Labels: map[string]string{
				"team": "myteam",
			},
			App: "yourapp",
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		assert.Equal(t, 2, len(templateData.Labels))
		assert.Equal(t, "yourapp", templateData.Labels["app"])
	})

	t.Run("SetsContainerRepositoryToImageRepositoryParam", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				ImageRepository: "myproject",
			},
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		assert.Equal(t, "myproject", templateData.Container.Repository)
	})

	t.Run("SetsContainerNameToImageNameParam", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				ImageName: "my-app",
			},
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		assert.Equal(t, "my-app", templateData.Container.Name)
	})

	t.Run("SetsContainerTagToImageTagParam", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				ImageTag: "1.0.0",
			},
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		assert.Equal(t, "1.0.0", templateData.Container.Tag)
	})

	t.Run("SetsServiceTypeToClusterIPIfVisibilityParamIsPrivate", func(t *testing.T) {

		params := Params{
			Visibility: "private",
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		assert.Equal(t, "ClusterIP", templateData.ServiceType)
	})

	t.Run("SetsServiceTypeToClusterIPIfVisibilityParamIsPublicWhitelist", func(t *testing.T) {

		params := Params{
			Visibility: "public-whitelist",
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		assert.Equal(t, "ClusterIP", templateData.ServiceType)
	})

	t.Run("SetsServiceTypeToNodePortIfVisibilityParamIsIap", func(t *testing.T) {

		params := Params{
			Visibility: "iap",
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		assert.Equal(t, "NodePort", templateData.ServiceType)
	})

	t.Run("SetsServiceTypeToLoadBalancerIfVisibilityParamIsPublic", func(t *testing.T) {

		params := Params{
			Visibility: "public",
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		assert.Equal(t, "LoadBalancer", templateData.ServiceType)
	})

	t.Run("SetsUseDNSAnnotationsOnIngressToTrueIfVisibilityParamIsPrivate", func(t *testing.T) {

		params := Params{
			Visibility: "private",
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		assert.True(t, templateData.UseDNSAnnotationsOnIngress)
	})

	t.Run("SetsUseDNSAnnotationsOnIngressToTrueIfVisibilityParamIsPublicWhitelist", func(t *testing.T) {

		params := Params{
			Visibility: "public-whitelist",
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		assert.True(t, templateData.UseDNSAnnotationsOnIngress)
	})

	t.Run("SetsUseDNSAnnotationsOnIngressToFalseIfVisibilityParamIsPublic", func(t *testing.T) {

		params := Params{
			Visibility: "public",
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		assert.False(t, templateData.UseDNSAnnotationsOnIngress)
	})

	t.Run("SetsUseDNSAnnotationsOnServiceToTrueIfVisibilityParamIsPublic", func(t *testing.T) {

		params := Params{
			Visibility: "public",
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		assert.True(t, templateData.UseDNSAnnotationsOnService)
	})

	t.Run("SetsUseDNSAnnotationsOnServiceToFalseIfVisibilityParamIsPrivate", func(t *testing.T) {

		params := Params{
			Visibility: "private",
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		assert.False(t, templateData.UseDNSAnnotationsOnService)
	})

	t.Run("SetsUseDNSAnnotationsOnServiceToFalseIfVisibilityParamIsPublicWhitelist", func(t *testing.T) {

		params := Params{
			Visibility: "public-whitelist",
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

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
		templateData := generateTemplateData(params, -1, "", "")

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
		templateData := generateTemplateData(params, -1, "", "")

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
		templateData := generateTemplateData(params, -1, "", "")

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
		templateData := generateTemplateData(params, -1, "", "")

		assert.Equal(t, "2048Mi", templateData.Container.MemoryLimit)
	})

	t.Run("SetsContainerPortToContainerPortParam", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				Port: 3080,
			},
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

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
		templateData := generateTemplateData(params, -1, "", "")

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
		templateData := generateTemplateData(params, -1, "", "")

		assert.Equal(t, "gke.estafette.io,gke-deploy.estafette.io", templateData.HostsJoined)
	})

	t.Run("SetsInternalHostsToInternalHostsParam", func(t *testing.T) {

		params := Params{
			InternalHosts: []string{
				"gke.estafette.io",
				"gke-deploy.estafette.io",
			},
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		assert.Equal(t, 2, len(templateData.InternalHosts))
		assert.Equal(t, "gke.estafette.io", templateData.InternalHosts[0])
		assert.Equal(t, "gke-deploy.estafette.io", templateData.InternalHosts[1])
	})

	t.Run("SetsInternalHostsJoinedToCommaSeparatedJoinOfInternalHostsParam", func(t *testing.T) {

		params := Params{
			InternalHosts: []string{
				"gke.estafette.io",
				"gke-deploy.estafette.io",
			},
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		assert.Equal(t, "gke.estafette.io,gke-deploy.estafette.io", templateData.InternalHostsJoined)
	})

	t.Run("SetsMinReplicasToAutoscaleMinReplicasParam", func(t *testing.T) {

		params := Params{
			Autoscale: AutoscaleParams{
				MinReplicas: 5,
			},
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		assert.Equal(t, 5, templateData.MinReplicas)
	})

	t.Run("SetsMaxReplicasToAutoscaleMaxReplicasParam", func(t *testing.T) {

		params := Params{
			Autoscale: AutoscaleParams{
				MaxReplicas: 16,
			},
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		assert.Equal(t, 16, templateData.MaxReplicas)
	})

	t.Run("SetsUseNginxIngressToTrueIfVisibilityIsPrivate", func(t *testing.T) {

		params := Params{
			Visibility: "private",
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		assert.True(t, templateData.UseNginxIngress)
	})

	t.Run("SetsUseNginxIngressToTrueIfVisibilityIsPublicWhitelist", func(t *testing.T) {

		params := Params{
			Visibility: "public-whitelist",
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		assert.True(t, templateData.UseNginxIngress)
	})

	t.Run("SetsUseNginxIngressToFalseIfVisibilityIsPublic", func(t *testing.T) {

		params := Params{
			Visibility: "public",
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		assert.False(t, templateData.UseNginxIngress)
	})

	t.Run("SetsUseNginxIngressToFalseIfVisibilityIsIap", func(t *testing.T) {

		params := Params{
			Visibility: "iap",
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		assert.False(t, templateData.UseNginxIngress)
	})

	t.Run("SetsUseGCEIngressToTrueIfVisibilityIsIap", func(t *testing.T) {

		params := Params{
			Visibility: "iap",
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		assert.True(t, templateData.UseGCEIngress)
	})

	t.Run("SetsUseGCEIngressToFalseIfVisibilityIsPublic", func(t *testing.T) {

		params := Params{
			Visibility: "public",
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		assert.False(t, templateData.UseGCEIngress)
	})

	t.Run("SetsUseGCEIngressToFalseIfVisibilityIsPrivate", func(t *testing.T) {

		params := Params{
			Visibility: "private",
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		assert.False(t, templateData.UseGCEIngress)
	})

	t.Run("SetsUseGCEIngressToFalseIfVisibilityIsPublicWhitelist", func(t *testing.T) {

		params := Params{
			Visibility: "public-whitelist",
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		assert.False(t, templateData.UseGCEIngress)
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
		templateData := generateTemplateData(params, -1, "", "")

		assert.Equal(t, "/liveness", templateData.Container.Liveness.Path)
	})

	t.Run("SetsLivenessPortToLivenessProbePortParam", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				LivenessProbe: ProbeParams{
					Port: 5001,
				},
			},
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		assert.Equal(t, 5001, templateData.Container.Liveness.Port)
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
		templateData := generateTemplateData(params, -1, "", "")

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
		templateData := generateTemplateData(params, -1, "", "")

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
		templateData := generateTemplateData(params, -1, "", "")

		assert.Equal(t, "/readiness", templateData.Container.Readiness.Path)
	})

	t.Run("SetsReadinessPortToReadinessProbePortParam", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				ReadinessProbe: ProbeParams{
					Port: 5002,
				},
			},
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		assert.Equal(t, 5002, templateData.Container.Readiness.Port)
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
		templateData := generateTemplateData(params, -1, "", "")

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
		templateData := generateTemplateData(params, -1, "", "")

		assert.Equal(t, 1, templateData.Container.Readiness.TimeoutSeconds)
	})

	t.Run("SetsEnvironmentVariablesToContainerEnvironmentVariablesParam", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				EnvironmentVariables: map[string]interface{}{
					"MY_CUSTOM_ENV":       "value1",
					"MY_OTHER_CUSTOM_ENV": "value2",
				},
			},
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		assert.Equal(t, "value1", templateData.Container.EnvironmentVariables["MY_CUSTOM_ENV"])
		assert.Equal(t, "value2", templateData.Container.EnvironmentVariables["MY_OTHER_CUSTOM_ENV"])
	})

	t.Run("AddsJaegerServiceNameToEnvironmentVariables", func(t *testing.T) {

		params := Params{
			App: "my-app",
			Container: ContainerParams{
				EnvironmentVariables: map[string]interface{}{},
			},
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		assert.Equal(t, "my-app", templateData.Container.EnvironmentVariables["JAEGER_SERVICE_NAME"])
	})

	t.Run("SetsMetricsPathToMetricsProbePathParam", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				Metrics: MetricsParams{
					Path: "/readiness",
				},
			},
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		assert.Equal(t, "/readiness", templateData.Container.Metrics.Path)
	})

	t.Run("SetsMetricsPortToMetricsPortParam", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				Metrics: MetricsParams{
					Port: 3080,
				},
			},
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		assert.Equal(t, 3080, templateData.Container.Metrics.Port)
	})

	t.Run("SetsMetricsScrapeToMetricsScrapeParam", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				Metrics: MetricsParams{
					Scrape: &trueValue,
				},
			},
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		assert.Equal(t, true, templateData.Container.Metrics.Scrape)
	})

	t.Run("SetsUseLifecyclePreStopSleepCommandToLifecyclePrestopSleepParam", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				Lifecycle: LifecycleParams{
					PrestopSleep: &trueValue,
				},
			},
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		assert.Equal(t, true, templateData.Container.UseLifecyclePreStopSleepCommand)
	})

	t.Run("SetsPreStopSleepSecondsToLifecyclePrestopSleepSecondsParam", func(t *testing.T) {

		sleepValue := 25

		params := Params{
			Container: ContainerParams{
				Lifecycle: LifecycleParams{
					PrestopSleepSeconds: &sleepValue,
				},
			},
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		assert.Equal(t, 25, templateData.Container.PreStopSleepSeconds)
	})

	t.Run("SidecarAddedToSidecarsCollection", func(t *testing.T) {

		params := Params{
			Sidecars: []*SidecarParams{
				&SidecarParams{
					Type: "openresty",
				},
			},
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		assert.Equal(t, 1, len(templateData.Sidecars))
	})

	t.Run("SetsSidecarTypeToSidecarType", func(t *testing.T) {

		params := Params{
			Sidecars: []*SidecarParams{
				&SidecarParams{
					Type: "openresty",
				},
			},
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		assert.Equal(t, "openresty", templateData.Sidecars[0].Type)
	})

	t.Run("SetsSidecarImageToSidecarImageParam", func(t *testing.T) {

		params := Params{
			Sidecars: []*SidecarParams{
				&SidecarParams{
					Image: "estafette/openresty-sidecar:1.13.6.1-alpine",
				},
			},
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		assert.Equal(t, "estafette/openresty-sidecar:1.13.6.1-alpine", templateData.Sidecars[0].Image)
	})

	t.Run("SetsSidecarHealthCheckPathToSidecarHealthCheckPathParam", func(t *testing.T) {

		params := Params{
			Sidecars: []*SidecarParams{
				&SidecarParams{
					HealthCheckPath: "/readiness",
				},
			},
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		assert.Equal(t, "/readiness", templateData.Sidecars[0].SidecarSpecificProperties["healthcheckpath"])
	})

	t.Run("SetsSidecarCPURequestToSidecarCPURequestParam", func(t *testing.T) {

		params := Params{
			Sidecars: []*SidecarParams{
				&SidecarParams{
					CPU: CPUParams{
						Request: "1200m",
					},
				},
			},
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		assert.Equal(t, "1200m", templateData.Sidecars[0].CPURequest)
	})

	t.Run("SetsSidecarCPULimitToSidecarCPULimitParam", func(t *testing.T) {

		params := Params{
			Sidecars: []*SidecarParams{
				&SidecarParams{
					CPU: CPUParams{
						Limit: "1500m",
					},
				},
			},
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		assert.Equal(t, "1500m", templateData.Sidecars[0].CPULimit)
	})

	t.Run("SetsSidecarMemoryRequestToSidecarMemoryRequestParam", func(t *testing.T) {

		params := Params{
			Sidecars: []*SidecarParams{
				&SidecarParams{
					Memory: MemoryParams{
						Request: "1024Mi",
					},
				},
			},
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		assert.Equal(t, "1024Mi", templateData.Sidecars[0].MemoryRequest)
	})

	t.Run("SetsSidecarMemoryLimitToSidecarMemoryLimitParam", func(t *testing.T) {

		params := Params{
			Sidecars: []*SidecarParams{
				&SidecarParams{
					Memory: MemoryParams{
						Limit: "2048Mi",
					},
				},
			},
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		assert.Equal(t, "2048Mi", templateData.Sidecars[0].MemoryLimit)
	})

	t.Run("SetsSidecarEnvironmentVariablesToSidecarEnvironmentVariablesParam", func(t *testing.T) {

		params := Params{
			Sidecars: []*SidecarParams{
				&SidecarParams{
					EnvironmentVariables: map[string]interface{}{
						"MY_CUSTOM_ENV":       "value1",
						"MY_OTHER_CUSTOM_ENV": "value2",
					},
				},
			},
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		// assert.Equal(t, 2, len(templateData.Sidecar.EnvironmentVariables))
		assert.Equal(t, "value1", templateData.Sidecars[0].EnvironmentVariables["MY_CUSTOM_ENV"])
		assert.Equal(t, "value2", templateData.Sidecars[0].EnvironmentVariables["MY_OTHER_CUSTOM_ENV"])
	})

	t.Run("SetsCloudSQLProxySpecificArgsToSidecarSpecificProperties", func(t *testing.T) {

		params := Params{
			Sidecars: []*SidecarParams{
				&SidecarParams{
					HealthCheckPath:                   "testHealthCheckPath",
					DbInstanceConnectionName:          "testDbInstanceConnectionName",
					SQLProxyPort:                      15,
					SQLProxyTerminationTimeoutSeconds: 16,
				},
			},
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		assert.Equal(t, 4, len(templateData.Sidecars[0].SidecarSpecificProperties))
		assert.Equal(t, "testHealthCheckPath", templateData.Sidecars[0].SidecarSpecificProperties["healthcheckpath"])
		assert.Equal(t, "testDbInstanceConnectionName", templateData.Sidecars[0].SidecarSpecificProperties["dbinstanceconnectionname"])
		assert.Equal(t, 15, templateData.Sidecars[0].SidecarSpecificProperties["sqlproxyport"])
		assert.Equal(t, 16, templateData.Sidecars[0].SidecarSpecificProperties["sqlproxyterminationtimeoutseconds"])
	})

	t.Run("SetsSecretsToSecretsParam", func(t *testing.T) {

		params := Params{
			Secrets: SecretsParams{
				Keys: map[string]interface{}{
					"secret-file-1.json": "c29tZSBzZWNyZXQgdmFsdWU=",
					"secret-file-2.yaml": "YW5vdGhlciBzZWNyZXQgdmFsdWU=",
				},
			},
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		assert.Equal(t, 2, len(templateData.Secrets))
		assert.Equal(t, "c29tZSBzZWNyZXQgdmFsdWU=", templateData.Secrets["secret-file-1.json"])
		assert.Equal(t, "YW5vdGhlciBzZWNyZXQgdmFsdWU=", templateData.Secrets["secret-file-2.yaml"])
	})

	t.Run("SetsMountApplicationSecretsToTrueIfSecretsParamLengthIsLargerThanZero", func(t *testing.T) {

		params := Params{
			Secrets: SecretsParams{
				Keys: map[string]interface{}{
					"secret-file-1.json": "c29tZSBzZWNyZXQgdmFsdWU=",
					"secret-file-2.yaml": "YW5vdGhlciBzZWNyZXQgdmFsdWU=",
				},
			},
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		assert.True(t, templateData.MountApplicationSecrets)
	})

	t.Run("SetsMountApplicationSecretsToFalseIfSecretsParamLengthIsZero", func(t *testing.T) {

		params := Params{
			Secrets: SecretsParams{
				Keys: map[string]interface{}{},
			},
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		assert.False(t, templateData.MountApplicationSecrets)
	})

	t.Run("SetsIngressPathToBasepathParam", func(t *testing.T) {

		params := Params{
			Basepath: "/",
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		assert.Equal(t, "/", templateData.IngressPath)
	})

	t.Run("AppendSlashToIngressPathIfBasepathParamDoesNotEndInSlash", func(t *testing.T) {

		params := Params{
			Basepath: "/api",
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		assert.Equal(t, "/api/", templateData.IngressPath)
	})

	t.Run("AppendSlashStarToIngressPathIfUseGCEIngressIsTrue", func(t *testing.T) {

		params := Params{
			Basepath:   "/api",
			Visibility: "iap",
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		assert.Equal(t, "/api/*", templateData.IngressPath)
	})

	t.Run("SetsInternalIngressPathToBasepathParam", func(t *testing.T) {

		params := Params{
			Basepath: "/",
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		assert.Equal(t, "/", templateData.InternalIngressPath)
	})

	t.Run("AppendSlashToInternalIngressPathIfBasepathParamDoesNotEndInSlash", func(t *testing.T) {

		params := Params{
			Basepath: "/api",
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		assert.Equal(t, "/api/", templateData.InternalIngressPath)
	})

	t.Run("DoNotAppendSlashStarToInternalIngressPathIfUseGCEIngressIsTrue", func(t *testing.T) {

		params := Params{
			Basepath:   "/api",
			Visibility: "iap",
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		assert.Equal(t, "/api/", templateData.InternalIngressPath)
	})

	t.Run("SetsMountPayloadLoggingToTrueIfEnablePayloadLoggingParamIsTrue", func(t *testing.T) {

		params := Params{
			EnablePayloadLogging: true,
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		assert.True(t, templateData.MountPayloadLogging)
	})

	t.Run("SetsMountPayloadLoggingToFalseIfEnablePayloadLoggingParamIsFalse", func(t *testing.T) {

		params := Params{
			EnablePayloadLogging: false,
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		assert.False(t, templateData.MountPayloadLogging)
	})

	t.Run("SetsAddSafeToEvictAnnotationToTrueIfEnablePayloadLoggingParamIsTrue", func(t *testing.T) {

		params := Params{
			EnablePayloadLogging: true,
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		assert.True(t, templateData.AddSafeToEvictAnnotation)
	})

	t.Run("SetsAddSafeToEvictAnnotationToFalseIfEnablePayloadLoggingParamIsFalse", func(t *testing.T) {

		params := Params{
			EnablePayloadLogging: false,
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		assert.False(t, templateData.AddSafeToEvictAnnotation)
	})

	t.Run("SetsRollingUpdateMaxSurgeToRollingUpdateMaxSurgeParam", func(t *testing.T) {

		params := Params{
			RollingUpdate: RollingUpdateParams{
				MaxSurge: "25%",
			},
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		assert.Equal(t, "25%", templateData.RollingUpdateMaxSurge)
	})

	t.Run("SetsRollingUpdateMaxSurgeToRollingUpdateMaxSurgeParam", func(t *testing.T) {

		params := Params{
			RollingUpdate: RollingUpdateParams{
				MaxUnavailable: "15%",
			},
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		assert.Equal(t, "15%", templateData.RollingUpdateMaxUnavailable)
	})

	t.Run("SetsBuildVersionToBuildVersionParam", func(t *testing.T) {

		params := Params{
			BuildVersion: "1.2.3",
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		assert.Equal(t, "1.2.3", templateData.BuildVersion)
	})

	t.Run("SetsPreferPreemptiblesgToTrueIfChaosProofParamIsTrue", func(t *testing.T) {

		params := Params{
			ChaosProof: true,
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		assert.True(t, templateData.PreferPreemptibles)
	})

	t.Run("SetsPreferPreemptiblesToFalseIfChaosProofParamIsFalse", func(t *testing.T) {

		params := Params{
			ChaosProof: false,
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		assert.False(t, templateData.PreferPreemptibles)
	})

	t.Run("SetsMountConfigmapToTrueIfConfigFilesParamsLengthIsLargerThanZero", func(t *testing.T) {

		params := Params{
			Configs: ConfigsParams{
				Files: []string{
					"config.json",
				},
			},
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		assert.True(t, templateData.MountConfigmap)
	})

	t.Run("SetsMountConfigmapToTrueIfInlineFilesParamsLengthIsLargerThanZero", func(t *testing.T) {

		params := Params{
			Configs: ConfigsParams{
				InlineFiles: map[string]string{
					"inline-config.properties": "enemies=aliens\nlives=3",
				},
			},
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		assert.True(t, templateData.MountConfigmap)
	})

	t.Run("SetsMountConfigmapToFalseIfConfigFilesAndInlineFilesParamsLengthAreZero", func(t *testing.T) {

		params := Params{
			Configs: ConfigsParams{
				Files:       []string{},
				InlineFiles: map[string]string{},
			},
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		assert.False(t, templateData.MountConfigmap)
	})

	t.Run("SetsConfigMountPathToConfigMountPathParam", func(t *testing.T) {

		params := Params{
			Configs: ConfigsParams{
				MountPath: "/configs",
			},
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		assert.Equal(t, "/configs", templateData.ConfigMountPath)
	})

	t.Run("SetsSecretMountPathToSecretMountPathParam", func(t *testing.T) {

		params := Params{
			Secrets: SecretsParams{
				MountPath: "/secrets",
			},
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		assert.Equal(t, "/secrets", templateData.SecretMountPath)
	})

	t.Run("SetsLimitTrustedIPRangesIfVisibilityParamIsPublic", func(t *testing.T) {

		params := Params{
			Visibility: "public",
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		assert.True(t, templateData.LimitTrustedIPRanges)
	})

	t.Run("SetsLimitTrustedIPRangesToTrueIfVisibilityParamIsPublic", func(t *testing.T) {

		params := Params{
			Visibility: "public",
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		assert.True(t, templateData.LimitTrustedIPRanges)
	})

	t.Run("SetsLimitTrustedIPRangesToFalseIfVisibilityParamIsIap", func(t *testing.T) {

		params := Params{
			Visibility: "iap",
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		assert.False(t, templateData.LimitTrustedIPRanges)
	})

	t.Run("SetsLimitTrustedIPRangesToFalseIfVisibilityParamIsPrivate", func(t *testing.T) {

		params := Params{
			Visibility: "private",
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		assert.False(t, templateData.LimitTrustedIPRanges)
	})

	t.Run("SetsTrustedIPRangesToTrustedIPRangesParams", func(t *testing.T) {

		params := Params{
			TrustedIPRanges: []string{
				"103.21.244.0/22",
				"103.22.200.0/22",
				"103.31.4.0/22",
				"104.16.0.0/12",
				"108.162.192.0/18",
				"131.0.72.0/22",
				"141.101.64.0/18",
				"162.158.0.0/15",
				"172.64.0.0/13",
				"173.245.48.0/20",
				"188.114.96.0/20",
				"190.93.240.0/20",
				"197.234.240.0/22",
				"198.41.128.0/17",
			},
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		assert.Equal(t, 14, len(templateData.TrustedIPRanges))
	})

	t.Run("SetsLocalManifestDataToAllLocalManifestDataCombined", func(t *testing.T) {

		params := Params{
			Manifests: ManifestsParams{
				Files: []string{
					"./gke/service.yaml",
					"./gke/ingress.yaml",
				},
				Data: map[string]interface{}{
					"property1": "value 1",
					"property2": "value 2",
					"property3": "value 3",
				},
			},
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		assert.Equal(t, 3, len(templateData.ManifestData))
		assert.Equal(t, "value 1", templateData.ManifestData["property1"])
		assert.Equal(t, "value 2", templateData.ManifestData["property2"])
		assert.Equal(t, "value 3", templateData.ManifestData["property3"])
	})

	t.Run("AppendsCanaryToNameWithTrackIfParamsTypeIsCanary", func(t *testing.T) {

		params := Params{
			App:    "myapp",
			Action: "deploy-canary",
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		assert.Equal(t, "myapp-canary", templateData.NameWithTrack)
	})

	t.Run("AppendsStableToNameWithTrackIfParamsTypeIsRollforward", func(t *testing.T) {

		params := Params{
			App:    "myapp",
			Action: "deploy-stable",
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		assert.Equal(t, "myapp-stable", templateData.NameWithTrack)
	})

	t.Run("DoesNotAppendTrackToNameWithTrackIfParamsTypeIsSimple", func(t *testing.T) {

		params := Params{
			App:    "myapp",
			Action: "deploy-simple",
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		assert.Equal(t, "myapp", templateData.NameWithTrack)
	})

	t.Run("SetsIncludeReleaseIDLabelToFalseIfReleaseIDIsEmpty", func(t *testing.T) {

		params := Params{
			App:    "myapp",
			Action: "deploy-simple",
		}
		releaseID := ""

		// act
		templateData := generateTemplateData(params, -1, releaseID, "")

		assert.False(t, templateData.IncludeReleaseIDLabel)
	})

	t.Run("SetsIncludeReleaseIDLabelToTrueIfReleaseIDIsNotEmpty", func(t *testing.T) {

		params := Params{
			App:    "myapp",
			Action: "deploy-simple",
		}
		releaseID := "1"

		// act
		templateData := generateTemplateData(params, -1, releaseID, "")

		assert.True(t, templateData.IncludeReleaseIDLabel)
	})

	t.Run("SetsReleaseIDLabelToReleaseID", func(t *testing.T) {

		params := Params{
			App:    "myapp",
			Action: "deploy-simple",
		}
		releaseID := "1"

		// act
		templateData := generateTemplateData(params, -1, releaseID, "")

		assert.Equal(t, "1", templateData.ReleaseIDLabel)
	})

	t.Run("SetsIncludeTriggeredByLabelToFalseIfTriggeredByIsEmpty", func(t *testing.T) {

		params := Params{
			App:    "myapp",
			Action: "deploy-simple",
		}
		triggeredBy := ""

		// act
		templateData := generateTemplateData(params, -1, "", triggeredBy)

		assert.False(t, templateData.IncludeTriggeredByLabel)
	})

	t.Run("SetsIncludeTriggeredByLabelToTrueIfTriggeredByIsNotEmpty", func(t *testing.T) {

		params := Params{
			App:    "myapp",
			Action: "deploy-simple",
		}
		triggeredBy := "user@estafette.io"

		// act
		templateData := generateTemplateData(params, -1, "", triggeredBy)

		assert.True(t, templateData.IncludeTriggeredByLabel)
	})

	t.Run("SetsTriggeredByLabelToTriggeredBySanitizedAsLabel", func(t *testing.T) {

		params := Params{
			App:    "myapp",
			Action: "deploy-simple",
		}
		triggeredBy := "user@estafette.io"

		// act
		templateData := generateTemplateData(params, -1, "", triggeredBy)

		assert.Equal(t, "user-estafette.io", templateData.TriggeredByLabel)
	})

	t.Run("SetsIncludeTrackLabelToFalseIfParamsTypeIsSimple", func(t *testing.T) {

		params := Params{
			App:    "myapp",
			Action: "deploy-simple",
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		assert.False(t, templateData.IncludeTrackLabel)
	})

	t.Run("SetsIncludeTrackLabelToTrueIfParamsTypeIsCanary", func(t *testing.T) {

		params := Params{
			App:    "myapp",
			Action: "deploy-canary",
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		assert.True(t, templateData.IncludeTrackLabel)
	})

	t.Run("SetsIncludeTrackLabelToTrueIfParamsTypeIsRollforward", func(t *testing.T) {

		params := Params{
			App:    "myapp",
			Action: "deploy-stable",
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		assert.True(t, templateData.IncludeTrackLabel)
	})

	t.Run("SetsTrackLabelToCanaryIfParamsTypeIsCanary", func(t *testing.T) {

		params := Params{
			App:    "myapp",
			Action: "deploy-canary",
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		assert.Equal(t, "canary", templateData.TrackLabel)
	})

	t.Run("SetsTrackLabelToStableIfParamsTypeIsRollforward", func(t *testing.T) {

		params := Params{
			App:    "myapp",
			Action: "deploy-stable",
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		assert.Equal(t, "stable", templateData.TrackLabel)
	})

	t.Run("SetsAdditionalVolumeMountsToVolumeMountsParam", func(t *testing.T) {

		params := Params{
			VolumeMounts: []VolumeMountParams{
				VolumeMountParams{
					Name:      "client-certs",
					MountPath: "/cockroach-certs",
					Volume: map[string]interface{}{
						"secret": map[string]interface{}{
							"secretName": "estafette.client.estafette",
							"items": []interface{}{
								map[string]interface{}{
									"key":  "key",
									"path": "key",
									"mode": 0600,
								},
								map[string]interface{}{
									"key":  "cert",
									"path": "cert",
								},
							},
						},
					},
				},
			},
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		assert.Equal(t, 1, len(templateData.AdditionalVolumeMounts))
		assert.Equal(t, "client-certs", templateData.AdditionalVolumeMounts[0].Name)
		assert.Equal(t, "/cockroach-certs", templateData.AdditionalVolumeMounts[0].MountPath)
		assert.Equal(t, "secret:\n  items:\n  - key: key\n    mode: 384\n    path: key\n  - key: cert\n    path: cert\n  secretName: estafette.client.estafette\n", templateData.AdditionalVolumeMounts[0].VolumeYAML)
	})

	t.Run("SetsAdditionalContainerPortsToContainerAdditionalPortsParam", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				AdditionalPorts: []*AdditionalPortParams{
					&AdditionalPortParams{
						Name:       "grpc",
						Port:       8085,
						Protocol:   "TCP",
						Visibility: "private",
					},
					&AdditionalPortParams{
						Name:       "grpc",
						Port:       8085,
						Protocol:   "UDP",
						Visibility: "public",
					},
				},
			},
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		assert.Equal(t, 2, len(templateData.AdditionalContainerPorts))
	})

	t.Run("SetsAdditionalServicePortsToContainerAdditionalPortsParamForPortsWithVisibilityEqualToVisibilityParam", func(t *testing.T) {

		params := Params{
			Visibility: "private",
			Container: ContainerParams{
				AdditionalPorts: []*AdditionalPortParams{
					&AdditionalPortParams{
						Name:       "grpc",
						Port:       8085,
						Protocol:   "TCP",
						Visibility: "private",
					},
					&AdditionalPortParams{
						Name:       "snmp",
						Port:       8086,
						Protocol:   "UDP",
						Visibility: "public",
					},
				},
			},
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		assert.Equal(t, 1, len(templateData.AdditionalServicePorts))
		assert.Equal(t, "grpc", templateData.AdditionalServicePorts[0].Name)
		assert.Equal(t, 8085, templateData.AdditionalServicePorts[0].Port)
		assert.Equal(t, "TCP", templateData.AdditionalServicePorts[0].Protocol)
	})

	t.Run("SetsOverrideDefaultWhitelistToTrueIfVisibilityEqualsPublicWhitelistAndWhitelistedIPSHasOneOrMoreItems", func(t *testing.T) {

		params := Params{
			Visibility: "public-whitelist",
			WhitelistedIPS: []string{
				"0.0.0.0/0",
			},
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		assert.True(t, templateData.OverrideDefaultWhitelist)
	})

	t.Run("SetsOverrideDefaultWhitelistToFalseIfVisibilityEqualsPublicWhitelistButWhitelistedIPSHasNoItems", func(t *testing.T) {

		params := Params{
			Visibility:     "public-whitelist",
			WhitelistedIPS: []string{},
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		assert.False(t, templateData.OverrideDefaultWhitelist)
	})

	t.Run("SetsOverrideDefaultWhitelistToFalseIfVisibilityIsPrivate", func(t *testing.T) {

		params := Params{
			Visibility: "private",
			WhitelistedIPS: []string{
				"0.0.0.0/0",
			},
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		assert.False(t, templateData.OverrideDefaultWhitelist)
	})

	t.Run("SetsOverrideDefaultWhitelistToFalseIfVisibilityIsPrivate", func(t *testing.T) {

		params := Params{
			Visibility: "iap",
			WhitelistedIPS: []string{
				"0.0.0.0/0",
			},
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		assert.False(t, templateData.OverrideDefaultWhitelist)
	})

	t.Run("SetsOverrideDefaultWhitelistToFalseIfVisibilityIsPrivate", func(t *testing.T) {

		params := Params{
			Visibility: "public",
			WhitelistedIPS: []string{
				"0.0.0.0/0",
			},
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		assert.False(t, templateData.OverrideDefaultWhitelist)
	})

	t.Run("SetsNginxIngressWhitelistToCommaSeparatedJoingOfWhitelistedIPS", func(t *testing.T) {

		params := Params{
			Visibility: "public-whitelist",
			WhitelistedIPS: []string{
				"10.0.0.0/8",
				"172.16.0.0/12",
				"192.168.0.0/16",
			},
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		assert.Equal(t, "10.0.0.0/8,172.16.0.0/12,192.168.0.0/16", templateData.NginxIngressWhitelist)
	})

	t.Run("SetsIncludeReplicasToTrueIfCurrentReplicasIsGreaterThanZero", func(t *testing.T) {

		params := Params{}

		// act
		templateData := generateTemplateData(params, 1, "", "")

		assert.True(t, templateData.IncludeReplicas)
	})

	t.Run("SetsIncludeReplicasToFalseIfCurrentReplicasIsZeroOrLess", func(t *testing.T) {

		params := Params{}

		// act
		templateData := generateTemplateData(params, 0, "", "")

		assert.False(t, templateData.IncludeReplicas)
	})

	t.Run("SetsReplicasToCurrentReplicasIfGreaterThanZero", func(t *testing.T) {

		params := Params{}

		// act
		templateData := generateTemplateData(params, 15, "", "")

		assert.Equal(t, 15, templateData.Replicas)
	})

	t.Run("SetsReplicasToMinReplicasIfReplicasIsZeroOrLess", func(t *testing.T) {

		params := Params{
			Autoscale: AutoscaleParams{
				MinReplicas: 3,
			},
		}

		// act
		templateData := generateTemplateData(params, 0, "", "")

		assert.Equal(t, 3, templateData.Replicas)
	})

	t.Run("SetsScheduleToScheduleParam", func(t *testing.T) {

		params := Params{
			Schedule: "*/5 * * * *",
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		assert.Equal(t, "*/5 * * * *", templateData.Schedule)
	})

	t.Run("SetsUseHpaScalerToAutoscalerSafetyEnabledParam", func(t *testing.T) {

		params := Params{
			Autoscale: AutoscaleParams{
				Safety: AutoscaleSafetyParams{
					Enabled: true,
				},
			},
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		assert.True(t, templateData.UseHpaScaler)
	})

	t.Run("SetsHpaScalerPromQueryToAutoscalerSafetyPromQueryParam", func(t *testing.T) {

		params := Params{
			Autoscale: AutoscaleParams{
				Safety: AutoscaleSafetyParams{
					PromQuery: "sum(rate(nginx_http_requests_total{app='my-app'}[5m])) by (app)",
				},
			},
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		assert.Equal(t, "sum(rate(nginx_http_requests_total{app='my-app'}[5m])) by (app)", templateData.HpaScalerPromQuery)
	})

	t.Run("SetsHpaScalerRequestsPerReplicaToAutoscalerSafetyPromQueryParam", func(t *testing.T) {

		params := Params{
			Autoscale: AutoscaleParams{
				Safety: AutoscaleSafetyParams{
					Ratio: "0.25",
				},
			},
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		assert.Equal(t, "0.25", templateData.HpaScalerRequestsPerReplica)
	})

	t.Run("SetsHpaScalerRequestsPerReplicaToAutoscalerSafetyPromQueryParam", func(t *testing.T) {

		params := Params{
			Autoscale: AutoscaleParams{
				Safety: AutoscaleSafetyParams{
					Delta: "-2.7584",
				},
			},
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		assert.Equal(t, "-2.7584", templateData.HpaScalerDelta)
	})

	t.Run("SetsHpaScalerRequestsPerReplicaToAutoscalerSafetyPromQueryParam", func(t *testing.T) {

		params := Params{
			Autoscale: AutoscaleParams{
				Safety: AutoscaleSafetyParams{
					ScaleDownRatio: "0.2",
				},
			},
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		assert.Equal(t, "0.2", templateData.HpaScalerScaleDownMaxRatio)
	})

	t.Run("SetsAllHostsToHostsAndInternalHostsAppended", func(t *testing.T) {

		params := Params{
			Hosts: []string{
				"ci.estafette.io",
			},
			InternalHosts: []string{
				"ci.internal.estafette.io",
			},
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		assert.Equal(t, 2, len(templateData.AllHosts))
		assert.Equal(t, "ci.estafette.io", templateData.AllHosts[0])
		assert.Equal(t, "ci.internal.estafette.io", templateData.AllHosts[1])
	})

	t.Run("SetsAllHostsJoinedToHostsAndInternalHostsAppendedSeparatedByComma", func(t *testing.T) {

		params := Params{
			Hosts: []string{
				"ci.estafette.io",
			},
			InternalHosts: []string{
				"ci.internal.estafette.io",
			},
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		assert.Equal(t, "ci.estafette.io,ci.internal.estafette.io", templateData.AllHostsJoined)
	})

	t.Run("SetsAllHostsToHostsAndInternalHostsAppendedWhenOnlyHostsAreSet", func(t *testing.T) {

		params := Params{
			Hosts: []string{
				"ci.estafette.io",
			},
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		assert.Equal(t, 1, len(templateData.AllHosts))
		assert.Equal(t, "ci.estafette.io", templateData.AllHosts[0])
	})

	t.Run("SetsAllHostsJoinedToHostsAndInternalHostsAppendedSeparatedByCommaWhenOnlyHostsAreSet", func(t *testing.T) {

		params := Params{
			Hosts: []string{
				"ci.estafette.io",
			},
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		assert.Equal(t, "ci.estafette.io", templateData.AllHostsJoined)
	})

	t.Run("SetsAllHostsToHostsAndInternalHostsAppendedWhenOnlyInternalHostsAreSet", func(t *testing.T) {

		params := Params{
			InternalHosts: []string{
				"ci.internal.estafette.io",
			},
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		assert.Equal(t, 1, len(templateData.AllHosts))
		assert.Equal(t, "ci.internal.estafette.io", templateData.AllHosts[0])
	})

	t.Run("SetsAllHostsJoinedToHostsAndInternalHostsAppendedSeparatedByCommaWhenOnlyInternalHostsAreSet", func(t *testing.T) {

		params := Params{
			InternalHosts: []string{
				"ci.internal.estafette.io",
			},
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		assert.Equal(t, "ci.internal.estafette.io", templateData.AllHostsJoined)
	})

	t.Run("SetsNginxIngressProxyConnectTimeoutToRequestTimeoutParam", func(t *testing.T) {

		params := Params{
			Request: RequestParams{
				Timeout: "75",
			},
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		assert.Equal(t, 75, templateData.NginxIngressProxyConnectTimeout)
	})

	t.Run("SetsNginxIngressProxyConnectTimeoutToRequestTimeoutParamWithSecondSuffix", func(t *testing.T) {

		params := Params{
			Request: RequestParams{
				Timeout: "75s",
			},
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		assert.Equal(t, 75, templateData.NginxIngressProxyConnectTimeout)
	})

	t.Run("SetsNginxIngressProxyConnectTimeoutTo75IfRequestTimeoutParamIsLargerThan75", func(t *testing.T) {

		params := Params{
			Request: RequestParams{
				Timeout: "180s",
			},
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		assert.Equal(t, 75, templateData.NginxIngressProxyConnectTimeout)
	})

	t.Run("SetsNginxIngressProxySendTimeoutToRequestTimeoutParam", func(t *testing.T) {

		params := Params{
			Request: RequestParams{
				Timeout: "300",
			},
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		assert.Equal(t, 300, templateData.NginxIngressProxySendTimeout)
	})

	t.Run("SetsNginxIngressProxySendTimeoutToRequestTimeoutParamWithSecondSuffix", func(t *testing.T) {

		params := Params{
			Request: RequestParams{
				Timeout: "300s",
			},
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		assert.Equal(t, 300, templateData.NginxIngressProxySendTimeout)
	})

	t.Run("SetsNginxIngressProxyReadTimeoutToRequestTimeoutParam", func(t *testing.T) {

		params := Params{
			Request: RequestParams{
				Timeout: "300",
			},
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		assert.Equal(t, 300, templateData.NginxIngressProxyReadTimeout)
	})

	t.Run("SetsNginxIngressProxyReadTimeoutToRequestTimeoutParamWithSecondSuffix", func(t *testing.T) {

		params := Params{
			Request: RequestParams{
				Timeout: "300s",
			},
		}

		// act
		templateData := generateTemplateData(params, -1, "", "")

		assert.Equal(t, 300, templateData.NginxIngressProxyReadTimeout)
	})

}
