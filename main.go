package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/alecthomas/kingpin"
	foundation "github.com/estafette/estafette-foundation"
	"github.com/rs/zerolog/log"
	"github.com/sethgrid/pester"
	"gopkg.in/yaml.v2"
)

var (
	appgroup  string
	app       string
	version   string
	branch    string
	revision  string
	buildDate string
	goVersion = runtime.Version()
)

var (
	// flags
	paramsJSON      = kingpin.Flag("params", "Extension parameters, created from custom properties.").Envar("ESTAFETTE_EXTENSION_CUSTOM_PROPERTIES").Required().String()
	paramsYAML      = kingpin.Flag("params-yaml", "Extension parameters, created from custom properties.").Envar("ESTAFETTE_EXTENSION_CUSTOM_PROPERTIES_YAML").Required().String()
	credentialsJSON = kingpin.Flag("credentials", "GKE credentials configured at service level, passed in to this trusted extension.").Envar("ESTAFETTE_CREDENTIALS_KUBERNETES_ENGINE").Required().String()

	// optional flags
	gitSource     = kingpin.Flag("git-source", "Repository source.").Envar("ESTAFETTE_GIT_SOURCE").String()
	gitOwner      = kingpin.Flag("git-owner", "Repository owner.").Envar("ESTAFETTE_GIT_OWNER").String()
	gitName       = kingpin.Flag("git-name", "Repository name, used as application name if not passed explicitly and app label not being set.").Envar("ESTAFETTE_GIT_NAME").String()
	gitBranch     = kingpin.Flag("git-branch", "Repository commit branch.").Envar("ESTAFETTE_GIT_BRANCH").String()
	gitRevision   = kingpin.Flag("git-revision", "Repository commit revisition.").Envar("ESTAFETTE_GIT_REVISION").String()
	appLabel      = kingpin.Flag("app-name", "App label, used as application name if not passed explicitly.").Envar("ESTAFETTE_LABEL_APP").String()
	buildVersion  = kingpin.Flag("build-version", "Version number, used if not passed explicitly.").Envar("ESTAFETTE_BUILD_VERSION").String()
	releaseName   = kingpin.Flag("release-name", "Name of the release section, which is used by convention to resolve the credentials.").Envar("ESTAFETTE_RELEASE_NAME").String()
	releaseAction = kingpin.Flag("release-action", "Name of the release action, to control the type of release.").Envar("ESTAFETTE_RELEASE_ACTION").String()
	releaseID     = kingpin.Flag("release-id", "ID of the release, to use as a label.").Envar("ESTAFETTE_RELEASE_ID").String()
	triggeredBy   = kingpin.Flag("triggered-by", "The user id of the person triggering the release.").Envar("ESTAFETTE_TRIGGER_MANUAL_USER_ID").String()

	assistTroubleshootingOnError = false
	paramsForTroubleshooting     = Params{}
)

