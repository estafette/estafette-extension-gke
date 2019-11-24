package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/alecthomas/kingpin"
	foundation "github.com/estafette/estafette-foundation"
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

	// put all estafette labels in map
	logInfo("Getting all estafette labels from envvars...")
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

	logInfo("Unmarshalling credentials parameter...")
	var credentialsParam CredentialsParam
	err := json.Unmarshal([]byte(*paramsJSON), &credentialsParam)
	if err != nil {
		log.Fatal("Failed unmarshalling credential parameter: ", err)
	}

	logInfo("Setting default for credential parameter...")
	credentialsParam.SetDefaults(*releaseName)

	logInfo("Validating required credential parameter...")
	valid, errors := credentialsParam.ValidateRequiredProperties()
	if !valid {
		log.Fatal("Not all valid fields are set: ", errors)
	}

	logInfo("Unmarshalling injected credentials...")
	var credentials []GKECredentials
	err = json.Unmarshal([]byte(*credentialsJSON), &credentials)
	if err != nil {
		log.Fatal("Failed unmarshalling injected credentials: ", err)
	}

	logInfo("Checking if credential %v exists...", credentialsParam.Credentials)
	credential := GetCredentialsByName(credentials, credentialsParam.Credentials)
	if credential == nil {
		log.Fatalf("Credential with name %v does not exist.", credentialsParam.Credentials)
	}

	var params Params
	if credential.AdditionalProperties.Defaults != nil {
		logInfo("Using defaults from credential %v...", credentialsParam.Credentials)
		// todo log just the specified defaults, not the entire parms object
		// defaultsAsYAML, err := yaml.Marshal(credential.AdditionalProperties.Defaults)
		// if err == nil {
		// 	log.Printf(string(defaultsAsYAML))
		// }
		params = *credential.AdditionalProperties.Defaults
	}

	logInfo("Unmarshalling parameters / custom properties...")
	err = yaml.Unmarshal([]byte(*paramsYAML), &params)
	if err != nil {
		log.Fatal("Failed unmarshalling parameters: ", err)
	}

	logInfo("Setting defaults for parameters that are not set in the manifest...")
	params.SetDefaults(*gitName, *appLabel, *buildVersion, *releaseName, *releaseAction, estafetteLabels)

	logInfo("Validating required parameters...")
	valid, errors, warnings := params.ValidateRequiredProperties()
	if !valid {
		log.Fatal("Not all valid fields are set: ", errors)
	}

	for _, warning := range warnings {
		log.Printf("Warning: %s", warning)
	}

	// check for visibility esp if openapi.yaml exists
	if _, err := os.Stat(params.EspOpenAPIYamlPath); params.Visibility == "esp" && os.IsNotExist(err) {
		log.Fatal("When using visibility: esp make sure to set clone: true and have openapi.yaml available in the working directory")
	}

	// replacing sidecar image tags with digest
	params.ReplaceSidecarTagsWithDigest()

	logInfo("Retrieving service account email from credentials...")
	var keyFileMap map[string]interface{}
	err = json.Unmarshal([]byte(credential.AdditionalProperties.ServiceAccountKeyfile), &keyFileMap)
	if err != nil {
		log.Fatal("Failed unmarshalling service account keyfile: ", err)
	}
	var saClientEmail string
	if saClientEmailIntfc, ok := keyFileMap["client_email"]; !ok {
		log.Fatal("Field client_email missing from service account keyfile")
	} else {
		if t, aok := saClientEmailIntfc.(string); !aok {
			log.Fatal("Field client_email not of type string")
		} else {
			saClientEmail = t
		}
	}

	logInfo("Storing gke credential %v on disk...", credentialsParam.Credentials)
	err = ioutil.WriteFile("/key-file.json", []byte(credential.AdditionalProperties.ServiceAccountKeyfile), 0600)
	if err != nil {
		log.Fatal("Failed writing service account keyfile: ", err)
	}

	logInfo("Authenticating to google cloud")
	foundation.RunCommandWithArgs("gcloud", []string{"auth", "activate-service-account", saClientEmail, "--key-file", "/key-file.json"})

	logInfo("Setting gcloud account to %v", saClientEmail)
	foundation.RunCommandWithArgs("gcloud", []string{"config", "set", "account", saClientEmail})

	logInfo("Setting gcloud project")
	foundation.RunCommandWithArgs("gcloud", []string{"config", "set", "project", credential.AdditionalProperties.Project})

	logInfo("Getting gke credentials for cluster %v", credential.AdditionalProperties.Cluster)
	clustersGetCredentialsArsgs := []string{"container", "clusters", "get-credentials", credential.AdditionalProperties.Cluster}
	if credential.AdditionalProperties.Zone != "" {
		clustersGetCredentialsArsgs = append(clustersGetCredentialsArsgs, "--zone", credential.AdditionalProperties.Zone)
	} else if credential.AdditionalProperties.Region != "" {
		clustersGetCredentialsArsgs = append(clustersGetCredentialsArsgs, "--region", credential.AdditionalProperties.Region)
	} else {
		log.Fatal("Credentials have no zone or region; at least one of them has to be defined")
	}
	foundation.RunCommandWithArgs("gcloud", clustersGetCredentialsArsgs)

	// combine templates
	tmpl, err := buildTemplates(params)
	if err != nil {
		log.Fatal("Failed building templates: ", err)
	}

	// pre-render config files if they exist
	params.Configs.RenderedFileContent = renderConfig(params)

	// checking number of replicas for existing deployment to make switching deployment type safe
	currentReplicas := params.Replicas
	if params.Kind == "deployment" {
		currentReplicas = getExistingNumberOfReplicas(params)
	}

	// generate the data required for rendering the templates
	templateData := generateTemplateData(params, currentReplicas, *gitSource, *gitOwner, *gitName, *gitBranch, *gitRevision, *releaseID, *triggeredBy)

	// render the template
	renderedTemplate, err := renderTemplate(tmpl, templateData)
	if err != nil {
		log.Fatal("Failed rendering templates: ", err)
	}

	if tmpl != nil {
		logInfo("Storing rendered manifest on disk...")
		err = ioutil.WriteFile("/kubernetes.yaml", renderedTemplate.Bytes(), 0600)
		if err != nil {
			log.Fatal("Failed writing manifest: ", err)
		}
	}

	kubectlApplyArgs := []string{"apply", "-f", "/kubernetes.yaml", "-n", templateData.Namespace}
	if tmpl != nil {
		// always perform a dryrun to ensure we're not ending up in a semi broken state where half of the templates is successfully applied and others not
		logInfo("Performing a dryrun to test the validity of the manifests...")
		foundation.RunCommandWithArgs("kubectl", append(kubectlApplyArgs, "--dry-run"))
	}

	if !params.DryRun {

		// ensure that from now on any error runs the troubleshooting assistant
		assistTroubleshootingOnError = true
		paramsForTroubleshooting = params

		if tmpl != nil {
			deployGoogleEndpointsServiceIfRequired(params)
			patchServiceIfRequired(params, templateData, templateData.Name, templateData.Namespace)
			patchDeploymentIfRequired(params, templateData.Name, templateData.Namespace)
			removePoddisruptionBudgetIfRequired(params, templateData.NameWithTrack, templateData.Namespace)
			removeIngressIfRequired(params, templateData, templateData.Name, templateData.Namespace)
			cleanupJobIfRequired(params, templateData, templateData.Name, templateData.Namespace)

			logInfo("Applying the manifests for real...")
			foundation.RunCommandWithArgs("kubectl", kubectlApplyArgs)

			if params.Kind == "deployment" {
				logInfo("Waiting for the deployment to finish...")
				foundation.RunCommandWithArgs("kubectl", []string{"rollout", "status", "deployment", templateData.NameWithTrack, "-n", templateData.Namespace})
			}
			if params.Kind == "statefulset" {
				logInfo("Waiting for the statefulset to finish...")
				foundation.RunCommandWithArgs("kubectl", []string{"rollout", "status", "statefulset", templateData.Name, "-n", templateData.Namespace})
			}
		}

		// clean up old stuff
		switch params.Kind {
		case "deployment":
			switch params.Action {
			case "deploy-canary":
				scaleCanaryDeployment(templateData.Name, templateData.Namespace, 1)
				deleteConfigsForParamsChange(params, templateData.NameWithTrack, templateData.Namespace)
				deleteSecretsForParamsChange(params, templateData.NameWithTrack, templateData.Namespace)
				break
			case "deploy-stable":
				scaleCanaryDeployment(templateData.Name, templateData.Namespace, 0)
				deleteResourcesForTypeSwitch(templateData.Name, templateData.Namespace)
				deleteConfigsForParamsChange(params, templateData.NameWithTrack, templateData.Namespace)
				deleteSecretsForParamsChange(params, templateData.NameWithTrack, templateData.Namespace)
				deleteServiceAccountSecretForParamsChange(params, templateData.GoogleCloudCredentialsAppName, templateData.Namespace)
				deleteIngressForVisibilityChange(templateData, templateData.Name, templateData.Namespace)
				removeEstafetteCloudflareAnnotations(templateData, templateData.Name, templateData.Namespace)
				removeBackendConfigAnnotation(templateData, templateData.Name, templateData.Namespace)
				deleteBackendConfigAndIAPOauthSecret(templateData, templateData.Name, templateData.Namespace)
				deleteHorizontalPodAutoscaler(params, templateData.NameWithTrack, templateData.Namespace)
				break
			case "rollback-canary":
				scaleCanaryDeployment(templateData.Name, templateData.Namespace, 0)
				break
			case "deploy-simple":
				deleteResourcesForTypeSwitch(fmt.Sprintf("%v-canary", templateData.Name), templateData.Namespace)
				deleteResourcesForTypeSwitch(fmt.Sprintf("%v-stable", templateData.Name), templateData.Namespace)
				deleteConfigsForParamsChange(params, templateData.Name, templateData.Namespace)
				deleteSecretsForParamsChange(params, templateData.Name, templateData.Namespace)
				deleteServiceAccountSecretForParamsChange(params, templateData.GoogleCloudCredentialsAppName, templateData.Namespace)
				deleteIngressForVisibilityChange(templateData, templateData.Name, templateData.Namespace)
				removeEstafetteCloudflareAnnotations(templateData, templateData.Name, templateData.Namespace)
				removeBackendConfigAnnotation(templateData, templateData.Name, templateData.Namespace)
				deleteBackendConfigAndIAPOauthSecret(templateData, templateData.Name, templateData.Namespace)
				deleteHorizontalPodAutoscaler(params, templateData.Name, templateData.Namespace)
				break
			}
			break
		case "statefulset":
			deleteConfigsForParamsChange(params, templateData.Name, templateData.Namespace)
			deleteSecretsForParamsChange(params, templateData.Name, templateData.Namespace)
			deleteServiceAccountSecretForParamsChange(params, templateData.GoogleCloudCredentialsAppName, templateData.Namespace)
			deleteIngressForVisibilityChange(templateData, templateData.Name, templateData.Namespace)
			removeEstafetteCloudflareAnnotations(templateData, templateData.Name, templateData.Namespace)
			removeBackendConfigAnnotation(templateData, templateData.Name, templateData.Namespace)
			deleteBackendConfigAndIAPOauthSecret(templateData, templateData.Name, templateData.Namespace)
			break
		}

		assistTroubleshooting()
	}
}

