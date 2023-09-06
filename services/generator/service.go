package generator

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"text/template"

	"github.com/estafette/estafette-extension-gke/api"
	yaml "gopkg.in/yaml.v2"
)

//go:generate mockgen -package=generator -destination ./mock.go -source=service.go
type Service interface {
	GenerateTemplateData(params api.Params, currentReplicas int, gitSource, gitOwner, gitName, gitBranch, gitRevision, releaseID, builderImageSHA, builderImageDate, triggeredBy string) api.TemplateData
	BuildSidecar(sidecar *api.SidecarParams, params api.Params) api.SidecarData
	AddEnvironmentVariableIfNotSet(environmentVariables map[string]interface{}, name, value string) map[string]interface{}
	IsSimpleEnvvarValue(i interface{}) bool
	ToYAML(v interface{}) string
	RenderToYAML(v interface{}, data interface{}) string
}

// NewService returns a new extension.Service
func NewService(ctx context.Context) (Service, error) {
	return &service{}, nil
}

type service struct {
}

func (s *service) GenerateTemplateData(params api.Params, currentReplicas int, gitSource, gitOwner, gitName, gitBranch, gitRevision, releaseID, builderImageSHA, builderImageDate, triggeredBy string) api.TemplateData {

	data := api.TemplateData{
		Name:                    params.App,
		NameWithTrack:           params.App,
		Namespace:               params.Namespace,
		Schedule:                params.Schedule,
		ConcurrencyPolicy:       params.ConcurrencyPolicy,
		RestartPolicy:           params.RestartPolicy,
		Completions:             params.Completions,
		Parallelism:             params.Parallelism,
		ProgressDeadlineSeconds: params.ProgressDeadlineSeconds,
		Labels:                  api.SanitizeLabels(params.Labels),
		PodLabels:               api.SanitizeLabels(params.Labels),
		AppLabelSelector:        api.SanitizeLabel(params.App),

		Hosts:               params.Hosts,
		HostsJoined:         strings.Join(params.Hosts, ","),
		InternalHosts:       params.InternalHosts,
		InternalHostsJoined: strings.Join(params.InternalHosts, ","),
		AllHosts:            append(params.Hosts, params.InternalHosts...),
		AllHostsJoined:      strings.Join(append(params.Hosts, params.InternalHosts...), ","),
		IngressPath:         params.Basepath,
		PathType:            "Prefix",
		InternalIngressPath: params.Basepath,
		AllowHTTP:           params.AllowHTTP,

		IncludeReplicas: currentReplicas > 0 || ((params.Autoscale.Enabled == nil || !*params.Autoscale.Enabled || params.StrategyType == "Recreate") && params.Replicas > 0),

		MinReplicas:         params.Autoscale.MinReplicas,
		MaxReplicas:         params.Autoscale.MaxReplicas,
		TargetCPUPercentage: params.Autoscale.CPUPercentage,

		UseHpaScaler:                params.Autoscale.Safety.Enabled,
		HpaScalerPromQuery:          params.Autoscale.Safety.PromQuery,
		HpaScalerRequestsPerReplica: params.Autoscale.Safety.Ratio,
		HpaScalerDelta:              params.Autoscale.Safety.Delta,
		HpaScalerScaleDownMaxRatio:  params.Autoscale.Safety.ScaleDownRatio,

		VpaUpdateMode: string(params.VerticalPodAutoscaler.UpdateMode),

		Secrets:                 params.Secrets.Keys,
		MountSslCertificate:     params.Kind == api.KindDeployment,
		MountApplicationSecrets: params.HasSecrets(),
		SecretMountPath:         params.Secrets.MountPath,
		MountConfigmap:          len(params.Configs.Files) > 0 || len(params.Configs.InlineFiles) > 0,
		ConfigMountPath:         params.Configs.MountPath,
		Tolerations:             []*map[string]interface{}{},
		Affinity:                []*map[string]interface{}{},
		PodSecurityContext:      params.PodSecurityContext,

		MountPayloadLogging:      params.EnablePayloadLogging,
		AddSafeToEvictAnnotation: params.EnablePayloadLogging,

		RollingUpdateMaxSurge:       params.RollingUpdate.MaxSurge,
		RollingUpdateMaxUnavailable: params.RollingUpdate.MaxUnavailable,

		PreferPreemptibles:            params.ChaosProof,
		UseWindowsNodes:               params.OperatingSystem == api.OperatingSystemWindows,
		MountServiceAccountSecret:     params.UseGoogleCloudCredentials || params.LegacyGoogleCloudServiceAccountKeyFile != "",
		UseLegacyServiceAccountKey:    params.LegacyGoogleCloudServiceAccountKeyFile != "",
		GoogleCloudCredentialsAppName: params.GoogleCloudCredentialsApp,
		GoogleCloudCredentialsLabels:  api.SanitizeLabels(params.Labels),

		PodManagementPolicy: params.PodManagementPolicy,
		StorageClass:        params.StorageClass,
		StorageSize:         params.StorageSize,
		StorageMountPath:    params.StorageMountPath,

		Container: api.ContainerData{
			Repository:      params.Container.ImageRepository,
			Name:            params.Container.ImageName,
			Tag:             params.Container.ImageTag,
			ImagePullPolicy: params.Container.ImagePullPolicy,
			Port:            params.Container.Port,
			PortGrpc:        params.Container.PortGrpc,

			CPURequest:    params.Container.CPU.Request,
			CPULimit:      params.Container.CPU.Limit,
			MemoryRequest: params.Container.Memory.Request,
			MemoryLimit:   params.Container.Memory.Limit,

			EnvironmentVariables:       params.Container.EnvironmentVariables,
			SecretEnvironmentVariables: params.Container.SecretEnvironmentVariables,

			ContainerSecurityContext: params.Container.ContainerSecurityContext,

			ContainerLifeCycle: params.Container.ContainerLifeCycle,

			Liveness: api.ProbeData{
				Path:                params.Container.LivenessProbe.Path,
				Port:                params.Container.LivenessProbe.Port,
				InitialDelaySeconds: params.Container.LivenessProbe.InitialDelaySeconds,
				TimeoutSeconds:      params.Container.LivenessProbe.TimeoutSeconds,
				PeriodSeconds:       params.Container.LivenessProbe.PeriodSeconds,
				FailureThreshold:    params.Container.LivenessProbe.FailureThreshold,
				SuccessThreshold:    params.Container.LivenessProbe.SuccessThreshold,
				IncludeOnContainer:  params.Container.LivenessProbe.Enabled != nil && *params.Container.LivenessProbe.Enabled,
			},
			Readiness: api.ProbeData{
				Path:                params.Container.ReadinessProbe.Path,
				Port:                params.Container.ReadinessProbe.Port,
				InitialDelaySeconds: params.Container.ReadinessProbe.InitialDelaySeconds,
				TimeoutSeconds:      params.Container.ReadinessProbe.TimeoutSeconds,
				PeriodSeconds:       params.Container.ReadinessProbe.PeriodSeconds,
				FailureThreshold:    params.Container.ReadinessProbe.FailureThreshold,
				SuccessThreshold:    params.Container.ReadinessProbe.SuccessThreshold,
			},
			Metrics: api.MetricsData{
				Path: params.Container.Metrics.Path,
				Port: params.Container.Metrics.Port,
			},
		},

		// IsSimpleEnvvarValue returns true if a value should be wrapped in 'value: ""', otherwise the interface should be outputted as yaml
		IsSimpleEnvvarValue: s.IsSimpleEnvvarValue,
		ToYAML:              s.ToYAML,
		RenderToYAML:        s.RenderToYAML,
	}

	if data.Secrets == nil {
		data.Secrets = make(map[string]interface{}, 0)
	}

	// add SecretEnvironmentVariables to secrets map, but do base64 encode the values
	for key, value := range params.Container.SecretEnvironmentVariables {
		data.Secrets[key] = base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%v", value)))
	}
	// add sidecar SecretEnvironmentVariables to secrets map, but do base64 encode the values
	for _, sc := range params.Sidecars {
		for key, value := range sc.SecretEnvironmentVariables {
			data.Secrets[key] = base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%v", value)))
		}
	}

	if params.BackoffLimit != nil {
		data.BackoffLimit = *params.BackoffLimit
	}

	if params.DisableServiceAccountKeyRotation != nil {
		data.DisableServiceAccountKeyRotation = *params.DisableServiceAccountKeyRotation
	}

	if data.MountServiceAccountSecret {
		data.Container.EnvironmentVariables = s.AddEnvironmentVariableIfNotSet(data.Container.EnvironmentVariables, "GOOGLE_APPLICATION_CREDENTIALS", "/gcp-service-account/service-account-key.json")
		if data.GoogleCloudCredentialsAppName != "" {
			data.GoogleCloudCredentialsLabels["app"] = data.GoogleCloudCredentialsAppName
		}

		data.LegacyServiceAccountKey = params.LegacyGoogleCloudServiceAccountKeyFile
	}

	// ensure the app label exists and is identical to the app label selector
	if data.AppLabelSelector != "" {
		data.Labels["app"] = data.AppLabelSelector
	}

	// set tracing service name
	data.Container.EnvironmentVariables = s.AddEnvironmentVariableIfNotSet(data.Container.EnvironmentVariables, "JAEGER_SERVICE_NAME", params.App)

	if params.Action == api.ActionDeployCanary || params.Action == api.ActionDiffCanary {
		data.Container.EnvironmentVariables = s.AddEnvironmentVariableIfNotSet(data.Container.EnvironmentVariables, "JAEGER_SAMPLER_TYPE", "probabilistic")
		data.Container.EnvironmentVariables = s.AddEnvironmentVariableIfNotSet(data.Container.EnvironmentVariables, "JAEGER_SAMPLER_PARAM", "0.1")
		data.Container.EnvironmentVariables = s.AddEnvironmentVariableIfNotSet(data.Container.EnvironmentVariables, "JAEGER_TAGS", "track=canary")
	} else {
		data.Container.EnvironmentVariables = s.AddEnvironmentVariableIfNotSet(data.Container.EnvironmentVariables, "JAEGER_SAMPLER_TYPE", "remote")
		data.Container.EnvironmentVariables = s.AddEnvironmentVariableIfNotSet(data.Container.EnvironmentVariables, "JAEGER_SAMPLER_PARAM", "0.001")
	}

	data.HasOpenrestySidecar = false
	for _, sidecarParams := range params.Sidecars {
		sidecar := s.BuildSidecar(sidecarParams, params)
		data.Sidecars = append(data.Sidecars, sidecar)
		if sidecar.Type == string(api.SidecarTypeOpenresty) {
			data.HasOpenrestySidecar = true
		}
	}
	if params.CustomSidecars != nil {
		data.HasCustomSidecars = true
		data.CustomSidecars = params.CustomSidecars
	}

	data.UseESP = params.Visibility == api.VisibilityESP || params.Visibility == api.VisibilityESPv2
	data.HasEspConfigID = params.EspConfigID != ""
	data.EspConfigID = params.EspConfigID
	if (params.Visibility == api.VisibilityESP || params.Visibility == api.VisibilityESPv2) && len(params.Hosts) > 0 {
		data.EspService = params.Hosts[0]
	}

	if data.PreferPreemptibles {
		data.HasTolerations = true
		data.Tolerations = append(data.Tolerations, &map[string]interface{}{
			"key":      "cloud.google.com/gke-preemptible",
			"operator": "Equal",
			"value":    "true",
			"effect":   "NoSchedule",
		})
	}
	if data.UseWindowsNodes {
		data.HasTolerations = true
		data.Tolerations = append(data.Tolerations, &map[string]interface{}{
			"key":      "node.kubernetes.io/os",
			"operator": "Equal",
			"value":    "windows",
			"effect":   "NoSchedule",
		})
	}

	if params.Tolerations != nil {
		data.HasTolerations = true
		data.Tolerations = append(data.Tolerations, params.Tolerations...)
	}

	if params.InitContainers != nil {
		data.HasInitContainers = true
		data.InitContainers = params.InitContainers
	}

	data.Container.Readiness.IncludeOnContainer = params.Container.ReadinessProbe.Enabled != nil && *params.Container.ReadinessProbe.Enabled && (!data.HasOpenrestySidecar || params.Container.ReadinessProbe.Port != params.Container.Port || params.Container.ReadinessProbe.Path != params.Sidecar.HealthCheckPath)

	// if container port is set to 443, we always use https named port
	data.UseHTTPS = data.HasOpenrestySidecar || params.Container.Port == 443

	// set request params on the nginx ingress
	requestTimeout, requestTimeoutConvertError := strconv.Atoi(strings.Trim(params.Request.Timeout, "s"))
	data.EspRequestTimeout = requestTimeout

	data.NginxIngressProxyConnectTimeout = requestTimeout
	if requestTimeoutConvertError != nil {
		data.NginxIngressProxyConnectTimeout = 60
	}
	if data.NginxIngressProxyConnectTimeout > 75 {
		data.NginxIngressProxyConnectTimeout = 75
	}
	data.NginxIngressProxySendTimeout = requestTimeout
	if requestTimeoutConvertError != nil {
		data.NginxIngressProxySendTimeout = 60
	}
	data.NginxIngressProxyReadTimeout = requestTimeout
	if requestTimeoutConvertError != nil {
		data.NginxIngressProxyReadTimeout = 60
	}
	data.NginxIngressBackendProtocol = params.Request.IngressBackendProtocol
	data.NginxIngressProxyBodySize = params.Request.MaxBodySize
	data.NginxIngressClientBodyBufferSize = params.Request.ClientBodyBufferSize
	data.NginxIngressProxyBufferSize = params.Request.ProxyBufferSize
	data.NginxIngressProxyBuffersNumber = strconv.Itoa(params.Request.ProxyBuffersNumber)
	data.SetsNginxIngressLoadBalanceAlgorithm = params.Request.LoadBalanceAlgorithm != ""
	data.NginxIngressLoadBalanceAlgorithm = params.Request.LoadBalanceAlgorithm
	data.NginxAuthTLSSecret = params.Request.AuthSecret
	data.NginxAuthTLSVerifyDepth = params.Request.VerifyDepth

	// set request params for gce ingress
	data.BackendConfigTimeout = requestTimeout
	if requestTimeoutConvertError != nil {
		data.BackendConfigTimeout = 30
	}

	if params.ProbeService != nil {
		data.UsePrometheusProbe = *params.ProbeService
	}
	if params.TopologyAwareHints != nil {
		data.UseTopologyAwareHints = *params.TopologyAwareHints
	}
	if params.Container.Metrics.Scrape != nil {
		data.Container.Metrics.Scrape = *params.Container.Metrics.Scrape
	}
	if params.Container.Lifecycle.PrestopSleep != nil {
		data.Container.UseLifecyclePreStopSleepCommand = *params.Container.Lifecycle.PrestopSleep
	}
	if params.Container.Lifecycle.PrestopSleepSeconds != nil {
		data.Container.PreStopSleepSeconds = *params.Container.Lifecycle.PrestopSleepSeconds
	}

	if currentReplicas > 0 {
		data.Replicas = currentReplicas
	} else if (params.Autoscale.Enabled != nil && !*params.Autoscale.Enabled) || params.StrategyType == "Recreate" || params.Replicas > data.MinReplicas {
		data.Replicas = params.Replicas
	} else {
		data.Replicas = data.MinReplicas
	}

	if params.BuildVersion != "" {
		data.PodLabels["version"] = api.SanitizeLabel(params.BuildVersion)
	}
	if releaseID != "" {
		data.PodLabels["estafette.io/release-id"] = api.SanitizeLabel(releaseID)
	}
	if triggeredBy != "" {
		data.PodLabels["estafette.io/triggered-by"] = api.SanitizeLabel(triggeredBy)
	}
	if gitSource != "" && gitOwner != "" && gitName != "" {
		data.PodLabels["estafette.io/git-repository"] = api.SanitizeLabel(fmt.Sprintf("%v/%v/%v", gitSource, gitOwner, gitName))
	}
	if gitBranch != "" {
		data.PodLabels["estafette.io/git-branch"] = api.SanitizeLabel(gitBranch)
	}
	if gitRevision != "" {
		data.PodLabels["estafette.io/git-revision"] = api.SanitizeLabel(gitRevision)
	}
	if builderImageSHA != "" {
		// grab only first 20 char of hash since it is not necessary to go beyond that
		data.Labels["estafette.io/builder-image-sha"] = api.SanitizeLabel(builderImageSHA)[0:19]
	}
	if builderImageDate != "" {
		builderImageDateTrimmed, err := api.GetTrimmedDate(builderImageDate)
		if err == nil {
			data.Labels["estafette.io/builder-image-date"] = api.SanitizeLabel(builderImageDateTrimmed)
		}
	}

	switch params.Action {
	case api.ActionDeploySimple,
		api.ActionDiffSimple:
		data.IncludeTrackLabel = false
	case api.ActionDeployCanary,
		api.ActionDiffCanary:
		data.NameWithTrack += "-canary"
		data.IncludeTrackLabel = true
		data.TrackLabel = "canary"
	case api.ActionDeployStable,
		api.ActionDiffStable:
		data.NameWithTrack += "-stable"
		data.IncludeTrackLabel = true
		data.TrackLabel = "stable"
	}

	switch params.StrategyType {
	case api.StrategyTypeRollingUpdate:
		data.StrategyType = string(params.StrategyType)
	case api.StrategyTypeRecreate:
		data.StrategyType = string(params.StrategyType)
	case api.StrategyTypeAtomicUpdate:
		data.StrategyType = string(api.StrategyTypeRollingUpdate)
	}

	if params.StrategyType == api.StrategyTypeAtomicUpdate && params.Action == api.ActionDeploySimple && params.AtomicID != "" {
		data.NameWithTrack += "-" + params.AtomicID
		data.IncludeAtomicIDSelector = true
		data.AtomicID = params.AtomicID

		data.Labels["estafette.io/atomic-id"] = params.AtomicID
		data.PodLabels["estafette.io/atomic-id"] = params.AtomicID
	}

	// set some additional labels similar to helm charts in order to unify alerting and dashboards
	data.PodLabels["app.kubernetes.io/name"] = data.Name
	data.PodLabels["app.kubernetes.io/instance"] = data.NameWithTrack
	data.PodLabels["app.kubernetes.io/version"] = api.SanitizeLabel(params.BuildVersion)
	data.PodLabels["app.kubernetes.io/managed-by"] = "estafette"

	data.ConfigmapFiles = params.Configs.RenderedFileContent

	data.ManifestData = map[string]interface{}{}
	for k, v := range params.Manifests.Data {
		data.ManifestData[k] = v
	}

	switch params.Visibility {
	case api.VisibilityPrivate:
		data.Service = api.ServiceData{
			ServiceType: string(api.ServiceTypeClusterIP),
			Name:        params.App,
		}
		data.UseNginxIngress = true
		data.UseGCEIngress = false
		data.UseDNSAnnotationsOnIngress = true
		data.UseCloudflareProxy = true
		data.LimitTrustedIPRanges = false
		data.OverrideDefaultWhitelist = false

	case api.VisibilityIAP:
		data.Service = api.ServiceData{
			ServiceType:                         string(api.ServiceTypeNodePort),
			Name:                                params.App,
			UseBackendConfigAnnotationOnService: true,
			UseNegAnnotationOnService:           params.ContainerNativeLoadBalancing,
		}
		data.UseNginxIngress = false
		data.UseGCEIngress = true
		data.UseDNSAnnotationsOnIngress = true
		data.UseCloudflareProxy = false
		data.LimitTrustedIPRanges = false
		data.OverrideDefaultWhitelist = false
		data.IapOauthCredentialsClientID = params.IapOauthCredentialsClientID
		data.IapOauthCredentialsClientSecret = params.IapOauthCredentialsClientSecret

	case api.VisibilityPublicWhitelist:
		data.Service = api.ServiceData{
			ServiceType: string(api.ServiceTypeClusterIP),
			Name:        params.App,
		}
		data.UseNginxIngress = true
		data.UseGCEIngress = false
		data.UseDNSAnnotationsOnIngress = true
		data.UseCloudflareProxy = true
		data.LimitTrustedIPRanges = false
		data.OverrideDefaultWhitelist = len(params.WhitelistedIPS) > 0
		data.NginxIngressWhitelist = strings.Join(params.WhitelistedIPS, ",")

	case api.VisibilityApigee:
		data.Service = api.ServiceData{
			ServiceType: string(api.ServiceTypeClusterIP),
			Name:        params.App,
		}
		data.UseNginxIngress = true
		data.UseGCEIngress = false
		data.UseDNSAnnotationsOnIngress = true
		data.UseCloudflareProxy = true // For private ingress. For Apigee it is hard-coded to be false.
		data.LimitTrustedIPRanges = false
		data.OverrideDefaultWhitelist = false
		for _, h := range params.Hosts {
			hparts := strings.Split(h, ".")
			hparts[0] = hparts[0] + "-" + params.ApigeeSuffix
			data.ApigeeHosts = append(data.ApigeeHosts, strings.Join(hparts, "."))
		}
		data.ApigeeHostsJoined = strings.Join(data.ApigeeHosts, ",")

	case api.VisibilityESP, api.VisibilityESPv2:
		if params.EspServiceTypeClusterIP {
			data.Service = api.ServiceData{
				ServiceType: string(api.ServiceTypeClusterIP),
				Name:        params.App + "-cluster-ip",
			}
			data.UseNginxIngress = true
			data.UseDNSAnnotationsOnIngress = true
			data.UseCloudflareProxy = true
			data.LimitTrustedIPRanges = false
			data.OverrideDefaultWhitelist = false
		} else {
			data.Service = api.ServiceData{
				ServiceType:                string(api.ServiceTypeLoadBalancer),
				Name:                       params.App,
				UseDNSAnnotationsOnService: true,
			}
			data.UseNginxIngress = false
			data.UseGCEIngress = false
			data.UseDNSAnnotationsOnIngress = false
			data.UseCloudflareProxy = true
			data.LimitTrustedIPRanges = true
			data.OverrideDefaultWhitelist = false
		}

	case api.VisibilityPublic:
		data.Service = api.ServiceData{
			ServiceType:                string(api.ServiceTypeLoadBalancer),
			Name:                       params.App,
			UseDNSAnnotationsOnService: true,
		}
		data.UseNginxIngress = false
		data.UseGCEIngress = false
		data.UseDNSAnnotationsOnIngress = false
		data.UseCloudflareProxy = true
		data.LimitTrustedIPRanges = true
		data.OverrideDefaultWhitelist = false
	}

	if params.WorkloadIdentity != nil {
		data.UseWorkloadIdentity = *params.WorkloadIdentity
	}

	if params.DNS.UseCloudflareEstafetteExtension != nil {
		data.UseCloudflareEstafetteExtension = *params.DNS.UseCloudflareEstafetteExtension
	}

	if params.DNS.UseExternalDNS != nil {
		data.UseExternalDNS = *params.DNS.UseExternalDNS
	}
	// add extra hosts for routing in ingress, without setting their dns records
	data.Hosts = append(data.Hosts, params.HostsRouteOnly...)
	data.InternalHosts = append(data.InternalHosts, params.InternalHostsRouteOnly...)

	if !strings.HasSuffix(data.IngressPath, "/") && !strings.HasSuffix(data.IngressPath, "*") {
		data.IngressPath += "/"
	}
	if data.UseGCEIngress && !strings.HasSuffix(data.IngressPath, "*") {
		data.IngressPath += "*"
	}

	if data.UseGCEIngress {
		data.PathType = "ImplementationSpecific"
	}
	if !strings.HasSuffix(data.InternalIngressPath, "/") && !strings.HasSuffix(data.InternalIngressPath, "*") {
		data.InternalIngressPath += "/"
	}

	data.TrustedIPRanges = params.TrustedIPRanges

	data.AdditionalVolumeMounts = []api.VolumeMountData{}
	for _, vm := range params.VolumeMounts {
		yamlBytes, err := yaml.Marshal(vm.Volume)
		if err == nil {
			data.AdditionalVolumeMounts = append(data.AdditionalVolumeMounts, api.VolumeMountData{
				Name:       vm.Name,
				MountPath:  vm.MountPath,
				VolumeYAML: string(yamlBytes),
			})
		}
	}
	data.MountAdditionalVolumes = len(data.AdditionalVolumeMounts) > 0

	data.AdditionalContainerPorts = []api.AdditionalPortData{}
	data.AdditionalServicePorts = []api.AdditionalPortData{}
	for _, ap := range params.Container.AdditionalPorts {
		additionalPortData := api.AdditionalPortData{
			Name:     ap.Name,
			Port:     ap.Port,
			Protocol: ap.Protocol,
		}
		data.AdditionalContainerPorts = append(data.AdditionalContainerPorts, additionalPortData)

		includeAsServicePort := ap.Visibility == params.Visibility

		if includeAsServicePort {
			data.AdditionalServicePorts = append(data.AdditionalServicePorts, additionalPortData)
		}
	}

	// Use certificate secret if it's specified
	if params.CertificateSecret != "" {
		data.UseCertificateSecret = true
		data.CertificateSecretName = params.CertificateSecret
	}

	data.MountVolumes = data.MountSslCertificate || data.MountApplicationSecrets || data.MountConfigmap || data.MountPayloadLogging || data.MountServiceAccountSecret || data.MountAdditionalVolumes

	if params.ImagePullSecretUser != "" && params.ImagePullSecretPassword != "" {
		data.HasImagePullSecret = true
		data.DockerConfig = map[string]map[string]map[string]string{
			"auths": {
				"https://index.docker.io/v1/": map[string]string{
					"username": params.ImagePullSecretUser,
					"password": params.ImagePullSecretPassword,
					"auth":     base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%v:%v", params.ImagePullSecretUser, params.ImagePullSecretPassword))),
				},
			},
		}
	}

	return data
}

