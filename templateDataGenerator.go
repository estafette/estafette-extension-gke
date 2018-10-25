package main

func generateTemplateData(params Params) TemplateData {

	data := TemplateData{
		Name:             params.App,
		Namespace:        params.Namespace,
		Labels:           params.Labels,
		AppLabelSelector: params.App,

		// Hosts               []string
		// HostsJoined         string
		// IngressPath         string
		// UseNginxIngress     bool
		// UseGCEIngress       bool
		// ServiceType         string
		// MinReplicas         int
		// MaxReplicas         int
		// TargetCPUPercentage int
		// PreferPreemptibles  bool

		Container: ContainerData{
			Repository: params.ImageRepository,
			Name:       params.ImageName,
			Tag:        params.ImageTag,

			// CPURequest    string
			// MemoryRequest string
			// CPULimit      string
			// MemoryLimit   string

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

	return data
}