func main() {

	// parse command line parameters
	kingpin.Parse()

	// init log format from envvar ESTAFETTE_LOG_FORMAT
	foundation.InitLoggingFromEnv(appgroup, app, version, branch, revision, buildDate)

	// create context to cancel commands on sigterm
	ctx := foundation.InitCancellationContext(context.Background())

	// put all estafette labels in map
	log.Info().Msg("Getting all estafette labels from envvars...")
	estafetteLabels := map[string]string{}
	for _, e := range os.Environ() {
		kvPair := strings.SplitN(e, "=", 2)

		if len(kvPair) == 2 {
			envvarName := kvPair[0]
			envvarValue := kvPair[1]

			if strings.HasPrefix(envvarName, "ESTAFETTE_LABEL_") && !strings.HasSuffix(envvarName, "_DNS_SAFE") {
				// strip prefix and convert to lowercase
				key := strings.ToLower(strings.Replace(envvarName, "ESTAFETTE_LABEL_", "", 1))
				estafetteLabels[key] = envvarValue
			}
		}
	}

	log.Info().Msg("Unmarshalling credentials parameter...")
	var credentialsParam CredentialsParam
	err := json.Unmarshal([]byte(*paramsJSON), &credentialsParam)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed unmarshalling credential parameter")
	}

	log.Info().Msg("Setting default for credential parameter...")
	credentialsParam.SetDefaults(*releaseName)

	log.Info().Msg("Validating required credential parameter...")
	valid, errors := credentialsParam.ValidateRequiredProperties()
	if !valid {
		log.Fatal().Msgf("Not all valid fields are set: %v", errors)
	}

	log.Info().Msg("Unmarshalling injected credentials...")
	var credentials []GKECredentials
	err = json.Unmarshal([]byte(*credentialsJSON), &credentials)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed unmarshalling injected credentials")
	}

	log.Info().Msgf("Checking if credential %v exists...", credentialsParam.Credentials)
	credential := GetCredentialsByName(credentials, credentialsParam.Credentials)
	if credential == nil {
		log.Fatal().Msgf("Credential with name %v does not exist.", credentialsParam.Credentials)
	}

	var params Params
	if credential.AdditionalProperties.Defaults != nil {
		log.Info().Msgf("Using defaults from credential %v...", credentialsParam.Credentials)
		// todo log just the specified defaults, not the entire parms object
		// defaultsAsYAML, err := yaml.Marshal(credential.AdditionalProperties.Defaults)
		// if err == nil {
		// 	log.Printf(string(defaultsAsYAML))
		// }
		params = *credential.AdditionalProperties.Defaults
	}

	log.Info().Msg("Unmarshalling parameters / custom properties...")
	err = yaml.Unmarshal([]byte(*paramsYAML), &params)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed unmarshalling parameters")
	}

	log.Info().Msg("Setting defaults for parameters that are not set in the manifest...")
	params.SetDefaults(*gitName, *appLabel, *buildVersion, *releaseName, *releaseAction, estafetteLabels)

	log.Info().Msg("Validating required parameters...")
	valid, errors, warnings := params.ValidateRequiredProperties()
	if !valid {
		log.Fatal().Msgf("Not all valid fields are set: %v", errors)
	}

	for _, warning := range warnings {
		log.Printf("Warning: %s", warning)
	}

	// check for visibility esp if openapi.yaml exists
	if _, err := os.Stat(params.EspOpenAPIYamlPath); params.Visibility == "esp" && os.IsNotExist(err) {
		log.Fatal().Err(err).Msg("When using visibility: esp make sure to set clone: true and have openapi.yaml available in the working directory")
	}

	// replacing sidecar image tags with digest
	params.ReplaceSidecarTagsWithDigest()

	log.Info().Msg("Retrieving service account email from credentials...")
	var keyFileMap map[string]interface{}
	err = json.Unmarshal([]byte(credential.AdditionalProperties.ServiceAccountKeyfile), &keyFileMap)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed unmarshalling service account keyfile")
	}
	var saClientEmail string
	if saClientEmailIntfc, ok := keyFileMap["client_email"]; !ok {
		log.Fatal().Msg("Field client_email missing from service account keyfile")
	} else {
		if t, aok := saClientEmailIntfc.(string); !aok {
			log.Fatal().Msg("Field client_email not of type string")
		} else {
			saClientEmail = t
		}
	}

	log.Info().Msgf("Storing gke credential %v on disk...", credentialsParam.Credentials)
	err = ioutil.WriteFile("/key-file.json", []byte(credential.AdditionalProperties.ServiceAccountKeyfile), 0600)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed writing service account keyfile")
	}

	log.Info().Msg("Authenticating to google cloud")
	foundation.RunCommandWithArgs(ctx, "gcloud", []string{"auth", "activate-service-account", saClientEmail, "--key-file", "/key-file.json"})

	log.Info().Msgf("Setting gcloud account to %v", saClientEmail)
	foundation.RunCommandWithArgs(ctx, "gcloud", []string{"config", "set", "account", saClientEmail})

	log.Info().Msg("Setting gcloud project")
	foundation.RunCommandWithArgs(ctx, "gcloud", []string{"config", "set", "project", credential.AdditionalProperties.Project})

	log.Info().Msgf("Getting gke credentials for cluster %v", credential.AdditionalProperties.Cluster)
	clustersGetCredentialsArsgs := []string{"container", "clusters", "get-credentials", credential.AdditionalProperties.Cluster}
	if credential.AdditionalProperties.Zone != "" {
		clustersGetCredentialsArsgs = append(clustersGetCredentialsArsgs, "--zone", credential.AdditionalProperties.Zone)
	} else if credential.AdditionalProperties.Region != "" {
		clustersGetCredentialsArsgs = append(clustersGetCredentialsArsgs, "--region", credential.AdditionalProperties.Region)
	} else {
		log.Fatal().Msg("Credentials have no zone or region; at least one of them has to be defined")
	}
	foundation.RunCommandWithArgs(ctx, "gcloud", clustersGetCredentialsArsgs)

	// combine templates
	tmpl, err := buildTemplates(params)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed building templates")
	}

	// pre-render config files if they exist
	params.Configs.RenderedFileContent = renderConfig(params)

	// checking number of replicas for existing deployment to make switching deployment type safe
	currentReplicas := params.Replicas
	if params.Kind == "deployment" || params.Kind == "headless-deployment" {
		currentReplicas = getExistingNumberOfReplicas(ctx, params)
	}

	// generate the data required for rendering the templates
	templateData := generateTemplateData(params, currentReplicas, *gitSource, *gitOwner, *gitName, *gitBranch, *gitRevision, *releaseID, *triggeredBy)

	// render the template
	renderedTemplate, err := renderTemplate(tmpl, templateData)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed rendering templates")
	}

	if tmpl != nil {
		log.Info().Msg("Storing rendered manifest on disk...")
		err = ioutil.WriteFile("/kubernetes.yaml", renderedTemplate.Bytes(), 0600)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed writing manifest")
		}
	}

	kubectlApplyArgs := []string{"apply", "-f", "/kubernetes.yaml", "-n", templateData.Namespace}
	if tmpl != nil {
		// always perform a dryrun to ensure we're not ending up in a semi broken state where half of the templates is successfully applied and others not
		log.Info().Msg("Performing a dryrun to test the validity of the manifests...")
		foundation.RunCommandWithArgs(ctx, "kubectl", append(kubectlApplyArgs, "--dry-run"))
	}

	if !params.DryRun {

		// ensure that from now on any error runs the troubleshooting assistant
		assistTroubleshootingOnError = true
		paramsForTroubleshooting = params

		if tmpl != nil {
			deployGoogleEndpointsServiceIfRequired(ctx, params)
			patchServiceIfRequired(ctx, params, templateData, templateData.Name, templateData.Namespace)
			patchDeploymentIfRequired(ctx, params, templateData.Name, templateData.Namespace)
			removePoddisruptionBudgetIfRequired(ctx, params, templateData.NameWithTrack, templateData.Namespace)
			removeIngressIfRequired(ctx, params, templateData, templateData.Name, templateData.Namespace)
			cleanupJobIfRequired(ctx, params, templateData, templateData.Name, templateData.Namespace)

			log.Info().Msg("Applying the manifests for real...")
			foundation.RunCommandWithArgs(ctx, "kubectl", kubectlApplyArgs)

			if params.Kind == "deployment" || params.Kind == "headless-deployment" {
				log.Info().Msg("Waiting for the deployment to finish...")
				err = foundation.RunCommandWithArgsExtended(ctx, "kubectl", []string{"rollout", "status", "deployment", templateData.NameWithTrack, "-n", templateData.Namespace})
			}
			if params.Kind == "statefulset" {
				log.Info().Msg("Waiting for the statefulset to finish...")
				err = foundation.RunCommandWithArgsExtended(ctx, "kubectl", []string{"rollout", "status", "statefulset", templateData.Name, "-n", templateData.Namespace})
			}
		}

		if err != nil {
			assistTroubleshooting(ctx, templateData, err)
		}

		// clean up old stuff
		switch params.Kind {
		case "deployment":
			switch params.Action {
			case "deploy-canary":
				scaleCanaryDeployment(ctx, templateData.Name, templateData.Namespace, 1)
				deleteConfigsForParamsChange(ctx, params, templateData.NameWithTrack, templateData.Namespace)
				deleteSecretsForParamsChange(ctx, params, templateData.NameWithTrack, templateData.Namespace)
				break
			case "deploy-stable":
				scaleCanaryDeployment(ctx, templateData.Name, templateData.Namespace, 0)
				deleteResourcesForTypeSwitch(ctx, templateData.Name, templateData.Namespace)
				deleteConfigsForParamsChange(ctx, params, templateData.NameWithTrack, templateData.Namespace)
				deleteSecretsForParamsChange(ctx, params, templateData.NameWithTrack, templateData.Namespace)
				deleteServiceAccountSecretForParamsChange(ctx, params, templateData.GoogleCloudCredentialsAppName, templateData.Namespace)
				deleteIngressForVisibilityChange(ctx, templateData, templateData.Name, templateData.Namespace)
				removeEstafetteCloudflareAnnotations(ctx, templateData, templateData.Name, templateData.Namespace)
				removeBackendConfigAnnotation(ctx, templateData, templateData.Name, templateData.Namespace)
				deleteBackendConfigAndIAPOauthSecret(ctx, templateData, templateData.Name, templateData.Namespace)
				deleteHorizontalPodAutoscaler(ctx, params, templateData.NameWithTrack, templateData.Namespace)
				break
			case "rollback-canary":
				scaleCanaryDeployment(ctx, templateData.Name, templateData.Namespace, 0)
				break
			case "deploy-simple":
				deleteResourcesForTypeSwitch(ctx, fmt.Sprintf("%v-canary", templateData.Name), templateData.Namespace)
				deleteResourcesForTypeSwitch(ctx, fmt.Sprintf("%v-stable", templateData.Name), templateData.Namespace)
				deleteConfigsForParamsChange(ctx, params, templateData.Name, templateData.Namespace)
				deleteSecretsForParamsChange(ctx, params, templateData.Name, templateData.Namespace)
				deleteServiceAccountSecretForParamsChange(ctx, params, templateData.GoogleCloudCredentialsAppName, templateData.Namespace)
				deleteIngressForVisibilityChange(ctx, templateData, templateData.Name, templateData.Namespace)
				removeEstafetteCloudflareAnnotations(ctx, templateData, templateData.Name, templateData.Namespace)
				removeBackendConfigAnnotation(ctx, templateData, templateData.Name, templateData.Namespace)
				deleteBackendConfigAndIAPOauthSecret(ctx, templateData, templateData.Name, templateData.Namespace)
				deleteHorizontalPodAutoscaler(ctx, params, templateData.Name, templateData.Namespace)
				break
			}
			break

		case "headless-deployment":
			switch params.Action {
			case "deploy-canary":
				scaleCanaryDeployment(ctx, templateData.Name, templateData.Namespace, 1)
				deleteConfigsForParamsChange(ctx, params, templateData.NameWithTrack, templateData.Namespace)
				deleteSecretsForParamsChange(ctx, params, templateData.NameWithTrack, templateData.Namespace)
				break
			case "deploy-stable":
				scaleCanaryDeployment(ctx, templateData.Name, templateData.Namespace, 0)
				deleteResourcesForTypeSwitch(ctx, templateData.Name, templateData.Namespace)
				deleteConfigsForParamsChange(ctx, params, templateData.NameWithTrack, templateData.Namespace)
				deleteSecretsForParamsChange(ctx, params, templateData.NameWithTrack, templateData.Namespace)
				deleteServiceAccountSecretForParamsChange(ctx, params, templateData.GoogleCloudCredentialsAppName, templateData.Namespace)
				deleteHorizontalPodAutoscaler(ctx, params, templateData.NameWithTrack, templateData.Namespace)
				break
			case "rollback-canary":
				scaleCanaryDeployment(ctx, templateData.Name, templateData.Namespace, 0)
				break
			case "deploy-simple":
				deleteResourcesForTypeSwitch(ctx, fmt.Sprintf("%v-canary", templateData.Name), templateData.Namespace)
				deleteResourcesForTypeSwitch(ctx, fmt.Sprintf("%v-stable", templateData.Name), templateData.Namespace)
				deleteConfigsForParamsChange(ctx, params, templateData.Name, templateData.Namespace)
				deleteSecretsForParamsChange(ctx, params, templateData.Name, templateData.Namespace)
				deleteServiceAccountSecretForParamsChange(ctx, params, templateData.GoogleCloudCredentialsAppName, templateData.Namespace)
				deleteHorizontalPodAutoscaler(ctx, params, templateData.Name, templateData.Namespace)
				break
			}
			break
		case "statefulset":
			deleteConfigsForParamsChange(ctx, params, templateData.Name, templateData.Namespace)
			deleteSecretsForParamsChange(ctx, params, templateData.Name, templateData.Namespace)
			deleteServiceAccountSecretForParamsChange(ctx, params, templateData.GoogleCloudCredentialsAppName, templateData.Namespace)
			deleteIngressForVisibilityChange(ctx, templateData, templateData.Name, templateData.Namespace)
			removeEstafetteCloudflareAnnotations(ctx, templateData, templateData.Name, templateData.Namespace)
			removeBackendConfigAnnotation(ctx, templateData, templateData.Name, templateData.Namespace)
			deleteBackendConfigAndIAPOauthSecret(ctx, templateData, templateData.Name, templateData.Namespace)
			break
		}

		assistTroubleshooting(ctx, templateData, err)
	}
}

