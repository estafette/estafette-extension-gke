package generator

import (
	"context"
	"testing"

	"github.com/estafette/estafette-extension-gke/api"
	"github.com/stretchr/testify/assert"
)

var (
	trueValue  = true
	falseValue = false
)

func TestGenerateTemplateData(t *testing.T) {

	t.Run("SetsNameToAppParam", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			App: "myapp",
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.Equal(t, "myapp", templateData.Name)
	})

	t.Run("SetsNamespaceToNamespaceParam", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Namespace: "mynamespace",
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.Equal(t, "mynamespace", templateData.Namespace)
	})

	t.Run("SetsLabelsToLabelsParam", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Labels: map[string]string{
				"app":  "myapp",
				"team": "myteam",
			},
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.Equal(t, 2, len(templateData.Labels))
		assert.Equal(t, "myapp", templateData.Labels["app"])
		assert.Equal(t, "myteam", templateData.Labels["team"])
	})

	t.Run("SetsAppLabelSelectorToAppParam", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			App: "myapp",
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.Equal(t, "myapp", templateData.AppLabelSelector)
	})

	t.Run("ReplacesAppLabelValueWithAppParamIfAppLabelExists", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Labels: map[string]string{
				"app":  "myapp",
				"team": "myteam",
			},
			App: "yourapp",
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.Equal(t, 2, len(templateData.Labels))
		assert.Equal(t, "yourapp", templateData.Labels["app"])
	})

	t.Run("AddsAppLabelValueWithAppParamIfAppLabelDoesNotExists", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Labels: map[string]string{
				"team": "myteam",
			},
			App: "yourapp",
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.Equal(t, 2, len(templateData.Labels))
		assert.Equal(t, "yourapp", templateData.Labels["app"])
	})

	t.Run("SetsContainerRepositoryToImageRepositoryParam", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Container: api.ContainerParams{
				ImageRepository: "myproject",
			},
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.Equal(t, "myproject", templateData.Container.Repository)
	})

	t.Run("SetsContainerNameToImageNameParam", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Container: api.ContainerParams{
				ImageName: "my-app",
			},
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.Equal(t, "my-app", templateData.Container.Name)
	})

	t.Run("SetsContainerTagToImageTagParam", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Container: api.ContainerParams{
				ImageTag: "1.0.0",
			},
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.Equal(t, "1.0.0", templateData.Container.Tag)
	})

	t.Run("SetsServiceTypeToClusterIPIfVisibilityParamIsPrivate", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Visibility: api.VisibilityPrivate,
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.Equal(t, "ClusterIP", templateData.ServiceType)
	})

	t.Run("SetsServiceTypeToClusterIPIfVisibilityParamIsPublicWhitelist", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Visibility: api.VisibilityPublicWhitelist,
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.Equal(t, "ClusterIP", templateData.ServiceType)
	})

	t.Run("SetsServiceTypeToNodePortIfVisibilityParamIsIap", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Visibility: api.VisibilityIAP,
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.Equal(t, "NodePort", templateData.ServiceType)
	})

	t.Run("SetsServiceTypeToLoadBalancerIfVisibilityParamIsPublic", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Visibility: api.VisibilityPublic,
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.Equal(t, "LoadBalancer", templateData.ServiceType)
	})

	t.Run("SetsServiceTypeToClusterIpIfVisibilityParamIsESP", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Visibility: api.VisibilityESP,
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")
		assert.Equal(t, "ClusterIP", templateData.ServiceType)
	})

	t.Run("SetsUseDNSAnnotationsOnIngressToTrueIfVisibilityParamIsPrivate", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Visibility: api.VisibilityPrivate,
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.True(t, templateData.UseDNSAnnotationsOnIngress)
	})

	t.Run("SetsUseDNSAnnotationsOnIngressToTrueIfVisibilityParamIsPublicWhitelist", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Visibility: api.VisibilityPublicWhitelist,
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.True(t, templateData.UseDNSAnnotationsOnIngress)
	})

	t.Run("SetsUseDNSAnnotationsOnIngressToFalseIfVisibilityParamIsPublic", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Visibility: api.VisibilityPublic,
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.False(t, templateData.UseDNSAnnotationsOnIngress)
	})

	t.Run("SetsUseDNSAnnotationsOnServiceToTrueIfVisibilityParamIsPublic", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Visibility: api.VisibilityPublic,
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.True(t, templateData.UseDNSAnnotationsOnService)
	})

	t.Run("SetsUseDNSAnnotationsOnServiceToFalseIfVisibilityParamIsPrivate", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Visibility: api.VisibilityPrivate,
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.False(t, templateData.UseDNSAnnotationsOnService)
	})

	t.Run("SetsUseDNSAnnotationsOnServiceToFalseIfVisibilityParamIsPublicWhitelist", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Visibility: api.VisibilityPublicWhitelist,
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.False(t, templateData.UseDNSAnnotationsOnService)
	})

	t.Run("SetsUseCloudflareProxyToTrueIfVisibilityParamIsPrivate", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Visibility: api.VisibilityPrivate,
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.True(t, templateData.UseCloudflareProxy)
	})

	t.Run("SetsUseCloudflareProxyToTrueIfVisibilityParamIsPublicWhitelist", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Visibility: api.VisibilityPublicWhitelist,
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.True(t, templateData.UseCloudflareProxy)
	})

	t.Run("SetsUseCloudflareProxyToTrueIfVisibilityParamIsPublic", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Visibility: api.VisibilityPublic,
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.True(t, templateData.UseCloudflareProxy)
	})

	t.Run("SetsUseCloudflareProxyToFalseIfVisibilityParamIsIAP", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Visibility: api.VisibilityIAP,
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.False(t, templateData.UseCloudflareProxy)
	})

	t.Run("SetsContainerCPURequestToCPURequestParam", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Container: api.ContainerParams{
				CPU: api.CPUParams{
					Request: "1200m",
				},
			},
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.Equal(t, "1200m", templateData.Container.CPURequest)
	})

	t.Run("SetsContainerCPULimitToCPULimitParam", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Container: api.ContainerParams{
				CPU: api.CPUParams{
					Limit: "1500m",
				},
			},
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.Equal(t, "1500m", templateData.Container.CPULimit)
	})

	t.Run("SetsContainerMemoryRequestToMemoryRequestParam", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Container: api.ContainerParams{
				Memory: api.MemoryParams{
					Request: "1024Mi",
				},
			},
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.Equal(t, "1024Mi", templateData.Container.MemoryRequest)
	})

	t.Run("SetsContainerMemoryLimitToMemoryLimitParam", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Container: api.ContainerParams{
				Memory: api.MemoryParams{
					Limit: "2048Mi",
				},
			},
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.Equal(t, "2048Mi", templateData.Container.MemoryLimit)
	})

	t.Run("SetsContainerPortToContainerPortParam", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Container: api.ContainerParams{
				Port: 3080,
			},
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.Equal(t, 3080, templateData.Container.Port)
	})

	t.Run("SetsHostsToHostsParam", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Hosts: []string{
				"gke.estafette.io",
				"gke-deploy.estafette.io",
			},
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.Equal(t, 2, len(templateData.Hosts))
		assert.Equal(t, "gke.estafette.io", templateData.Hosts[0])
		assert.Equal(t, "gke-deploy.estafette.io", templateData.Hosts[1])
	})

	t.Run("SetsHostsJoinedToCommaSeparatedJoinOfHostsParam", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Hosts: []string{
				"gke.estafette.io",
				"gke-deploy.estafette.io",
			},
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.Equal(t, "gke.estafette.io,gke-deploy.estafette.io", templateData.HostsJoined)
	})

	t.Run("SetsInternalHostsToInternalHostsParam", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			InternalHosts: []string{
				"gke.estafette.io",
				"gke-deploy.estafette.io",
			},
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.Equal(t, 2, len(templateData.InternalHosts))
		assert.Equal(t, "gke.estafette.io", templateData.InternalHosts[0])
		assert.Equal(t, "gke-deploy.estafette.io", templateData.InternalHosts[1])
	})

	t.Run("SetsInternalHostsJoinedToCommaSeparatedJoinOfInternalHostsParam", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			InternalHosts: []string{
				"gke.estafette.io",
				"gke-deploy.estafette.io",
			},
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.Equal(t, "gke.estafette.io,gke-deploy.estafette.io", templateData.InternalHostsJoined)
	})

	t.Run("SetsMinReplicasToAutoscaleMinReplicasParam", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Autoscale: api.AutoscaleParams{
				MinReplicas: 5,
			},
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.Equal(t, 5, templateData.MinReplicas)
	})

	t.Run("SetsMaxReplicasToAutoscaleMaxReplicasParam", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Autoscale: api.AutoscaleParams{
				MaxReplicas: 16,
			},
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.Equal(t, 16, templateData.MaxReplicas)
	})

	t.Run("SetsUseNginxIngressToTrueIfVisibilityIsPrivate", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Visibility: api.VisibilityPrivate,
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.True(t, templateData.UseNginxIngress)
	})

	t.Run("SetsUseNginxIngressToTrueIfVisibilityIsPublicWhitelist", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Visibility: api.VisibilityPublicWhitelist,
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.True(t, templateData.UseNginxIngress)
	})

	t.Run("SetsUseNginxIngressToFalseIfVisibilityIsPublic", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Visibility: api.VisibilityPublic,
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.False(t, templateData.UseNginxIngress)
	})

	t.Run("SetsUseNginxIngressToFalseIfVisibilityIsIap", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Visibility: api.VisibilityIAP,
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.False(t, templateData.UseNginxIngress)
	})

	t.Run("SetsUseGCEIngressToTrueIfVisibilityIsIap", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Visibility: api.VisibilityIAP,
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.True(t, templateData.UseGCEIngress)
	})

	t.Run("SetsUseGCEIngressToFalseIfVisibilityIsPublic", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Visibility: api.VisibilityPublic,
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.False(t, templateData.UseGCEIngress)
	})

	t.Run("SetsUseGCEIngressToFalseIfVisibilityIsPrivate", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Visibility: api.VisibilityPrivate,
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.False(t, templateData.UseGCEIngress)
	})

	t.Run("SetsUseGCEIngressToFalseIfVisibilityIsPublicWhitelist", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Visibility: api.VisibilityPublicWhitelist,
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.False(t, templateData.UseGCEIngress)
	})

	t.Run("SetsLivenessPathToLivenessProbePathParam", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Container: api.ContainerParams{
				LivenessProbe: api.ProbeParams{
					Path: "/liveness",
				},
			},
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.Equal(t, "/liveness", templateData.Container.Liveness.Path)
	})

	t.Run("SetsLivenessPortToLivenessProbePortParam", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Container: api.ContainerParams{
				LivenessProbe: api.ProbeParams{
					Port: 5001,
				},
			},
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.Equal(t, 5001, templateData.Container.Liveness.Port)
	})

	t.Run("SetsLivenessInitialDelaySecondsToLivenessProbeInitialDelaySecondsParam", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Container: api.ContainerParams{
				LivenessProbe: api.ProbeParams{
					InitialDelaySeconds: 30,
				},
			},
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.Equal(t, 30, templateData.Container.Liveness.InitialDelaySeconds)
	})

	t.Run("SetsLivenessTimeoutSecondsToLivenessProbeTimeoutSecondsParam", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Container: api.ContainerParams{
				LivenessProbe: api.ProbeParams{
					TimeoutSeconds: 1,
				},
			},
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.Equal(t, 1, templateData.Container.Liveness.TimeoutSeconds)
	})

	t.Run("SetsLivenessFailureThresholdToLivenessProbeFailureThresholdParam", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Container: api.ContainerParams{
				LivenessProbe: api.ProbeParams{
					FailureThreshold: 2,
				},
			},
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.Equal(t, 2, templateData.Container.Liveness.FailureThreshold)
	})

	t.Run("SetsLivenessSuccessThresholdToLivenessProbeSuccessThresholdParam", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Container: api.ContainerParams{
				LivenessProbe: api.ProbeParams{
					SuccessThreshold: 7,
				},
			},
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.Equal(t, 7, templateData.Container.Liveness.SuccessThreshold)
	})

	t.Run("SetsReadinessPathToReadinessProbePathParam", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Container: api.ContainerParams{
				ReadinessProbe: api.ProbeParams{
					Path: "/readiness",
				},
			},
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.Equal(t, "/readiness", templateData.Container.Readiness.Path)
	})

	t.Run("SetsReadinessPortToReadinessProbePortParam", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Container: api.ContainerParams{
				ReadinessProbe: api.ProbeParams{
					Port: 5002,
				},
			},
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.Equal(t, 5002, templateData.Container.Readiness.Port)
	})

	t.Run("SetsReadinessInitialDelaySecondsToReadinessProbeInitialDelaySecondsParam", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Container: api.ContainerParams{
				ReadinessProbe: api.ProbeParams{
					InitialDelaySeconds: 30,
				},
			},
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.Equal(t, 30, templateData.Container.Readiness.InitialDelaySeconds)
	})

	t.Run("SetsReadinessTimeoutSecondsToReadinessProbeTimeoutSecondsParam", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Container: api.ContainerParams{
				ReadinessProbe: api.ProbeParams{
					TimeoutSeconds: 1,
				},
			},
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.Equal(t, 1, templateData.Container.Readiness.TimeoutSeconds)
	})

	t.Run("SetsReadinessFailureThresholdToReadinessProbeFailureThresholdParam", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Container: api.ContainerParams{
				ReadinessProbe: api.ProbeParams{
					FailureThreshold: 6,
				},
			},
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.Equal(t, 6, templateData.Container.Readiness.FailureThreshold)
	})

	t.Run("SetsReadinessSuccessThresholdToReadinessProbeSuccessThresholdParam", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Container: api.ContainerParams{
				ReadinessProbe: api.ProbeParams{
					SuccessThreshold: 3,
				},
			},
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.Equal(t, 3, templateData.Container.Readiness.SuccessThreshold)
	})

	t.Run("SetsEnvironmentVariablesToContainerEnvironmentVariablesParam", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Container: api.ContainerParams{
				EnvironmentVariables: map[string]interface{}{
					"MY_CUSTOM_ENV":       "value1",
					"MY_OTHER_CUSTOM_ENV": "value2",
				},
			},
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.Equal(t, "value1", templateData.Container.EnvironmentVariables["MY_CUSTOM_ENV"])
		assert.Equal(t, "value2", templateData.Container.EnvironmentVariables["MY_OTHER_CUSTOM_ENV"])
	})

	t.Run("AddsJaegerServiceNameToEnvironmentVariables", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			App: "my-app",
			Container: api.ContainerParams{
				EnvironmentVariables: map[string]interface{}{},
			},
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.Equal(t, "my-app", templateData.Container.EnvironmentVariables["JAEGER_SERVICE_NAME"])
	})

	t.Run("SetsMetricsPathToMetricsProbePathParam", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Container: api.ContainerParams{
				Metrics: api.MetricsParams{
					Path: "/readiness",
				},
			},
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.Equal(t, "/readiness", templateData.Container.Metrics.Path)
	})

	t.Run("SetsMetricsPortToMetricsPortParam", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Container: api.ContainerParams{
				Metrics: api.MetricsParams{
					Port: 3080,
				},
			},
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.Equal(t, 3080, templateData.Container.Metrics.Port)
	})

	t.Run("SetsMetricsScrapeToMetricsScrapeParam", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Container: api.ContainerParams{
				Metrics: api.MetricsParams{
					Scrape: &trueValue,
				},
			},
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.Equal(t, true, templateData.Container.Metrics.Scrape)
	})

	t.Run("SetsUseLifecyclePreStopSleepCommandToLifecyclePrestopSleepParam", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Container: api.ContainerParams{
				Lifecycle: api.LifecycleParams{
					PrestopSleep: &trueValue,
				},
			},
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.Equal(t, true, templateData.Container.UseLifecyclePreStopSleepCommand)
	})

	t.Run("SetsPreStopSleepSecondsToLifecyclePrestopSleepSecondsParam", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		sleepValue := 25

		params := api.Params{
			Container: api.ContainerParams{
				Lifecycle: api.LifecycleParams{
					PrestopSleepSeconds: &sleepValue,
				},
			},
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.Equal(t, 25, templateData.Container.PreStopSleepSeconds)
	})

	t.Run("SidecarAddedToSidecarsCollection", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Sidecars: []*api.SidecarParams{
				{
					Type: api.SidecarTypeOpenresty,
				},
			},
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.Equal(t, 1, len(templateData.Sidecars))
	})

	t.Run("SetsSidecarTypeToSidecarTypeAsString", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Sidecars: []*api.SidecarParams{
				{
					Type: api.SidecarTypeOpenresty,
				},
			},
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.Equal(t, "openresty", templateData.Sidecars[0].Type)
	})

	t.Run("SetsSidecarImageToSidecarImageParam", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Sidecars: []*api.SidecarParams{
				&api.SidecarParams{
					Image: "estafette/openresty-sidecar:1.13.6.1-alpine",
				},
			},
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.Equal(t, "estafette/openresty-sidecar:1.13.6.1-alpine", templateData.Sidecars[0].Image)
	})

	t.Run("SetsSidecarHealthCheckPathToSidecarHealthCheckPathParam", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Sidecars: []*api.SidecarParams{
				&api.SidecarParams{
					HealthCheckPath: "/readiness",
				},
			},
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.Equal(t, "/readiness", templateData.Sidecars[0].SidecarSpecificProperties["healthcheckpath"])
	})

	t.Run("SetsSidecarCPURequestToSidecarCPURequestParam", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Sidecars: []*api.SidecarParams{
				&api.SidecarParams{
					CPU: api.CPUParams{
						Request: "1200m",
					},
				},
			},
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.Equal(t, "1200m", templateData.Sidecars[0].CPURequest)
	})

	t.Run("SetsSidecarCPULimitToSidecarCPULimitParam", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Sidecars: []*api.SidecarParams{
				&api.SidecarParams{
					CPU: api.CPUParams{
						Limit: "1500m",
					},
				},
			},
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.Equal(t, "1500m", templateData.Sidecars[0].CPULimit)
	})

	t.Run("SetsSidecarMemoryRequestToSidecarMemoryRequestParam", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Sidecars: []*api.SidecarParams{
				&api.SidecarParams{
					Memory: api.MemoryParams{
						Request: "1024Mi",
					},
				},
			},
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.Equal(t, "1024Mi", templateData.Sidecars[0].MemoryRequest)
	})

	t.Run("SetsSidecarMemoryLimitToSidecarMemoryLimitParam", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Sidecars: []*api.SidecarParams{
				&api.SidecarParams{
					Memory: api.MemoryParams{
						Limit: "2048Mi",
					},
				},
			},
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.Equal(t, "2048Mi", templateData.Sidecars[0].MemoryLimit)
	})

	t.Run("SetsSidecarEnvironmentVariablesToSidecarEnvironmentVariablesParam", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Sidecars: []*api.SidecarParams{
				&api.SidecarParams{
					EnvironmentVariables: map[string]interface{}{
						"MY_CUSTOM_ENV":       "value1",
						"MY_OTHER_CUSTOM_ENV": "value2",
					},
				},
			},
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		// assert.Equal(t, 2, len(templateData.Sidecar.EnvironmentVariables))
		assert.Equal(t, "value1", templateData.Sidecars[0].EnvironmentVariables["MY_CUSTOM_ENV"])
		assert.Equal(t, "value2", templateData.Sidecars[0].EnvironmentVariables["MY_OTHER_CUSTOM_ENV"])
	})

	t.Run("SetsCloudSQLProxySpecificArgsToSidecarSpecificProperties", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Sidecars: []*api.SidecarParams{
				&api.SidecarParams{
					HealthCheckPath:                   "testHealthCheckPath",
					DbInstanceConnectionName:          "testDbInstanceConnectionName",
					SQLProxyPort:                      15,
					SQLProxyTerminationTimeoutSeconds: 16,
				},
			},
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.Equal(t, 4, len(templateData.Sidecars[0].SidecarSpecificProperties))
		assert.Equal(t, "testHealthCheckPath", templateData.Sidecars[0].SidecarSpecificProperties["healthcheckpath"])
		assert.Equal(t, "testDbInstanceConnectionName", templateData.Sidecars[0].SidecarSpecificProperties["dbinstanceconnectionname"])
		assert.Equal(t, 15, templateData.Sidecars[0].SidecarSpecificProperties["sqlproxyport"])
		assert.Equal(t, 16, templateData.Sidecars[0].SidecarSpecificProperties["sqlproxyterminationtimeoutseconds"])
	})

	t.Run("SetsContainerLifecycle", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Container: api.ContainerParams{
				ContainerLifeCycle: map[string]interface{}{
					"preStop": map[string]interface{}{
						"exec": map[string]interface{}{
							"command": []string{
								"/bin/sh",
								"-c",
								"echo Hello from the preStop handler > /usr/share/message",
							},
						},
					},
				},
			},
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.Equal(t, map[string]interface{}{"preStop": map[string]interface{}{"exec": map[string]interface{}{"command": []string{"/bin/sh", "-c", "echo Hello from the preStop handler > /usr/share/message"}}}}, templateData.Container.ContainerLifeCycle)
	})

	t.Run("SetsSecretsToSecretsParam", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Secrets: api.SecretsParams{
				Keys: map[string]interface{}{
					"secret-file-1.json": "c29tZSBzZWNyZXQgdmFsdWU=",
					"secret-file-2.yaml": "YW5vdGhlciBzZWNyZXQgdmFsdWU=",
				},
			},
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.Equal(t, 2, len(templateData.Secrets))
		assert.Equal(t, "c29tZSBzZWNyZXQgdmFsdWU=", templateData.Secrets["secret-file-1.json"])
		assert.Equal(t, "YW5vdGhlciBzZWNyZXQgdmFsdWU=", templateData.Secrets["secret-file-2.yaml"])
	})

	t.Run("SetsMountApplicationSecretsToTrueIfSecretsParamLengthIsLargerThanZero", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Secrets: api.SecretsParams{
				Keys: map[string]interface{}{
					"secret-file-1.json": "c29tZSBzZWNyZXQgdmFsdWU=",
					"secret-file-2.yaml": "YW5vdGhlciBzZWNyZXQgdmFsdWU=",
				},
			},
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.True(t, templateData.MountApplicationSecrets)
	})

	t.Run("SetsMountApplicationSecretsToFalseIfSecretsParamLengthIsZero", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Secrets: api.SecretsParams{
				Keys: map[string]interface{}{},
			},
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.False(t, templateData.MountApplicationSecrets)
	})

	t.Run("SetsIngressPathToBasepathParam", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Basepath: "/",
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.Equal(t, "/", templateData.IngressPath)
	})

	t.Run("AppendSlashToIngressPathIfBasepathParamDoesNotEndInSlash", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Basepath: "/api",
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.Equal(t, "/api/", templateData.IngressPath)
	})

	t.Run("AppendSlashStarToIngressPathIfUseGCEIngressIsTrue", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Basepath:   "/api",
			Visibility: api.VisibilityIAP,
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.Equal(t, "/api/*", templateData.IngressPath)
	})

	t.Run("SetsInternalIngressPathToBasepathParam", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Basepath: "/",
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.Equal(t, "/", templateData.InternalIngressPath)
	})

	t.Run("AppendSlashToInternalIngressPathIfBasepathParamDoesNotEndInSlash", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Basepath: "/api",
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.Equal(t, "/api/", templateData.InternalIngressPath)
	})

	t.Run("DoNotAppendSlashStarToInternalIngressPathIfUseGCEIngressIsTrue", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Basepath:   "/api",
			Visibility: api.VisibilityIAP,
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.Equal(t, "/api/", templateData.InternalIngressPath)
	})

	t.Run("SetsMountPayloadLoggingToTrueIfEnablePayloadLoggingParamIsTrue", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			EnablePayloadLogging: true,
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.True(t, templateData.MountPayloadLogging)
	})

	t.Run("SetsMountPayloadLoggingToFalseIfEnablePayloadLoggingParamIsFalse", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			EnablePayloadLogging: false,
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.False(t, templateData.MountPayloadLogging)
	})

	t.Run("SetsAddSafeToEvictAnnotationToTrueIfEnablePayloadLoggingParamIsTrue", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			EnablePayloadLogging: true,
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.True(t, templateData.AddSafeToEvictAnnotation)
	})

	t.Run("SetsAddSafeToEvictAnnotationToFalseIfEnablePayloadLoggingParamIsFalse", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			EnablePayloadLogging: false,
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.False(t, templateData.AddSafeToEvictAnnotation)
	})

	t.Run("SetsRollingUpdateMaxSurgeToRollingUpdateMaxSurgeParam", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			RollingUpdate: api.RollingUpdateParams{
				MaxSurge: "25%",
			},
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.Equal(t, "25%", templateData.RollingUpdateMaxSurge)
	})

	t.Run("SetsRollingUpdateMaxSurgeToRollingUpdateMaxSurgeParam", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			RollingUpdate: api.RollingUpdateParams{
				MaxUnavailable: "15%",
			},
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.Equal(t, "15%", templateData.RollingUpdateMaxUnavailable)
	})

	t.Run("AddsBuildVersionLabelSetToBuildVersionParam", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			BuildVersion: "1.2.3",
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.Equal(t, "1.2.3", templateData.PodLabels["version"])
	})

	t.Run("SetsPreferPreemptiblesToTrueIfChaosProofParamIsTrue", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			ChaosProof: true,
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.True(t, templateData.PreferPreemptibles)
	})

	t.Run("SetsPreferPreemptiblesToFalseIfChaosProofParamIsFalse", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			ChaosProof: false,
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.False(t, templateData.PreferPreemptibles)
	})

	t.Run("SetsHasTolerationsToTrueIfChaosProofParamIsTrue", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			ChaosProof: true,
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.True(t, templateData.HasTolerations)
	})

	t.Run("AddsPreemptibleTolerationToTolerationsIfChaosProofParamIsTrue", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			ChaosProof: true,
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.Equal(t, 1, len(templateData.Tolerations))
		assert.Equal(t, &map[string]interface{}{
			"key":      "cloud.google.com/gke-preemptible",
			"operator": "Equal",
			"value":    "true",
			"effect":   "NoSchedule",
		}, templateData.Tolerations[0])
	})

	t.Run("AddsPreemptibleTolerationAndOtherTolerationsIfChaosProofParamIsTrueAndTolerationsAreSet", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			ChaosProof: true,
			Tolerations: []*map[string]interface{}{
				{
					"key":      "role",
					"operator": "Equal",
					"value":    "tooling",
					"effect":   "NoSchedule",
				},
			},
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.Equal(t, 2, len(templateData.Tolerations))
		assert.Equal(t, &map[string]interface{}{
			"key":      "cloud.google.com/gke-preemptible",
			"operator": "Equal",
			"value":    "true",
			"effect":   "NoSchedule",
		}, templateData.Tolerations[0])
		assert.Equal(t, &map[string]interface{}{
			"key":      "role",
			"operator": "Equal",
			"value":    "tooling",
			"effect":   "NoSchedule",
		}, templateData.Tolerations[1])
	})

	t.Run("SetsMountConfigmapToTrueIfConfigFilesParamsLengthIsLargerThanZero", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Configs: api.ConfigsParams{
				Files: []string{
					"config.json",
				},
			},
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.True(t, templateData.MountConfigmap)
	})

	t.Run("SetsMountConfigmapToTrueIfInlineFilesParamsLengthIsLargerThanZero", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Configs: api.ConfigsParams{
				InlineFiles: map[string]string{
					"inline-config.properties": "enemies=aliens\nlives=3",
				},
			},
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.True(t, templateData.MountConfigmap)
	})

	t.Run("SetsMountConfigmapToFalseIfConfigFilesAndInlineFilesParamsLengthAreZero", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Configs: api.ConfigsParams{
				Files:       []string{},
				InlineFiles: map[string]string{},
			},
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.False(t, templateData.MountConfigmap)
	})

	t.Run("SetsConfigMountPathToConfigMountPathParam", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Configs: api.ConfigsParams{
				MountPath: "/configs",
			},
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.Equal(t, "/configs", templateData.ConfigMountPath)
	})

	t.Run("SetsSecretMountPathToSecretMountPathParam", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Secrets: api.SecretsParams{
				MountPath: "/secrets",
			},
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.Equal(t, "/secrets", templateData.SecretMountPath)
	})

	t.Run("SetsLimitTrustedIPRangesIfVisibilityParamIsPublic", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Visibility: api.VisibilityPublic,
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.True(t, templateData.LimitTrustedIPRanges)
	})

	t.Run("SetsLimitTrustedIPRangesToTrueIfVisibilityParamIsPublic", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Visibility: api.VisibilityPublic,
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.True(t, templateData.LimitTrustedIPRanges)
	})

	t.Run("SetsLimitTrustedIPRangesToFalseIfVisibilityParamIsIap", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Visibility: api.VisibilityIAP,
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.False(t, templateData.LimitTrustedIPRanges)
	})

	t.Run("SetsLimitTrustedIPRangesToFalseIfVisibilityParamIsPrivate", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Visibility: api.VisibilityPrivate,
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.False(t, templateData.LimitTrustedIPRanges)
	})

	t.Run("SetsTrustedIPRangesToTrustedIPRangesParams", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
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
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.Equal(t, 14, len(templateData.TrustedIPRanges))
	})

	t.Run("SetsLocalManifestDataToAllLocalManifestDataCombined", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Manifests: api.ManifestsParams{
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
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.Equal(t, 3, len(templateData.ManifestData))
		assert.Equal(t, "value 1", templateData.ManifestData["property1"])
		assert.Equal(t, "value 2", templateData.ManifestData["property2"])
		assert.Equal(t, "value 3", templateData.ManifestData["property3"])
	})

	t.Run("AppendsCanaryToNameWithTrackIfParamsTypeIsCanary", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			App:    "myapp",
			Action: api.ActionDeployCanary,
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.Equal(t, "myapp-canary", templateData.NameWithTrack)
	})

	t.Run("AppendsStableToNameWithTrackIfParamsTypeIsRollforward", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			App:    "myapp",
			Action: api.ActionDeployStable,
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.Equal(t, "myapp-stable", templateData.NameWithTrack)
	})

	t.Run("DoesNotAppendTrackToNameWithTrackIfParamsTypeIsSimple", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			App:    "myapp",
			Action: api.ActionDeploySimple,
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.Equal(t, "myapp", templateData.NameWithTrack)
	})

	t.Run("DoesNotAddReleaseIDLabelToPodLabelsIfReleaseIDIsEmpty", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			App:    "myapp",
			Action: api.ActionDeploySimple,
		}
		releaseID := ""

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", releaseID, "")

		assert.Equal(t, "", templateData.PodLabels["estafette.io/release-id"])
	})

	t.Run("AddsReleaseIDLabelToPodLabelsIfReleaseIDIsNotEmpty", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			App:    "myapp",
			Action: api.ActionDeploySimple,
		}
		releaseID := "1"

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", releaseID, "")

		assert.Equal(t, "1", templateData.PodLabels["estafette.io/release-id"])
	})

	t.Run("SetsIncludeTriggeredByLabelToFalseIfTriggeredByIsEmpty", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			App:    "myapp",
			Action: api.ActionDeploySimple,
		}
		triggeredBy := ""

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", triggeredBy)

		assert.Equal(t, "", templateData.PodLabels["estafette.io/triggered-by"])
	})

	t.Run("SetsIncludeTriggeredByLabelToTrueIfTriggeredByIsNotEmpty", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			App:    "myapp",
			Action: api.ActionDeploySimple,
		}
		triggeredBy := "user@estafette.io"

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", triggeredBy)

		assert.Equal(t, "user-at-estafette.io", templateData.PodLabels["estafette.io/triggered-by"])
	})

	t.Run("SetsIncludeTrackLabelToFalseIfParamsTypeIsSimple", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			App:    "myapp",
			Action: api.ActionDeploySimple,
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.False(t, templateData.IncludeTrackLabel)
	})

	t.Run("SetsIncludeTrackLabelToTrueIfParamsTypeIsCanary", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			App:    "myapp",
			Action: api.ActionDeployCanary,
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.True(t, templateData.IncludeTrackLabel)
	})

	t.Run("SetsIncludeTrackLabelToTrueIfParamsTypeIsRollforward", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			App:    "myapp",
			Action: api.ActionDeployStable,
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.True(t, templateData.IncludeTrackLabel)
	})

	t.Run("SetsTrackLabelToCanaryIfParamsTypeIsCanary", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			App:    "myapp",
			Action: api.ActionDeployCanary,
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.Equal(t, "canary", templateData.TrackLabel)
	})

	t.Run("SetsTrackLabelToStableIfParamsTypeIsRollforward", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			App:    "myapp",
			Action: api.ActionDeployStable,
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.Equal(t, "stable", templateData.TrackLabel)
	})

	t.Run("SetsAdditionalVolumeMountsToVolumeMountsParam", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			VolumeMounts: []api.VolumeMountParams{
				api.VolumeMountParams{
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
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.Equal(t, 1, len(templateData.AdditionalVolumeMounts))
		assert.Equal(t, "client-certs", templateData.AdditionalVolumeMounts[0].Name)
		assert.Equal(t, "/cockroach-certs", templateData.AdditionalVolumeMounts[0].MountPath)
		assert.Equal(t, "secret:\n  items:\n  - key: key\n    mode: 384\n    path: key\n  - key: cert\n    path: cert\n  secretName: estafette.client.estafette\n", templateData.AdditionalVolumeMounts[0].VolumeYAML)
	})

	t.Run("SetsAdditionalContainerPortsToContainerAdditionalPortsParam", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Container: api.ContainerParams{
				AdditionalPorts: []*api.AdditionalPortParams{
					&api.AdditionalPortParams{
						Name:       "grpc",
						Port:       8085,
						Protocol:   "TCP",
						Visibility: api.VisibilityPrivate,
					},
					&api.AdditionalPortParams{
						Name:       "grpc",
						Port:       8085,
						Protocol:   "UDP",
						Visibility: api.VisibilityPublic,
					},
				},
			},
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.Equal(t, 2, len(templateData.AdditionalContainerPorts))
	})

	t.Run("SetsAdditionalServicePortsToContainerAdditionalPortsParamForPortsWithVisibilityEqualToVisibilityParam", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Visibility: api.VisibilityPrivate,
			Container: api.ContainerParams{
				AdditionalPorts: []*api.AdditionalPortParams{
					&api.AdditionalPortParams{
						Name:       "grpc",
						Port:       8085,
						Protocol:   "TCP",
						Visibility: api.VisibilityPrivate,
					},
					&api.AdditionalPortParams{
						Name:       "snmp",
						Port:       8086,
						Protocol:   "UDP",
						Visibility: api.VisibilityPublic,
					},
				},
			},
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.Equal(t, 1, len(templateData.AdditionalServicePorts))
		assert.Equal(t, "grpc", templateData.AdditionalServicePorts[0].Name)
		assert.Equal(t, 8085, templateData.AdditionalServicePorts[0].Port)
		assert.Equal(t, "TCP", templateData.AdditionalServicePorts[0].Protocol)
	})

	t.Run("SetsOverrideDefaultWhitelistToTrueIfVisibilityEqualsPublicWhitelistAndWhitelistedIPSHasOneOrMoreItems", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Visibility: api.VisibilityPublicWhitelist,
			WhitelistedIPS: []string{
				"0.0.0.0/0",
			},
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.True(t, templateData.OverrideDefaultWhitelist)
	})

	t.Run("SetsOverrideDefaultWhitelistToFalseIfVisibilityEqualsPublicWhitelistButWhitelistedIPSHasNoItems", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Visibility:     api.VisibilityPublicWhitelist,
			WhitelistedIPS: []string{},
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.False(t, templateData.OverrideDefaultWhitelist)
	})

	t.Run("SetsOverrideDefaultWhitelistToFalseIfVisibilityIsPrivate", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Visibility: api.VisibilityPrivate,
			WhitelistedIPS: []string{
				"0.0.0.0/0",
			},
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.False(t, templateData.OverrideDefaultWhitelist)
	})

	t.Run("SetsOverrideDefaultWhitelistToFalseIfVisibilityIsPrivate", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Visibility: api.VisibilityIAP,
			WhitelistedIPS: []string{
				"0.0.0.0/0",
			},
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.False(t, templateData.OverrideDefaultWhitelist)
	})

	t.Run("SetsOverrideDefaultWhitelistToFalseIfVisibilityIsPrivate", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Visibility: api.VisibilityPublic,
			WhitelistedIPS: []string{
				"0.0.0.0/0",
			},
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.False(t, templateData.OverrideDefaultWhitelist)
	})

	t.Run("SetsNginxIngressWhitelistToCommaSeparatedJoingOfWhitelistedIPS", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Visibility: api.VisibilityPublicWhitelist,
			WhitelistedIPS: []string{
				"10.0.0.0/8",
				"172.16.0.0/12",
				"192.168.0.0/16",
			},
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.Equal(t, "10.0.0.0/8,172.16.0.0/12,192.168.0.0/16", templateData.NginxIngressWhitelist)
	})

	t.Run("SetsIncludeReplicasToTrueIfCurrentReplicasIsGreaterThanZero", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{}

		// act
		templateData := service.GenerateTemplateData(params, 1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.True(t, templateData.IncludeReplicas)
	})

	t.Run("SetsIncludeReplicasToFalseIfCurrentReplicasIsZeroOrLess", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{}

		// act
		templateData := service.GenerateTemplateData(params, 0, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.False(t, templateData.IncludeReplicas)
	})

	t.Run("SetsReplicasToCurrentReplicasIfGreaterThanZero", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{}

		// act
		templateData := service.GenerateTemplateData(params, 15, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.Equal(t, 15, templateData.Replicas)
	})

	t.Run("SetsReplicasToReplicasParamIfCurrentReplicasIsZeroOrLessAndAutoscaleIsDisabled", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Replicas: 1,
			Autoscale: api.AutoscaleParams{
				Enabled:     &falseValue,
				MinReplicas: 3,
			},
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.Equal(t, 1, templateData.Replicas)
	})

	t.Run("SetsReplicasToReplicasParamIfCurrentReplicasIsZeroOrLessAndReplicasParamIsLargerThanMinReplicas", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Replicas: 5,
			Autoscale: api.AutoscaleParams{
				Enabled:     &falseValue,
				MinReplicas: 3,
			},
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.Equal(t, 5, templateData.Replicas)
	})

	t.Run("SetsReplicasToMinReplicasIfCurrentReplicasIsZeroOrLess", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Autoscale: api.AutoscaleParams{
				MinReplicas: 3,
			},
		}

		// act
		templateData := service.GenerateTemplateData(params, 0, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.Equal(t, 3, templateData.Replicas)
	})

	t.Run("SetsScheduleToScheduleParam", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Schedule: "*/5 * * * *",
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.Equal(t, "*/5 * * * *", templateData.Schedule)
	})

	t.Run("SetsUseHpaScalerToAutoscalerSafetyEnabledParam", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Autoscale: api.AutoscaleParams{
				Safety: api.AutoscaleSafetyParams{
					Enabled: true,
				},
			},
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.True(t, templateData.UseHpaScaler)
	})

	t.Run("SetsHpaScalerPromQueryToAutoscalerSafetyPromQueryParam", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Autoscale: api.AutoscaleParams{
				Safety: api.AutoscaleSafetyParams{
					PromQuery: "sum(rate(nginx_http_requests_total{app='my-app'}[5m])) by (app)",
				},
			},
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.Equal(t, "sum(rate(nginx_http_requests_total{app='my-app'}[5m])) by (app)", templateData.HpaScalerPromQuery)
	})

	t.Run("SetsHpaScalerRequestsPerReplicaToAutoscalerSafetyPromQueryParam", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Autoscale: api.AutoscaleParams{
				Safety: api.AutoscaleSafetyParams{
					Ratio: "0.25",
				},
			},
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.Equal(t, "0.25", templateData.HpaScalerRequestsPerReplica)
	})

	t.Run("SetsHpaScalerRequestsPerReplicaToAutoscalerSafetyPromQueryParam", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Autoscale: api.AutoscaleParams{
				Safety: api.AutoscaleSafetyParams{
					Delta: "-2.7584",
				},
			},
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.Equal(t, "-2.7584", templateData.HpaScalerDelta)
	})

	t.Run("SetsHpaScalerRequestsPerReplicaToAutoscalerSafetyPromQueryParam", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Autoscale: api.AutoscaleParams{
				Safety: api.AutoscaleSafetyParams{
					ScaleDownRatio: "0.2",
				},
			},
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.Equal(t, "0.2", templateData.HpaScalerScaleDownMaxRatio)
	})

	t.Run("SetsAllHostsToHostsAndInternalHostsAppended", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Hosts: []string{
				"ci.estafette.io",
			},
			InternalHosts: []string{
				"ci.internal.estafette.io",
			},
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.Equal(t, 2, len(templateData.AllHosts))
		assert.Equal(t, "ci.estafette.io", templateData.AllHosts[0])
		assert.Equal(t, "ci.internal.estafette.io", templateData.AllHosts[1])
	})

	t.Run("SetsAllHostsJoinedToHostsAndInternalHostsAppendedSeparatedByComma", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Hosts: []string{
				"ci.estafette.io",
			},
			InternalHosts: []string{
				"ci.internal.estafette.io",
			},
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.Equal(t, "ci.estafette.io,ci.internal.estafette.io", templateData.AllHostsJoined)
	})

	t.Run("SetsAllHostsToHostsAndInternalHostsAppendedWhenOnlyHostsAreSet", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Hosts: []string{
				"ci.estafette.io",
			},
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.Equal(t, 1, len(templateData.AllHosts))
		assert.Equal(t, "ci.estafette.io", templateData.AllHosts[0])
	})

	t.Run("SetsAllHostsJoinedToHostsAndInternalHostsAppendedSeparatedByCommaWhenOnlyHostsAreSet", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Hosts: []string{
				"ci.estafette.io",
			},
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.Equal(t, "ci.estafette.io", templateData.AllHostsJoined)
	})

	t.Run("SetsAllHostsToHostsAndInternalHostsAppendedWhenOnlyInternalHostsAreSet", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			InternalHosts: []string{
				"ci.internal.estafette.io",
			},
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.Equal(t, 1, len(templateData.AllHosts))
		assert.Equal(t, "ci.internal.estafette.io", templateData.AllHosts[0])
	})

	t.Run("SetsAllHostsJoinedToHostsAndInternalHostsAppendedSeparatedByCommaWhenOnlyInternalHostsAreSet", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			InternalHosts: []string{
				"ci.internal.estafette.io",
			},
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.Equal(t, "ci.internal.estafette.io", templateData.AllHostsJoined)
	})

	t.Run("SetsNginxIngressProxyConnectTimeoutToRequestTimeoutParam", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Request: api.RequestParams{
				Timeout: "75",
			},
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.Equal(t, 75, templateData.NginxIngressProxyConnectTimeout)
	})

	t.Run("SetsNginxIngressProxyConnectTimeoutToRequestTimeoutParamWithSecondSuffix", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Request: api.RequestParams{
				Timeout: "75s",
			},
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.Equal(t, 75, templateData.NginxIngressProxyConnectTimeout)
	})

	t.Run("SetsNginxIngressProxyConnectTimeoutTo75IfRequestTimeoutParamIsLargerThan75", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Request: api.RequestParams{
				Timeout: "180s",
			},
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.Equal(t, 75, templateData.NginxIngressProxyConnectTimeout)
	})

	t.Run("SetsNginxIngressProxySendTimeoutToRequestTimeoutParam", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Request: api.RequestParams{
				Timeout: "300",
			},
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.Equal(t, 300, templateData.NginxIngressProxySendTimeout)
	})

	t.Run("SetsNginxIngressProxySendTimeoutToRequestTimeoutParamWithSecondSuffix", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Request: api.RequestParams{
				Timeout: "300s",
			},
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.Equal(t, 300, templateData.NginxIngressProxySendTimeout)
	})

	t.Run("SetsNginxIngressProxyReadTimeoutToRequestTimeoutParam", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Request: api.RequestParams{
				Timeout: "300",
			},
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.Equal(t, 300, templateData.NginxIngressProxyReadTimeout)
	})

	t.Run("SetsNginxIngressProxyReadTimeoutToRequestTimeoutParamWithSecondSuffix", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Request: api.RequestParams{
				Timeout: "300s",
			},
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.Equal(t, 300, templateData.NginxIngressProxyReadTimeout)
	})

	t.Run("SetsCertificateSecret", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			CertificateSecret: "shared-wildcard-letsencrypt-certificate",
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.True(t, templateData.UseCertificateSecret)
	})

	t.Run("SetsServiceTypeToClusterIPIfVisibilityParamIsApigee", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Visibility: api.VisibilityApigee,
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.Equal(t, "ClusterIP", templateData.ServiceType)
	})

	t.Run("SetsUseCloudflareProxyToTrueIfVisibilityParamIsApigee", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Visibility: api.VisibilityApigee,
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.Equal(t, true, templateData.UseCloudflareProxy)
	})

	t.Run("SetsUseGCEIngressToFalseIfVisibilityParamIsApigee", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Visibility: api.VisibilityApigee,
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.Equal(t, false, templateData.UseGCEIngress)
	})

	t.Run("SetsNginxAuthTLSSecretIfVisibilityParamIsApigee", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Visibility: api.VisibilityApigee,
			Request: api.RequestParams{
				AuthSecret: "protected/some-secret",
			},
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.Equal(t, "protected/some-secret", templateData.NginxAuthTLSSecret)
	})

	t.Run("SetsVerifyDepthIfVisibilityParamIsApigee", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Visibility: api.VisibilityApigee,
			Request: api.RequestParams{
				VerifyDepth: 5,
			},
		}

		// act
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.Equal(t, 5, templateData.NginxAuthTLSVerifyDepth)
	})

	t.Run("SetsApigeeHostsIfVisibilityParamIsApigee", func(t *testing.T) {

		ctx := context.Background()
		service, err := NewService(ctx)
		assert.Nil(t, err)

		params := api.Params{
			Visibility: api.VisibilityApigee,
			Hosts:      []string{"google.com", "estafette.io", "test-app"},
		}

		// act
		params.SetDefaults("github.com", "estafette", "estafette-extension-gke", "sample-app", "0.1.0", "test", "deploy", "", nil)
		templateData := service.GenerateTemplateData(params, -1, "github.com", "estafette", "estafette-extension-gke", "master", "02770946ad015b34da9e9980007bf81308c41aec", "", "")

		assert.Equal(t, []string{"google-apigee.com", "estafette-apigee.io", "test-app-apigee"}, templateData.ApigeeHosts)
		assert.Equal(t, "google-apigee.com,estafette-apigee.io,test-app-apigee", templateData.ApigeeHostsJoined)
	})
}
