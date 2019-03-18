package main

import (
	"fmt"
	"regexp"
	"strings"
)

// Params is used to parameterize the deployment, set from custom properties in the manifest
type Params struct {
	// control params
	Action          string          `json:"action,omitempty"`
	Kind            string          `json:"kind,omitempty"`
	DryRun          bool            `json:"dryrun,omitempty"`
	BuildVersion    string          `json:"-"`
	ChaosProof      bool            `json:"chaosproof,omitempty"`
	Manifests       ManifestsParams `json:"manifests,omitempty"`
	TrustedIPRanges []string        `json:"trustedips,omitempty"`

	// app params
	App            string              `json:"app,omitempty"`
	Namespace      string              `json:"namespace,omitempty"`
	Schedule       string              `json:"schedule,omitempty"`
	Labels         map[string]string   `json:"labels,omitempty"`
	Visibility     string              `json:"visibility,omitempty"`
	WhitelistedIPS []string            `json:"whitelist,omitempty"`
	Hosts          []string            `json:"hosts,omitempty"`
	InternalHosts  []string            `json:"internalhosts,omitempty"`
	Basepath       string              `json:"basepath,omitempty"`
	Autoscale      AutoscaleParams     `json:"autoscale,omitempty"`
	Request        RequestParams       `json:"request,omitempty"`
	Secrets        SecretsParams       `json:"secrets,omitempty"`
	Configs        ConfigsParams       `json:"configs,omitempty"`
	VolumeMounts   []VolumeMountParams `json:"volumemounts,omitempty"`

	EnablePayloadLogging      bool `json:"enablePayloadLogging,omitempty"`
	UseGoogleCloudCredentials bool `json:"useGoogleCloudCredentials,omitempty"`

	// container params
	Container     ContainerParams     `json:"container,omitempty"`
	Sidecar       SidecarParams       `json:"sidecar,omitempty"`
	RollingUpdate RollingUpdateParams `json:"rollingupdate,omitempty"`
}

// ContainerParams defines the container image to deploy
type ContainerParams struct {
	ImageRepository      string                 `json:"repository,omitempty"`
	ImageName            string                 `json:"name,omitempty"`
	ImageTag             string                 `json:"tag,omitempty"`
	Port                 int                    `json:"port,omitempty"`
	EnvironmentVariables map[string]interface{} `json:"env,omitempty"`

	CPU            CPUParams       `json:"cpu,omitempty"`
	Memory         MemoryParams    `json:"memory,omitempty"`
	LivenessProbe  ProbeParams     `json:"liveness,omitempty"`
	ReadinessProbe ProbeParams     `json:"readiness,omitempty"`
	Metrics        MetricsParams   `json:"metrics,omitempty"`
	Lifecycle      LifecycleParams `json:"lifecycle,omitempty"`

	AdditionalPorts []*AdditionalPortParams `json:"additionalports,omitempty"`
}

// AdditionalPortParams provides information about any additional ports exposed and accessible via a service
type AdditionalPortParams struct {
	Name       string `json:"name,omitempty"`
	Port       int    `json:"port,omitempty"`
	Protocol   string `json:"protocol,omitempty"`
	Visibility string `json:"visibility,omitempty"`
}

// CPUParams sets cpu request and limit values
type CPUParams struct {
	Request string `json:"request,omitempty"`
	Limit   string `json:"limit,omitempty"`
}

// MemoryParams sets memory request and limit values
type MemoryParams struct {
	Request string `json:"request,omitempty"`
	Limit   string `json:"limit,omitempty"`
}

// AutoscaleParams controls autoscaling
type AutoscaleParams struct {
	MinReplicas   int `json:"min,omitempty"`
	MaxReplicas   int `json:"max,omitempty"`
	CPUPercentage int `json:"cpu,omitempty"`
}