func assistTroubleshooting() {
	if assistTroubleshootingOnError {
		logInfo("Showing current ingresses, services, configmaps, secrets, deployments, jobs, cronjobs, poddisruptionbudgets, horizontalpodautoscalers, pods, endpoints for app=%v...", paramsForTroubleshooting.App)
		foundation.RunCommandWithArgsExtended("kubectl", []string{"get", "ing,svc,cm,secret,deploy,job,cronjob,sts,pdb,hpa,po,ep", "-l", fmt.Sprintf("app=%v", paramsForTroubleshooting.App), "-n", paramsForTroubleshooting.Namespace})

		if paramsForTroubleshooting.Action == "deploy-canary" {
			logInfo("Showing logs for canary deployment...")
			foundation.RunCommandWithArgsExtended("kubectl", []string{"logs", "-l", fmt.Sprintf("app=%v,track=canary", paramsForTroubleshooting.App), "-n", paramsForTroubleshooting.Namespace, "-c", paramsForTroubleshooting.App, "--tail", "50"})
		}

		// logInfo("Showing kubernetes events with the word %v in it...", paramsForTroubleshooting.App)
		// c1 := exec.Command("kubectl", "get", "events", "--sort-by=.metadata.creationTimestamp", "-n", paramsForTroubleshooting.Namespace)
		// c2 := exec.Command("grep", paramsForTroubleshooting.App)

		// r, w := io.Pipe()
		// c1.Stdout = w
		// c2.Stdin = r

		// var b2 bytes.Buffer
		// c2.Stdout = &b2

		// c1.Start()
		// c2.Start()
		// c1.Wait()
		// w.Close()
		// c2.Wait()
		// io.Copy(os.Stdout, &b2)
	}
}