func assistTroubleshooting(ctx context.Context, templateData TemplateData, err error) {
	if assistTroubleshootingOnError {
		log.Info().Msgf("Showing current ingresses, services, configmaps, secrets, deployments, jobs, cronjobs, poddisruptionbudgets, horizontalpodautoscalers, pods, endpoints for app=%v...", paramsForTroubleshooting.App)
		foundation.RunCommandWithArgsExtended(ctx, "kubectl", []string{"get", "ing,svc,cm,secret,deploy,job,cronjob,sts,pdb,hpa,po,ep", "-l", fmt.Sprintf("app=%v", paramsForTroubleshooting.App), "-n", paramsForTroubleshooting.Namespace})

		if err != nil {
			log.Info().Msg("Rollout failed, trying to show logs...")
			if *releaseID != "" {
				_ = foundation.RunCommandWithArgsExtended(ctx, "kubectl", []string{"logs", "-l", fmt.Sprintf("app=%v,estafette.io/release-id=%v", templateData.AppLabelSelector, sanitizeLabel(*releaseID)), "-n", templateData.Namespace, "--all-containers"})
			} else if *buildVersion != "" {
				_ = foundation.RunCommandWithArgsExtended(ctx, "kubectl", []string{"logs", "-l", fmt.Sprintf("app=%v,version=%v", templateData.AppLabelSelector, sanitizeLabel(*buildVersion)), "-n", templateData.Namespace, "--all-containers"})
			}
		} else if paramsForTroubleshooting.Action == "deploy-canary" {
			log.Info().Msg("Showing logs for canary deployment...")
			foundation.RunCommandWithArgsExtended(ctx, "kubectl", []string{"logs", "-l", fmt.Sprintf("app=%v,track=canary", paramsForTroubleshooting.App), "-n", paramsForTroubleshooting.Namespace, "-c", paramsForTroubleshooting.App, "--tail", "50"})
		}

		foundation.HandleError(err)
	}
}

