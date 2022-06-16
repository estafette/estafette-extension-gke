package api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/rs/zerolog/log"
)

// Params is used to parameterize the deployment, set from custom properties in the manifest
type Params struct {
	// control params
	Action                  ActionType      `json:"action,omitempty" yaml:"action,omitempty"`
	Kind                    Kind            `json:"kind,omitempty" yaml:"kind,omitempty"`
	DryRun                  bool            `json:"dryrun,omitempty" yaml:"dryrun,omitempty"`
	ProgressDeadlineSeconds int             `json:"progressDeadlineSeconds,omitempty" yaml:"progressDeadlineSeconds,omitempty"`
	BuildVersion            string          `json:"-" yaml:"-"`
	ChaosProof              bool            `json:"chaosproof,omitempty" yaml:"chaosproof,omitempty"`
	OperatingSystem         OperatingSystem `json:"os,omitempty" yaml:"os,omitempty"`
	Manifests               ManifestsParams `json:"manifests,omitempty" yaml:"manifests,omitempty"`
	TrustedIPRanges         []string        `json:"trustedips,omitempty" yaml:"trustedips,omitempty"`

	// app params
	App                             string                 `json:"app,omitempty" yaml:"app,omitempty"`
	Namespace                       string                 `json:"namespace,omitempty" yaml:"namespace,omitempty"`
	Schedule                        string                 `json:"schedule,omitempty" yaml:"schedule,omitempty"`
	RestartPolicy                   string                 `json:"restartPolicy,omitempty" yaml:"restartPolicy,omitempty"`
	Completions                     int                    `json:"completions,omitempty" yaml:"completions,omitempty"`
	Parallelism                     int                    `json:"parallelism,omitempty" yaml:"parallelism,omitempty"`
	BackoffLimit                    *int                   `json:"backoffLimit,omitempty" yaml:"backoffLimit,omitempty"`
	ConcurrencyPolicy               string                 `json:"concurrencypolicy,omitempty" yaml:"concurrencypolicy,omitempty"`
	PodManagementPolicy             string                 `json:"podManagementpolicy,omitempty" yaml:"podManagementpolicy,omitempty"`
	Replicas                        int                    `json:"replicas,omitempty" yaml:"replicas,omitempty"`
	StorageClass                    string                 `json:"storageclass,omitempty" yaml:"storageclass,omitempty"`
	StorageSize                     string                 `json:"storagesize,omitempty" yaml:"storagesize,omitempty"`
	StorageMountPath                string                 `json:"storagemountpath,omitempty" yaml:"storagemountpath,omitempty"`
	Labels                          map[string]string      `json:"labels,omitempty" yaml:"labels,omitempty"`
	Visibility                      Visibility             `json:"visibility,omitempty" yaml:"visibility,omitempty"`
	ContainerNativeLoadBalancing    bool                   `json:"containerNativeLoadBalancing,omitempty" yaml:"containerNativeLoadBalancing,omitempty"`
	IapOauthCredentialsClientID     string                 `json:"iapOauthClientID,omitempty" yaml:"iapOauthClientID,omitempty"`
	IapOauthCredentialsClientSecret string                 `json:"iapOauthClientSecret,omitempty" yaml:"iapOauthClientSecret,omitempty"`
	EspEndpointsProjectID           string                 `json:"espEndpointsProjectID,omitempty" yaml:"espEndpointsProjectID,omitempty"`
	EspConfigID                     string                 `json:"espConfigID,omitempty" yaml:"espConfigID,omitempty"`
	EspOpenAPIYamlPath              string                 `json:"espOpenapiYamlPath,omitempty" yaml:"espOpenapiYamlPath,omitempty"`
	EspServiceTypeClusterIP         bool                   `json:"espServiceTypeClusterIP,omitempty" yaml:"espServiceTypeClusterIP,omitempty"`
	WhitelistedIPS                  []string               `json:"whitelist,omitempty" yaml:"whitelist,omitempty"`
	Hosts                           []string               `json:"hosts,omitempty" yaml:"hosts,omitempty"`
	HostsRouteOnly                  []string               `json:"hostsrouteonly,omitempty" yaml:"hostsrouteonly,omitempty"`
	InternalHosts                   []string               `json:"internalhosts,omitempty" yaml:"internalhosts,omitempty"`
	InternalHostsRouteOnly          []string               `json:"internalhostsrouteonly,omitempty" yaml:"internalhostsrouteonly,omitempty"`
	ApigeeSuffix                    string                 `json:"apigeesuffix,omitempty" yaml:"apigeesuffix,omitempty"`
	Basepath                        string                 `json:"basepath,omitempty" yaml:"basepath,omitempty"`
	Autoscale                       AutoscaleParams        `json:"autoscale,omitempty" yaml:"autoscale,omitempty"`
	VerticalPodAutoscaler           VPAParams              `json:"vpa,omitempty" yaml:"vpa,omitempty"`
	Request                         RequestParams          `json:"request,omitempty" yaml:"request,omitempty"`
	Secrets                         SecretsParams          `json:"secrets,omitempty" yaml:"secrets,omitempty"`
	Configs                         ConfigsParams          `json:"configs,omitempty" yaml:"configs,omitempty"`
	VolumeMounts                    []VolumeMountParams    `json:"volumemounts,omitempty" yaml:"volumemounts,omitempty"`
	CertificateSecret               string                 `json:"certificatesecret,omitempty" yaml:"certificatesecret,omitempty"`
	AllowHTTP                       bool                   `json:"allowhttp,omitempty" yaml:"allowhttp,omitempty"`
	EnablePayloadLogging            bool                   `json:"enablePayloadLogging,omitempty" yaml:"enablePayloadLogging,omitempty"`
	UseGoogleCloudCredentials       bool                   `json:"useGoogleCloudCredentials,omitempty" yaml:"useGoogleCloudCredentials,omitempty"`
	WorkloadIdentity                *bool                  `json:"workloadIdentity,omitempty" yaml:"workloadIdentity,omitempty"`
	PodSecurityContext              map[string]interface{} `json:"securityContext,omitempty" yaml:"securityContext,omitempty"`

	DisableServiceAccountKeyRotation       *bool                     `json:"disableServiceAccountKeyRotation,omitempty" yaml:"disableServiceAccountKeyRotation,omitempty"`
	LegacyGoogleCloudServiceAccountKeyFile string                    `json:"legacyGoogleCloudServiceAccountKeyFile,omitempty" yaml:"legacyGoogleCloudServiceAccountKeyFile,omitempty"`
	GoogleCloudCredentialsApp              string                    `json:"googleCloudCredentialsApp,omitempty" yaml:"googleCloudCredentialsApp,omitempty"`
	ProbeService                           *bool                     `json:"probeService,omitempty" yaml:"probeService,omitempty"`
	TopologyAwareHints                     *bool                     `json:"topologyAwareHints,omitempty" yaml:"topologyAwareHints,omitempty"`
	Tolerations                            []*map[string]interface{} `json:"tolerations,omitempty" yaml:"tolerations,omitempty"`

	// container params
	Container              ContainerParams           `json:"container,omitempty" yaml:"container,omitempty"`
	InjectHTTPProxySidecar *bool                     `json:"injecthttpproxysidecar,omitempty" yaml:"injecthttpproxysidecar,omitempty"`
	InitContainers         []*map[string]interface{} `json:"initcontainers,omitempty" yaml:"initcontainers,omitempty"`
	Sidecar                SidecarParams             `json:"sidecar,omitempty" yaml:"sidecar,omitempty"`
	Sidecars               []*SidecarParams          `json:"sidecars,omitempty" yaml:"sidecars,omitempty"`
	CustomSidecars         []*map[string]interface{} `json:"customsidecars,omitempty" yaml:"customsidecars,omitempty"`
	StrategyType           StrategyType              `json:"strategytype,omitempty" yaml:"strategytype,omitempty"`
	AtomicID               string                    `json:"-" yaml:"-"`
	RollingUpdate          RollingUpdateParams       `json:"rollingupdate,omitempty" yaml:"rollingupdate,omitempty"`

	// set default image for sidecars
	DefaultOpenrestySidecarImage     string `json:"defaultOpenrestySidecarImage,omitempty" yaml:"defaultOpenrestySidecarImage,omitempty"`
	DefaultESPSidecarImage           string `json:"defaultESPSidecarImage,omitempty" yaml:"defaultESPSidecarImage,omitempty"`
	DefaultESPv2SidecarImage         string `json:"defaultESPv2SidecarImage,omitempty" yaml:"defaultESPv2SidecarImage,omitempty"`
	DefaultCloudSQLProxySidecarImage string `json:"defaultCloudSQLProxySidecarImage,omitempty" yaml:"defaultCloudSQLProxySidecarImage,omitempty"`

	// params for image pull secret
	ImagePullSecretUser     string `json:"imagePullSecretUser,omitempty" yaml:"imagePullSecretUser,omitempty"`
	ImagePullSecretPassword string `json:"imagePullSecretPassword,omitempty" yaml:"imagePullSecretPassword,omitempty"`
}