func scaleCanaryDeployment(name, namespace string, replicas int) {
	logInfo("Scaling canary deployment to %v replicas...", replicas)
	foundation.RunCommandWithArgs("kubectl", []string{"scale", "deploy", fmt.Sprintf("%v-canary", name), "-n", namespace, fmt.Sprintf("--replicas=%v", replicas)})
}

func deleteResourcesForTypeSwitch(name, namespace string) {
	// clean up resources in case a switch from simple to canary releases or vice versa has been made
	logInfo("Deleting simple type deployment, configmap, secret, hpa and pdb...")
	foundation.RunCommandWithArgs("kubectl", []string{"delete", "deploy", name, "-n", namespace, "--ignore-not-found=true"})
	foundation.RunCommandWithArgs("kubectl", []string{"delete", "configmap", fmt.Sprintf("%v-configs", name), "-n", namespace, "--ignore-not-found=true"})
	foundation.RunCommandWithArgs("kubectl", []string{"delete", "secret", fmt.Sprintf("%v-secrets", name), "-n", namespace, "--ignore-not-found=true"})
	foundation.RunCommandWithArgs("kubectl", []string{"delete", "hpa", name, "-n", namespace, "--ignore-not-found=true"})
	foundation.RunCommandWithArgs("kubectl", []string{"delete", "pdb", name, "-n", namespace, "--ignore-not-found=true"})
}