func (s *service) BuildSidecar(sidecar *api.SidecarParams, params api.Params) api.SidecarData {
	builtSidecar := api.SidecarData{
		Type:                       string(sidecar.Type),
		Image:                      sidecar.Image,
		CPURequest:                 sidecar.CPU.Request,
		CPULimit:                   sidecar.CPU.Limit,
		MemoryRequest:              sidecar.Memory.Request,
		MemoryLimit:                sidecar.Memory.Limit,
		EnvironmentVariables:       sidecar.EnvironmentVariables,
		SecretEnvironmentVariables: sidecar.SecretEnvironmentVariables,
		HasEnvironmentVariables:    len(sidecar.EnvironmentVariables) > 0 || len(sidecar.SecretEnvironmentVariables) > 0,
		SidecarSpecificProperties: map[string]interface{}{
			"healthcheckpath":                   sidecar.HealthCheckPath,
			"dbinstanceconnectionname":          sidecar.DbInstanceConnectionName,
			"sqlproxyport":                      sidecar.SQLProxyPort,
			"sqlproxyterminationtimeoutseconds": sidecar.SQLProxyTerminationTimeoutSeconds,
		},
	}

	if sidecar.Type == api.SidecarTypeOpenresty {
		builtSidecar.EnvironmentVariables = s.AddEnvironmentVariableIfNotSet(builtSidecar.EnvironmentVariables, "KEEPALIVE_TIMEOUT", params.Request.KeepaliveTimeout)
		builtSidecar.EnvironmentVariables = s.AddEnvironmentVariableIfNotSet(builtSidecar.EnvironmentVariables, "SEND_TIMEOUT", params.Request.Timeout)
		builtSidecar.EnvironmentVariables = s.AddEnvironmentVariableIfNotSet(builtSidecar.EnvironmentVariables, "CLIENT_BODY_TIMEOUT", params.Request.Timeout)
		builtSidecar.EnvironmentVariables = s.AddEnvironmentVariableIfNotSet(builtSidecar.EnvironmentVariables, "CLIENT_HEADER_TIMEOUT", params.Request.Timeout)
		builtSidecar.EnvironmentVariables = s.AddEnvironmentVariableIfNotSet(builtSidecar.EnvironmentVariables, "PROXY_CONNECT_TIMEOUT", params.Request.Timeout)
		builtSidecar.EnvironmentVariables = s.AddEnvironmentVariableIfNotSet(builtSidecar.EnvironmentVariables, "PROXY_SEND_TIMEOUT", params.Request.Timeout)
		builtSidecar.EnvironmentVariables = s.AddEnvironmentVariableIfNotSet(builtSidecar.EnvironmentVariables, "PROXY_READ_TIMEOUT", params.Request.Timeout)
		builtSidecar.EnvironmentVariables = s.AddEnvironmentVariableIfNotSet(builtSidecar.EnvironmentVariables, "CLIENT_MAX_BODY_SIZE", params.Request.MaxBodySize)
		builtSidecar.EnvironmentVariables = s.AddEnvironmentVariableIfNotSet(builtSidecar.EnvironmentVariables, "CLIENT_BODY_BUFFER_SIZE", params.Request.ClientBodyBufferSize)
		builtSidecar.EnvironmentVariables = s.AddEnvironmentVariableIfNotSet(builtSidecar.EnvironmentVariables, "PROXY_BUFFER_SIZE", params.Request.ProxyBufferSize)
		builtSidecar.EnvironmentVariables = s.AddEnvironmentVariableIfNotSet(builtSidecar.EnvironmentVariables, "PROXY_BUFFERS_SIZE", params.Request.ProxyBufferSize)
		builtSidecar.EnvironmentVariables = s.AddEnvironmentVariableIfNotSet(builtSidecar.EnvironmentVariables, "PROXY_BUFFERS_NUMBER", strconv.Itoa(params.Request.ProxyBuffersNumber))

		if params.Container.Lifecycle.PrestopSleep != nil && *params.Container.Lifecycle.PrestopSleep && params.Container.Lifecycle.PrestopSleepSeconds != nil {
			builtSidecar.EnvironmentVariables = s.AddEnvironmentVariableIfNotSet(builtSidecar.EnvironmentVariables, "GRACEFUL_SHUTDOWN_DELAY_SECONDS", strconv.Itoa(*params.Container.Lifecycle.PrestopSleepSeconds))
		}

		if params.Visibility == api.VisibilityESP || params.Visibility == api.VisibilityESPv2 {
			builtSidecar.EnvironmentVariables = s.AddEnvironmentVariableIfNotSet(builtSidecar.EnvironmentVariables, "ENFORCE_HTTPS", "false")
		}
	}

	if sidecar.CustomProperties != nil {
		yamlBytes, err := yaml.Marshal(sidecar.CustomProperties)
		if err == nil {
			builtSidecar.CustomPropertiesYAML = string(yamlBytes)
			builtSidecar.HasCustomProperties = true
		}
	}

	return builtSidecar
}

func (s *service) AddEnvironmentVariableIfNotSet(environmentVariables map[string]interface{}, name, value string) map[string]interface{} {

	if environmentVariables == nil {
		environmentVariables = map[string]interface{}{}
	}
	if _, ok := environmentVariables[name]; !ok {
		environmentVariables[name] = value
	}

	return environmentVariables
}

func (s *service) IsSimpleEnvvarValue(i interface{}) bool {
	switch i.(type) {
	case int:
		return true
	case float64:
		return true
	case string:
		return true
	case bool:
		return true
	}

	return false
}

func (s *service) ToYAML(v interface{}) string {
	data, err := yaml.Marshal(v)
	if err != nil {
		// Swallow errors inside of a template.
		return ""
	}
	return string(data)
}

func (s *service) RenderToYAML(v interface{}, data interface{}) string {

	value := s.ToYAML(v)

	tmpl, err := template.New("renderToYAML").Parse(value)
	if err != nil {
		return value
	}

	var renderedTemplate bytes.Buffer
	err = tmpl.Execute(&renderedTemplate, data)
	if err != nil {
		return value
	}

	return renderedTemplate.String()
}