// ContainerParams defines the container image to deploy
type ContainerParams struct {
	ImageRepository            string                 `json:"repository,omitempty" yaml:"repository,omitempty"`
	ImageName                  string                 `json:"name,omitempty" yaml:"name,omitempty"`
	ImageTag                   string                 `json:"tag,omitempty" yaml:"tag,omitempty"`
	ImagePullPolicy            string                 `json:"imagePullPolicy,omitempty" yaml:"imagePullPolicy,omitempty"`
	Port                       int                    `json:"port,omitempty" yaml:"port,omitempty"`
	PortGrpc                   int                    `json:"portGrpc,omitempty" yaml:"portGrpc,omitempty"`
	EnvironmentVariables       map[string]interface{} `json:"env,omitempty" yaml:"env,omitempty"`
	SecretEnvironmentVariables map[string]interface{} `json:"secretEnv,omitempty" yaml:"secretEnv,omitempty"`

	CPU            CPUParams       `json:"cpu,omitempty" yaml:"cpu,omitempty"`
	Memory         MemoryParams    `json:"memory,omitempty" yaml:"memory,omitempty"`
	LivenessProbe  ProbeParams     `json:"liveness,omitempty" yaml:"liveness,omitempty"`
	ReadinessProbe ProbeParams     `json:"readiness,omitempty" yaml:"readiness,omitempty"`
	Metrics        MetricsParams   `json:"metrics,omitempty" yaml:"metrics,omitempty"`
	Lifecycle      LifecycleParams `json:"lifecycle,omitempty" yaml:"lifecycle,omitempty"`

	AdditionalPorts []*AdditionalPortParams `json:"additionalports,omitempty" yaml:"additionalports,omitempty"`

	ContainerSecurityContext map[string]interface{} `json:"securityContext,omitempty" yaml:"securityContext,omitempty"`

	ContainerLifeCycle map[string]interface{} `json:"containerLifecycle,omitempty" yaml:"containerLifecycle,omitempty"`
}

// AdditionalPortParams provides information about any additional ports exposed and accessible via a service
type AdditionalPortParams struct {
	Name       string     `json:"name,omitempty" yaml:"name,omitempty"`
	Port       int        `json:"port,omitempty" yaml:"port,omitempty"`
	Protocol   string     `json:"protocol,omitempty" yaml:"protocol,omitempty"`
	Visibility Visibility `json:"visibility,omitempty" yaml:"visibility,omitempty"`
}

// CPUParams sets cpu request and limit values
type CPUParams struct {
	Request string `json:"request,omitempty" yaml:"request,omitempty"`
	Limit   string `json:"limit,omitempty" yaml:"limit,omitempty"`
}

// MemoryParams sets memory request and limit values
type MemoryParams struct {
	Request string `json:"request,omitempty" yaml:"request,omitempty"`
	Limit   string `json:"limit,omitempty" yaml:"limit,omitempty"`
}

// AutoscaleParams controls autoscaling
type AutoscaleParams struct {
	Enabled       *bool                 `json:"enabled,omitempty" yaml:"enabled,omitempty"`
	MinReplicas   int                   `json:"min,omitempty" yaml:"min,omitempty"`
	MaxReplicas   int                   `json:"max,omitempty" yaml:"max,omitempty"`
	CPUPercentage int                   `json:"cpu,omitempty" yaml:"cpu,omitempty"`
	Safety        AutoscaleSafetyParams `json:"safety,omitempty" yaml:"safety,omitempty"`
}

type VPAParams struct {
	Enabled    *bool      `json:"enabled,omitempty" yaml:"enabled,omitempty"`
	UpdateMode UpdateMode `json:"updateMode,omitempty" yaml:"updateMode,omitempty"`
}

// AutoscaleSafetyParams configures the autoscaler to use estafette-hpa-scaler as a safety net
type AutoscaleSafetyParams struct {
	Enabled        bool   `json:"enabled,omitempty" yaml:"enabled,omitempty"`
	PromQuery      string `json:"promquery,omitempty" yaml:"promquery,omitempty"`
	Ratio          string `json:"ratio,omitempty" yaml:"ratio,omitempty"`
	Delta          string `json:"delta,omitempty" yaml:"delta,omitempty"`
	ScaleDownRatio string `json:"scaledownratio,omitempty" yaml:"scaledownratio,omitempty"`
}

// RequestParams controls timeouts, max body size, etc
type RequestParams struct {
	IngressBackendProtocol string `json:"ingressbackendprotocol,omitempty" yaml:"ingressbackendprotocol,omitempty"`
	Timeout                string `json:"timeout,omitempty" yaml:"timeout,omitempty"`
	MaxBodySize            string `json:"maxbodysize,omitempty" yaml:"maxbodysize,omitempty"`
	ProxyBufferSize        string `json:"proxybuffersize,omitempty" yaml:"proxybuffersize,omitempty"`
	ProxyBuffersNumber     int    `json:"proxybuffersnumber,omitempty" yaml:"proxybuffersnumber,omitempty"`
	ClientBodyBufferSize   string `json:"clientbodybuffersize,omitempty" yaml:"clientbodybuffersize,omitempty"`
	LoadBalanceAlgorithm   string `json:"loadbalance,omitempty" yaml:"loadbalance,omitempty"`
	AuthSecret             string `json:"authsecret,omitempty" yaml:"authsecret,omitempty"`
	VerifyDepth            int    `json:"verifydepth,omitempty" yaml:"verifydepth,omitempty"`
}