func deleteConfigsForParamsChange(params Params, name, namespace string) {
	if len(params.Configs.Files) == 0 && len(params.Configs.InlineFiles) == 0 {
		logInfo("Deleting application configs if it exists, because no configs are specified...")
		foundation.RunCommandWithArgs("kubectl", []string{"delete", "configmap", fmt.Sprintf("%v-configs", name), "-n", namespace, "--ignore-not-found=true"})
	}
}

func deleteSecretsForParamsChange(params Params, name, namespace string) {
	if len(params.Secrets.Keys) == 0 {
		logInfo("Deleting application secrets if it exists, because no secrets are specified...")
		foundation.RunCommandWithArgs("kubectl", []string{"delete", "secret", fmt.Sprintf("%v-secrets", name), "-n", namespace, "--ignore-not-found=true"})
	}
}

func deleteServiceAccountSecretForParamsChange(params Params, name, namespace string) {
	if !params.UseGoogleCloudCredentials {
		logInfo("Deleting service account secret if it exists, because no use of service account is specified...")
		foundation.RunCommandWithArgs("kubectl", []string{"delete", "secret", fmt.Sprintf("%v-gcp-service-account", name), "-n", namespace, "--ignore-not-found=true"})
	}
}

func deleteIngressForVisibilityChange(templateData TemplateData, name, namespace string) {
	if !templateData.UseNginxIngress && !templateData.UseGCEIngress {
		// public uses service of type loadbalancer and doesn't need ingress
		logInfo("Deleting ingress if it exists, which is used for visibility private, iap or public-whitelist...")
		foundation.RunCommandWithArgs("kubectl", []string{"delete", "ingress", name, "-n", namespace, "--ignore-not-found=true"})
	}
}

