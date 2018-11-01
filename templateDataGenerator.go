package main

import (
	"strings"
)

func generateTemplateData(params Params) TemplateData {

	data := TemplateData{
		Name:             params.App,
		Namespace:        params.Namespace,
		Labels:           params.Labels,
		AppLabelSelector: params.App,

		Hosts:       params.Hosts,
		HostsJoined: strings.Join(params.Hosts, ","),

		MinReplicas:         params.Autoscale.MinReplicas,
		MaxReplicas:         params.Autoscale.MaxReplicas,
		TargetCPUPercentage: params.Autoscale.CPUPercentage,

		Secrets: params.Secrets,

		// IngressPath         string
		// PreferPreemptibles  bool

		Container: ContainerData{
			Repository: params.Container.ImageRepository,
			Name:       params.Container.ImageName,
			Tag:        params.Container.ImageTag,
			Port:       params.Container.Port,

			CPURequest:    params.Container.CPU.Request,
			CPULimit:      params.Container.CPU.Limit,
			MemoryRequest: params.Container.Memory.Request,
			MemoryLimit:   params.Container.Memory.Limit,

			EnvironmentVariables: params.Container.EnvironmentVariables,

			Liveness: ProbeData{
				Path:                params.Container.LivenessProbe.Path,
				InitialDelaySeconds: params.Container.LivenessProbe.InitialDelaySeconds,
				TimeoutSeconds:      params.Container.LivenessProbe.TimeoutSeconds,
			},
			Readiness: ProbeData{
				Path:                params.Container.ReadinessProbe.Path,
				InitialDelaySeconds: params.Container.ReadinessProbe.InitialDelaySeconds,
				TimeoutSeconds:      params.Container.ReadinessProbe.TimeoutSeconds,
			},
			Metrics: MetricsData{
				Scrape: params.Container.Metrics.Scrape,
				Path:   params.Container.Metrics.Path,
				Port:   params.Container.Metrics.Port,
			},
		},

		Sidecar: SidecarData{
			UseOpenrestySidecar: params.Sidecar.Type == "openresty",

			Image:         params.Sidecar.Image,
			CPURequest:    params.Sidecar.CPU.Request,
			CPULimit:      params.Sidecar.CPU.Limit,
			MemoryRequest: params.Sidecar.Memory.Request,
			MemoryLimit:   params.Sidecar.Memory.Limit,

			EnvironmentVariables: params.Sidecar.EnvironmentVariables,
		},
	}

	if params.Visibility == "private" {
		data.ServiceType = "ClusterIP"
		data.UseNginxIngress = true
		data.UseDNSAnnotationsOnIngress = true
		data.UseDNSAnnotationsOnService = false
	} else if params.Visibility == "public" {
		data.ServiceType = "LoadBalancer"
		data.UseNginxIngress = false
		data.UseDNSAnnotationsOnIngress = false
		data.UseDNSAnnotationsOnService = true
	}

	return data
}
