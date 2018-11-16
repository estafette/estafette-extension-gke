package main

import (
	"strings"
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
				InitialDelaySeconds: params.Container.LivenessProbe.InitialDelaySeconds,
				TimeoutSeconds:      params.Container.LivenessProbe.TimeoutSeconds,
			},
			Readiness: ProbeData{
				Path:                params.Container.ReadinessProbe.Path,
				InitialDelaySeconds: params.Container.ReadinessProbe.InitialDelaySeconds,
				TimeoutSeconds:      params.Container.ReadinessProbe.TimeoutSeconds,
			},
			Metrics: MetricsData{
				Scrape: params.Container.Metrics.Scrape,
				Path:   params.Container.Metrics.Path,
				Port:   params.Container.Metrics.Port,
			},
		},

		Sidecar: SidecarData{
			UseOpenrestySidecar: params.Sidecar.Type == "openresty",

			Image:         params.Sidecar.Image,
			CPURequest:    params.Sidecar.CPU.Request,
			CPULimit:      params.Sidecar.CPU.Limit,
			MemoryRequest: params.Sidecar.Memory.Request,
			MemoryLimit:   params.Sidecar.Memory.Limit,

			EnvironmentVariables: params.Sidecar.EnvironmentVariables,
		},
	}

	switch params.Type {
	case "simple":
		data.IncludeTrackLabel = false
	case "canary":
		data.NameWithTrack += "-canary"
		data.IncludeTrackLabel = true
		data.TrackLabel = "canary"
		data.MinReplicas = 1
		data.MaxReplicas = 1
	case "rollforward":
		data.NameWithTrack += "-stable"
		data.IncludeTrackLabel = true
		data.TrackLabel = "stable"
	case "rollback":
		data.NameWithTrack += "-canary"
		data.MinReplicas = 0
		data.MaxReplicas = 0
	}

	data.ConfigmapFiles = params.Configs.RenderedFileContent

	data.ManifestData = map[string]string{}
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

	return data
}