func deleteBackendConfigAndIAPOauthSecret(templateData TemplateData, name, namespace string) {
	if !templateData.UseBackendConfigAnnotationOnService {
		logInfo("Deleting iap oauth secret if it exists, because visibility is not set to iap...")
		foundation.RunCommandWithArgs("kubectl", []string{"delete", "secret", fmt.Sprintf("%v--iap-oauth-credentials", name), "-n", namespace, "--ignore-not-found=true"})
		logInfo("Deleting iap backend config if it exists, because visibility is not set to iap...")
		foundation.RunCommandWithArgs("kubectl", []string{"delete", "backendconfig", name, "-n", namespace, "--ignore-not-found=true"})
	}
}

func removePoddisruptionBudgetIfRequired(params Params, name, namespace string) {
	if params.Kind == "deployment" && (params.Action == "deploy-simple" || params.Action == "deploy-stable") {
		// if there's a pdb that doesn't use maxUnavailable: 1 remove it so a new one can be created with correct settings
		deletePoddisruptionBudget := false
		maxUnavailable, err := getCommandOutput("kubectl", []string{"get", "pdb", name, "-n", namespace, "-o=jsonpath={.spec.maxUnavailable}"})
		if err == nil {
			maxUnavailableInt, err := strconv.Atoi(maxUnavailable)
			if err == nil {
				if maxUnavailableInt != 1 {
					logInfo("MaxUnavailable from pdb %v is %v instead of 1", name, maxUnavailableInt)
					deletePoddisruptionBudget = true
				}
			} else {
				logInfo("Failed reading maxUnavailable from pdb %v: %v", name, err)
				deletePoddisruptionBudget = true
			}
		} else {
			logInfo("Failed retrieving pdb %v: %v", name, err)
			deletePoddisruptionBudget = true
		}

		if deletePoddisruptionBudget {
			foundation.RunCommandWithArgs("kubectl", []string{"delete", "pdb", name, "-n", namespace, "--ignore-not-found=true"})
		} else {
			logInfo("Poddisruptionbudet %v is fine, not removing it", name)
		}
	}
}

func removeIngressIfRequired(params Params, templateData TemplateData, name, namespace string) {
	if params.Kind == "deployment" && (params.Action == "deploy-simple" || params.Action == "deploy-canary" || params.Action == "deploy-stable") {
		if templateData.UseNginxIngress {
			// check if ingress exists and has kubernetes.io/ingress.class: gce, then delete it because of https://github.com/kubernetes/ingress-gce/issues/481
			ingressClass, err := getCommandOutput("kubectl", []string{"get", "ing", name, "-n", namespace, "-o=go-template={{index .metadata.annotations \"kubernetes.io/ingress.class\"}}"})
			if err == nil {
				if ingressClass == "gce" {
					// delete the ingress so all related load balancers, etc get deleted
					logInfo("Deleting ingress so the gce ingress controller removes the related load balancer...")
					foundation.RunCommandWithArgs("kubectl", []string{"delete", "ingress", name, "-n", namespace, "--ignore-not-found=true"})
				} else {
					logInfo("Ingress %v already has kubernetes.io/ingress.class: %v annotation, no need to delete the ingress", name, ingressClass)
				}
			} else {
				logInfo("Ingress %v or kubernetes.io/ingress.class annotation doesn't exist, no need to delete the ingress: %v", name, err)
			}
		} else if templateData.UseGCEIngress {
			// check if ingress exists and has kubernetes.io/ingress.class: gce, then delete it to ensure there's no nginx ingress annotations lingering around
			ingressClass, err := getCommandOutput("kubectl", []string{"get", "ing", name, "-n", namespace, "-o=go-template={{index .metadata.annotations \"kubernetes.io/ingress.class\"}}"})
			if err == nil {
				if ingressClass == "nginx" {
					// delete the ingress so all related nginx ingress config gets deleted
					logInfo("Deleting ingress so the nginx ingress controller removes related config...")
					foundation.RunCommandWithArgs("kubectl", []string{"delete", "ingress", name, "-n", namespace, "--ignore-not-found=true"})
				} else {
					logInfo("Ingress %v already has kubernetes.io/ingress.class: %v annotation, no need to delete the ingress", name, ingressClass)
				}
			} else {
				logInfo("Ingress %v or kubernetes.io/ingress.class annotation doesn't exist, no need to delete the ingress: %v", name, err)
			}
		}
	}
}