// ProbeParams sets params for liveness or readiness probe
type ProbeParams struct {
	Enabled             *bool  `json:"enabled,omitempty" yaml:"enabled,omitempty"`
	Path                string `json:"path,omitempty" yaml:"path,omitempty"`
	Port                int    `json:"port,omitempty" yaml:"port,omitempty"`
	InitialDelaySeconds int    `json:"delay,omitempty" yaml:"delay,omitempty"`
	TimeoutSeconds      int    `json:"timeout,omitempty" yaml:"timeout,omitempty"`
	PeriodSeconds       int    `json:"period,omitempty" yaml:"period,omitempty"`
	FailureThreshold    int    `json:"failureThreshold,omitempty" yaml:"failureThreshold,omitempty"`
	SuccessThreshold    int    `json:"successThreshold,omitempty" yaml:"successThreshold,omitempty"`
}

// MetricsParams sets params for scraping prometheus metrics
type MetricsParams struct {
	Scrape *bool  `json:"scrape,omitempty" yaml:"scrape,omitempty"`
	Path   string `json:"path,omitempty" yaml:"path,omitempty"`
	Port   int    `json:"port,omitempty" yaml:"port,omitempty"`
}

// LifecycleParams sets params for lifecycle commands
type LifecycleParams struct {
	PrestopSleep        *bool `json:"prestopsleep,omitempty" yaml:"prestopsleep,omitempty"`
	PrestopSleepSeconds *int  `json:"prestopsleepseconds,omitempty" yaml:"prestopsleepseconds,omitempty"`
}

// SidecarParams sets params for sidecar injection
type SidecarParams struct {
	Type                              SidecarType            `json:"type,omitempty" yaml:"type,omitempty"`
	Image                             string                 `json:"image,omitempty" yaml:"image,omitempty"`
	EnvironmentVariables              map[string]interface{} `json:"env,omitempty" yaml:"env,omitempty"`
	SecretEnvironmentVariables        map[string]interface{} `json:"secretEnv,omitempty" yaml:"secretEnv,omitempty"`
	CPU                               CPUParams              `json:"cpu,omitempty" yaml:"cpu,omitempty"`
	Memory                            MemoryParams           `json:"memory,omitempty" yaml:"memory,omitempty"`
	HealthCheckPath                   string                 `json:"healthcheckpath,omitempty" yaml:"healthcheckpath,omitempty"`
	DbInstanceConnectionName          string                 `json:"dbinstanceconnectionname,omitempty" yaml:"dbinstanceconnectionname,omitempty"`
	SQLProxyPort                      int                    `json:"sqlproxyport,omitempty" yaml:"sqlproxyport,omitempty"`
	SQLProxyTerminationTimeoutSeconds int                    `json:"sqlproxyterminationtimeoutseconds,omitempty" yaml:"sqlproxyterminationtimeoutseconds,omitempty"`
	CustomProperties                  map[string]interface{} `yaml:",inline"`
}

// RollingUpdateParams sets params for controlling rolling update speed
type RollingUpdateParams struct {
	MaxSurge       string `json:"maxsurge,omitempty" yaml:"maxsurge,omitempty"`
	MaxUnavailable string `json:"maxunavailable,omitempty" yaml:"maxunavailable,omitempty"`
	Timeout        string `json:"timeout,omitempty" yaml:"timeout,omitempty"`
}

// ManifestsParams can be used to override or add additional manifests located in the application repository
type ManifestsParams struct {
	Files []string               `json:"files,omitempty" yaml:"files,omitempty"`
	Data  map[string]interface{} `json:"data,omitempty" yaml:"data,omitempty"`
}

// SecretsParams allows secrets to be set dynamically for the application
type SecretsParams struct {
	Keys      map[string]interface{} `json:"keys,omitempty" yaml:"keys,omitempty"`
	MountPath string                 `json:"mountpath,omitempty" yaml:"mountpath,omitempty"`
}

// ConfigsParams allows configs to be set dynamically for the application
type ConfigsParams struct {
	Files               []string               `json:"files,omitempty" yaml:"files,omitempty"`
	Data                map[string]interface{} `json:"data,omitempty" yaml:"data,omitempty"`
	InlineFiles         map[string]string      `json:"inline,omitempty" yaml:"inline,omitempty"`
	MountPath           string                 `json:"mountpath,omitempty" yaml:"mountpath,omitempty"`
	RenderedFileContent map[string]string      `json:"-" yaml:"-"`
}

// VolumeMountParams allows additional mounts for already existing volumes, secrets, etc
type VolumeMountParams struct {
	Name      string                 `json:"name,omitempty" yaml:"name,omitempty"`
	MountPath string                 `json:"mountpath,omitempty" yaml:"mountpath,omitempty"`
	Volume    map[string]interface{} `json:"volume,omitempty" yaml:"volume,omitempty"`
}