// RequestParams controls timeouts, max body size, etc
type RequestParams struct {
	Timeout              string `json:"timeout,omitempty"`
	MaxBodySize          string `json:"maxbodysize,omitempty"`
	ProxyBufferSize      string `json:"proxybuffersize,omitempty"`
	ProxyBuffersNumber   int    `json:"proxybuffersnumber,omitempty"`
	ClientBodyBufferSize string `json:"clientbodybuffersize,omitempty"`
}

// ProbeParams sets params for liveness or readiness probe
type ProbeParams struct {
	Path                string `json:"path,omitempty"`
	Port                int    `json:"port,omitempty"`
	InitialDelaySeconds int    `json:"delay,omitempty"`
	TimeoutSeconds      int    `json:"timeout,omitempty"`
}

// MetricsParams sets params for scraping prometheus metrics
type MetricsParams struct {
	Scrape *bool  `json:"scrape,omitempty"`
	Path   string `json:"path,omitempty"`
	Port   int    `json:"port,omitempty"`
}

// LifecycleParams sets params for lifecycle commands
type LifecycleParams struct {
	PrestopSleep        *bool `json:"prestopsleep,omitempty"`
	PrestopSleepSeconds *int  `json:"prestopsleepseconds,omitempty"`
}

// SidecarParams sets params for sidecar injection
type SidecarParams struct {
	Type                 string                 `json:"type,omitempty"`
	Image                string                 `json:"image,omitempty"`
	EnvironmentVariables map[string]interface{} `json:"env,omitempty"`
	CPU                  CPUParams              `json:"cpu,omitempty"`
	Memory               MemoryParams           `json:"memory,omitempty"`
	HealthCheckPath      string                 `json:"healthcheckpath,omitempty"`
}

// RollingUpdateParams sets params for controlling rolling update speed
type RollingUpdateParams struct {
	MaxSurge       string `json:"maxsurge,omitempty"`
	MaxUnavailable string `json:"maxunavailable,omitempty"`
	Timeout        string `json:"timeout,omitempty"`
}

// ManifestsParams can be used to override or add additional manifests located in the application repository
type ManifestsParams struct {
	Files []string               `json:"files,omitempty"`
	Data  map[string]interface{} `json:"data,omitempty"`
}

// SecretsParams allows secrets to be set dynamically for the application
type SecretsParams struct {
	Keys      map[string]interface{} `json:"keys,omitempty"`
	MountPath string                 `json:"mountpath,omitempty"`
}

// ConfigsParams allows configs to be set dynamically for the application
type ConfigsParams struct {
	Files               []string               `json:"files,omitempty"`
	Data                map[string]interface{} `json:"data,omitempty"`
	InlineFiles         map[string]string      `json:"inline,omitempty"`
	MountPath           string                 `json:"mountpath,omitempty"`
	RenderedFileContent map[string]string      `json:"-"`
}

// VolumeMountParams allows additional mounts for already existing volumes, secrets, etc
type VolumeMountParams struct {
	Name      string                 `json:"name,omitempty"`
	MountPath string                 `json:"mountpath,omitempty"`
	Volume    map[string]interface{} `json:"volume,omitempty"`
}

