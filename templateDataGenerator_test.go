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

	t.Run("SetsServiceTypeToNodePortIfVisibilityParamIsIap", func(t *testing.T) {

		params := Params{
			Visibility: "iap",
		}

		// act
		templateData := generateTemplateData(params)

		assert.Equal(t, "NodePort", templateData.ServiceType)
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

	t.Run("SetsUseNginxIngressToFalseIfVisibilityIsIap", func(t *testing.T) {

		params := Params{
			Visibility: "iap",
		}

		// act
		templateData := generateTemplateData(params)

		assert.False(t, templateData.UseNginxIngress)
	})

	t.Run("SetsUseGCEIngressToTrueIfVisibilityIsIap", func(t *testing.T) {

		params := Params{
			Visibility: "iap",
		}

		// act
		templateData := generateTemplateData(params)

		assert.True(t, templateData.UseGCEIngress)
	})

	t.Run("SetsUseGCEIngressToFalseIfVisibilityIsPublic", func(t *testing.T) {

		params := Params{
			Visibility: "public",
		}

		// act
		templateData := generateTemplateData(params)

		assert.False(t, templateData.UseGCEIngress)
	})

	t.Run("SetsUseGCEIngressToFalseIfVisibilityIsPrivate", func(t *testing.T) {

		params := Params{
			Visibility: "private",
		}

		// act
		templateData := generateTemplateData(params)

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
		templateData := generateTemplateData(params)

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
		templateData := generateTemplateData(params)

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

	t.Run("SetsReadinessPortToReadinessProbePortParam", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				ReadinessProbe: ProbeParams{
					Port: 5002,
				},
			},
		}

		// act
		templateData := generateTemplateData(params)

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
				EnvironmentVariables: map[string]interface{}{
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

	t.Run("SetsMetricsPathToMetricsProbePathParam", func(t *testing.T) {

		params := Params{
			Container: ContainerParams{
				Metrics: MetricsParams{
					Path: "/readiness",
				},
			},
		}

		// act
		templateData := generateTemplateData(params)

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
		templateData := generateTemplateData(params)

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
		templateData := generateTemplateData(params)

		assert.Equal(t, true, templateData.Container.Metrics.Scrape)
	})

	t.Run("SetsSidecarUseOpenrestySidecarToTrueIfSidecarTypeParamEqualsOpenresty", func(t *testing.T) {

		params := Params{
			Sidecar: SidecarParams{
				Type: "openresty",
			},
		}

		// act
		templateData := generateTemplateData(params)

		assert.True(t, templateData.Sidecar.UseOpenrestySidecar)
	})

	t.Run("SetsSidecarImageToSidecarImageParam", func(t *testing.T) {

		params := Params{
			Sidecar: SidecarParams{
				Image: "estafette/openresty-sidecar:1.13.6.1-alpine",
			},
		}

		// act
		templateData := generateTemplateData(params)

		assert.Equal(t, "estafette/openresty-sidecar:1.13.6.1-alpine", templateData.Sidecar.Image)
	})

	t.Run("SetsSidecarHealthCheckPathToSidecarHealthCheckPathParam", func(t *testing.T) {

		params := Params{
			Sidecar: SidecarParams{
				HealthCheckPath: "/readiness",
			},
		}

		// act
		templateData := generateTemplateData(params)

		assert.Equal(t, "/readiness", templateData.Sidecar.HealthCheckPath)
	})

	t.Run("SetsSidecarCPURequestToSidecarCPURequestParam", func(t *testing.T) {

		params := Params{
			Sidecar: SidecarParams{
				CPU: CPUParams{
					Request: "1200m",
				},
			},
		}

		// act
		templateData := generateTemplateData(params)

		assert.Equal(t, "1200m", templateData.Sidecar.CPURequest)
	})

	t.Run("SetsSidecarCPULimitToSidecarCPULimitParam", func(t *testing.T) {

		params := Params{
			Sidecar: SidecarParams{
				CPU: CPUParams{
					Limit: "1500m",
				},
			},
		}

		// act
		templateData := generateTemplateData(params)

		assert.Equal(t, "1500m", templateData.Sidecar.CPULimit)
	})

	t.Run("SetsSidecarMemoryRequestToSidecarMemoryRequestParam", func(t *testing.T) {

		params := Params{
			Sidecar: SidecarParams{
				Memory: MemoryParams{
					Request: "1024Mi",
				},
			},
		}

		// act
		templateData := generateTemplateData(params)

		assert.Equal(t, "1024Mi", templateData.Sidecar.MemoryRequest)
	})

	t.Run("SetsSidecarMemoryLimitToSidecarMemoryLimitParam", func(t *testing.T) {

		params := Params{
			Sidecar: SidecarParams{
				Memory: MemoryParams{
					Limit: "2048Mi",
				},
			},
		}

		// act
		templateData := generateTemplateData(params)

		assert.Equal(t, "2048Mi", templateData.Sidecar.MemoryLimit)
	})

	t.Run("SetsSidecarEnvironmentVariablesToSidecarEnvironmentVariablesParam", func(t *testing.T) {

		params := Params{
			Sidecar: SidecarParams{
				EnvironmentVariables: map[string]interface{}{
					"MY_CUSTOM_ENV":       "value1",
					"MY_OTHER_CUSTOM_ENV": "value2",
				},
			},
		}

		// act
		templateData := generateTemplateData(params)

		assert.Equal(t, 2, len(templateData.Sidecar.EnvironmentVariables))
		assert.Equal(t, "value1", templateData.Sidecar.EnvironmentVariables["MY_CUSTOM_ENV"])
		assert.Equal(t, "value2", templateData.Sidecar.EnvironmentVariables["MY_OTHER_CUSTOM_ENV"])
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
		templateData := generateTemplateData(params)

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
		templateData := generateTemplateData(params)

		assert.True(t, templateData.MountApplicationSecrets)
	})

	t.Run("SetsMountApplicationSecretsToFalseIfSecretsParamLengthIsZero", func(t *testing.T) {

		params := Params{
			Secrets: SecretsParams{
				Keys: map[string]interface{}{},
			},
		}

		// act
		templateData := generateTemplateData(params)

		assert.False(t, templateData.MountApplicationSecrets)
	})

	t.Run("SetsIngressPathToBasepathParam", func(t *testing.T) {

		params := Params{
			Basepath: "/",
		}

		// act
		templateData := generateTemplateData(params)

		assert.Equal(t, "/", templateData.IngressPath)
	})

	t.Run("AppendSlashToIngressPathToIfBasepathParamDoesNotEndInSlash", func(t *testing.T) {

		params := Params{
			Basepath: "/api",
		}

		// act
		templateData := generateTemplateData(params)

		assert.Equal(t, "/api/", templateData.IngressPath)
	})

	t.Run("AppendSlashStarToIngressPathToIfUseGCEIngressIsTrue", func(t *testing.T) {

		params := Params{
			Basepath:   "/api",
			Visibility: "iap",
		}

		// act
		templateData := generateTemplateData(params)

		assert.Equal(t, "/api/*", templateData.IngressPath)
	})

	t.Run("SetsMountPayloadLoggingToTrueIfEnablePayloadLoggingParamIsTrue", func(t *testing.T) {

		params := Params{
			EnablePayloadLogging: true,
		}

		// act
		templateData := generateTemplateData(params)

		assert.True(t, templateData.MountPayloadLogging)
	})

	t.Run("SetsMountPayloadLoggingToFalseIfEnablePayloadLoggingParamIsFalse", func(t *testing.T) {

		params := Params{
			EnablePayloadLogging: false,
		}

		// act
		templateData := generateTemplateData(params)

		assert.False(t, templateData.MountPayloadLogging)
	})

	t.Run("SetsRollingUpdateMaxSurgeToRollingUpdateMaxSurgeParam", func(t *testing.T) {

		params := Params{
			RollingUpdate: RollingUpdateParams{
				MaxSurge: "25%",
			},
		}

		// act
		templateData := generateTemplateData(params)

		assert.Equal(t, "25%", templateData.RollingUpdateMaxSurge)
	})

	t.Run("SetsRollingUpdateMaxSurgeToRollingUpdateMaxSurgeParam", func(t *testing.T) {

		params := Params{
			RollingUpdate: RollingUpdateParams{
				MaxUnavailable: "15%",
			},
		}

		// act
		templateData := generateTemplateData(params)

		assert.Equal(t, "15%", templateData.RollingUpdateMaxUnavailable)
	})

	t.Run("SetsBuildVersionToBuildVersionParam", func(t *testing.T) {

		params := Params{
			BuildVersion: "1.2.3",
		}

		// act
		templateData := generateTemplateData(params)

		assert.Equal(t, "1.2.3", templateData.BuildVersion)
	})

	t.Run("SetsPreferPreemptiblesgToTrueIfChaosProofParamIsTrue", func(t *testing.T) {

		params := Params{
			ChaosProof: true,
		}

		// act
		templateData := generateTemplateData(params)

		assert.True(t, templateData.PreferPreemptibles)
	})

	t.Run("SetsPreferPreemptiblesToFalseIfChaosProofParamIsFalse", func(t *testing.T) {

		params := Params{
			ChaosProof: false,
		}

		// act
		templateData := generateTemplateData(params)

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
		templateData := generateTemplateData(params)

		assert.True(t, templateData.MountConfigmap)
	})

	t.Run("SetsMountConfigmapToFalseIfConfigFilesParamsLengthIsZero", func(t *testing.T) {

		params := Params{
			Configs: ConfigsParams{
				Files: []string{},
			},
		}

		// act
		templateData := generateTemplateData(params)

		assert.False(t, templateData.MountConfigmap)
	})

	t.Run("SetsConfigMountPathToConfigMountPathParam", func(t *testing.T) {

		params := Params{
			Configs: ConfigsParams{
				MountPath: "/configs",
			},
		}

		// act
		templateData := generateTemplateData(params)

		assert.Equal(t, "/configs", templateData.ConfigMountPath)
	})

	t.Run("SetsSecretMountPathToSecretMountPathParam", func(t *testing.T) {

		params := Params{
			Secrets: SecretsParams{
				MountPath: "/secrets",
			},
		}

		// act
		templateData := generateTemplateData(params)

		assert.Equal(t, "/secrets", templateData.SecretMountPath)
	})

	t.Run("SetsLimitTrustedIPRangesIfVisibilityParamIsPublic", func(t *testing.T) {

		params := Params{
			Visibility: "public",
		}

		// act
		templateData := generateTemplateData(params)

		assert.True(t, templateData.LimitTrustedIPRanges)
	})

	t.Run("SetsLimitTrustedIPRangesToTrueIfVisibilityParamIsPublic", func(t *testing.T) {

		params := Params{
			Visibility: "public",
		}

		// act
		templateData := generateTemplateData(params)

		assert.True(t, templateData.LimitTrustedIPRanges)
	})

	t.Run("SetsLimitTrustedIPRangesToFalseIfVisibilityParamIsIap", func(t *testing.T) {

		params := Params{
			Visibility: "iap",
		}

		// act
		templateData := generateTemplateData(params)

		assert.False(t, templateData.LimitTrustedIPRanges)
	})

	t.Run("SetsLimitTrustedIPRangesToFalseIfVisibilityParamIsPrivate", func(t *testing.T) {

		params := Params{
			Visibility: "private",
		}

		// act
		templateData := generateTemplateData(params)

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
		templateData := generateTemplateData(params)

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
		templateData := generateTemplateData(params)

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
		templateData := generateTemplateData(params)

		assert.Equal(t, "myapp-canary", templateData.NameWithTrack)
	})

	t.Run("AppendsStableToNameWithTrackIfParamsTypeIsRollforward", func(t *testing.T) {

		params := Params{
			App:    "myapp",
			Action: "deploy-stable",
		}

		// act
		templateData := generateTemplateData(params)

		assert.Equal(t, "myapp-stable", templateData.NameWithTrack)
	})

	t.Run("DoesNotAppendTrackToNameWithTrackIfParamsTypeIsSimple", func(t *testing.T) {

		params := Params{
			App:    "myapp",
			Action: "deploy-simple",
		}

		// act
		templateData := generateTemplateData(params)

		assert.Equal(t, "myapp", templateData.NameWithTrack)
	})

	t.Run("SetsIncludeTrackLabelToFalseIfParamsTypeIsSimple", func(t *testing.T) {

		params := Params{
			App:    "myapp",
			Action: "deploy-simple",
		}

		// act
		templateData := generateTemplateData(params)

		assert.False(t, templateData.IncludeTrackLabel)
	})

	t.Run("SetsIncludeTrackLabelToTrueIfParamsTypeIsCanary", func(t *testing.T) {

		params := Params{
			App:    "myapp",
			Action: "deploy-canary",
		}

		// act
		templateData := generateTemplateData(params)

		assert.True(t, templateData.IncludeTrackLabel)
	})

	t.Run("SetsIncludeTrackLabelToTrueIfParamsTypeIsRollforward", func(t *testing.T) {

		params := Params{
			App:    "myapp",
			Action: "deploy-stable",
		}

		// act
		templateData := generateTemplateData(params)

		assert.True(t, templateData.IncludeTrackLabel)
	})

	t.Run("SetsTrackLabelToCanaryIfParamsTypeIsCanary", func(t *testing.T) {

		params := Params{
			App:    "myapp",
			Action: "deploy-canary",
		}

		// act
		templateData := generateTemplateData(params)

		assert.Equal(t, "canary", templateData.TrackLabel)
	})

	t.Run("SetsTrackLabelToStableIfParamsTypeIsRollforward", func(t *testing.T) {

		params := Params{
			App:    "myapp",
			Action: "deploy-stable",
		}

		// act
		templateData := generateTemplateData(params)

		assert.Equal(t, "stable", templateData.TrackLabel)
	})
}