// SetDefaults fills in empty fields with convention-based defaults
func (p *Params) SetDefaults(gitSource, gitOwner, gitName, appLabel, buildVersion, releaseName string, releaseAction ActionType, releaseID string, estafetteLabels map[string]string) {

	p.BuildVersion = buildVersion

	// default action to deploy-simple unless it's either specified on the stage or passed in as a release action
	if releaseAction != ActionUnknown {
		p.Action = ActionType(releaseAction)
	} else if p.Action == "" && releaseAction == ActionUnknown {
		p.Action = ActionDeploySimple
	}

	// default kind to deployment
	if p.Kind == KindUnknown {
		p.Kind = KindDeployment
	}

	// default operating system to linux
	if p.OperatingSystem == OperatingSystemUnknown {
		p.OperatingSystem = OperatingSystemLinux
	}

	// default app to estafette app label if no override in stage params
	if p.App == "" && appLabel == "" && gitName != "" {
		p.App = gitName
	}
	if p.App == "" && appLabel != "" {
		p.App = appLabel
	}

	if p.ProgressDeadlineSeconds <= 0 {
		p.ProgressDeadlineSeconds = 600
	}

	// default DisableServiceAccountKeyRotation to true for avoiding unintended side-effects of key rotation
	if p.DisableServiceAccountKeyRotation == nil {
		trueValue := true
		p.DisableServiceAccountKeyRotation = &trueValue
	}

	// default GoogleCloudCredentialsApp to App if empty
	if p.GoogleCloudCredentialsApp == "" {
		p.GoogleCloudCredentialsApp = p.App
	}

	// default image name to estafette app label if no override in stage params
	if p.Container.ImageName == "" && p.App != "" {
		p.Container.ImageName = p.App
	}

	// default image tag to estafette build version if no override in stage params
	if p.Container.ImageTag == "" && buildVersion != "" {
		p.Container.ImageTag = buildVersion
	}

	// default image pull policy to IfNotPresent
	if p.Container.ImagePullPolicy == "" {
		p.Container.ImagePullPolicy = "IfNotPresent"
	}

	// default labels to estafette labels if no override in stage params
	if p.Labels == nil {
		p.Labels = map[string]string{}
	}
	if len(p.Labels) == 0 && estafetteLabels != nil && len(estafetteLabels) != 0 {
		p.Labels = estafetteLabels
	}
	// ensure the app label is set and equals the app label or app override in stage params if present
	if p.App != "" {
		p.Labels["app"] = p.App
	}
	if gitSource != "" && gitOwner != "" && gitName != "" {
		// add estafette.io/pipeline label to labels set on all resources for linking catalog entities back to pipelines

		pipeline := fmt.Sprintf("%v/%v/%v", gitSource, gitOwner, gitName)
		pipelineBase64 := base64.StdEncoding.EncodeToString([]byte(pipeline))

		p.Labels["estafette.io/pipeline"] = SanitizeLabel(pipeline)
		p.Labels["estafette.io/pipeline-base64"] = SanitizeLabel(pipelineBase64)
	}

	// default visibility to private if no override in stage params
	if p.Visibility == VisibilityUnknown {
		p.Visibility = VisibilityPrivate
	}

	// set default workloadIdentity
	if p.WorkloadIdentity == nil {
		falseValue := false
		p.WorkloadIdentity = &falseValue
	}

	// set cpu defaults
	cpuRequestIsEmpty := p.Container.CPU.Request == ""
	if cpuRequestIsEmpty {
		if p.Container.CPU.Limit != "" {
			p.Container.CPU.Request = p.Container.CPU.Limit
		} else {
			p.Container.CPU.Request = "100m"
		}
	}

	// set memory defaults
	memoryRequestIsEmpty := p.Container.Memory.Request == ""
	if memoryRequestIsEmpty {
		if p.Container.Memory.Limit != "" {
			p.Container.Memory.Request = p.Container.Memory.Limit
		} else {
			p.Container.Memory.Request = "128Mi"
		}
	}
	if p.Container.Memory.Limit == "" {
		if !memoryRequestIsEmpty {
			p.Container.Memory.Limit = p.Container.Memory.Request
		} else {
			p.Container.Memory.Limit = "128Mi"
		}
	}

	// set container port defaults
	if p.Container.Port <= 0 {
		p.Container.Port = 5000
	}

	if p.Container.PortGrpc <= 0 {
		// gRPC is optional, we have to opt-in by setting the portgrpc field of the deployment.
		p.Container.PortGrpc = 0
	}

	// set additional ports defaults
	if len(p.Container.AdditionalPorts) > 0 {
		for _, ap := range p.Container.AdditionalPorts {
			if ap.Protocol == "" {
				ap.Protocol = "TCP"
			}
			if ap.Visibility == VisibilityUnknown {
				ap.Visibility = p.Visibility
			}
		}
	}

	// set autoscale defaults
	if p.Autoscale.Enabled == nil {
		trueValue := true
		p.Autoscale.Enabled = &trueValue
	}
	if p.Autoscale.MinReplicas <= 0 {
		p.Autoscale.MinReplicas = 3
	}
	if p.Autoscale.MaxReplicas <= 0 {
		p.Autoscale.MaxReplicas = 100
	}
	if p.Autoscale.CPUPercentage <= 0 {
		p.Autoscale.CPUPercentage = 80
	}

	if p.Autoscale.Safety.PromQuery == "" {
		p.Autoscale.Safety.PromQuery = fmt.Sprintf("sum(rate(nginx_http_requests_total{app='%v'}[5m])) by (app)", p.App)
	}
	if p.Autoscale.Safety.Ratio == "" {
		p.Autoscale.Safety.Ratio = "1"
	}
	if p.Autoscale.Safety.ScaleDownRatio == "" {
		p.Autoscale.Safety.ScaleDownRatio = "1"
	}

	// set vpa defaults
	if p.VerticalPodAutoscaler.Enabled == nil {
		falseValue := false
		p.VerticalPodAutoscaler.Enabled = &falseValue
	}
	if p.VerticalPodAutoscaler.UpdateMode == UpdateModeUnknown {
		p.VerticalPodAutoscaler.UpdateMode = UpdateModeOff
	}

	// set request defaults
	if p.Request.IngressBackendProtocol == "" {
		p.Request.IngressBackendProtocol = "HTTPS"
	}
	if p.Request.Timeout == "" {
		p.Request.Timeout = "60s"
	}
	if p.Request.MaxBodySize == "" {
		p.Request.MaxBodySize = "128m"
	}
	if p.Request.ProxyBufferSize == "" {
		p.Request.ProxyBufferSize = "4k"
	}
	if p.Request.ProxyBuffersNumber <= 0 {
		p.Request.ProxyBuffersNumber = 4
	}

	if p.Request.ClientBodyBufferSize == "" {
		p.Request.ClientBodyBufferSize = "8k"
	}

	// set liveness probe defaults
	if p.Container.LivenessProbe.Enabled == nil {
		trueValue := true
		p.Container.LivenessProbe.Enabled = &trueValue
	}
	if p.Container.LivenessProbe.Path == "" {
		p.Container.LivenessProbe.Path = "/liveness"
	}
	if p.Container.LivenessProbe.Port <= 0 {
		p.Container.LivenessProbe.Port = p.Container.Port
	}
	if p.Container.LivenessProbe.InitialDelaySeconds <= 0 {
		p.Container.LivenessProbe.InitialDelaySeconds = 30
	}
	if p.Container.LivenessProbe.TimeoutSeconds <= 0 {
		p.Container.LivenessProbe.TimeoutSeconds = 1
	}
	if p.Container.LivenessProbe.PeriodSeconds <= 0 {
		p.Container.LivenessProbe.PeriodSeconds = 10
	}
	if p.Container.LivenessProbe.FailureThreshold <= 0 {
		p.Container.LivenessProbe.FailureThreshold = 3
	}
	if p.Container.LivenessProbe.SuccessThreshold <= 0 {
		p.Container.LivenessProbe.SuccessThreshold = 1
	}

	// set readiness probe defaults
	if p.Container.ReadinessProbe.Enabled == nil {
		if p.Kind == KindHeadlessDeployment {
			falseValue := false
			p.Container.ReadinessProbe.Enabled = &falseValue
		} else {
			trueValue := true
			p.Container.ReadinessProbe.Enabled = &trueValue
		}
	}
	if p.Container.ReadinessProbe.Path == "" {
		p.Container.ReadinessProbe.Path = "/readiness"
	}
	if p.Container.ReadinessProbe.Port <= 0 {
		p.Container.ReadinessProbe.Port = p.Container.Port
	}
	if p.Container.ReadinessProbe.TimeoutSeconds <= 0 {
		p.Container.ReadinessProbe.TimeoutSeconds = 1
	}
	if p.Container.ReadinessProbe.PeriodSeconds <= 0 {
		p.Container.ReadinessProbe.PeriodSeconds = 10
	}
	if p.Container.ReadinessProbe.FailureThreshold <= 0 {
		p.Container.ReadinessProbe.FailureThreshold = 3
	}
	if p.Container.ReadinessProbe.SuccessThreshold <= 0 {
		p.Container.ReadinessProbe.SuccessThreshold = 1
	}
	if p.ProbeService == nil {
		if p.Visibility == VisibilityESP || p.Visibility == VisibilityESPv2 {
			falseValue := false
			p.ProbeService = &falseValue
		} else {
			trueValue := true
			p.ProbeService = &trueValue
		}
	}

	// set metrics defaults
	if p.Container.Metrics.Path == "" {
		p.Container.Metrics.Path = "/metrics"
	}
	if p.Container.Metrics.Port <= 0 {
		p.Container.Metrics.Port = p.Container.Port
	}
	if p.Container.Metrics.Scrape == nil {
		trueValue := true
		p.Container.Metrics.Scrape = &trueValue
	}

	// set lifecycle defaults
	if p.Container.Lifecycle.PrestopSleep == nil {
		if p.OperatingSystem == OperatingSystemWindows {
			falseValue := false
			p.Container.Lifecycle.PrestopSleep = &falseValue
		} else {
			trueValue := true
			p.Container.Lifecycle.PrestopSleep = &trueValue
		}
	}
	if p.Container.Lifecycle.PrestopSleepSeconds == nil {
		defaultSleepValue := 20
		p.Container.Lifecycle.PrestopSleepSeconds = &defaultSleepValue
	}

	if p.Container.ContainerLifeCycle != nil {
		p.Container.Lifecycle.PrestopSleep = nil
		p.Container.Lifecycle.PrestopSleepSeconds = nil
		p.Container.Lifecycle = LifecycleParams{}
	}

	if p.InjectHTTPProxySidecar == nil {
		trueValue := true
		p.InjectHTTPProxySidecar = &trueValue
	}

	// if deprecated sidecar is still used add it to the sidecars list for backwards compatibility
	if p.Sidecar.Type != "" && p.Sidecar.Type != "none" {
		p.Sidecars = append([]*SidecarParams{&p.Sidecar}, p.Sidecars...)
	}

	// check if an openresty sidecar is in the list
	openrestySidecarSpecifiedInList := false
	for _, sidecar := range p.Sidecars {
		if sidecar.Type == SidecarTypeOpenresty {
			openrestySidecarSpecifiedInList = true
		}
	}

	// inject an openresty sidecar in the sidecars list if it isn't there yet for deployments
	if *p.InjectHTTPProxySidecar && !openrestySidecarSpecifiedInList && p.Kind == KindDeployment {
		openrestySidecar := SidecarParams{Type: SidecarTypeOpenresty}
		p.initializeSidecarDefaults(&openrestySidecar)

		p.Sidecars = append(p.Sidecars, &openrestySidecar)
	}

	if p.Visibility == VisibilityESP || p.Visibility == VisibilityESPv2 {
		if p.EspOpenAPIYamlPath == "" {
			p.EspOpenAPIYamlPath = "openapi.yaml"
		}

		// check if an esp sidecar is in the list
		espSidecarSpecifiedInList := false
		espv2SidecarSpecifiedInList := false
		for _, sidecar := range p.Sidecars {
			if sidecar.Type == SidecarTypeESP {
				espSidecarSpecifiedInList = true
			}
			if sidecar.Type == SidecarTypeESPv2 {
				espv2SidecarSpecifiedInList = true
			}
		}

		// inject an esp sidecar in the sidecars list if it isn't there yet for deployments
		if *p.InjectHTTPProxySidecar && p.Visibility == VisibilityESP && !espSidecarSpecifiedInList && p.Kind == KindDeployment {
			espSidecar := SidecarParams{Type: SidecarTypeESP}
			p.initializeSidecarDefaults(&espSidecar)

			p.Sidecars = append(p.Sidecars, &espSidecar)
		}

		// inject an espv2 sidecar in the sidecars list if it isn't there yet for deployments
		if *p.InjectHTTPProxySidecar && p.Visibility == VisibilityESPv2 && !espv2SidecarSpecifiedInList && p.Kind == KindDeployment {
			espSidecar := SidecarParams{Type: SidecarTypeESPv2}
			p.initializeSidecarDefaults(&espSidecar)

			p.Sidecars = append(p.Sidecars, &espSidecar)
		}
	}

	if p.Visibility == VisibilityApigee {
		if p.Request.VerifyDepth <= 0 {
			p.Request.VerifyDepth = 3
		}
		if p.ApigeeSuffix == "" {
			p.ApigeeSuffix = "apigee"
		}
	}

	if p.DefaultOpenrestySidecarImage == "" {
		p.DefaultOpenrestySidecarImage = "estafette/openresty-sidecar@sha256:f13a8412ed89cb8fe3a5fe2f1955e1f16665f7d7bfadc83c94d7880301dd3e32"
	}
	if p.DefaultESPSidecarImage == "" {
		p.DefaultESPSidecarImage = "gcr.io/endpoints-release/endpoints-runtime:1.57.0"
	}
	if p.DefaultESPv2SidecarImage == "" {
		p.DefaultESPv2SidecarImage = "gcr.io/endpoints-release/endpoints-runtime:2.29.1"
	}
	if p.DefaultCloudSQLProxySidecarImage == "" {
		p.DefaultCloudSQLProxySidecarImage = "eu.gcr.io/cloudsql-docker/gce-proxy:1.24.0"
	}

	for i := range p.Sidecars {
		p.initializeSidecarDefaults(p.Sidecars[i])
	}

	// default basepath to /
	if p.Basepath == "" {
		p.Basepath = "/"
	}

	// defaults for rollingupdate
	if p.StrategyType == StrategyTypeUnknown {
		p.StrategyType = StrategyTypeRollingUpdate
	}
	if p.StrategyType == StrategyTypeAtomicUpdate {
		p.AtomicID = releaseID
		if p.AtomicID == "" {
			p.AtomicID = p.BuildVersion
		}
	}

	if p.RollingUpdate.MaxSurge == "" {
		p.RollingUpdate.MaxSurge = "25%"
	}
	if p.RollingUpdate.MaxUnavailable == "" {
		p.RollingUpdate.MaxUnavailable = "0"
	}
	if p.RollingUpdate.Timeout == "" {
		p.RollingUpdate.Timeout = "5m"
	}

	if p.Replicas == 0 && p.StrategyType == StrategyTypeRecreate {
		p.Replicas = 1
	}

	// set mountpaths for configs and secrets
	if p.Configs.MountPath == "" {
		p.Configs.MountPath = "/configs"
	}
	if p.Secrets.MountPath == "" {
		p.Secrets.MountPath = "/secrets"
	}

	// default trusted ip ranges to cloudflare's ips from https://www.cloudflare.com/ips-v4
	if len(p.TrustedIPRanges) == 0 {
		p.TrustedIPRanges = []string{
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
		}
	}

	if p.Kind == KindCronJob {
		if p.ConcurrencyPolicy == "" {
			p.ConcurrencyPolicy = "Allow"
		}
	}

	if p.RestartPolicy == "" {
		p.RestartPolicy = "OnFailure"
	}
	if p.Completions <= 0 {
		p.Completions = 1
	}
	if p.Parallelism <= 0 {
		p.Parallelism = 1
	}
	if p.BackoffLimit == nil || *p.BackoffLimit < 0 {
		defaultBackoffLimit := 6
		p.BackoffLimit = &defaultBackoffLimit
	}

	if p.Kind == KindStatefulset {
		if p.PodManagementPolicy == "" {
			p.PodManagementPolicy = "Parallel"
		}
		if p.StorageClass == "" {
			p.StorageClass = "standard"
		}
		if p.StorageSize == "" {
			p.StorageSize = "1Gi"
		}
		if p.StorageMountPath == "" {
			p.StorageMountPath = "/data"
		}
	}
}

