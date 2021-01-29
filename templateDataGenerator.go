package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"text/template"

	yaml "gopkg.in/yaml.v2"
)

func generateTemplateData(params Params, currentReplicas int, gitSource, gitOwner, gitName, gitBranch, gitRevision, releaseID, triggeredBy string) TemplateData {

	data := TemplateData{
		Name:              params.App,
		NameWithTrack:     params.App,
		Namespace:         params.Namespace,
		Schedule:          params.Schedule,
		ConcurrencyPolicy: params.ConcurrencyPolicy,
		RestartPolicy:     params.RestartPolicy,
		Completions:       params.Completions,
		Parallelism:       params.Parallelism,
		Labels:            sanitizeLabels(params.Labels),
		PodLabels:         sanitizeLabels(params.Labels),
		AppLabelSelector:  sanitizeLabel(params.App),

		Hosts:               params.Hosts,
		HostsJoined:         strings.Join(params.Hosts, ","),
		InternalHosts:       params.InternalHosts,
		InternalHostsJoined: strings.Join(params.InternalHosts, ","),
		AllHosts:            append(params.Hosts, params.InternalHosts...),
		AllHostsJoined:      strings.Join(append(params.Hosts, params.InternalHosts...), ","),
		IngressPath:         params.Basepath,
		InternalIngressPath: params.Basepath,

		IncludeReplicas: currentReplicas > 0 || ((params.Autoscale.Enabled == nil || !*params.Autoscale.Enabled || params.StrategyType == "Recreate") && params.Replicas > 0),

		MinReplicas:         params.Autoscale.MinReplicas,
		MaxReplicas:         params.Autoscale.MaxReplicas,
		TargetCPUPercentage: params.Autoscale.CPUPercentage,

		UseHpaScaler:                params.Autoscale.Safety.Enabled,
		HpaScalerPromQuery:          params.Autoscale.Safety.PromQuery,
		HpaScalerRequestsPerReplica: params.Autoscale.Safety.Ratio,
		HpaScalerDelta:              params.Autoscale.Safety.Delta,
		HpaScalerScaleDownMaxRatio:  params.Autoscale.Safety.ScaleDownRatio,

		Secrets:                 params.Secrets.Keys,
		MountSslCertificate:     params.Kind == KindDeployment,
		MountApplicationSecrets: len(params.Secrets.Keys) > 0,
		SecretMountPath:         params.Secrets.MountPath,
		MountConfigmap:          len(params.Configs.Files) > 0 || len(params.Configs.InlineFiles) > 0,
		ConfigMountPath:         params.Configs.MountPath,
		Tolerations:             []*map[string]interface{}{},

		MountPayloadLogging:      params.EnablePayloadLogging,
		AddSafeToEvictAnnotation: params.EnablePayloadLogging,

		StrategyType:                params.StrategyType,
		RollingUpdateMaxSurge:       params.RollingUpdate.MaxSurge,
		RollingUpdateMaxUnavailable: params.RollingUpdate.MaxUnavailable,

		PreferPreemptibles:            params.ChaosProof,
		UseWindowsNodes:               params.OperatingSystem == OperatingSystemWindows,
		MountServiceAccountSecret:     params.UseGoogleCloudCredentials || params.LegacyGoogleCloudServiceAccountKeyFile != "",
		UseLegacyServiceAccountKey:    params.LegacyGoogleCloudServiceAccountKeyFile != "",
		GoogleCloudCredentialsAppName: params.GoogleCloudCredentialsApp,
		GoogleCloudCredentialsLabels:  sanitizeLabels(params.Labels),

		PodManagementPolicy: params.PodManagementPolicy,
		StorageClass:        params.StorageClass,
		StorageSize:         params.StorageSize,
		StorageMountPath:    params.StorageMountPath,

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
				PeriodSeconds:       params.Container.LivenessProbe.PeriodSeconds,
				FailureThreshold:    params.Container.LivenessProbe.FailureThreshold,
				SuccessThreshold:    params.Container.LivenessProbe.SuccessThreshold,
				IncludeOnContainer:  params.Container.LivenessProbe.Enabled != nil && *params.Container.LivenessProbe.Enabled,
			},
			Readiness: ProbeData{
				Path:                params.Container.ReadinessProbe.Path,
				Port:                params.Container.ReadinessProbe.Port,
				InitialDelaySeconds: params.Container.ReadinessProbe.InitialDelaySeconds,
				TimeoutSeconds:      params.Container.ReadinessProbe.TimeoutSeconds,
				PeriodSeconds:       params.Container.ReadinessProbe.PeriodSeconds,
				FailureThreshold:    params.Container.ReadinessProbe.FailureThreshold,
				SuccessThreshold:    params.Container.ReadinessProbe.SuccessThreshold,
			},
			Metrics: MetricsData{
				Path: params.Container.Metrics.Path,
				Port: params.Container.Metrics.Port,
			},
		},

		// IsSimpleEnvvarValue returns true if a value should be wrapped in 'value: ""', otherwise the interface should be outputted as yaml
		IsSimpleEnvvarValue: isSimpleEnvvarValue,
		ToYAML:              toYAML,
		RenderToYAML:        renderToYAML,
	}

	if params.BackoffLimit != nil {
		data.BackoffLimit = *params.BackoffLimit
	}

	if params.DisableServiceAccountKeyRotation != nil {
		data.DisableServiceAccountKeyRotation = *params.DisableServiceAccountKeyRotation
	}

	if data.MountServiceAccountSecret {
		data.Container.EnvironmentVariables = addEnvironmentVariableIfNotSet(data.Container.EnvironmentVariables, "GOOGLE_APPLICATION_CREDENTIALS", "/gcp-service-account/service-account-key.json")
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
	data.Container.EnvironmentVariables = addEnvironmentVariableIfNotSet(data.Container.EnvironmentVariables, "JAEGER_SERVICE_NAME", params.App)

	if params.Action == ActionDeployCanary || params.Action == ActionDiffCanary {
		data.Container.EnvironmentVariables = addEnvironmentVariableIfNotSet(data.Container.EnvironmentVariables, "JAEGER_SAMPLER_TYPE", "probabilistic")
		data.Container.EnvironmentVariables = addEnvironmentVariableIfNotSet(data.Container.EnvironmentVariables, "JAEGER_SAMPLER_PARAM", "0.1")
		data.Container.EnvironmentVariables = addEnvironmentVariableIfNotSet(data.Container.EnvironmentVariables, "JAEGER_TAGS", "track=canary")
	} else {
		data.Container.EnvironmentVariables = addEnvironmentVariableIfNotSet(data.Container.EnvironmentVariables, "JAEGER_SAMPLER_TYPE", "remote")
		data.Container.EnvironmentVariables = addEnvironmentVariableIfNotSet(data.Container.EnvironmentVariables, "JAEGER_SAMPLER_PARAM", "0.001")
	}

	data.HasOpenrestySidecar = false
	for _, sidecarParams := range params.Sidecars {
		sidecar := buildSidecar(sidecarParams, params)
		data.Sidecars = append(data.Sidecars, sidecar)
		if sidecar.Type == string(SidecarTypeOpenresty) {
			data.HasOpenrestySidecar = true
		}
	}
	if params.CustomSidecars != nil {
		data.HasCustomSidecars = true
		data.CustomSidecars = params.CustomSidecars
	}

	data.UseESP = params.Visibility == VisibilityESP
	data.HasEspConfigID = params.EspConfigID != ""
	data.EspConfigID = params.EspConfigID
	if params.Visibility == VisibilityESP && len(params.Hosts) > 0 {
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
		data.PodLabels["version"] = sanitizeLabel(params.BuildVersion)
	}
	if releaseID != "" {
		data.PodLabels["estafette.io/release-id"] = sanitizeLabel(releaseID)
	}
	if triggeredBy != "" {
		data.PodLabels["estafette.io/triggered-by"] = sanitizeLabel(triggeredBy)
	}
	if gitSource != "" && gitOwner != "" && gitName != "" {
		data.PodLabels["estafette.io/git-repository"] = sanitizeLabel(fmt.Sprintf("%v/%v/%v", gitSource, gitOwner, gitName))
	}
	if gitBranch != "" {
		data.PodLabels["estafette.io/git-branch"] = sanitizeLabel(gitBranch)
	}
	if gitRevision != "" {
		data.PodLabels["estafette.io/git-revision"] = sanitizeLabel(gitRevision)
	}

	switch params.Action {
	case ActionDeploySimple,
		ActionDiffSimple:
		data.IncludeTrackLabel = false
	case ActionDeployCanary,
		ActionDiffCanary:
		data.NameWithTrack += "-canary"
		data.IncludeTrackLabel = true
		data.TrackLabel = "canary"
	case ActionDeployStable,
		ActionDiffStable:
		data.NameWithTrack += "-stable"
		data.IncludeTrackLabel = true
		data.TrackLabel = "stable"
	}

	// set some additional labels similar to helm charts in order to unify alerting and dashboards
	data.PodLabels["app.kubernetes.io/name"] = data.Name
	data.PodLabels["app.kubernetes.io/instance"] = data.NameWithTrack
	data.PodLabels["app.kubernetes.io/version"] = sanitizeLabel(params.BuildVersion)
	data.PodLabels["app.kubernetes.io/managed-by"] = "estafette"

	data.ConfigmapFiles = params.Configs.RenderedFileContent

	data.ManifestData = map[string]interface{}{}
	for k, v := range params.Manifests.Data {
		data.ManifestData[k] = v
	}

	switch params.Visibility {
	case VisibilityPrivate:
		data.ServiceType = "ClusterIP"
		data.UseNginxIngress = true
		data.UseGCEIngress = false
		data.UseDNSAnnotationsOnIngress = true
		data.UseDNSAnnotationsOnService = false
		data.UseCloudflareProxy = true
		data.UseBackendConfigAnnotationOnService = false
		data.UseNegAnnotationOnService = false
		data.LimitTrustedIPRanges = false
		data.OverrideDefaultWhitelist = false

	case VisibilityIAP:
		data.ServiceType = "NodePort"
		data.UseNginxIngress = false
		data.UseGCEIngress = true
		data.UseDNSAnnotationsOnIngress = true
		data.UseDNSAnnotationsOnService = false
		data.UseCloudflareProxy = false
		data.UseBackendConfigAnnotationOnService = true
		data.UseNegAnnotationOnService = params.ContainerNativeLoadBalancing
		data.LimitTrustedIPRanges = false
		data.OverrideDefaultWhitelist = false
		data.IapOauthCredentialsClientID = params.IapOauthCredentialsClientID
		data.IapOauthCredentialsClientSecret = params.IapOauthCredentialsClientSecret

	case VisibilityPublicWhitelist:
		data.ServiceType = "ClusterIP"
		data.UseNginxIngress = true
		data.UseGCEIngress = false
		data.UseDNSAnnotationsOnIngress = true
		data.UseDNSAnnotationsOnService = false
		data.UseCloudflareProxy = true
		data.UseBackendConfigAnnotationOnService = false
		data.UseNegAnnotationOnService = false
		data.LimitTrustedIPRanges = false
		data.OverrideDefaultWhitelist = len(params.WhitelistedIPS) > 0
		data.NginxIngressWhitelist = strings.Join(params.WhitelistedIPS, ",")

	case VisibilityApigee:
		data.ServiceType = "ClusterIP"
		data.UseNginxIngress = true
		data.UseGCEIngress = false
		data.UseDNSAnnotationsOnIngress = true
		data.UseDNSAnnotationsOnService = false
		data.UseCloudflareProxy = true // For private ingress. For Apigee it is hard-coded to be false.
		data.UseBackendConfigAnnotationOnService = false
		data.UseNegAnnotationOnService = false
		data.LimitTrustedIPRanges = false
		data.OverrideDefaultWhitelist = false
		for _, h := range params.Hosts {
			hparts := strings.Split(h, ".")
			hparts[0] = hparts[0] + "-" + params.ApigeeSuffix
			data.ApigeeHosts = append(data.ApigeeHosts, strings.Join(hparts, "."))
		}
		data.ApigeeHostsJoined = strings.Join(data.ApigeeHosts, ",")

	case VisibilityESP:
		data.ServiceType = "LoadBalancer"
		data.UseNginxIngress = false
		data.UseGCEIngress = false
		data.UseDNSAnnotationsOnIngress = false
		data.UseDNSAnnotationsOnService = true
		data.UseCloudflareProxy = true
		data.LimitTrustedIPRanges = true
		data.OverrideDefaultWhitelist = false

	case VisibilityPublic:
		data.ServiceType = "LoadBalancer"
		data.UseNginxIngress = false
		data.UseGCEIngress = false
		data.UseDNSAnnotationsOnIngress = false
		data.UseDNSAnnotationsOnService = true
		data.UseCloudflareProxy = true
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

func buildSidecar(sidecar *SidecarParams, params Params) SidecarData {
	builtSidecar := SidecarData{
		Type:                    string(sidecar.Type),
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

	if sidecar.Type == SidecarTypeOpenresty {
		builtSidecar.EnvironmentVariables = addEnvironmentVariableIfNotSet(builtSidecar.EnvironmentVariables, "SEND_TIMEOUT", params.Request.Timeout)
		builtSidecar.EnvironmentVariables = addEnvironmentVariableIfNotSet(builtSidecar.EnvironmentVariables, "CLIENT_BODY_TIMEOUT", params.Request.Timeout)
		builtSidecar.EnvironmentVariables = addEnvironmentVariableIfNotSet(builtSidecar.EnvironmentVariables, "CLIENT_HEADER_TIMEOUT", params.Request.Timeout)
		builtSidecar.EnvironmentVariables = addEnvironmentVariableIfNotSet(builtSidecar.EnvironmentVariables, "PROXY_CONNECT_TIMEOUT", params.Request.Timeout)
		builtSidecar.EnvironmentVariables = addEnvironmentVariableIfNotSet(builtSidecar.EnvironmentVariables, "PROXY_SEND_TIMEOUT", params.Request.Timeout)
		builtSidecar.EnvironmentVariables = addEnvironmentVariableIfNotSet(builtSidecar.EnvironmentVariables, "PROXY_READ_TIMEOUT", params.Request.Timeout)
		builtSidecar.EnvironmentVariables = addEnvironmentVariableIfNotSet(builtSidecar.EnvironmentVariables, "CLIENT_MAX_BODY_SIZE", params.Request.MaxBodySize)
		builtSidecar.EnvironmentVariables = addEnvironmentVariableIfNotSet(builtSidecar.EnvironmentVariables, "CLIENT_BODY_BUFFER_SIZE", params.Request.ClientBodyBufferSize)
		builtSidecar.EnvironmentVariables = addEnvironmentVariableIfNotSet(builtSidecar.EnvironmentVariables, "PROXY_BUFFER_SIZE", params.Request.ProxyBufferSize)
		builtSidecar.EnvironmentVariables = addEnvironmentVariableIfNotSet(builtSidecar.EnvironmentVariables, "PROXY_BUFFERS_SIZE", params.Request.ProxyBufferSize)
		builtSidecar.EnvironmentVariables = addEnvironmentVariableIfNotSet(builtSidecar.EnvironmentVariables, "PROXY_BUFFERS_NUMBER", strconv.Itoa(params.Request.ProxyBuffersNumber))

		if params.Container.Lifecycle.PrestopSleep != nil && *params.Container.Lifecycle.PrestopSleep && params.Container.Lifecycle.PrestopSleepSeconds != nil {
			builtSidecar.EnvironmentVariables = addEnvironmentVariableIfNotSet(builtSidecar.EnvironmentVariables, "GRACEFUL_SHUTDOWN_DELAY_SECONDS", strconv.Itoa(*params.Container.Lifecycle.PrestopSleepSeconds))
		}

		if params.Visibility == VisibilityESP {
			builtSidecar.EnvironmentVariables = addEnvironmentVariableIfNotSet(builtSidecar.EnvironmentVariables, "ENFORCE_HTTPS", "false")
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

	// replace @ with -at-
	reg := regexp.MustCompile(`@+`)
	value = reg.ReplaceAllString(value, "-at-")

	// replace all invalid characters with a hyphen
	reg = regexp.MustCompile(`[^a-zA-Z0-9-_.]+`)
	value = reg.ReplaceAllString(value, "-")

	// replace double hyphens with a single one
	value = strings.Replace(value, "--", "-", -1)

	// ensure it starts with an alphanumeric character
	reg = regexp.MustCompile(`^[^a-zA-Z0-9]+`)
	value = reg.ReplaceAllString(value, "")

	// maximize length at 63 characters
	if len(value) > 63 {
		value = value[:63]
	}

	// ensure it ends with an alphanumeric character
	reg = regexp.MustCompile(`[^a-zA-Z0-9]+$`)
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

func isSimpleEnvvarValue(i interface{}) bool {
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

func toYAML(v interface{}) string {
	data, err := yaml.Marshal(v)
	if err != nil {
		// Swallow errors inside of a template.
		return ""
	}
	return string(data)
}

func renderToYAML(v interface{}, data interface{}) string {

	value := toYAML(v)

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