func scaleCanaryDeployment(ctx context.Context, name, namespace string, replicas int) {
	log.Info().Msgf("Scaling canary deployment to %v replicas...", replicas)
	foundation.RunCommandWithArgs(ctx, "kubectl", []string{"scale", "deploy", fmt.Sprintf("%v-canary", name), "-n", namespace, fmt.Sprintf("--replicas=%v", replicas)})
}

func deleteResourcesForTypeSwitch(ctx context.Context, name, namespace string) {
	// clean up resources in case a switch from simple to canary releases or vice versa has been made
	log.Info().Msg("Deleting simple type deployment, configmap, secret, hpa and pdb...")
	foundation.RunCommandWithArgs(ctx, "kubectl", []string{"delete", "deploy", name, "-n", namespace, "--ignore-not-found=true"})
	foundation.RunCommandWithArgs(ctx, "kubectl", []string{"delete", "configmap", fmt.Sprintf("%v-configs", name), "-n", namespace, "--ignore-not-found=true"})
	foundation.RunCommandWithArgs(ctx, "kubectl", []string{"delete", "secret", fmt.Sprintf("%v-secrets", name), "-n", namespace, "--ignore-not-found=true"})
	foundation.RunCommandWithArgs(ctx, "kubectl", []string{"delete", "hpa", name, "-n", namespace, "--ignore-not-found=true"})
	foundation.RunCommandWithArgs(ctx, "kubectl", []string{"delete", "pdb", name, "-n", namespace, "--ignore-not-found=true"})
}