func (p *Params) HasSecrets() bool {
	if len(p.Secrets.Keys) > 0 {
		return true
	}

	if len(p.Container.SecretEnvironmentVariables) > 0 {
		return true
	}

	for _, sc := range p.Sidecars {
		if len(sc.SecretEnvironmentVariables) > 0 {
			return true
		}
	}

	return false
}

func (p *Params) initializeSidecarDefaults(sidecar *SidecarParams) {
	switch sidecar.Type {
	case SidecarTypeOpenresty:
		if sidecar.Image == "" {
			sidecar.Image = p.DefaultOpenrestySidecarImage
		}
		if sidecar.HealthCheckPath == "" {
			sidecar.HealthCheckPath = p.Container.ReadinessProbe.Path
		}
	case SidecarTypeESP:
		if sidecar.Image == "" {
			sidecar.Image = p.DefaultESPSidecarImage
		}
	case SidecarTypeESPv2:
		if sidecar.Image == "" {
			sidecar.Image = p.DefaultESPv2SidecarImage
		}
	case SidecarTypeCloudSQLProxy:
		if sidecar.Image == "" {
			sidecar.Image = p.DefaultCloudSQLProxySidecarImage
		}
		if sidecar.SQLProxyPort <= 0 {
			sidecar.SQLProxyPort = 5432
		}
		if sidecar.SQLProxyTerminationTimeoutSeconds <= 0 {
			sidecar.SQLProxyTerminationTimeoutSeconds = 60
		}
	}

	// set sidecar cpu defaults
	sidecarCPURequestIsEmpty := sidecar.CPU.Request == ""
	if sidecarCPURequestIsEmpty {
		if sidecar.CPU.Limit != "" {
			sidecar.CPU.Request = sidecar.CPU.Limit
		} else {
			sidecar.CPU.Request = "50m"
		}
	}

	// set sidecar memory defaults
	sidecarMemoryRequestIsEmpty := sidecar.Memory.Request == ""
	if sidecarMemoryRequestIsEmpty {
		if sidecar.Memory.Limit != "" {
			sidecar.Memory.Request = sidecar.Memory.Limit
		} else {
			sidecar.Memory.Request = "30Mi"
		}
	}
	if sidecar.Memory.Limit == "" {
		if !sidecarMemoryRequestIsEmpty {
			sidecar.Memory.Limit = sidecar.Memory.Request
		} else {
			sidecar.Memory.Limit = "50Mi"
		}
	}
}

