package main

import (
	"regexp"
	"strconv"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

func generateTemplateData(params Params, currentReplicas int, releaseID, triggeredBy string) TemplateData {

	data := TemplateData{
		BuildVersion: params.BuildVersion,

		Name:              params.App,
		NameWithTrack:     params.App,
		Namespace:         params.Namespace,
		Schedule:          params.Schedule,
		ConcurrencyPolicy: params.ConcurrencyPolicy,
		Labels:            sanitizeLabels(params.Labels),
		AppLabelSelector:  sanitizeLabel(params.App),

		Hosts:               params.Hosts,
		HostsJoined:         strings.Join(params.Hosts, ","),
		InternalHosts:       params.InternalHosts,
		InternalHostsJoined: strings.Join(params.InternalHosts, ","),
		AllHosts:            append(params.Hosts, params.InternalHosts...),
		AllHostsJoined:      strings.Join(append(params.Hosts, params.InternalHosts...), ","),
		IngressPath:         params.Basepath,
		InternalIngressPath: params.Basepath,

		IncludeReplicas: currentReplicas > 0,

		MinReplicas:         params.Autoscale.MinReplicas,
		MaxReplicas:         params.Autoscale.MaxReplicas,
		TargetCPUPercentage: params.Autoscale.CPUPercentage,

		UseHpaScaler:                params.Autoscale.Safety.Enabled,
		HpaScalerPromQuery:          params.Autoscale.Safety.PromQuery,
		HpaScalerRequestsPerReplica: params.Autoscale.Safety.Ratio,
		HpaScalerDelta:              params.Autoscale.Safety.Delta,
		HpaScalerScaleDownMaxRatio:  params.Autoscale.Safety.ScaleDownRatio,

		Secrets:                 params.Secrets.Keys,
		MountApplicationSecrets: len(params.Secrets.Keys) > 0,
		SecretMountPath:         params.Secrets.MountPath,
		MountConfigmap:          len(params.Configs.Files) > 0 || len(params.Configs.InlineFiles) > 0,
		ConfigMountPath:         params.Configs.MountPath,

		MountPayloadLogging:      params.EnablePayloadLogging,
		AddSafeToEvictAnnotation: params.EnablePayloadLogging,

		RollingUpdateMaxSurge:       params.RollingUpdate.MaxSurge,
		RollingUpdateMaxUnavailable: params.RollingUpdate.MaxUnavailable,

		PreferPreemptibles:               params.ChaosProof,
		MountServiceAccountSecret:        params.UseGoogleCloudCredentials,
		GoogleCloudCredentialsAppName:    params.GoogleCloudCredentialsApp,
		DisableServiceAccountKeyRotation: params.DisableServiceAccountKeyRotation,

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
	}

	if params.UseGoogleCloudCredentials {
		data.Container.EnvironmentVariables = addEnvironmentVariableIfNotSet(data.Container.EnvironmentVariables, "GOOGLE_APPLICATION_CREDENTIALS", "/gcp-service-account/service-account-key.json")
	}

	// ensure the app label exists and is identical to the app label selector
	if data.AppLabelSelector != "" {
		data.Labels["app"] = data.AppLabelSelector
	}

	// set tracing service name
	data.Container.EnvironmentVariables = addEnvironmentVariableIfNotSet(data.Container.EnvironmentVariables, "JAEGER_SERVICE_NAME", params.App)

	for _, sidecarParams := range params.Sidecars {
		sidecar := buildSidecar(sidecarParams, params.Request)
		data.Sidecars = append(data.Sidecars, sidecar)

		logInfo("Added sidecar of type %v to data.Sidecars of length %v", sidecarParams.Type, len(data.Sidecars))
	}

	// set request params on the nginx ingress
	requestTimeout, requestTimeoutConvertError := strconv.Atoi(strings.Trim(params.Request.Timeout, "s"))

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
	data.NginxIngressProxyBodySize = params.Request.MaxBodySize
	data.NginxIngressClientBodyBufferSize = params.Request.ClientBodyBufferSize
	data.NginxIngressProxyBufferSize = params.Request.ProxyBufferSize
	data.NginxIngressProxyBuffersNumber = strconv.Itoa(params.Request.ProxyBuffersNumber)

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
	} else {
		data.Replicas = data.MinReplicas
	}

	if releaseID != "" {
		data.IncludeReleaseIDLabel = true
		data.ReleaseIDLabel = sanitizeLabel(releaseID)
	}

	if triggeredBy != "" {
		data.IncludeTriggeredByLabel = true
		data.TriggeredByLabel = sanitizeLabel(triggeredBy)
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

	switch params.Visibility {
	case "private":
		data.ServiceType = "ClusterIP"
		data.UseNginxIngress = true
		data.UseGCEIngress = false
		data.UseDNSAnnotationsOnIngress = true
		data.UseDNSAnnotationsOnService = false
		data.UseBackendConfigAnnotationOnService = false
		data.LimitTrustedIPRanges = false
		data.OverrideDefaultWhitelist = false

	case "iap":
		data.ServiceType = "NodePort"
		data.UseNginxIngress = false
		data.UseGCEIngress = true
		data.UseDNSAnnotationsOnIngress = true
		data.UseDNSAnnotationsOnService = false
		data.UseBackendConfigAnnotationOnService = true
		data.LimitTrustedIPRanges = false
		data.OverrideDefaultWhitelist = false
		data.IapOauthCredentialsClientID = params.IapOauthCredentialsClientID
		data.IapOauthCredentialsClientSecret = params.IapOauthCredentialsClientSecret

	case "public-whitelist":
		data.ServiceType = "ClusterIP"
		data.UseNginxIngress = true
		data.UseGCEIngress = false
		data.UseDNSAnnotationsOnIngress = true
		data.UseDNSAnnotationsOnService = false
		data.UseBackendConfigAnnotationOnService = false
		data.LimitTrustedIPRanges = false
		data.OverrideDefaultWhitelist = len(params.WhitelistedIPS) > 0
		data.NginxIngressWhitelist = strings.Join(params.WhitelistedIPS, ",")

	case "public":
		data.ServiceType = "LoadBalancer"
		data.UseNginxIngress = false
		data.UseGCEIngress = false
		data.UseDNSAnnotationsOnIngress = false
		data.UseDNSAnnotationsOnService = true
		data.LimitTrustedIPRanges = true
		data.OverrideDefaultWhitelist = false
	}

	if !strings.HasSuffix(data.IngressPath, "/") && !strings.HasSuffix(data.IngressPath, "*") {
		data.IngressPath += "/"
	}
	if data.UseGCEIngress && !strings.HasSuffix(data.IngressPath, "*") {
		data.IngressPath += "*"
	}
	if !strings.HasSuffix(data.InternalIngressPath, "/") && !strings.HasSuffix(data.InternalIngressPath, "*") {
		data.InternalIngressPath += "/"
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
	data.MountAdditionalVolumes = len(data.AdditionalVolumeMounts) > 0

	data.AdditionalContainerPorts = []AdditionalPortData{}
	data.AdditionalServicePorts = []AdditionalPortData{}
	for _, ap := range params.Container.AdditionalPorts {
		additionalPortData := AdditionalPortData{
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

	return data
}

func buildSidecar(sidecar *SidecarParams, request RequestParams) SidecarData {
	builtSidecar := SidecarData{
		Type:                    sidecar.Type,
		Image:                   sidecar.Image,
		CPURequest:              sidecar.CPU.Request,
		CPULimit:                sidecar.CPU.Limit,
		MemoryRequest:           sidecar.Memory.Request,
		MemoryLimit:             sidecar.Memory.Limit,
		EnvironmentVariables:    sidecar.EnvironmentVariables,
		HasEnvironmentVariables: len(sidecar.EnvironmentVariables) > 0,
		SidecarSpecificProperties: map[string]interface{}{
			"healthcheckpath":                   sidecar.HealthCheckPath,
			"dbinstanceconnectionname":          sidecar.DbInstanceConnectionName,
			"sqlproxyport":                      sidecar.SQLProxyPort,
			"sqlproxyterminationtimeoutseconds": sidecar.SQLProxyTerminationTimeoutSeconds,
		},
	}

	if builtSidecar.Type == "openresty" {
		builtSidecar.EnvironmentVariables = addEnvironmentVariableIfNotSet(builtSidecar.EnvironmentVariables, "SEND_TIMEOUT", request.Timeout)
		builtSidecar.EnvironmentVariables = addEnvironmentVariableIfNotSet(builtSidecar.EnvironmentVariables, "CLIENT_BODY_TIMEOUT", request.Timeout)
		builtSidecar.EnvironmentVariables = addEnvironmentVariableIfNotSet(builtSidecar.EnvironmentVariables, "CLIENT_HEADER_TIMEOUT", request.Timeout)
		builtSidecar.EnvironmentVariables = addEnvironmentVariableIfNotSet(builtSidecar.EnvironmentVariables, "PROXY_CONNECT_TIMEOUT", request.Timeout)
		builtSidecar.EnvironmentVariables = addEnvironmentVariableIfNotSet(builtSidecar.EnvironmentVariables, "PROXY_SEND_TIMEOUT", request.Timeout)
		builtSidecar.EnvironmentVariables = addEnvironmentVariableIfNotSet(builtSidecar.EnvironmentVariables, "PROXY_READ_TIMEOUT", request.Timeout)
		builtSidecar.EnvironmentVariables = addEnvironmentVariableIfNotSet(builtSidecar.EnvironmentVariables, "CLIENT_MAX_BODY_SIZE", request.MaxBodySize)
		builtSidecar.EnvironmentVariables = addEnvironmentVariableIfNotSet(builtSidecar.EnvironmentVariables, "CLIENT_BODY_BUFFER_SIZE", request.ClientBodyBufferSize)
		builtSidecar.EnvironmentVariables = addEnvironmentVariableIfNotSet(builtSidecar.EnvironmentVariables, "PROXY_BUFFER_SIZE", request.ProxyBufferSize)
		builtSidecar.EnvironmentVariables = addEnvironmentVariableIfNotSet(builtSidecar.EnvironmentVariables, "PROXY_BUFFERS_SIZE", request.ProxyBufferSize)
		builtSidecar.EnvironmentVariables = addEnvironmentVariableIfNotSet(builtSidecar.EnvironmentVariables, "PROXY_BUFFERS_NUMBER", strconv.Itoa(request.ProxyBuffersNumber))
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

func addEnvironmentVariableIfNotSet(environmentVariables map[string]interface{}, name, value string) map[string]interface{} {

	if environmentVariables == nil {
		environmentVariables = map[string]interface{}{}
	}
	if _, ok := environmentVariables[name]; !ok {
		environmentVariables[name] = value
	}

	return environmentVariables
}

// a valid label must be an empty string or consist of alphanumeric characters, '-', '_' or '.', and must start and end with an alphanumeric character (e.g. 'MyValue',  or 'my_value',  or '12345', regex used for validation is '(([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9])?')
func sanitizeLabel(value string) string {

	// Valid label values must be 63 characters or less and must be empty or begin and end with an alphanumeric character ([a-z0-9A-Z])
	// with dashes (-), underscores (_), dots (.), and alphanumerics between.

	// replace all invalid characters with a hyphen
	reg := regexp.MustCompile(`[^a-zA-Z0-9-_.]+`)
	value = reg.ReplaceAllString(value, "-")

	// replace double hyphens with a single one
	value = strings.Replace(value, "--", "-", -1)

	// ensure it starts with an alphanumeric character
	reg = regexp.MustCompile(`^[-_.]+`)
	value = reg.ReplaceAllString(value, "")

	// maximize length at 63 characters
	if len(value) > 63 {
		value = value[:63]
	}

	// ensure it ends with an alphanumeric character
	reg = regexp.MustCompile(`[-_.]+$`)
	value = reg.ReplaceAllString(value, "")

	return value
}

func sanitizeLabels(labels map[string]string) (sanitizedLabels map[string]string) {
	sanitizedLabels = make(map[string]string, len(labels))
	for k, v := range labels {
		sanitizedLabels[k] = sanitizeLabel(v)
	}
	return
}