func deployGoogleEndpointsServiceIfRequired(params Params) {
	if params.Kind == "deployment" && params.Visibility == "esp" && (params.Action == "deploy-simple" || params.Action == "deploy-canary") {
		foundation.RunCommandWithArgs("gcloud", []string{"endpoints", "services", "deploy", params.EspOpenAPIYamlPath})
	}
}

func patchServiceIfRequired(params Params, templateData TemplateData, name, namespace string) {
	if params.Kind == "deployment" && templateData.ServiceType == "ClusterIP" {
		serviceType, err := getCommandOutput("kubectl", []string{"get", "service", name, "-n", namespace, "-o=jsonpath={.spec.type}"})
		if err != nil {
			logInfo("Failed retrieving service type: %v", err)
		}
		if err == nil && (serviceType == "NodePort" || serviceType == "LoadBalancer") {
			logInfo("Service is of type %v, patching it...", serviceType)

			// brute force patch the service
			err = foundation.RunCommandWithArgsExtended("kubectl", []string{"patch", "service", name, "-n", namespace, "--type", "json", "--patch", "[{\"op\": \"remove\", \"path\": \"/spec/loadBalancerSourceRanges\"},{\"op\": \"remove\", \"path\": \"/spec/externalTrafficPolicy\"}, {\"op\": \"remove\", \"path\": \"/spec/ports/0/nodePort\"}, {\"op\": \"remove\", \"path\": \"/spec/ports/1/nodePort\"}, {\"op\": \"replace\", \"path\": \"/spec/type\", \"value\": \"ClusterIP\"}]"})
			if err != nil {
				err = foundation.RunCommandWithArgsExtended("kubectl", []string{"patch", "service", name, "-n", namespace, "--type", "json", "--patch", "[{\"op\": \"remove\", \"path\": \"/spec/externalTrafficPolicy\"}, {\"op\": \"remove\", \"path\": \"/spec/ports/0/nodePort\"}, {\"op\": \"remove\", \"path\": \"/spec/ports/1/nodePort\"}, {\"op\": \"replace\", \"path\": \"/spec/type\", \"value\": \"ClusterIP\"}]"})
			}
			if err != nil {
				log.Fatal(fmt.Sprintf("Failed patching service to change from %v to ClusterIP: ", serviceType), err)
			}
		} else {
			logInfo("Service is of type %v, no need to patch it", serviceType)
		}
	}
}

func cleanupJobIfRequired(params Params, templateData TemplateData, name, namespace string) {
	if params.Kind == "job" {
		err := foundation.RunCommandWithArgsExtended("kubectl", []string{"delete", "job", name, "-n", namespace, "--ignore-not-found=true"})
		if err != nil {
			logInfo("Deleting job %v failed: %v", name, err)
		}
	}
	if params.Kind == "cronjob" {
		err := foundation.RunCommandWithArgsExtended("kubectl", []string{"delete", "cronjob", name, "-n", namespace, "--ignore-not-found=true"})
		if err != nil {
			logInfo("Deleting cronjob %v failed: %v", name, err)
		}
	}
}

func getExistingNumberOfReplicas(params Params) int {
	if params.Kind == "deployment" {
		deploymentName := ""
		if params.Action == "deploy-simple" {
			deploymentName = params.App + "-stable"
		} else if params.Action == "deploy-stable" {
			deploymentName = params.App
		}
		if deploymentName != "" {
			replicas, err := getCommandOutput("kubectl", []string{"get", "deploy", deploymentName, "-n", params.Namespace, "-o=jsonpath={.spec.replicas}"})
			if err != nil {
				logInfo("Failed retrieving replicas for %v: %v ignoring setting replicas since there's no switch for deployment type...", deploymentName, err)
				return -1
			}
			replicasInt, err := strconv.Atoi(replicas)
			if err != nil {
				logInfo("Failed converting replicas value %v for %v: %v ignoring setting replicas since there's no switch for deployment type...", replicas, deploymentName, err)
				return -1
			}
			logInfo("Retrieved number of replicas for %v is %v; using it to set correct number of replicas switching deployment type...", deploymentName, replicasInt)
			return replicasInt
		}
	}

	return -1
}

