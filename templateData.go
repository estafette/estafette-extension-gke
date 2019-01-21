package main

// TemplateData contains the root data for rendering the Kubernetes manifests
type TemplateData struct {
	Name                        string
	NameWithTrack               string
	Namespace                   string
	Labels                      map[string]string
	AppLabelSelector            string
	Hosts                       []string
	HostsJoined                 string
	IngressPath                 string
	UseIngress                  bool
	UseNginxIngress             bool
	UseGCEIngress               bool
	UseDNSAnnotationsOnIngress  bool
	UseDNSAnnotationsOnService  bool
	ServiceType                 string
	MinReplicas                 int
	MaxReplicas                 int
	TargetCPUPercentage         int
	PreferPreemptibles          bool
	Container                   ContainerData
	Sidecar                     SidecarData
	MountApplicationSecrets     bool
	Secrets                     map[string]interface{}
	SecretMountPath             string
	MountConfigmap              bool
	ConfigmapFiles              map[string]string
	ConfigMountPath             string
	MountPayloadLogging         bool
	RollingUpdateMaxSurge       string
	RollingUpdateMaxUnavailable string
	BuildVersion                string
	LimitTrustedIPRanges        bool
	TrustedIPRanges             []string
	ManifestData                map[string]interface{}
	IncludeTrackLabel           bool
	IncludeReleaseIDLabel       bool
	ReleaseIDLabel              string
	TrackLabel                  string
	AddSafeToEvictAnnotation    bool
	AdditionalVolumeMounts      []VolumeMountData
	AdditionalContainerPorts    []AdditionalPortData
	AdditionalServicePorts      []AdditionalPortData
	OverrideDefaultWhitelist    bool
	NginxIngressWhitelist       string
	IncludeReplicas             bool
	Replicas                    int
}

// ContainerData has data specific to the application container
type ContainerData struct {
	Repository                      string
	Name                            string
	Tag                             string
	CPURequest                      string
	MemoryRequest                   string
	CPULimit                        string
	MemoryLimit                     string
	Port                            int
	EnvironmentVariables            map[string]interface{}
	Liveness                        ProbeData
	Readiness                       ProbeData
	Metrics                         MetricsData
	UseLifecyclePreStopSleepCommand bool
	PreStopSleepSeconds             int
}

// ProbeData has data specific to liveness and readiness probes
type ProbeData struct {
	Path                string
	Port                int
	InitialDelaySeconds int
	TimeoutSeconds      int
	IncludeOnContainer  bool
}

// MetricsData has data to configure prometheus metrics scraping
type MetricsData struct {
	Scrape bool
	Path   string
	Port   int
}

// SidecarData configures the injected sidecar
type SidecarData struct {
	UseOpenrestySidecar  bool
	Image                string
	HealthCheckPath      string
	EnvironmentVariables map[string]interface{}
	CPURequest           string
	MemoryRequest        string
	CPULimit             string
	MemoryLimit          string
}

// VolumeMountData configures additional volume mounts for shared secrets, existing volumes, etc
type VolumeMountData struct {
	Name       string
	MountPath  string
	VolumeYAML string
}

// AdditionalPortData provides information about extra ports on the container
type AdditionalPortData struct {
	Name     string
	Port     int
	Protocol string
}
