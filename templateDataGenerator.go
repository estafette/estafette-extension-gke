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

		// IngressPath         string
		// UseNginxIngress     bool
		// UseGCEIngress       bool
		// PreferPreemptibles  bool

		Container: ContainerData{
			Repository: params.Container.ImageRepository,
			Name:       params.Container.ImageName,
			Tag:        params.Container.ImageTag,
			Port:       params.Container.Port,

			CPURequest:    params.CPU.Request,
			CPULimit:      params.CPU.Limit,
			MemoryRequest: params.Memory.Request,
			MemoryLimit:   params.Memory.Limit,

			Liveness: ProbeData{
				Path:                params.LivenessProbe.Path,
				InitialDelaySeconds: params.LivenessProbe.InitialDelaySeconds,
				TimeoutSeconds:      params.LivenessProbe.TimeoutSeconds,
			},
			Readiness: ProbeData{
				Path:                params.ReadinessProbe.Path,
				InitialDelaySeconds: params.ReadinessProbe.InitialDelaySeconds,
				TimeoutSeconds:      params.ReadinessProbe.TimeoutSeconds,
			},
		},
	}

	if params.Visibility == "private" {
		data.ServiceType = "ClusterIP"
		data.UseNginxIngress = true
	} else if params.Visibility == "public" {
		data.ServiceType = "LoadBalancer"
		data.UseNginxIngress = false
	}

	return data
}