func patchDeploymentIfRequired(params Params, name, namespace string) {
	if params.Kind == "deployment" && params.Action == "deploy-simple" {
		selectorLabels, err := getCommandOutput("kubectl", []string{"get", "deploy", name, "-n", namespace, "-o=jsonpath={.spec.selector.matchLabels}"})
		if err != nil {
			logInfo("Failed retrieving deployment selector labels: %v", err)
		}
		if err == nil && selectorLabels != fmt.Sprintf("map[app:%v]", name) {
			logInfo("Deployment selector labels %v not correct, patching it...", selectorLabels)

			// patch the deployment
			err = foundation.RunCommandWithArgsExtended("kubectl", []string{"patch", "deploy", name, "-n", namespace, "--type", "json", "--patch", fmt.Sprintf("[{\"op\": \"replace\", \"path\": \"/spec/selector/matchLabels\", \"value\": {\"app\":\"%v\"}}]", name)})
			if err != nil {
				log.Fatal(fmt.Sprintf("Failed patching deployment to change selector labels from %v to app=%v: ", selectorLabels, name), err)
			}
		} else {
			logInfo("Deployment selector labels %v are correct, not patching", selectorLabels)
		}
	}
}

func removeEstafetteCloudflareAnnotations(templateData TemplateData, name, namespace string) {
	if !templateData.UseDNSAnnotationsOnService {
		// ingress is used and has the estafette.io/cloudflare annotations, so they should be removed from the service
		logInfo("Removing estafette.io/cloudflare annotations on the service if they exists, since they're now set on the ingress instead...")
		foundation.RunCommandWithArgs("kubectl", []string{"annotate", "svc", name, "-n", namespace, "estafette.io/cloudflare-dns-"})
		foundation.RunCommandWithArgs("kubectl", []string{"annotate", "svc", name, "-n", namespace, "estafette.io/cloudflare-proxy-"})
		foundation.RunCommandWithArgs("kubectl", []string{"annotate", "svc", name, "-n", namespace, "estafette.io/cloudflare-hostnames-"})
		foundation.RunCommandWithArgs("kubectl", []string{"annotate", "svc", name, "-n", namespace, "estafette.io/cloudflare-state-"})
	}
}

func removeBackendConfigAnnotation(templateData TemplateData, name, namespace string) {
	if !templateData.UseBackendConfigAnnotationOnService {
		// iap is not used, so the beta.cloud.google.com/backend-config annotations should be removed from the service
		logInfo("Removing beta.cloud.google.com/backend-config annotations on the service if they exists, since visibility is not set to iap...")
		foundation.RunCommandWithArgs("kubectl", []string{"annotate", "svc", name, "-n", namespace, "beta.cloud.google.com/backend-config-"})
	}
}

func deleteHorizontalPodAutoscaler(params Params, name, namespace string) {
	if params.Kind == "deployment" && (params.Autoscale.Enable == nil || !*params.Autoscale.Enable) && (params.Action == "deploy-simple" || params.Action == "deploy-stable") {
		logInfo("Deleting HorizontalPodAutoscaler %v, since autoscaling is disabled...", name)
		foundation.RunCommandWithArgs("kubectl", []string{"delete", "hpa", name, "-n", namespace, "--ignore-not-found=true"})
	}
}

func getCommandOutput(command string, args []string) (string, error) {
	logInfo("Getting output for command '%v %v'...", command, strings.Join(args, " "))
	output, err := exec.Command(command, args...).Output()

	return string(output), err
}

func logInfo(message string, args ...interface{}) {
	formattedMessage := fmt.Sprintf(message, args...)
	log.Printf("%v\n\n", formattedMessage)
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
