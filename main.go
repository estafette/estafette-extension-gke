package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"runtime"
	"strconv"
	"time"

	"github.com/alecthomas/kingpin"
	"github.com/estafette/estafette-extension-gke/api"
	"github.com/estafette/estafette-extension-gke/clients/credentials"
	"github.com/estafette/estafette-extension-gke/clients/gcp"
	"github.com/estafette/estafette-extension-gke/clients/parameters"
	"github.com/estafette/estafette-extension-gke/services/builder"
	"github.com/estafette/estafette-extension-gke/services/generator"
	foundation "github.com/estafette/estafette-foundation"
	"github.com/rs/zerolog/log"
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
	credentialsPath = kingpin.Flag("credentials-path", "Path to GKE credentials configured at service level, passed in to this trusted extension.").Default("/credentials/kubernetes_engine.json").String()

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
	paramsForTroubleshooting     = api.Params{}
)

func main() {

	// parse command line parameters
	kingpin.Parse()

	// init log format from envvar ESTAFETTE_LOG_FORMAT
	foundation.InitLoggingFromEnv(foundation.NewApplicationInfo(appgroup, app, version, branch, revision, buildDate))

	// create context to cancel commands on sigterm
	ctx := foundation.InitCancellationContext(context.Background())

	credentialsClient, err := credentials.NewClient(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed creating credentials.Client")
	}

	credential, err := credentialsClient.Init(ctx, *paramsJSON, *releaseName, *credentialsPath)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed initializing credentials")
	}

	parametersClient, err := parameters.NewClient(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed creating parameters.Client")
	}

	params, err := parametersClient.Init(ctx, *paramsYAML, credential, *gitSource, *gitOwner, *gitName, *appLabel, *buildVersion, *releaseName, *releaseAction, *releaseID)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed initializing parameters")
	}

	gcpClient, err := gcp.NewClient(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed creating gcp.Client")
	}

	builderService, err := builder.NewService(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed creating builder.Service")
	}

	generatorService, err := generator.NewService(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed creating generator.Service")
	}

	_, err = gcpClient.LoadGKEClusterKubeConfig(ctx, credential)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed creating kube config for gke cluster")
	}

	// combine templates
	tmpl, err := builderService.BuildTemplates(params, true)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed building templates")
	}

	tmplNoPDB, err := builderService.BuildTemplates(params, false)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed building templates without poddisruptionbudget")
	}

	// pre-render config files if they exist
	params.Configs.RenderedFileContent = builderService.RenderConfig(params)
	if params.Kind == api.KindConfigToFile {
		// write files to working directory
		for filename, data := range params.Configs.RenderedFileContent {
			ioutil.WriteFile(filename, []byte(data), 0600)
		}

		return
	}

	// checking number of replicas for existing deployment to make switching deployment type safe
	currentReplicas := params.Replicas
	if params.Kind == api.KindDeployment || params.Kind == api.KindHeadlessDeployment {
		currentReplicas = getExistingNumberOfReplicas(ctx, params)
	}

	// generate the data required for rendering the templates
	templateData := generatorService.GenerateTemplateData(params, currentReplicas, *gitSource, *gitOwner, *gitName, *gitBranch, *gitRevision, *releaseID, *triggeredBy)

	// render the template
	renderedTemplate, err := builderService.RenderTemplate(tmpl, templateData, true)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed rendering templates")
	}
	renderedNoPDBTemplate, err := builderService.RenderTemplate(tmplNoPDB, templateData, false)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed rendering templates without poddisruptionbudget")
	}

	if tmpl != nil {
		log.Info().Msg("Storing rendered manifest on disk...")
		err = ioutil.WriteFile("/kubernetes.yaml", renderedTemplate.Bytes(), 0600)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed writing manifest")
		}
	}

	if tmplNoPDB != nil {
		log.Info().Msg("Storing rendered manifest without poddisruptionbudget on disk...")
		err = ioutil.WriteFile("/kubernetes-no-pdb.yaml", renderedNoPDBTemplate.Bytes(), 0600)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed writing manifest without poddisruptionbudget")
		}
	}

	if tmpl != nil {
		// visibility public is deprecated, so fail if creating new public service
		failIfCreatingNewPublicService(ctx, params, templateData, templateData.Name, templateData.Namespace)

		// fix resources before server-side dry-run to avoid failure
		cleanupJobIfRequired(ctx, params, templateData, templateData.Name, templateData.Namespace)
		patchServiceIfRequired(ctx, params, templateData, templateData.Name, templateData.Namespace)
		patchDeploymentIfRequired(ctx, params, templateData.Name, templateData.Namespace)

		// always perform a dryrun to ensure we're not ending up in a semi broken state where half of the templates is successfully applied and others not
		// await https://github.com/kubernetes/kubernetes/issues/83562 to switch back to server-side dry-run and not fail for new namespaces
		log.Info().Msg("Performing a dryrun to test the validity of the manifests...")
		foundation.RunCommandWithArgs(ctx, "kubectl", []string{"apply", "-f", "/kubernetes-no-pdb.yaml", "-n", templateData.Namespace, "--dry-run=client"})

		log.Info().Msg("Performing a diff to show what's changed...")
		_ = foundation.RunCommandWithArgsExtended(ctx, "kubectl", []string{"diff", "-f", "/kubernetes-no-pdb.yaml", "-n", templateData.Namespace})
	}

	if !params.DryRun && params.Action != api.ActionDiffSimple && params.Action != api.ActionDiffCanary && params.Action != api.ActionDiffStable {

		// ensure that from now on any error runs the troubleshooting assistant
		assistTroubleshootingOnError = true
		paramsForTroubleshooting = params

		if tmpl != nil {
			deployGoogleEndpointsServiceIfRequired(ctx, gcpClient, params)
			removePoddisruptionBudgetIfRequired(ctx, params, templateData.NameWithTrack, templateData.Namespace)
			removeIngressIfRequired(ctx, params, templateData, templateData.Name, templateData.Namespace)

			log.Info().Msg("Applying the manifests for real...")
			foundation.RunCommandWithArgs(ctx, "kubectl", []string{"apply", "-f", "/kubernetes.yaml", "-n", templateData.Namespace})

			if params.Kind == api.KindDeployment || params.Kind == api.KindHeadlessDeployment {
				log.Info().Msg("Waiting for the deployment to finish...")
				err = foundation.RunCommandWithArgsExtended(ctx, "kubectl", []string{"rollout", "status", "deployment", templateData.NameWithTrack, "-n", templateData.Namespace})
			}
			if params.Kind == api.KindStatefulset {
				log.Info().Msg("Waiting for the statefulset to finish...")
				err = foundation.RunCommandWithArgsExtended(ctx, "kubectl", []string{"rollout", "status", "statefulset", templateData.Name, "-n", templateData.Namespace})
			}
		}

		if err != nil {
			assistTroubleshooting(ctx, templateData, err)
		}

		handleAtomicUpdate(ctx, builderService, params, templateData)

		// clean up old stuff
		switch params.Kind {
		case api.KindDeployment:
			switch params.Action {
			case api.ActionDeployCanary:
				scaleCanaryDeployment(ctx, templateData.Name, templateData.Namespace, 1)
				deleteConfigsForParamsChange(ctx, params, templateData.NameWithTrack, templateData.Namespace)
				deleteSecretsForParamsChange(ctx, params, templateData.NameWithTrack, templateData.Namespace)
				break
			case api.ActionDeployStable:
				scaleCanaryDeployment(ctx, templateData.Name, templateData.Namespace, 0)
				deleteResourcesForTypeSwitch(ctx, templateData.Name, templateData.Namespace)
				deleteConfigsForParamsChange(ctx, params, templateData.NameWithTrack, templateData.Namespace)
				deleteSecretsForParamsChange(ctx, params, templateData.NameWithTrack, templateData.Namespace)
				deleteServiceAccountSecretForParamsChange(ctx, params, templateData.GoogleCloudCredentialsAppName, templateData.Namespace)
				deleteIngressForVisibilityChange(ctx, templateData, templateData.Name, templateData.Namespace)
				removeEstafetteCloudflareAnnotations(ctx, templateData, templateData.Name, templateData.Namespace)
				removeBackendConfigAnnotation(ctx, templateData, templateData.Name, templateData.Namespace)
				removeNegAnnotation(ctx, templateData, templateData.Name, templateData.Namespace)
				deleteBackendConfigAndIAPOauthSecret(ctx, templateData, templateData.Name, templateData.Namespace)
				deleteHorizontalPodAutoscaler(ctx, params, templateData.NameWithTrack, templateData.Namespace)
				break
			case api.ActionRollbackCanary:
				scaleCanaryDeployment(ctx, templateData.Name, templateData.Namespace, 0)
				break
			case api.ActionRestartCanary:
				restartDeployment(ctx, fmt.Sprintf("%v-canary", templateData.Name), templateData.Namespace)
				break
			case api.ActionRestartStable:
				restartDeployment(ctx, fmt.Sprintf("%v-stable", templateData.Name), templateData.Namespace)
				break
			case api.ActionRestartSimple:
				restartDeployment(ctx, templateData.Name, templateData.Namespace)
				break
			case api.ActionDeploySimple:
				deleteResourcesForTypeSwitch(ctx, fmt.Sprintf("%v-canary", templateData.Name), templateData.Namespace)
				deleteResourcesForTypeSwitch(ctx, fmt.Sprintf("%v-stable", templateData.Name), templateData.Namespace)
				deleteConfigsForParamsChange(ctx, params, templateData.Name, templateData.Namespace)
				deleteSecretsForParamsChange(ctx, params, templateData.Name, templateData.Namespace)
				deleteServiceAccountSecretForParamsChange(ctx, params, templateData.GoogleCloudCredentialsAppName, templateData.Namespace)
				deleteIngressForVisibilityChange(ctx, templateData, templateData.Name, templateData.Namespace)
				removeEstafetteCloudflareAnnotations(ctx, templateData, templateData.Name, templateData.Namespace)
				removeBackendConfigAnnotation(ctx, templateData, templateData.Name, templateData.Namespace)
				removeNegAnnotation(ctx, templateData, templateData.Name, templateData.Namespace)
				deleteBackendConfigAndIAPOauthSecret(ctx, templateData, templateData.Name, templateData.Namespace)
				deleteHorizontalPodAutoscaler(ctx, params, templateData.Name, templateData.Namespace)
				break
			}
			break

		case api.KindHeadlessDeployment:
			switch params.Action {
			case api.ActionDeployCanary:
				scaleCanaryDeployment(ctx, templateData.Name, templateData.Namespace, 1)
				deleteConfigsForParamsChange(ctx, params, templateData.NameWithTrack, templateData.Namespace)
				deleteSecretsForParamsChange(ctx, params, templateData.NameWithTrack, templateData.Namespace)
				break
			case api.ActionDeployStable:
				scaleCanaryDeployment(ctx, templateData.Name, templateData.Namespace, 0)
				deleteResourcesForTypeSwitch(ctx, templateData.Name, templateData.Namespace)
				deleteConfigsForParamsChange(ctx, params, templateData.NameWithTrack, templateData.Namespace)
				deleteSecretsForParamsChange(ctx, params, templateData.NameWithTrack, templateData.Namespace)
				deleteServiceAccountSecretForParamsChange(ctx, params, templateData.GoogleCloudCredentialsAppName, templateData.Namespace)
				deleteHorizontalPodAutoscaler(ctx, params, templateData.NameWithTrack, templateData.Namespace)
				break
			case api.ActionRollbackCanary:
				scaleCanaryDeployment(ctx, templateData.Name, templateData.Namespace, 0)
				break
			case api.ActionRestartCanary:
				restartDeployment(ctx, fmt.Sprintf("%v-canary", templateData.Name), templateData.Namespace)
				break
			case api.ActionRestartStable:
				restartDeployment(ctx, fmt.Sprintf("%v-stable", templateData.Name), templateData.Namespace)
				break
			case api.ActionRestartSimple:
				restartDeployment(ctx, templateData.Name, templateData.Namespace)
				break
			case api.ActionDeploySimple:
				deleteResourcesForTypeSwitch(ctx, fmt.Sprintf("%v-canary", templateData.Name), templateData.Namespace)
				deleteResourcesForTypeSwitch(ctx, fmt.Sprintf("%v-stable", templateData.Name), templateData.Namespace)
				deleteConfigsForParamsChange(ctx, params, templateData.Name, templateData.Namespace)
				deleteSecretsForParamsChange(ctx, params, templateData.Name, templateData.Namespace)
				deleteServiceAccountSecretForParamsChange(ctx, params, templateData.GoogleCloudCredentialsAppName, templateData.Namespace)
				deleteHorizontalPodAutoscaler(ctx, params, templateData.Name, templateData.Namespace)
				break
			}
			break
		case api.KindStatefulset:
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

func assistTroubleshooting(ctx context.Context, templateData api.TemplateData, err error) {
	if assistTroubleshootingOnError {
		log.Info().Msgf("Showing current ingresses, services, configmaps, secrets, deployments, jobs, cronjobs, poddisruptionbudgets, horizontalpodautoscalers, pods, endpoints for app=%v...", paramsForTroubleshooting.App)
		err = foundation.RunCommandWithArgsExtended(ctx, "kubectl", []string{"get", "ing,svc,cm,secret,deploy,job,cronjob,sts,pdb,hpa,po,ep", "-l", fmt.Sprintf("app=%v", paramsForTroubleshooting.App), "-n", paramsForTroubleshooting.Namespace})

		if err != nil {
			log.Info().Msg("Rollout failed, trying to show logs...")
			if *releaseID != "" {
				_ = foundation.RunCommandWithArgsExtended(ctx, "kubectl", []string{"logs", "-l", fmt.Sprintf("app=%v,estafette.io/release-id=%v", templateData.AppLabelSelector, api.SanitizeLabel(*releaseID)), "-n", templateData.Namespace, "--all-containers"})
			} else if *buildVersion != "" {
				_ = foundation.RunCommandWithArgsExtended(ctx, "kubectl", []string{"logs", "-l", fmt.Sprintf("app=%v,version=%v", templateData.AppLabelSelector, api.SanitizeLabel(*buildVersion)), "-n", templateData.Namespace, "--all-containers"})
			}
		} else if paramsForTroubleshooting.Action == api.ActionDeployCanary {
			log.Info().Msg("Showing logs for canary deployment...")
			_ = foundation.RunCommandWithArgsExtended(ctx, "kubectl", []string{"logs", "-l", fmt.Sprintf("app=%v,track=canary", paramsForTroubleshooting.App), "-n", paramsForTroubleshooting.Namespace, "-c", paramsForTroubleshooting.App, "--tail", "50"})
		}

		foundation.HandleError(err)
	}
}

func scaleCanaryDeployment(ctx context.Context, name, namespace string, replicas int) {
	log.Info().Msgf("Scaling canary deployment to %v replicas...", replicas)
	foundation.RunCommandWithArgs(ctx, "kubectl", []string{"scale", "deploy", fmt.Sprintf("%v-canary", name), "-n", namespace, fmt.Sprintf("--replicas=%v", replicas)})
}

func restartDeployment(ctx context.Context, name, namespace string) {
	log.Info().Msgf("Restarting deployment rollout...")
	foundation.RunCommandWithArgs(ctx, "kubectl", []string{"rollout", "restart", "deployment", name, "-n", namespace})
	foundation.RunCommandWithArgs(ctx, "kubectl", []string{"rollout", "status", "deployment", name, "-n", namespace})
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

func deleteConfigsForParamsChange(ctx context.Context, params api.Params, name, namespace string) {
	if len(params.Configs.Files) == 0 && len(params.Configs.InlineFiles) == 0 {
		log.Info().Msg("Deleting application configs if it exists, because no configs are specified...")
		foundation.RunCommandWithArgs(ctx, "kubectl", []string{"delete", "configmap", fmt.Sprintf("%v-configs", name), "-n", namespace, "--ignore-not-found=true"})
	}
}

func deleteSecretsForParamsChange(ctx context.Context, params api.Params, name, namespace string) {
	if len(params.Secrets.Keys) == 0 {
		log.Info().Msg("Deleting application secrets if it exists, because no secrets are specified...")
		foundation.RunCommandWithArgs(ctx, "kubectl", []string{"delete", "secret", fmt.Sprintf("%v-secrets", name), "-n", namespace, "--ignore-not-found=true"})
	}
}

func deleteServiceAccountSecretForParamsChange(ctx context.Context, params api.Params, name, namespace string) {
	if !params.UseGoogleCloudCredentials && params.LegacyGoogleCloudServiceAccountKeyFile == "" {
		log.Info().Msg("Deleting service account secret if it exists, because no use of service account is specified...")
		foundation.RunCommandWithArgs(ctx, "kubectl", []string{"delete", "secret", fmt.Sprintf("%v-gcp-service-account", name), "-n", namespace, "--ignore-not-found=true"})
	}
}

func deleteIngressForVisibilityChange(ctx context.Context, templateData api.TemplateData, name, namespace string) {
	if !templateData.UseNginxIngress && !templateData.UseGCEIngress {
		// public uses service of type loadbalancer and doesn't need ingress
		log.Info().Msg("Deleting ingress if it exists, which is used for visibility private, iap or public-whitelist...")
		foundation.RunCommandWithArgs(ctx, "kubectl", []string{"delete", "ingress", name, "-n", namespace, "--ignore-not-found=true"})
	}
}

func deleteBackendConfigAndIAPOauthSecret(ctx context.Context, templateData api.TemplateData, name, namespace string) {
	if !templateData.UseBackendConfigAnnotationOnService {
		log.Info().Msg("Deleting iap oauth secret if it exists, because visibility is not set to iap...")
		foundation.RunCommandWithArgs(ctx, "kubectl", []string{"delete", "secret", fmt.Sprintf("%v--iap-oauth-credentials", name), "-n", namespace, "--ignore-not-found=true"})
		log.Info().Msg("Deleting iap backend config if it exists, because visibility is not set to iap...")
		foundation.RunCommandWithArgs(ctx, "kubectl", []string{"delete", "backendconfig", name, "-n", namespace, "--ignore-not-found=true"})
	}
}

func removePoddisruptionBudgetIfRequired(ctx context.Context, params api.Params, name, namespace string) {
	if (params.Kind == api.KindDeployment || params.Kind == api.KindHeadlessDeployment) && (params.Action == api.ActionDeploySimple || params.Action == api.ActionDeployStable) {
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

func removeIngressIfRequired(ctx context.Context, params api.Params, templateData api.TemplateData, name, namespace string) {
	if params.Kind == api.KindDeployment && (params.Action == api.ActionDeploySimple || params.Action == api.ActionDeployCanary || params.Action == api.ActionDeployStable) {
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

func deployGoogleEndpointsServiceIfRequired(ctx context.Context, gcpClient gcp.Client, params api.Params) {
	if params.Kind == api.KindDeployment && params.Visibility == api.VisibilityESP && (params.Action == api.ActionDeploySimple || params.Action == api.ActionDeployCanary) {
		err := gcpClient.DeployGoogleCloudEndpoints(ctx, params)
		if err != nil {
			log.Fatal().Err(err).Msgf("Failed deploying endpoints service in project %v", params.EspEndpointsProjectID)
		}
	}
}

func failIfCreatingNewPublicService(ctx context.Context, params api.Params, templateData api.TemplateData, name, namespace string) {
	if params.Kind == api.KindDeployment && params.Visibility == api.VisibilityPublic {
		serviceType, err := foundation.GetCommandWithArgsOutput(ctx, "kubectl", []string{"get", "service", name, "-n", namespace, "-o=jsonpath={.spec.type}"})
		// fail if creating new public service or updating to public
		if err != nil {
			log.Fatal().Err(err).Msgf("Creating new public service is no longer supported, please use visibility esp or apigee.")
		} else if serviceType != "LoadBalancer" {
			log.Fatal().Msgf("Changing service visibility to public is no longer supported, please use visibility esp or apigee.")
		}
	}
}

func patchServiceIfRequired(ctx context.Context, params api.Params, templateData api.TemplateData, name, namespace string) {
	if params.Kind == api.KindDeployment && templateData.ServiceType == "ClusterIP" {
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

func cleanupJobIfRequired(ctx context.Context, params api.Params, templateData api.TemplateData, name, namespace string) {
	if params.Kind == api.KindJob {
		err := foundation.RunCommandWithArgsExtended(ctx, "kubectl", []string{"delete", "job", name, "-n", namespace, "--ignore-not-found=true"})
		if err != nil {
			log.Info().Msgf("Deleting job %v failed: %v", name, err)
		}
	}
	if params.Kind == api.KindCronJob {
		err := foundation.RunCommandWithArgsExtended(ctx, "kubectl", []string{"delete", "cronjob", name, "-n", namespace, "--ignore-not-found=true"})
		if err != nil {
			log.Info().Msgf("Deleting cronjob %v failed: %v", name, err)
		}
	}
}

func getExistingNumberOfReplicas(ctx context.Context, params api.Params) int {
	if params.Kind == api.KindDeployment || params.Kind == api.KindHeadlessDeployment {
		if params.StrategyType == api.StrategyTypeAtomicUpdate {
			replicas, err := foundation.GetCommandWithArgsOutput(ctx, "kubectl", []string{"get", "deploy", "-l", fmt.Sprintf("app in (%v),estafette.io/atomic-id,estafette.io/atomic-id notin (%v)", api.SanitizeLabel(params.App), params.AtomicID), "-n", params.Namespace, "--sort-by=.metadata.creationTimestamp", "-o=jsonpath={.items[-1:].spec.replicas}"})
			if err != nil {
				log.Info().Err(err).Msg("Failed retrieving replicas for previous atomic deployments. Ignoring setting replicas since there's no switch for deployment type...")
				return -1
			}
			replicasInt, err := strconv.Atoi(replicas)
			if err != nil {
				log.Info().Err(err).Msgf("Failed converting replicas value %v for previous atomic deployments. Ignoring setting replicas since there's no switch for deployment type...", replicas)
				return -1
			}
			log.Info().Msgf("Retrieved number of replicas for previous atomic deployments is %v; using it to set correct number of replicas switching deployment type...", replicasInt)
			return replicasInt
		}

		deploymentName := ""
		if params.Action == api.ActionDeploySimple || params.Action == api.ActionDiffSimple {
			deploymentName = params.App + "-stable"
		} else if params.Action == api.ActionDeployStable || params.Action == api.ActionDiffStable {
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

func patchDeploymentIfRequired(ctx context.Context, params api.Params, name, namespace string) {
	if (params.Kind == api.KindDeployment || params.Kind == api.KindHeadlessDeployment) && params.Action == api.ActionDeploySimple {
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

func removeEstafetteCloudflareAnnotations(ctx context.Context, templateData api.TemplateData, name, namespace string) {
	if !templateData.UseDNSAnnotationsOnService {
		// ingress is used and has the estafette.io/cloudflare annotations, so they should be removed from the service
		log.Info().Msg("Removing estafette.io/cloudflare annotations on the service if they exists, since they're now set on the ingress instead...")
		foundation.RunCommandWithArgs(ctx, "kubectl", []string{"annotate", "svc", name, "-n", namespace, "estafette.io/cloudflare-dns-"})
		foundation.RunCommandWithArgs(ctx, "kubectl", []string{"annotate", "svc", name, "-n", namespace, "estafette.io/cloudflare-proxy-"})
		foundation.RunCommandWithArgs(ctx, "kubectl", []string{"annotate", "svc", name, "-n", namespace, "estafette.io/cloudflare-hostnames-"})
		foundation.RunCommandWithArgs(ctx, "kubectl", []string{"annotate", "svc", name, "-n", namespace, "estafette.io/cloudflare-state-"})
	}
}

func removeBackendConfigAnnotation(ctx context.Context, templateData api.TemplateData, name, namespace string) {
	if !templateData.UseBackendConfigAnnotationOnService {
		// iap is not used, so the beta.cloud.google.com/backend-config annotations should be removed from the service
		log.Info().Msg("Removing beta.cloud.google.com/backend-config annotations on the service if they exists, since visibility is not set to iap...")
		foundation.RunCommandWithArgs(ctx, "kubectl", []string{"annotate", "svc", name, "-n", namespace, "beta.cloud.google.com/backend-config-"})
	}
}

func removeNegAnnotation(ctx context.Context, templateData api.TemplateData, name, namespace string) {
	if !templateData.UseNegAnnotationOnService {
		// cloud native load balancing is not used, so the beta.cloud.google.com/backend-config annotations should be removed from the service
		log.Info().Msg("Removing cloud.google.com/neg annotations on the service if they exists, since visibility is not set to iap or containerNativeLoadBalancing is set to fals...")
		foundation.RunCommandWithArgs(ctx, "kubectl", []string{"annotate", "svc", name, "-n", namespace, "cloud.google.com/neg-"})
	}
}

func deleteHorizontalPodAutoscaler(ctx context.Context, params api.Params, name, namespace string) {
	if (params.Kind == api.KindDeployment || params.Kind == api.KindHeadlessDeployment) && (params.Autoscale.Enabled == nil || !*params.Autoscale.Enabled) && (params.Action == api.ActionDeploySimple || params.Action == api.ActionDeployStable) {
		log.Info().Msgf("Deleting HorizontalPodAutoscaler %v, since autoscaling is disabled...", name)
		foundation.RunCommandWithArgs(ctx, "kubectl", []string{"delete", "hpa", name, "-n", namespace, "--ignore-not-found=true"})
	}
}

func handleAtomicUpdate(ctx context.Context, builderService builder.Service, params api.Params, templateData api.TemplateData) {
	if params.StrategyType != api.StrategyTypeAtomicUpdate {
		return
	}

	// update service in order to point to new deployment
	log.Info().Msgf("Updating service selector to use the latest atomic id...")
	atomicServiceTmpl, err := builderService.GetAtomicUpdateServiceTemplate()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed building service template")
	}

	renderedTemplate, err := builderService.RenderTemplate(atomicServiceTmpl, templateData, true)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed rendering templates")
	}

	log.Info().Msg("Storing rendered service manifest on disk...")
	err = ioutil.WriteFile("/service.yaml", renderedTemplate.Bytes(), 0600)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed writing manifest")
	}

	log.Info().Msg("Applying the service manifest...")
	foundation.RunCommandWithArgs(ctx, "kubectl", []string{"apply", "-f", "/service.yaml", "-n", templateData.Namespace})

	// wait a bit to drain traffic to old deployment
	sleepTime := 30
	log.Info().Msgf("Waiting for %v seconds to drain traffic to previous deployment(s)...", sleepTime)
	time.Sleep(time.Duration(sleepTime) * time.Second)

	// clean up old deployments, configmaps, secrets, hpa, pdb
	log.Info().Msg("Cleaning up previous deployments, configmaps, secrets, hpas and pdbs...")
	foundation.RunCommandWithArgs(ctx, "kubectl", []string{"delete", "deploy,hpa,pdb", "-l", fmt.Sprintf("app in (%v),estafette.io/atomic-id,estafette.io/atomic-id notin (%v)", api.SanitizeLabel(params.App), params.AtomicID), "-n", templateData.Namespace, "--ignore-not-found=true"})
	foundation.RunCommandWithArgs(ctx, "kubectl", []string{"delete", "configmap,secret", "-l", fmt.Sprintf("app in (%v),type in (application),estafette.io/atomic-id,estafette.io/atomic-id notin (%v)", api.SanitizeLabel(params.App), params.AtomicID), "-n", templateData.Namespace, "--ignore-not-found=true"})
	if templateData.IncludeTrackLabel {
		foundation.RunCommandWithArgs(ctx, "kubectl", []string{"delete", "deploy,hpa,pdb", "-l", fmt.Sprintf("app in (%v),!estafette.io/atomic-id,track in (%v)", api.SanitizeLabel(params.App), templateData.TrackLabel), "-n", templateData.Namespace, "--ignore-not-found=true"})
		foundation.RunCommandWithArgs(ctx, "kubectl", []string{"delete", "configmap,secret", "-l", fmt.Sprintf("app in (%v),type in (application),!estafette.io/atomic-id,track in (%v)", api.SanitizeLabel(params.App), templateData.TrackLabel), "-n", templateData.Namespace, "--ignore-not-found=true"})
	} else {
		foundation.RunCommandWithArgs(ctx, "kubectl", []string{"delete", "deploy,hpa,pdb", "-l", fmt.Sprintf("app in (%v),!estafette.io/atomic-id,!track", api.SanitizeLabel(params.App)), "-n", templateData.Namespace, "--ignore-not-found=true"})
		foundation.RunCommandWithArgs(ctx, "kubectl", []string{"delete", "configmap,secret", "-l", fmt.Sprintf("app in (%v),type in (application),!estafette.io/atomic-id,!track", api.SanitizeLabel(params.App)), "-n", templateData.Namespace, "--ignore-not-found=true"})
	}
}