func deleteConfigsForParamsChange(ctx context.Context, params Params, name, namespace string) {
	if len(params.Configs.Files) == 0 && len(params.Configs.InlineFiles) == 0 {
		log.Info().Msg("Deleting application configs if it exists, because no configs are specified...")
		foundation.RunCommandWithArgs(ctx, "kubectl", []string{"delete", "configmap", fmt.Sprintf("%v-configs", name), "-n", namespace, "--ignore-not-found=true"})
	}
}

func deleteSecretsForParamsChange(ctx context.Context, params Params, name, namespace string) {
	if len(params.Secrets.Keys) == 0 {
		log.Info().Msg("Deleting application secrets if it exists, because no secrets are specified...")
		foundation.RunCommandWithArgs(ctx, "kubectl", []string{"delete", "secret", fmt.Sprintf("%v-secrets", name), "-n", namespace, "--ignore-not-found=true"})
	}
}

func deleteServiceAccountSecretForParamsChange(ctx context.Context, params Params, name, namespace string) {
	if !params.UseGoogleCloudCredentials {
		log.Info().Msg("Deleting service account secret if it exists, because no use of service account is specified...")
		foundation.RunCommandWithArgs(ctx, "kubectl", []string{"delete", "secret", fmt.Sprintf("%v-gcp-service-account", name), "-n", namespace, "--ignore-not-found=true"})
	}
}

func deleteIngressForVisibilityChange(ctx context.Context, templateData TemplateData, name, namespace string) {
	if !templateData.UseNginxIngress && !templateData.UseGCEIngress {
		// public uses service of type loadbalancer and doesn't need ingress
		log.Info().Msg("Deleting ingress if it exists, which is used for visibility private, iap or public-whitelist...")
		foundation.RunCommandWithArgs(ctx, "kubectl", []string{"delete", "ingress", name, "-n", namespace, "--ignore-not-found=true"})
	}
}