// ValidateRequiredProperties checks whether all needed properties are set
func (p *Params) ValidateRequiredProperties() (bool, []error, []string) {

	errors := []error{}
	warnings := []string{}

	// validate app params
	if p.App == "" {
		errors = append(errors, fmt.Errorf("Application name is required; either define an app label or use app property on this stage"))
	}
	if p.Namespace == "" {
		errors = append(errors, fmt.Errorf("Namespace is required; either use credentials with a defaultNamespace or set it via namespace property on this stage"))
	}

	if p.Action == ActionRollbackCanary || p.Kind == KindConfig || p.Kind == KindConfigToFile {
		// the above properties are all you need for a rollback
		return len(errors) == 0, errors, warnings
	}

	// validate container params
	if p.Container.ImageRepository == "" {
		errors = append(errors, fmt.Errorf("Image repository is required; set it via container.repository property on this stage"))
	}
	if p.Container.ImageName == "" {
		errors = append(errors, fmt.Errorf("Image name is required; set it via container.name property on this stage"))
	}
	if p.Container.ImageTag == "" {
		errors = append(errors, fmt.Errorf("Image tag is required; set it via container.tag property on this stage"))
	}
	if p.Container.ImagePullPolicy == "" {
		errors = append(errors, fmt.Errorf("Image pull policy is required; set it via container.imagePullPolicy property on this stage; allowed values are IfNotPresent or Always"))
	}

	// validate cpu params
	if p.Container.CPU.Request == "" {
		errors = append(errors, fmt.Errorf("Cpu request is required; set it via container.cpu.request property on this stage"))
	}

	// validate memory params
	if p.Container.Memory.Request == "" {
		errors = append(errors, fmt.Errorf("Memory request is required; set it via container.memory.request property on this stage"))
	}
	if p.Container.Memory.Limit == "" {
		errors = append(errors, fmt.Errorf("Memory limit is required; set it via container.memory.limit property on this stage"))
	}

	// validate params for rollingupdate
	if p.StrategyType == StrategyTypeUnknown {
		errors = append(errors, fmt.Errorf("StrategyType is required; set it via strategytype property on this stage; valid values are RollingUpdate, Recreate or AtomicUpdate"))
	}
	if p.StrategyType == StrategyTypeAtomicUpdate && p.Action != ActionDeploySimple {
		errors = append(errors, fmt.Errorf("StrategyType: AtomicUpdate can't be used in combination with other actions than deploy-simple as this would allow multiple versions to be served. Please use action: deploy-simple"))
	}
	if p.RollingUpdate.MaxSurge == "" {
		errors = append(errors, fmt.Errorf("Rollingupdate max surge is required; set it via rollingupdate.maxsurge property on this stage"))
	}
	if p.RollingUpdate.MaxSurge == "" {
		errors = append(errors, fmt.Errorf("Rollingupdate max surge is required; set it via rollingupdate.maxsurge property on this stage"))
	}
	if p.RollingUpdate.MaxUnavailable == "" {
		errors = append(errors, fmt.Errorf("Rollingupdate max unavailable is required; set it via rollingupdate.maxunavailable property on this stage"))
	}

	if p.Kind == KindJob || p.Kind == KindCronJob {
		if p.Kind == KindCronJob {
			if p.Schedule == "" {
				errors = append(errors, fmt.Errorf("Schedule is required for a cronjob; set it via schedule property on this stage"))
			}

			if p.ConcurrencyPolicy != "Allow" && p.ConcurrencyPolicy != "Forbid" && p.ConcurrencyPolicy != "Replace" {
				errors = append(errors, fmt.Errorf("ConcurrencyPolicy is invalid; allowed values for concurrencypolicy property are Allow, Forbid or Replace"))
			}
		}

		// the above properties are all you need for a worker
		return len(errors) == 0, errors, warnings
	}

	if p.Kind == KindStatefulset {
		if p.PodManagementPolicy != "OrderedReady" && p.PodManagementPolicy != "Parallel" {
			errors = append(errors, fmt.Errorf("PodManagementPolicy is required for a statefulset; allowed values for podmanagementpolicy property are OrderedReady or Parallel"))
		}
		if p.StorageClass == "" {
			errors = append(errors, fmt.Errorf("StorageClass is required for a statefulset; set it via storageclass property on this stage"))
		}
		if p.StorageSize == "" {
			errors = append(errors, fmt.Errorf("StorageSize is required for a statefulset; set it via storagesize property on this stage"))
		}
		if p.StorageMountPath == "" {
			errors = append(errors, fmt.Errorf("StorageMountPath is required for a statefulset; set it via storagemountpath property on this stage"))
		}
	}
	// validate params with respect to incoming requests
	if p.Kind == KindDeployment {
		if p.Visibility == VisibilityUnknown || (p.Visibility != VisibilityPrivate && p.Visibility != VisibilityPublic && p.Visibility != VisibilityIAP && p.Visibility != VisibilityESP && p.Visibility != VisibilityESPv2 && p.Visibility != VisibilityPublicWhitelist && p.Visibility != VisibilityApigee) {
			errors = append(errors, fmt.Errorf("Visibility property is required; set it via visibility property on this stage; allowed values are private, iap, esp, public-whitelist, public or apigee"))
		}
		if p.Visibility == VisibilityPublic {
			warnings = append(warnings, "Visibility public is deprecated, please use esp or apigee.")
		}
		if p.Visibility == VisibilityIAP && p.IapOauthCredentialsClientID == "" {
			errors = append(errors, fmt.Errorf("With visibility 'iap' property iapOauthClientID is required; set it via iapOauthClientID property on this stage"))
		}
		if p.Visibility == VisibilityIAP && p.IapOauthCredentialsClientSecret == "" {
			errors = append(errors, fmt.Errorf("With visibility 'iap' property iapOauthClientSecret is required; set it via iapOauthClientSecret property on this stage"))
		}

		if (p.Visibility == VisibilityESP || p.Visibility == VisibilityESPv2) && (!p.UseGoogleCloudCredentials && (p.WorkloadIdentity == nil || !*p.WorkloadIdentity)) {
			errors = append(errors, fmt.Errorf("With visibility 'esp' property useGoogleCloudCredentials or workloadIdentity is required; set useGoogleCloudCredentials: true or workloadIdentity: true on this stage"))
		}
		if (p.Visibility == VisibilityESP || p.Visibility == VisibilityESPv2) && ((p.DisableServiceAccountKeyRotation == nil || !*p.DisableServiceAccountKeyRotation) && (p.WorkloadIdentity == nil || !*p.WorkloadIdentity)) {
			errors = append(errors, fmt.Errorf("With visibility 'esp' property disableServiceAccountKeyRotation is required; set disableServiceAccountKeyRotation: true on this stage"))
		}
		if (p.Visibility == VisibilityESP || p.Visibility == VisibilityESPv2) && (p.EspEndpointsProjectID == "") {
			errors = append(errors, fmt.Errorf("With visibility 'esp' property espEndpointsProjectID is required; provide id of the 'endpoints' project"))
		}
		if (p.Visibility == VisibilityESP || p.Visibility == VisibilityESPv2) && p.EspOpenAPIYamlPath == "" {
			errors = append(errors, fmt.Errorf("With visibility 'esp' property espOpenapiYamlPath is required; set espOpenapiYamlPath to the path towards openapi.yaml"))
		}
		if p.EspServiceTypeClusterIP && (p.Visibility != VisibilityESP && p.Visibility != VisibilityESPv2) {
			errors = append(errors, fmt.Errorf("With EspServiceTypeClusterIP set to true, visibility needs to be set to 'esp' or 'espv2'"))
		}
		if (p.Visibility == VisibilityESP || p.Visibility == VisibilityESPv2) && len(p.Hosts) < 1 {
			errors = append(errors, fmt.Errorf("With visibility 'esp' property at least one host is required. Set it via hosts array property on this stage"))
		}

		if p.Visibility == VisibilityApigee && p.Request.AuthSecret == "" {
			errors = append(errors, fmt.Errorf("With visibility 'apigee' property authsecret is required; set it via authsecret property for request on this stage"))
		}

		if len(p.Hosts) == 0 {
			errors = append(errors, fmt.Errorf("At least one host is required; set it via hosts array property on this stage"))
		}
		for _, host := range p.Hosts {
			if len(host) > 253 {
				errors = append(errors, fmt.Errorf("Host %v is longer than the allowed 253 characters, which is invalid for DNS; please shorten your host", host))
				break
			}

			matchesInvalidChars, _ := regexp.MatchString("[^a-zA-Z0-9-.]", host)
			if matchesInvalidChars {
				errors = append(errors, fmt.Errorf("Host %v has invalid characters; only a-z, 0-9, - and . are allowed; please fix your host", host))
			}

			hostLabels := strings.Split(host, ".")
			for _, label := range hostLabels {
				if len(label) > 63 {
					errors = append(errors, fmt.Errorf("Host %v has label %v - the parts between dots - that is longer than the allowed 63 characters, which is invalid for DNS; please shorten your host label", host, label))
				}
			}
		}

		for _, host := range p.InternalHosts {
			if len(host) > 253 {
				errors = append(errors, fmt.Errorf("Internal host %v is longer than the allowed 253 characters, which is invalid for DNS; please shorten your host", host))
				break
			}

			matchesInvalidChars, _ := regexp.MatchString("[^a-zA-Z0-9-.]", host)
			if matchesInvalidChars {
				errors = append(errors, fmt.Errorf("Internal host %v has invalid characters; only a-z, 0-9, - and . are allowed; please fix your host", host))
			}

			hostLabels := strings.Split(host, ".")
			for _, label := range hostLabels {
				if len(label) > 63 {
					errors = append(errors, fmt.Errorf("Internal host %v has label %v - the parts between dots - that is longer than the allowed 63 characters, which is invalid for DNS; please shorten your host label", host, label))
				}
			}
		}
	}

	if p.Basepath == "" {
		errors = append(errors, fmt.Errorf("Basepath property is required; set it via basepath property on this stage"))
	}
	if p.Container.Port <= 0 {
		errors = append(errors, fmt.Errorf("Container port must be larger than zero; set it via container.port property on this stage"))
	}

	// validate autoscale params
	if p.Autoscale.MinReplicas <= 0 {
		errors = append(errors, fmt.Errorf("Autoscaling min replicas must be larger than zero; set it via autoscale.min property on this stage"))
	}
	if p.Autoscale.MaxReplicas <= 0 {
		errors = append(errors, fmt.Errorf("Autoscaling max replicas must be larger than zero; set it via autoscale.max property on this stage"))
	}
	if p.Autoscale.CPUPercentage <= 0 {
		errors = append(errors, fmt.Errorf("Autoscaling cpu percentage must be larger than zero; set it via autoscale.cpu property on this stage"))
	}

	// validate liveness params
	if p.Container.LivenessProbe.Path == "" {
		errors = append(errors, fmt.Errorf("Liveness path is required; set it via container.liveness.path property on this stage"))
	}
	if p.Container.LivenessProbe.Port <= 0 {
		errors = append(errors, fmt.Errorf("Liveness port must be larger than zero; set it via container.liveness.port property on this stage"))
	}
	if p.Container.LivenessProbe.InitialDelaySeconds <= 0 {
		errors = append(errors, fmt.Errorf("Liveness initial delay must be larger than zero; set it via container.liveness.delay property on this stage"))
	}
	if p.Container.LivenessProbe.TimeoutSeconds <= 0 {
		errors = append(errors, fmt.Errorf("Liveness timeout must be larger than zero; set it via container.liveness.timeout property on this stage"))
	}
	if p.Container.LivenessProbe.PeriodSeconds <= 0 {
		errors = append(errors, fmt.Errorf("Liveness period must be larger than zero; set it via container.liveness.period property on this stage"))
	}

	// validate readiness params
	if p.Container.ReadinessProbe.Path == "" {
		errors = append(errors, fmt.Errorf("Readiness path is required; set it via container.readiness.path property on this stage"))
	}
	if p.Container.ReadinessProbe.Port <= 0 {
		errors = append(errors, fmt.Errorf("Readiness port must be larger than zero; set it via container.readiness.port property on this stage"))
	}
	if p.Container.ReadinessProbe.TimeoutSeconds <= 0 {
		errors = append(errors, fmt.Errorf("Readiness timeout must be larger than zero; set it via container.readiness.timeout property on this stage"))
	}
	if p.Container.ReadinessProbe.PeriodSeconds <= 0 {
		errors = append(errors, fmt.Errorf("Readiness period must be larger than zero; set it via container.liveness.period property on this stage"))
	}

	// validate metrics params
	if p.Container.Metrics.Scrape == nil {
		errors = append(errors, fmt.Errorf("Metrics scrape is required; set it via container.metrics.scrape property on this stage; allowed values are true or false"))
	}
	if p.Container.Metrics.Scrape != nil && *p.Container.Metrics.Scrape {
		if p.Container.Metrics.Path == "" {
			errors = append(errors, fmt.Errorf("Metrics path is required; set it via container.metrics.path property on this stage"))
		}
		if p.Container.Metrics.Port <= 0 {
			errors = append(errors, fmt.Errorf("Metrics port must be larger than zero; set it via container.metrics.port property on this stage"))
		}
	}

	// The "sidecar" field is deprecated, so it can be empty. But if it's specified, then we validate it.
	if p.Sidecar.Type != "" && p.Sidecar.Type != "none" {
		errors = p.validateSidecar(&p.Sidecar, errors)
		warnings = append(warnings, "The sidecar field is deprecated, the sidecars list should be used instead.")
	}

	// check if openresty was defined as deprecated sidecar type
	hasOpenrestySidecar := p.Sidecar.Type == SidecarTypeOpenresty

	// validate sidecars params
	for _, sidecar := range p.Sidecars {
		errors = p.validateSidecar(sidecar, errors)
		if sidecar.Type == SidecarTypeOpenresty {
			hasOpenrestySidecar = true
		}
	}

	// openresty sidecar cannot be added in combination with port 443
	if hasOpenrestySidecar && p.Container.Port == 443 {
		errors = append(errors, fmt.Errorf("Container port can't be 443 if an openresty sidecar is injected"))
	}

	// validate load balance algorithm
	if p.Request.LoadBalanceAlgorithm != "" && p.Request.LoadBalanceAlgorithm != "ewma" && p.Request.LoadBalanceAlgorithm != "round_robin" {
		errors = append(errors, fmt.Errorf("Load balance algorithm is invalid; leave it empty or set request.loadbalance property on this stage to 'ewma' or 'round_robin'"))
	}

	// check for visibility esp if openapi.yaml exists
	if _, err := os.Stat(p.EspOpenAPIYamlPath); (p.Visibility == VisibilityESP || p.Visibility == VisibilityESPv2) && os.IsNotExist(err) {
		errors = append(errors, fmt.Errorf("When using visibility: esp make sure to set clone: true and have openapi.yaml available in the working directory"))
	}

	return len(errors) == 0, errors, warnings
}

