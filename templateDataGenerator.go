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
		// Container           ContainerData

	}

	return data
}