func deleteBackendConfigAndIAPOauthSecret(ctx context.Context, templateData TemplateData, name, namespace string) {
	if !templateData.UseBackendConfigAnnotationOnService {
		log.Info().Msg("Deleting iap oauth secret if it exists, because visibility is not set to iap...")
		foundation.RunCommandWithArgs(ctx, "kubectl", []string{"delete", "secret", fmt.Sprintf("%v--iap-oauth-credentials", name), "-n", namespace, "--ignore-not-found=true"})
		log.Info().Msg("Deleting iap backend config if it exists, because visibility is not set to iap...")
		foundation.RunCommandWithArgs(ctx, "kubectl", []string{"delete", "backendconfig", name, "-n", namespace, "--ignore-not-found=true"})
	}
}

func removePoddisruptionBudgetIfRequired(ctx context.Context, params Params, name, namespace string) {
	if (params.Kind == "deployment" || params.Kind == "headless-deployment") && (params.Action == "deploy-simple" || params.Action == "deploy-stable") {
		// if there's a pdb that doesn't use maxUnavailable: 1 remove it so a new one can be created with correct settings
		deletePoddisruptionBudget := false
		maxUnavailable, err := foundation.GetCommandWithArgsOutput(ctx, "kubectl", []string{"get", "pdb", name, "-n", namespace, "-o=jsonpath={.spec.maxUnavailable}"})
		if err == nil {
			maxUnavailableInt, err := strconv.Atoi(maxUnavailable)
			if err == nil {
				if maxUnavailableInt != 1 {
					log.Info().Msgf("MaxUnavailable from pdb %v is %v instead of 1", name, maxUnavailableInt)
					deletePoddisruptionBudget = true
				}
			} else {
				log.Info().Msgf("Failed reading maxUnavailable from pdb %v: %v", name, err)
				deletePoddisruptionBudget = true
			}
		} else {
			log.Info().Msgf("Failed retrieving pdb %v: %v", name, err)
			deletePoddisruptionBudget = true
		}

		if deletePoddisruptionBudget {
			foundation.RunCommandWithArgs(ctx, "kubectl", []string{"delete", "pdb", name, "-n", namespace, "--ignore-not-found=true"})
		} else {
			log.Info().Msgf("Poddisruptionbudet %v is fine, not removing it", name)
		}
	}
}

func removeIngressIfRequired(ctx context.Context, params Params, templateData TemplateData, name, namespace string) {
	if params.Kind == "deployment" && (params.Action == "deploy-simple" || params.Action == "deploy-canary" || params.Action == "deploy-stable") {
		if templateData.UseNginxIngress {
			// check if ingress exists and has kubernetes.io/ingress.class: gce, then delete it because of https://github.com/kubernetes/ingress-gce/issues/481
			ingressClass, err := foundation.GetCommandWithArgsOutput(ctx, "kubectl", []string{"get", "ing", name, "-n", namespace, "-o=go-template={{index .metadata.annotations \"kubernetes.io/ingress.class\"}}"})
			if err == nil {
				if ingressClass == "gce" {
					// delete the ingress so all related load balancers, etc get deleted
					log.Info().Msg("Deleting ingress so the gce ingress controller removes the related load balancer...")
					foundation.RunCommandWithArgs(ctx, "kubectl", []string{"delete", "ingress", name, "-n", namespace, "--ignore-not-found=true"})
				} else {
					log.Info().Msgf("Ingress %v already has kubernetes.io/ingress.class: %v annotation, no need to delete the ingress", name, ingressClass)
				}
			} else {
				log.Info().Msgf("Ingress %v or kubernetes.io/ingress.class annotation doesn't exist, no need to delete the ingress: %v", name, err)
			}
		} else if templateData.UseGCEIngress {
			// check if ingress exists and has kubernetes.io/ingress.class: gce, then delete it to ensure there's no nginx ingress annotations lingering around
			ingressClass, err := foundation.GetCommandWithArgsOutput(ctx, "kubectl", []string{"get", "ing", name, "-n", namespace, "-o=go-template={{index .metadata.annotations \"kubernetes.io/ingress.class\"}}"})
			if err == nil {
				if ingressClass == "nginx" {
					// delete the ingress so all related nginx ingress config gets deleted
					log.Info().Msg("Deleting ingress so the nginx ingress controller removes related config...")
					foundation.RunCommandWithArgs(ctx, "kubectl", []string{"delete", "ingress", name, "-n", namespace, "--ignore-not-found=true"})
				} else {
					log.Info().Msgf("Ingress %v already has kubernetes.io/ingress.class: %v annotation, no need to delete the ingress", name, ingressClass)
				}
			} else {
				log.Info().Msgf("Ingress %v or kubernetes.io/ingress.class annotation doesn't exist, no need to delete the ingress: %v", name, err)
			}
		}
	}
}

