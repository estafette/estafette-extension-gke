package main

import (
	"fmt"
)

// Params is used to parameterize the deployment, set from custom properties in the manifest
type Params struct {
	// which credentials to use
	Credentials string `json:"credentials,omitempty"`

	// application common properties
	App       string            `json:"app,omitempty"`
	Namespace string            `json:"namespace,omitempty"`
	Labels    map[string]string `json:"labels,omitempty"`

	// container specific properties
	Container ContainerParams `json:"container,omitempty"`

	// misc
	Visibility string          `json:"visibility,omitempty"`
	Hosts      []string        `json:"hosts,omitempty"`
	Autoscale  AutoscaleParams `json:"autoscale,omitempty"`

	// resources
	CPU            CPUParams    `json:"cpu,omitempty"`
	Memory         MemoryParams `json:"memory,omitempty"`
	LivenessProbe  ProbeParams  `json:"liveness,omitempty"`
	ReadinessProbe ProbeParams  `json:"readiness,omitempty"`

	// used for seeing the rendered template without executing it but testing it with a dryrun
	DryRun bool `json:"dryrun,string,omitempty"`

	// IngressPath         string
	// UseNginxIngress     bool
	// UseGCEIngress       bool
	// ServiceType         string
	// PreferPreemptibles  bool
}

// ContainerParams defines the container image to deploy
type ContainerParams struct {
	ImageRepository string `json:"repository,omitempty"`
	ImageName       string `json:"name,omitempty"`
	ImageTag        string `json:"tag,omitempty"`
	Port            int    `json:"port,string,omitempty"`
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
	MinReplicas   int `json:"min,string,omitempty"`
	MaxReplicas   int `json:"max,string,omitempty"`
	CPUPercentage int `json:"cpu,string,omitempty"`
}

// ProbeParams sets params for liveness or readiness probe
type ProbeParams struct {
	Path                string `json:"path,omitempty"`
	InitialDelaySeconds int    `json:"delay,string,omitempty"`
	TimeoutSeconds      int    `json:"timeout,string,omitempty"`
}