func (p *Params) validateSidecar(sidecar *SidecarParams, errors []error) []error {
	switch sidecar.Type {
	case SidecarTypeOpenresty:
		break
	case SidecarTypeCloudSQLProxy:
		if sidecar.DbInstanceConnectionName == "" {
			errors = append(errors, fmt.Errorf("The name of the DB instance used by this Cloud SQL Proxy is required; set it via sidecar.dbinstanceconnectionname property on this stage"))
		}
		if sidecar.SQLProxyPort == 0 {
			errors = append(errors, fmt.Errorf("The port on which the Cloud SQL Proxy listens is required; set it via sidecar.sqlproxyport property on this stage"))
		}
	case SidecarTypeUnknown:
		errors = append(errors, fmt.Errorf("The sidecar type is empty; set a type"))
	}

	if sidecar.Image == "" {
		errors = append(errors, fmt.Errorf("Sidecar image is required; set it via sidecar.image property on this stage"))
	}

	// validate sidecar cpu params
	if sidecar.CPU.Request == "" {
		errors = append(errors, fmt.Errorf("Sidecar cpu request is required; set it via sidecar.cpu.request property on this stage"))
	}

	// validate sidecar memory params
	if sidecar.Memory.Request == "" {
		errors = append(errors, fmt.Errorf("Sidecar memory request is required; set it via sidecar.memory.request property on this stage"))
	}
	if sidecar.Memory.Limit == "" {
		errors = append(errors, fmt.Errorf("Sidecar memory limit is required; set it via sidecar.memory.limit property on this stage"))
	}

	return errors
}

