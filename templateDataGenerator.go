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

		// IngressPath         string
		// UseNginxIngress     bool
		// UseGCEIngress       bool
		// ServiceType         string
		// MinReplicas         int
		// MaxReplicas         int
		// TargetCPUPercentage int
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
				// Path                string
				// InitialDelaySeconds int
				// TimeoutSeconds      int
			},
			Readiness: ProbeData{
				// Path                string
				// InitialDelaySeconds int
				// TimeoutSeconds      int
			},
		},
	}

	if params.Visibility == "private" {
		data.ServiceType = "ClusterIP"
	}

	return data
}
