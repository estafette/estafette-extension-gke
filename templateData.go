package main

// TemplateData contains the root data for rendering the Kubernetes manifests
type TemplateData struct {
	Name                       string
	Namespace                  string
	Labels                     map[string]string
	AppLabelSelector           string
	Hosts                      []string
	HostsJoined                string
	IngressPath                string
	UseIngress                 bool
	UseNginxIngress            bool
	UseGCEIngress              bool
	UseDNSAnnotationsOnIngress bool
	UseDNSAnnotationsOnService bool
	ServiceType                string
	MinReplicas                int
	MaxReplicas                int
	TargetCPUPercentage        int
	PreferPreemptibles         bool
	Container                  ContainerData
}

// ContainerData has data specific to the application container
type ContainerData struct {
	Repository           string
	Name                 string
	Tag                  string
	CPURequest           string
	MemoryRequest        string
	CPULimit             string
	MemoryLimit          string
	Port                 int
	EnvironmentVariables map[string]string
	Liveness             ProbeData
	Readiness            ProbeData
	Metrics              MetricsData
}

// ProbeData has data specific to liveness and readiness probes
type ProbeData struct {
	Path                string
	InitialDelaySeconds int
	TimeoutSeconds      int
}

// MetricsData has data to configure prometheus metrics scraping
type MetricsData struct {
	Path string
	Port int
}