func deployGoogleEndpointsServiceIfRequired(ctx context.Context, params Params) {
	if params.Kind == "deployment" && params.Visibility == "esp" && (params.Action == "deploy-simple" || params.Action == "deploy-canary") {
		foundation.RunCommandWithArgs(ctx, "gcloud", []string{"endpoints", "services", "deploy", params.EspOpenAPIYamlPath})
	}
}

func patchServiceIfRequired(ctx context.Context, params Params, templateData TemplateData, name, namespace string) {
	if params.Kind == "deployment" && templateData.ServiceType == "ClusterIP" {
		serviceType, err := foundation.GetCommandWithArgsOutput(ctx, "kubectl", []string{"get", "service", name, "-n", namespace, "-o=jsonpath={.spec.type}"})
		if err != nil {
			log.Info().Msgf("Failed retrieving service type: %v", err)
		}
		if err == nil && (serviceType == "NodePort" || serviceType == "LoadBalancer") {
			log.Info().Msgf("Service is of type %v, patching it...", serviceType)

			// brute force patch the service
			err = foundation.RunCommandWithArgsExtended(ctx, "kubectl", []string{"patch", "service", name, "-n", namespace, "--type", "json", "--patch", "[{\"op\": \"remove\", \"path\": \"/spec/loadBalancerSourceRanges\"},{\"op\": \"remove\", \"path\": \"/spec/externalTrafficPolicy\"}, {\"op\": \"remove\", \"path\": \"/spec/ports/0/nodePort\"}, {\"op\": \"remove\", \"path\": \"/spec/ports/1/nodePort\"}, {\"op\": \"replace\", \"path\": \"/spec/type\", \"value\": \"ClusterIP\"}]"})
			if err != nil {
				err = foundation.RunCommandWithArgsExtended(ctx, "kubectl", []string{"patch", "service", name, "-n", namespace, "--type", "json", "--patch", "[{\"op\": \"remove\", \"path\": \"/spec/externalTrafficPolicy\"}, {\"op\": \"remove\", \"path\": \"/spec/ports/0/nodePort\"}, {\"op\": \"remove\", \"path\": \"/spec/ports/1/nodePort\"}, {\"op\": \"replace\", \"path\": \"/spec/type\", \"value\": \"ClusterIP\"}]"})
			}
			if err != nil {
				log.Fatal().Err(err).Msg(fmt.Sprintf("Failed patching service to change from %v to ClusterIP", serviceType))
			}
		} else {
			log.Info().Msgf("Service is of type %v, no need to patch it", serviceType)
		}
	}
}

func cleanupJobIfRequired(ctx context.Context, params Params, templateData TemplateData, name, namespace string) {
	if params.Kind == "job" {
		err := foundation.RunCommandWithArgsExtended(ctx, "kubectl", []string{"delete", "job", name, "-n", namespace, "--ignore-not-found=true"})
		if err != nil {
			log.Info().Msgf("Deleting job %v failed: %v", name, err)
		}
	}
	if params.Kind == "cronjob" {
		err := foundation.RunCommandWithArgsExtended(ctx, "kubectl", []string{"delete", "cronjob", name, "-n", namespace, "--ignore-not-found=true"})
		if err != nil {
			log.Info().Msgf("Deleting cronjob %v failed: %v", name, err)
		}
	}
}

func getExistingNumberOfReplicas(ctx context.Context, params Params) int {
	if params.Kind == "deployment" || params.Kind == "headless-deployment" {
		deploymentName := ""
		if params.Action == "deploy-simple" {
			deploymentName = params.App + "-stable"
		} else if params.Action == "deploy-stable" {
			deploymentName = params.App
		}
		if deploymentName != "" {
			replicas, err := foundation.GetCommandWithArgsOutput(ctx, "kubectl", []string{"get", "deploy", deploymentName, "-n", params.Namespace, "-o=jsonpath={.spec.replicas}"})
			if err != nil {
				log.Info().Msgf("Failed retrieving replicas for %v: %v ignoring setting replicas since there's no switch for deployment type...", deploymentName, err)
				return -1
			}
			replicasInt, err := strconv.Atoi(replicas)
			if err != nil {
				log.Info().Msgf("Failed converting replicas value %v for %v: %v ignoring setting replicas since there's no switch for deployment type...", replicas, deploymentName, err)
				return -1
			}
			log.Info().Msgf("Retrieved number of replicas for %v is %v; using it to set correct number of replicas switching deployment type...", deploymentName, replicasInt)
			return replicasInt
		}
	}

	return -1
}