// SetDefaults fills in empty fields with convention-based defaults
func (p *Params) SetDefaults(appLabel, buildVersion, releaseName, releaseAction string, estafetteLabels map[string]string) {

	p.BuildVersion = buildVersion

	// default action to deploy-simple unless it's either specified on the stage or passed in as a release action
	if releaseAction != "" {
		p.Action = releaseAction
	} else if p.Action == "" && releaseAction == "" {
		p.Action = "deploy-simple"
	}

	// default kind to deployment
	if p.Kind == "" {
		p.Kind = "deployment"
	}

	// default app to estafette app label if no override in stage params
	if p.App == "" && appLabel != "" {
		p.App = appLabel
	}

	// default image name to estafette app label if no override in stage params
	if p.Container.ImageName == "" && p.App != "" {
		p.Container.ImageName = p.App
	}

	// default image tag to estafette build version if no override in stage params
	if p.Container.ImageTag == "" && buildVersion != "" {
		p.Container.ImageTag = buildVersion
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

	// default visibility to private if no override in stage params
	if p.Visibility == "" {
		p.Visibility = "private"
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
	if p.Container.CPU.Limit == "" {
		if !cpuRequestIsEmpty {
			p.Container.CPU.Limit = p.Container.CPU.Request
		} else {
			p.Container.CPU.Limit = "125m"
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

	// set additional ports defaults
	if len(p.Container.AdditionalPorts) > 0 {
		for _, ap := range p.Container.AdditionalPorts {
			if ap.Protocol == "" {
				ap.Protocol = "TCP"
			}
			if ap.Visibility == "" {
				ap.Visibility = p.Visibility
			}
		}
	}

	// set autoscale defaults
	if p.Autoscale.MinReplicas <= 0 {
		p.Autoscale.MinReplicas = 3
	}
	if p.Autoscale.MaxReplicas <= 0 {
		p.Autoscale.MaxReplicas = 100
	}
	if p.Autoscale.CPUPercentage <= 0 {
		p.Autoscale.CPUPercentage = 80
	}

	// set request defaults
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

	// set readiness probe defaults
	if p.Container.ReadinessProbe.Path == "" {
		p.Container.ReadinessProbe.Path = "/readiness"
	}
	if p.Container.ReadinessProbe.Port <= 0 {
		p.Container.ReadinessProbe.Port = p.Container.Port
	}
	if p.Container.ReadinessProbe.TimeoutSeconds <= 0 {
		p.Container.ReadinessProbe.TimeoutSeconds = 1
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
		trueValue := true
		p.Container.Lifecycle.PrestopSleep = &trueValue
	}
	if p.Container.Lifecycle.PrestopSleepSeconds == nil {
		defaultSleepValue := 15
		p.Container.Lifecycle.PrestopSleepSeconds = &defaultSleepValue
	}

	// set sidecar defaults
	if p.Sidecar.Type == "" {
		switch p.Kind {
		case "job":
			p.Sidecar.Type = "none"

		default:
			p.Sidecar.Type = "openresty"
		}
	}
	if p.Sidecar.Image == "" {
		p.Sidecar.Image = "estafette/openresty-sidecar:1.13.6.2-alpine"
	}
	if p.Sidecar.HealthCheckPath == "" {
		p.Sidecar.HealthCheckPath = p.Container.ReadinessProbe.Path
	}

	// set sidecar cpu defaults
	sidecarCPURequestIsEmpty := p.Sidecar.CPU.Request == ""
	if sidecarCPURequestIsEmpty {
		if p.Sidecar.CPU.Limit != "" {
			p.Sidecar.CPU.Request = p.Sidecar.CPU.Limit
		} else {
			p.Sidecar.CPU.Request = "50m"
		}
	}
	if p.Sidecar.CPU.Limit == "" {
		if !sidecarCPURequestIsEmpty {
			p.Sidecar.CPU.Limit = p.Sidecar.CPU.Request
		} else {
			p.Sidecar.CPU.Limit = "75m"
		}
	}

	// set sidecar memory defaults
	sidecarMemoryRequestIsEmpty := p.Sidecar.Memory.Request == ""
	if sidecarMemoryRequestIsEmpty {
		if p.Sidecar.Memory.Limit != "" {
			p.Sidecar.Memory.Request = p.Sidecar.Memory.Limit
		} else {
			p.Sidecar.Memory.Request = "30Mi"
		}
	}
	if p.Sidecar.Memory.Limit == "" {
		if !sidecarMemoryRequestIsEmpty {
			p.Sidecar.Memory.Limit = p.Sidecar.Memory.Request
		} else {
			p.Sidecar.Memory.Limit = "50Mi"
		}
	}

	// default basepath to /
	if p.Basepath == "" {
		p.Basepath = "/"
	}

	// defaults for rollingupdate
	if p.RollingUpdate.MaxSurge == "" {
		p.RollingUpdate.MaxSurge = "25%"
	}
	if p.RollingUpdate.MaxUnavailable == "" {
		p.RollingUpdate.MaxUnavailable = "0"
	}
	if p.RollingUpdate.Timeout == "" {
		p.RollingUpdate.Timeout = "5m"
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
}

// ValidateRequiredProperties checks whether all needed properties are set
func (p *Params) ValidateRequiredProperties() (bool, []error) {

	errors := []error{}

	// validate app params
	if p.App == "" {
		errors = append(errors, fmt.Errorf("Application name is required; either define an app label or use app property on this stage"))
	}
	if p.Namespace == "" {
		errors = append(errors, fmt.Errorf("Namespace is required; either use credentials with a defaultNamespace or set it via namespace property on this stage"))
	}

	if p.Action == "rollback-canary" {
		// the above properties are all you need for a rollback
		return len(errors) == 0, errors
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

	// validate cpu params
	if p.Container.CPU.Request == "" {
		errors = append(errors, fmt.Errorf("Cpu request is required; set it via container.cpu.request property on this stage"))
	}
	if p.Container.CPU.Limit == "" {
		errors = append(errors, fmt.Errorf("Cpu limit is required; set it via container.cpu.limit property on this stage"))
	}

	// validate memory params
	if p.Container.Memory.Request == "" {
		errors = append(errors, fmt.Errorf("Memory request is required; set it via container.memory.request property on this stage"))
	}
	if p.Container.Memory.Limit == "" {
		errors = append(errors, fmt.Errorf("Memory limit is required; set it via container.memory.limit property on this stage"))
	}

	// defaults for rollingupdate
	if p.RollingUpdate.MaxSurge == "" {
		errors = append(errors, fmt.Errorf("Rollingupdate max surge is required; set it via rollingupdate.maxsurge property on this stage"))
	}
	if p.RollingUpdate.MaxUnavailable == "" {
		errors = append(errors, fmt.Errorf("Rollingupdate max unavailable is required; set it via rollingupdate.maxunavailable property on this stage"))
	}

	if p.Kind == "cronjob" && p.Schedule == "" {
		errors = append(errors, fmt.Errorf("Schedule is required for a cronjob; set it via schedule property on this stage"))
	}

	if p.Kind == "job" || p.Kind == "cronjob" {
		// the above properties are all you need for a worker
		return len(errors) == 0, errors
	}

	// validate params with respect to incoming requests
	if p.Visibility == "" || (p.Visibility != "private" && p.Visibility != "public" && p.Visibility != "iap" && p.Visibility != "public-whitelist") {
		errors = append(errors, fmt.Errorf("Visibility property is required; set it via visibility property on this stage; allowed values are private, iap, public-whitelist or public"))
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

	// validate sidecar params
	if p.Sidecar.Type == "" {
		errors = append(errors, fmt.Errorf("Sidecar type is required; set it via sidecar.type property on this stage; allowed values are openresty"))
	}
	if p.Sidecar.Image == "" {
		errors = append(errors, fmt.Errorf("Sidecar image is required; set it via sidecar.image property on this stage"))
	}

	// validate sidecar cpu params
	if p.Sidecar.CPU.Request == "" {
		errors = append(errors, fmt.Errorf("Sidecar cpu request is required; set it via sidecar.cpu.request property on this stage"))
	}
	if p.Sidecar.CPU.Limit == "" {
		errors = append(errors, fmt.Errorf("Sidecar cpu limit is required; set it via sidecar.cpu.limit property on this stage"))
	}

	// validate sidecar memory params
	if p.Sidecar.Memory.Request == "" {
		errors = append(errors, fmt.Errorf("Sidecar memory request is required; set it via sidecar.memory.request property on this stage"))
	}
	if p.Sidecar.Memory.Limit == "" {
		errors = append(errors, fmt.Errorf("Sidecar memory limit is required; set it via sidecar.memory.limit property on this stage"))
	}

	return len(errors) == 0, errors
}
