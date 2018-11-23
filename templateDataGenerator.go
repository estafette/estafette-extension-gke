package main

import (
	"strings"

	"gopkg.in/yaml.v2"
)

func generateTemplateData(params Params) TemplateData {

	data := TemplateData{
		BuildVersion: params.BuildVersion,

		Name:             params.App,
		NameWithTrack:    params.App,
		Namespace:        params.Namespace,
		Labels:           params.Labels,
		AppLabelSelector: params.App,

		Hosts:       params.Hosts,
		HostsJoined: strings.Join(params.Hosts, ","),
		IngressPath: params.Basepath,

		MinReplicas:         params.Autoscale.MinReplicas,
		MaxReplicas:         params.Autoscale.MaxReplicas,
		TargetCPUPercentage: params.Autoscale.CPUPercentage,

		Secrets:                 params.Secrets.Keys,
		MountApplicationSecrets: len(params.Secrets.Keys) > 0,
		SecretMountPath:         params.Secrets.MountPath,
		MountConfigmap:          len(params.Configs.Files) > 0,
		ConfigMountPath:         params.Configs.MountPath,

		MountPayloadLogging: params.EnablePayloadLogging,

		RollingUpdateMaxSurge:       params.RollingUpdate.MaxSurge,
		RollingUpdateMaxUnavailable: params.RollingUpdate.MaxUnavailable,

		PreferPreemptibles: params.ChaosProof,

		Container: ContainerData{
			Repository: params.Container.ImageRepository,
			Name:       params.Container.ImageName,
			Tag:        params.Container.ImageTag,
			Port:       params.Container.Port,

			CPURequest:    params.Container.CPU.Request,
			CPULimit:      params.Container.CPU.Limit,
			MemoryRequest: params.Container.Memory.Request,
			MemoryLimit:   params.Container.Memory.Limit,

			EnvironmentVariables: params.Container.EnvironmentVariables,

			Liveness: ProbeData{
				Path:                params.Container.LivenessProbe.Path,
				Port:                params.Container.LivenessProbe.Port,
				InitialDelaySeconds: params.Container.LivenessProbe.InitialDelaySeconds,
				TimeoutSeconds:      params.Container.LivenessProbe.TimeoutSeconds,
				IncludeOnContainer:  true,
			},
			Readiness: ProbeData{
				Path:                params.Container.ReadinessProbe.Path,
				Port:                params.Container.ReadinessProbe.Port,
				InitialDelaySeconds: params.Container.ReadinessProbe.InitialDelaySeconds,
				TimeoutSeconds:      params.Container.ReadinessProbe.TimeoutSeconds,
				IncludeOnContainer:  params.Sidecar.Type != "openresty" || params.Container.ReadinessProbe.Port != params.Container.Port || params.Container.ReadinessProbe.Path != params.Sidecar.HealthCheckPath,
			},
			Metrics: MetricsData{
				Path: params.Container.Metrics.Path,
				Port: params.Container.Metrics.Port,
			},
		},

		Sidecar: SidecarData{
			UseOpenrestySidecar: params.Sidecar.Type == "openresty",

			Image:           params.Sidecar.Image,
			HealthCheckPath: params.Sidecar.HealthCheckPath,
			CPURequest:      params.Sidecar.CPU.Request,
			CPULimit:        params.Sidecar.CPU.Limit,
			MemoryRequest:   params.Sidecar.Memory.Request,
			MemoryLimit:     params.Sidecar.Memory.Limit,

			EnvironmentVariables: params.Sidecar.EnvironmentVariables,
		},
	}

	if params.Container.Metrics.Scrape != nil {
		data.Container.Metrics.Scrape = *params.Container.Metrics.Scrape
	}

	switch params.Action {
	case "deploy-simple":
		data.IncludeTrackLabel = false
	case "deploy-canary":
		data.NameWithTrack += "-canary"
		data.IncludeTrackLabel = true
		data.TrackLabel = "canary"
	case "deploy-stable":
		data.NameWithTrack += "-stable"
		data.IncludeTrackLabel = true
		data.TrackLabel = "stable"
	}

	data.ConfigmapFiles = params.Configs.RenderedFileContent

	data.ManifestData = map[string]interface{}{}
	for k, v := range params.Manifests.Data {
		data.ManifestData[k] = v
	}

	if params.Visibility == "private" {
		data.ServiceType = "ClusterIP"
		data.UseNginxIngress = true
		data.UseGCEIngress = false
		data.UseDNSAnnotationsOnIngress = true
		data.UseDNSAnnotationsOnService = false
		data.LimitTrustedIPRanges = false
	} else if params.Visibility == "iap" {
		data.ServiceType = "NodePort"
		data.UseNginxIngress = false
		data.UseGCEIngress = true
		data.UseDNSAnnotationsOnIngress = true
		data.UseDNSAnnotationsOnService = false
		data.LimitTrustedIPRanges = false
	} else if params.Visibility == "public" {
		data.ServiceType = "LoadBalancer"
		data.UseNginxIngress = false
		data.UseGCEIngress = false
		data.UseDNSAnnotationsOnIngress = false
		data.UseDNSAnnotationsOnService = true
		data.LimitTrustedIPRanges = true
	}

	if !strings.HasSuffix(data.IngressPath, "/") && !strings.HasSuffix(data.IngressPath, "*") {
		data.IngressPath += "/"
	}
	if data.UseGCEIngress && !strings.HasSuffix(data.IngressPath, "*") {
		data.IngressPath += "*"
	}

	data.TrustedIPRanges = params.TrustedIPRanges

	data.AdditionalVolumeMounts = []VolumeMountData{}
	for _, vm := range params.VolumeMounts {
		yamlBytes, err := yaml.Marshal(vm.Volume)
		if err == nil {
			data.AdditionalVolumeMounts = append(data.AdditionalVolumeMounts, VolumeMountData{
				Name:       vm.Name,
				MountPath:  vm.MountPath,
				VolumeYAML: string(yamlBytes),
			})
		}
	}

	return data
}