// SetDefaults fills in empty fields with convention-based defaults
func (p *Params) SetDefaults(appLabel, buildVersion, releaseName string, estafetteLabels map[string]string) {

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

	// default credentials to release name if no override in stage params
	if p.Credentials == "" && releaseName != "" {
		p.Credentials = fmt.Sprintf("gke-%v", releaseName)
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
	cpuRequestIsEmpty := p.CPU.Request == ""
	if cpuRequestIsEmpty {
		if p.CPU.Limit != "" {
			p.CPU.Request = p.CPU.Limit
		} else {
			p.CPU.Request = "100m"
		}
	}
	if p.CPU.Limit == "" {
		if !cpuRequestIsEmpty {
			p.CPU.Limit = p.CPU.Request
		} else {
			p.CPU.Limit = "125m"
		}
	}

	// set memory defaults
	memoryRequestIsEmpty := p.Memory.Request == ""
	if memoryRequestIsEmpty {
		if p.Memory.Limit != "" {
			p.Memory.Request = p.Memory.Limit
		} else {
			p.Memory.Request = "128Mi"
		}
	}
	if p.Memory.Limit == "" {
		if !memoryRequestIsEmpty {
			p.Memory.Limit = p.Memory.Request
		} else {
			p.Memory.Limit = "128Mi"
		}
	}

	// set container port defaults
	if p.Container.Port <= 0 {
		p.Container.Port = 5000
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

	// set probe defaults
	if p.LivenessProbe.Path == "" {
		p.LivenessProbe.Path = "/liveness"
	}
	if p.LivenessProbe.InitialDelaySeconds <= 0 {
		p.LivenessProbe.InitialDelaySeconds = 30
	}
	if p.LivenessProbe.TimeoutSeconds <= 0 {
		p.LivenessProbe.TimeoutSeconds = 1
	}

	if p.ReadinessProbe.Path == "" {
		p.ReadinessProbe.Path = "/readiness"
	}
	if p.ReadinessProbe.TimeoutSeconds <= 0 {
		p.ReadinessProbe.TimeoutSeconds = 1
	}
}

// SetDefaultsFromCredentials sets defaults based on the credentials fetched with first-run defaults
func (p *Params) SetDefaultsFromCredentials(credentials GKECredentials) {

	// default namespace to credential default namespace if no override in stage params
	if p.Namespace == "" && credentials.AdditionalProperties.DefaultNamespace != "" {
		p.Namespace = credentials.AdditionalProperties.DefaultNamespace
	}

	// default image repository to credential project if no override in stage params
	if p.Container.ImageRepository == "" && credentials.AdditionalProperties.Project != "" {
		p.Container.ImageRepository = credentials.AdditionalProperties.Project
	}
}

// ValidateRequiredProperties checks whether all needed properties are set
func (p *Params) ValidateRequiredProperties() (bool, []error) {

	errors := []error{}

	if p.App == "" {
		errors = append(errors, fmt.Errorf("Application name is required; either define an app label or use app property on this stage"))
	}
	if p.Namespace == "" {
		errors = append(errors, fmt.Errorf("Namespace is required; either use credentials with a defaultNamespace or set it via namespace property on this stage"))
	}
	if p.Container.ImageRepository == "" {
		errors = append(errors, fmt.Errorf("Image repository is required; set it via container.repository property on this stage"))
	}
	if p.Container.ImageName == "" {
		errors = append(errors, fmt.Errorf("Image name is required; set it via container.name property on this stage"))
	}
	if p.Container.ImageTag == "" {
		errors = append(errors, fmt.Errorf("Image tag is required; set it via container.tag property on this stage"))
	}
	if p.Container.Port <= 0 {
		errors = append(errors, fmt.Errorf("Container port must be larger than zero; set it via container.port property on this stage"))
	}
	if p.Credentials == "" {
		errors = append(errors, fmt.Errorf("Credentials property is required; set it via credentials property on this stage"))
	}
	if p.Visibility == "" || (p.Visibility != "private" && p.Visibility != "public") {
		errors = append(errors, fmt.Errorf("Visibility property is required; set it via visibility property on this stage; allowed values are private or public"))
	}
	if p.CPU.Request == "" {
		errors = append(errors, fmt.Errorf("Cpu request is required; set it via cpu.request property on this stage"))
	}
	if p.CPU.Limit == "" {
		errors = append(errors, fmt.Errorf("Cpu limit is required; set it via cpu.limit property on this stage"))
	}
	if p.Memory.Request == "" {
		errors = append(errors, fmt.Errorf("Memory request is required; set it via memory.request property on this stage"))
	}
	if p.Memory.Limit == "" {
		errors = append(errors, fmt.Errorf("Memory limit is required; set it via memory.limit property on this stage"))
	}
	if len(p.Hosts) == 0 {
		errors = append(errors, fmt.Errorf("At least one host is required; set it via hosts array property on this stage"))
	}
	if p.Autoscale.MinReplicas <= 0 {
		errors = append(errors, fmt.Errorf("Autoscaling min replicas must be larger than zero; set it via autoscale.min property on this stage"))
	}
	if p.Autoscale.MaxReplicas <= 0 {
		errors = append(errors, fmt.Errorf("Autoscaling max replicas must be larger than zero; set it via autoscale.max property on this stage"))
	}
	if p.Autoscale.CPUPercentage <= 0 {
		errors = append(errors, fmt.Errorf("Autoscaling cpu percentage must be larger than zero; set it via autoscale.cpu property on this stage"))
	}

	if p.LivenessProbe.Path == "" {
		errors = append(errors, fmt.Errorf("Liveness path is required; set it via liveness.path property on this stage"))
	}
	if p.LivenessProbe.InitialDelaySeconds <= 0 {
		errors = append(errors, fmt.Errorf("Liveness initial delay must be larger than zero; set it via liveness.delay property on this stage"))
	}
	if p.LivenessProbe.TimeoutSeconds <= 0 {
		errors = append(errors, fmt.Errorf("Liveness timeout must be larger than zero; set it via liveness.timeout property on this stage"))
	}

	if p.ReadinessProbe.Path == "" {
		errors = append(errors, fmt.Errorf("Readiness path is required; set it via readiness.path property on this stage"))
	}
	if p.ReadinessProbe.TimeoutSeconds <= 0 {
		errors = append(errors, fmt.Errorf("Readiness timeout must be larger than zero; set it via readiness.timeout property on this stage"))
	}

	return len(errors) == 0, errors
}
