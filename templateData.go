package main

// TemplateData contains the root data for rendering the Kubernetes manifests
type TemplateData struct {
	Name                                 string
	NameWithTrack                        string
	Namespace                            string
	Schedule                             string
	RestartPolicy                        string
	Completions                          int
	Parallelism                          int
	BackoffLimit                         int
	ConcurrencyPolicy                    string
	Labels                               map[string]string
	PodLabels                            map[string]string
	AppLabelSelector                     string
	Hosts                                []string
	HostsJoined                          string
	InternalHosts                        []string
	InternalHostsJoined                  string
	AllHosts                             []string
	AllHostsJoined                       string
	IngressPath                          string
	InternalIngressPath                  string
	UseIngress                           bool
	UseNginxIngress                      bool
	UseGCEIngress                        bool
	UseDNSAnnotationsOnIngress           bool
	UseCloudflareProxy                   bool
	UseDNSAnnotationsOnService           bool
	UseBackendConfigAnnotationOnService  bool
	UsePrometheusProbe                   bool
	ServiceType                          string
	MinReplicas                          int
	MaxReplicas                          int
	TargetCPUPercentage                  int
	UseHpaScaler                         bool
	HpaScalerPromQuery                   string
	HpaScalerRequestsPerReplica          string
	HpaScalerDelta                       string
	HpaScalerScaleDownMaxRatio           string
	PreferPreemptibles                   bool
	Container                            ContainerData
	Sidecars                             []SidecarData
	HasInitContainers                    bool
	InitContainers                       []*map[string]interface{}
	HasOpenrestySidecar                  bool
	UseESP                               bool
	HasEspConfigID                       bool
	EspConfigID                          string
	MountApplicationSecrets              bool
	Secrets                              map[string]interface{}
	SecretMountPath                      string
	MountConfigmap                       bool
	ConfigmapFiles                       map[string]string
	ConfigMountPath                      string
	MountPayloadLogging                  bool
	MountServiceAccountSecret            bool
	DisableServiceAccountKeyRotation     bool
	GoogleCloudCredentialsAppName        string
	StrategyType                         string
	RollingUpdateMaxSurge                string
	RollingUpdateMaxUnavailable          string
	LimitTrustedIPRanges                 bool
	TrustedIPRanges                      []string
	ManifestData                         map[string]interface{}
	IncludeTrackLabel                    bool
	TrackLabel                           string
	AddSafeToEvictAnnotation             bool
	MountAdditionalVolumes               bool
	AdditionalVolumeMounts               []VolumeMountData
	AdditionalContainerPorts             []AdditionalPortData
	AdditionalServicePorts               []AdditionalPortData
	OverrideDefaultWhitelist             bool
	NginxIngressWhitelist                string
	NginxIngressClientBodyBufferSize     string
	NginxIngressProxyConnectTimeout      int
	NginxIngressProxySendTimeout         int
	NginxIngressProxyReadTimeout         int
	NginxIngressProxyBodySize            string
	NginxIngressProxyBufferSize          string
	NginxIngressProxyBuffersNumber       string
	SetsNginxIngressLoadBalanceAlgorithm bool
	NginxIngressLoadBalanceAlgorithm     string

	IncludeReplicas                 bool
	Replicas                        int
	PodManagementPolicy             string
	StorageClass                    string
	StorageSize                     string
	StorageMountPath                string
	IapOauthCredentialsClientID     string
	IapOauthCredentialsClientSecret string
	IsSimpleEnvvarValue             func(interface{}) bool
	ToYAML                          func(interface{}) string
	UseCertificateSecret            bool
	CertificateSecretName           string
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
	PeriodSeconds       int
	IncludeOnContainer  bool
	FailureThreshold    int
	SuccessThreshold    int
}

// MetricsData has data to configure prometheus metrics scraping
type MetricsData struct {
	Scrape bool
	Path   string
	Port   int
}

// SidecarData configures the injected sidecar
type SidecarData struct {
	Type                      string
	Image                     string
	EnvironmentVariables      map[string]interface{}
	HasEnvironmentVariables   bool
	CPURequest                string
	MemoryRequest             string
	CPULimit                  string
	MemoryLimit               string
	SidecarSpecificProperties map[string]interface{}
	HasCustomProperties       bool
	CustomPropertiesYAML      string
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