func patchDeploymentIfRequired(ctx context.Context, params Params, name, namespace string) {
	if (params.Kind == "deployment" || params.Kind == "headless-deployment") && params.Action == "deploy-simple" {
		selectorLabels, err := foundation.GetCommandWithArgsOutput(ctx, "kubectl", []string{"get", "deploy", name, "-n", namespace, "-o=jsonpath={.spec.selector.matchLabels}"})
		if err != nil {
			log.Info().Msgf("Failed retrieving deployment selector labels: %v", err)
		}
		if err == nil && selectorLabels != fmt.Sprintf("map[app:%v]", name) {
			log.Info().Msgf("Deployment selector labels %v not correct, patching it...", selectorLabels)

			// patch the deployment
			err = foundation.RunCommandWithArgsExtended(ctx, "kubectl", []string{"patch", "deploy", name, "-n", namespace, "--type", "json", "--patch", fmt.Sprintf("[{\"op\": \"replace\", \"path\": \"/spec/selector/matchLabels\", \"value\": {\"app\":\"%v\"}}]", name)})
			if err != nil {
				log.Fatal().Err(err).Msg(fmt.Sprintf("Failed patching deployment to change selector labels from %v to app=%v", selectorLabels, name))
			}
		} else {
			log.Info().Msgf("Deployment selector labels %v are correct, not patching", selectorLabels)
		}
	}
}

func removeEstafetteCloudflareAnnotations(ctx context.Context, templateData TemplateData, name, namespace string) {
	if !templateData.UseDNSAnnotationsOnService {
		// ingress is used and has the estafette.io/cloudflare annotations, so they should be removed from the service
		log.Info().Msg("Removing estafette.io/cloudflare annotations on the service if they exists, since they're now set on the ingress instead...")
		foundation.RunCommandWithArgs(ctx, "kubectl", []string{"annotate", "svc", name, "-n", namespace, "estafette.io/cloudflare-dns-"})
		foundation.RunCommandWithArgs(ctx, "kubectl", []string{"annotate", "svc", name, "-n", namespace, "estafette.io/cloudflare-proxy-"})
		foundation.RunCommandWithArgs(ctx, "kubectl", []string{"annotate", "svc", name, "-n", namespace, "estafette.io/cloudflare-hostnames-"})
		foundation.RunCommandWithArgs(ctx, "kubectl", []string{"annotate", "svc", name, "-n", namespace, "estafette.io/cloudflare-state-"})
	}
}

func removeBackendConfigAnnotation(ctx context.Context, templateData TemplateData, name, namespace string) {
	if !templateData.UseBackendConfigAnnotationOnService {
		// iap is not used, so the beta.cloud.google.com/backend-config annotations should be removed from the service
		log.Info().Msg("Removing beta.cloud.google.com/backend-config annotations on the service if they exists, since visibility is not set to iap...")
		foundation.RunCommandWithArgs(ctx, "kubectl", []string{"annotate", "svc", name, "-n", namespace, "beta.cloud.google.com/backend-config-"})
	}
}

func deleteHorizontalPodAutoscaler(ctx context.Context, params Params, name, namespace string) {
	if (params.Kind == "deployment" || params.Kind == "headless-deployment") && (params.Autoscale.Enabled == nil || !*params.Autoscale.Enabled) && (params.Action == "deploy-simple" || params.Action == "deploy-stable") {
		log.Info().Msgf("Deleting HorizontalPodAutoscaler %v, since autoscaling is disabled...", name)
		foundation.RunCommandWithArgs(ctx, "kubectl", []string{"delete", "hpa", name, "-n", namespace, "--ignore-not-found=true"})
	}
}

func httpRequestBody(method, url string, headers map[string]string) string {
	client := pester.New()
	client.MaxRetries = 3
	client.Backoff = pester.ExponentialJitterBackoff
	client.KeepLog = true
	client.Timeout = time.Second * 5
	request, err := http.NewRequest(method, url, nil)
	if err != nil {
		return ""
	}

	for k, v := range headers {
		request.Header.Add(k, v)
	}

	// perform actual request
	response, err := client.Do(request)
	if err != nil {
		return ""
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return ""
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return ""
	}

	return string(body)
}

func httpRequestHeader(method, url string, headers map[string]string, responseHeader string) string {
	client := pester.New()
	client.MaxRetries = 3
	client.Backoff = pester.ExponentialJitterBackoff
	client.KeepLog = true
	client.Timeout = time.Second * 5
	request, err := http.NewRequest(method, url, nil)
	if err != nil {
		return ""
	}

	for k, v := range headers {
		request.Header.Add(k, v)
	}

	// perform actual request
	response, err := client.Do(request)
	if err != nil {
		return ""
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return ""
	}

	return response.Header.Get(responseHeader)
}
