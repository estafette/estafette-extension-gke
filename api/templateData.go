package api

// TemplateData contains the root data for rendering the Kubernetes manifests
type TemplateData struct {
	Name                       string
	NameWithTrack              string
	Namespace                  string
	Schedule                   string
	RestartPolicy              string
	Completions                int
	Parallelism                int
	BackoffLimit               int
	ProgressDeadlineSeconds    int
	ConcurrencyPolicy          string
	Labels                     map[string]string
	PodLabels                  map[string]string
	AppLabelSelector           string
	Hosts                      []string
	HostsJoined                string
	InternalHosts              []string
	InternalHostsJoined        string
	AllHosts                   []string
	AllHostsJoined             string
	ApigeeHosts                []string
	ApigeeHostsJoined          string
	IngressPath                string
	InternalIngressPath        string
	UseIngress                 bool
	UseNginxIngress            bool
	UseGCEIngress              bool
	PathType                   string
	UseDNSAnnotationsOnIngress bool
	UseCloudflareProxy         bool

	Service            ServiceData
	UsePrometheusProbe bool

	MinReplicas                          int
	MaxReplicas                          int
	TargetCPUPercentage                  int
	UseHpaScaler                         bool
	HpaScalerPromQuery                   string
	HpaScalerRequestsPerReplica          string
	HpaScalerDelta                       string
	HpaScalerScaleDownMaxRatio           string
	VpaUpdateMode                        string
	PreferPreemptibles                   bool
	UseWindowsNodes                      bool
	Container                            ContainerData
	Sidecars                             []SidecarData
	HasCustomSidecars                    bool
	CustomSidecars                       []*map[string]interface{}
	HasInitContainers                    bool
	InitContainers                       []*map[string]interface{}
	PodSecurityContext                   map[string]interface{}
	HasOpenrestySidecar                  bool
	UseESP                               bool
	UseWorkloadIdentity                  bool
	HasEspConfigID                       bool
	EspConfigID                          string
	EspService                           string
	EspRequestTimeout                    int
	MountVolumes                         bool
	MountSslCertificate                  bool
	MountApplicationSecrets              bool
	Secrets                              map[string]interface{}
	SecretMountPath                      string
	MountConfigmap                       bool
	ConfigmapFiles                       map[string]string
	ConfigMountPath                      string
	MountPayloadLogging                  bool
	MountServiceAccountSecret            bool
	DisableServiceAccountKeyRotation     bool
	UseLegacyServiceAccountKey           bool
	LegacyServiceAccountKey              string
	GoogleCloudCredentialsAppName        string
	GoogleCloudCredentialsLabels         map[string]string
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
	UseHTTPS                             bool
	AllowHTTP                            bool
	BackendConfigTimeout                 int
	NginxAuthTLSSecret                   string
	NginxAuthTLSVerifyDepth              int
	Tolerations                          []*map[string]interface{}
	HasTolerations                       bool

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
	RenderToYAML                    func(v interface{}, data interface{}) string
	UseCertificateSecret            bool
	CertificateSecretName           string

	HasImagePullSecret bool
	DockerConfig       map[string]map[string]map[string]string

	IncludeAtomicIDSelector bool
	AtomicID                string
}

// ServiceData has data specific to service
type ServiceData struct {
	ServiceType                         ServiceType
	Name                                string
	UseDNSAnnotationsOnService          bool   `default:"false"`
	UseBackendConfigAnnotationOnService bool   `default:"false"`
	UseNegAnnotationOnService           bool   `default:"false"`
	LimitTrustedIPRanges                bool   `default:"false"`
	NameSuffix                          string `default:""`
}

// ContainerData has data specific to the application container
type ContainerData struct {
	Repository                      string
	Name                            string
	Tag                             string
	ImagePullPolicy                 string
	CPURequest                      string
	MemoryRequest                   string
	CPULimit                        string
	MemoryLimit                     string
	Port                            int
	EnvironmentVariables            map[string]interface{}
	SecretEnvironmentVariables      map[string]interface{}
	Liveness                        ProbeData
	Readiness                       ProbeData
	Metrics                         MetricsData
	UseLifecyclePreStopSleepCommand bool
	PreStopSleepSeconds             int
	ContainerSecurityContext        map[string]interface{}
	ContainerLifeCycle              map[string]interface{}
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
	Type                       string
	Image                      string
	EnvironmentVariables       map[string]interface{}
	SecretEnvironmentVariables map[string]interface{}
	HasEnvironmentVariables    bool
	CPURequest                 string
	MemoryRequest              string
	CPULimit                   string
	MemoryLimit                string
	SidecarSpecificProperties  map[string]interface{}
	HasCustomProperties        bool
	CustomPropertiesYAML       string
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