// ReplaceSidecarTagsWithDigest replaces image tags for sidecars with a digest
func (p *Params) ReplaceSidecarTagsWithDigest() {

	// see if there's a sidecar of type openresty
	for _, s := range p.Sidecars {
		log.Info().Msgf("Replacing sidecar %v image tag with digest...", s.Type)

		if s.Image == "" {
			continue
		}

		imageDigestParts := strings.Split(s.Image, "@")
		if len(imageDigestParts) > 1 {
			// already uses a digest, skip replacement
			continue
		}

		imageParts := strings.Split(s.Image, ":")
		repository := imageParts[0]
		tag := "latest"
		if len(imageParts) > 1 {
			tag = imageParts[1]
		}

		// get docker hub api token
		tokenJSON := httpRequestBody("GET", fmt.Sprintf("https://auth.docker.io/token?scope=repository:%v:pull&service=registry.docker.io", repository), map[string]string{})
		if tokenJSON == "" {
			continue
		}

		type TokenObject struct {
			Token string `json:"token"`
		}

		tokenObject := TokenObject{}
		err := json.Unmarshal([]byte(tokenJSON), &tokenObject)
		if err != nil {
			continue
		}
		if tokenObject.Token == "" {
			continue
		}

		digest := httpRequestHeader("HEAD", fmt.Sprintf("https://index.docker.io/v2/%v/manifests/%v", repository, tag), map[string]string{
			"Accept":        "application/vnd.docker.distribution.manifest.v2+json",
			"Authorization": fmt.Sprintf("Bearer %v", tokenObject.Token),
		}, "Docker-Content-Digest")

		if len(digest) == 0 {
			continue
		}

		s.Image = fmt.Sprintf("%v@%v", repository, digest)

		log.Info().Msgf("Successfully replaced tag %v with digest %v...", tag, digest)
	}
}
