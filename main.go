package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/alecthomas/kingpin"
)

var (
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
	credentialsJSON = kingpin.Flag("credentials", "GKE credentials configured at service level, passed in to this trusted extension.").Envar("ESTAFETTE_CREDENTIALS_KUBERNETES_ENGINE").Required().String()

	// optional flags
	appLabel      = kingpin.Flag("app-name", "App label, used as application name if not passed explicitly.").Envar("ESTAFETTE_LABEL_APP").String()
	buildVersion  = kingpin.Flag("build-version", "Version number, used if not passed explicitly.").Envar("ESTAFETTE_BUILD_VERSION").String()
	releaseName   = kingpin.Flag("release-name", "Name of the release section, which is used by convention to resolve the credentials.").Envar("ESTAFETTE_RELEASE_ACTION").String()
	releaseAction = kingpin.Flag("release-action", "Name of the release action, to control the type of release.").Envar("ESTAFETTE_RELEASE_NAME").String()
)

func main() {

	// parse command line parameters
	kingpin.Parse()

	// log to stdout and hide timestamp
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	// log startup message
	log.Printf("Starting %v version %v...", app, version)

	// put all estafette labels in map
	log.Printf("Getting all estafette labels from envvars...")
	estafetteLabels := map[string]string{}
	for _, e := range os.Environ() {
		kvPair := strings.SplitN(e, "=", 2)

		if len(kvPair) == 2 {
			envvarName := kvPair[0]
			envvarValue := kvPair[1]

			if strings.HasPrefix(envvarName, "ESTAFETTE_LABEL_") {
				// strip prefix and convert to lowercase
				key := strings.ToLower(strings.Replace(envvarName, "ESTAFETTE_LABEL_", "", 1))
				estafetteLabels[key] = envvarValue
			}
		}
	}

	log.Printf("Unmarshalling parameters / custom properties...")
	var params Params
	err := json.Unmarshal([]byte(*paramsJSON), &params)
	if err != nil {
		log.Fatal("Failed unmarshalling parameters: ", err)
	}

	log.Printf("Setting defaults for parameters that are not set in the manifest...")
	params.SetDefaults(*appLabel, *buildVersion, *releaseName, *releaseAction, estafetteLabels)

	log.Printf("Unmarshalling credentials...")
	var credentials []GKECredentials
	err = json.Unmarshal([]byte(*credentialsJSON), &credentials)
	if err != nil {
		log.Fatal("Failed unmarshalling credentials: ", err)
	}

	log.Printf("Checking if credential %v exists...", params.Credentials)
	credential := GetCredentialsByName(credentials, params.Credentials)
	if credential == nil {
		log.Fatalf("Credential with name %v does not exist.", params.Credentials)
	}

	log.Printf("Setting default namespace from credentials in case the parameter is not set in the manifest...")
	params.SetDefaultsFromCredentials(*credential)

	log.Printf("Validating required parameters...")
	valid, errors := params.ValidateRequiredProperties()
	if !valid {
		log.Fatal("Not all valid fields are set: ", errors)
	}

	// combine templates
	tmpl, err := buildTemplates(params)
	if err != nil {
		log.Fatal("Failed building templates: ", err)
	}

	// pre-render config files if they exist
	params.Configs.RenderedFileContent = renderConfig(params)

	// generate the data required for rendering the templates
	templateData := generateTemplateData(params)

	// render the template
	renderedTemplate, err := renderTemplate(tmpl, templateData)
	if err != nil {
		log.Fatal("Failed rendering templates: ", err)
	}

	if tmpl != nil {
		log.Printf("Storing rendered manifest on disk...\n")
		err = ioutil.WriteFile("/kubernetes.yaml", renderedTemplate.Bytes(), 0600)
		if err != nil {
			log.Fatal("Failed writing manifest: ", err)
		}
	}

	log.Printf("Retrieving service account email from credentials...\n")
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

	log.Printf("Storing gke credential %v on disk...\n", params.Credentials)
	err = ioutil.WriteFile("/key-file.json", []byte(credential.AdditionalProperties.ServiceAccountKeyfile), 0600)
	if err != nil {
		log.Fatal("Failed writing service account keyfile: ", err)
	}

	log.Printf("Authenticating to google cloud\n")
	runCommand("gcloud", []string{"auth", "activate-service-account", saClientEmail, "--key-file", "/key-file.json"})

	log.Printf("Setting gcloud account\n")
	runCommand("gcloud", []string{"config", "set", "account", saClientEmail})

	log.Printf("Setting gcloud project\n")
	runCommand("gcloud", []string{"config", "set", "project", credential.AdditionalProperties.Project})

	log.Printf("Getting gke credentials for cluster %v\n", credential.AdditionalProperties.Cluster)
	clustersGetCredentialsArsgs := []string{"container", "clusters", "get-credentials", credential.AdditionalProperties.Cluster}
	if credential.AdditionalProperties.Zone != "" {
		clustersGetCredentialsArsgs = append(clustersGetCredentialsArsgs, "--zone", credential.AdditionalProperties.Zone)
	} else if credential.AdditionalProperties.Region != "" {
		clustersGetCredentialsArsgs = append(clustersGetCredentialsArsgs, "--region", credential.AdditionalProperties.Region)
	} else {
		log.Fatal("Credentials have no zone or region; at least one of them has to be defined")
	}
	runCommand("gcloud", clustersGetCredentialsArsgs)

	kubectlApplyArgs := []string{"apply", "-f", "/kubernetes.yaml", "-n", templateData.Namespace}
	if tmpl != nil {
		// always perform a dryrun to ensure we're not ending up in a semi broken state where half of the templates is successfully applied and others not
		log.Printf("Performing a dryrun to test the validity of the manifests...\n")
		runCommand("kubectl", append(kubectlApplyArgs, "--dry-run"))
	}

	if !params.DryRun {
		if tmpl != nil {
			log.Printf("Applying the manifests for real...\n")
			runCommand("kubectl", kubectlApplyArgs)

			log.Printf("Waiting for the deployment to finish...\n")
			runCommand("kubectl", []string{"rollout", "status", "deployment", templateData.NameWithTrack, "-n", templateData.Namespace})
		}

		// clean up old stuff
		switch params.Action {
		case "deploy-canary":
			scaleCanaryDeployment(templateData.Name, templateData.Namespace, 1)
			deleteConfigsForParamsChange(params, fmt.Sprintf("%v-canary", templateData.Name), templateData.Namespace)
			deleteSecretsForParamsChange(params, fmt.Sprintf("%v-canary", templateData.Name), templateData.Namespace)
			break
		case "deploy-stable":
			scaleCanaryDeployment(templateData.Name, templateData.Namespace, 0)
			deleteResourcesForTypeSwitch(templateData.Name, templateData.Namespace)
			deleteConfigsForParamsChange(params, fmt.Sprintf("%v-stable", templateData.Name), templateData.Namespace)
			deleteSecretsForParamsChange(params, fmt.Sprintf("%v-stable", templateData.Name), templateData.Namespace)
			deleteIngressForVisibilityChange(params, templateData.Name, templateData.Namespace)
			break
		case "rollback-canary":
			scaleCanaryDeployment(templateData.Name, templateData.Namespace, 0)
			break
		case "deploy-simple":
			deleteResourcesForTypeSwitch(fmt.Sprintf("%v-canary", templateData.Name), templateData.Namespace)
			deleteResourcesForTypeSwitch(fmt.Sprintf("%v-stable", templateData.Name), templateData.Namespace)
			deleteConfigsForParamsChange(params, templateData.Name, templateData.Namespace)
			deleteSecretsForParamsChange(params, templateData.Name, templateData.Namespace)
			deleteIngressForVisibilityChange(params, templateData.Name, templateData.Namespace)
			break
		}

		log.Printf("Showing current secrets, configmaps, horizontalpodautoscalers, services, ingresses, deployments and pods for app=%v...\n", params.App)
		runCommand("kubectl", []string{"get", "secret,cm,hpa,pdb,svc,ing,deploy,po", "-l", fmt.Sprintf("app=%v", params.App), "-n", params.Namespace})
	}
}

func scaleCanaryDeployment(name, namespace string, replicas int) {
	log.Printf("Scaling canary deployment to %v replicas...\n", replicas)
	runCommand("kubectl", []string{"scale", "deploy", fmt.Sprintf("%v-canary", name), "-n", namespace, fmt.Sprintf("--replicas=%v", replicas)})
}

func deleteResourcesForTypeSwitch(name, namespace string) {
	// clean up resources in case a switch from simple to canary releases or vice versa has been made
	log.Printf("Deleting simple type deployment, configmap, secret, hpa and pdb...\n")
	runCommand("kubectl", []string{"delete", "deploy", name, "-n", namespace, "--ignore-not-found=true"})
	runCommand("kubectl", []string{"delete", "configmap", fmt.Sprintf("%v-configs", name), "-n", namespace, "--ignore-not-found=true"})
	runCommand("kubectl", []string{"delete", "secret", fmt.Sprintf("%v-secrets", name), "-n", namespace, "--ignore-not-found=true"})
	runCommand("kubectl", []string{"delete", "hpa", name, "-n", namespace, "--ignore-not-found=true"})
	runCommand("kubectl", []string{"delete", "pdb", name, "-n", namespace, "--ignore-not-found=true"})
}

func deleteConfigsForParamsChange(params Params, name, namespace string) {
	if len(params.Configs.Files) == 0 {
		log.Printf("Deleting application configs if it exists, because no configs are specified...\n")
		runCommand("kubectl", []string{"delete", "configmap", fmt.Sprintf("%v-configs", name), "-n", namespace, "--ignore-not-found=true"})
	}
}

func deleteSecretsForParamsChange(params Params, name, namespace string) {
	if len(params.Secrets.Keys) == 0 {
		log.Printf("Deleting application secrets if it exists, because no secrets are specified...\n")
		runCommand("kubectl", []string{"delete", "secret", fmt.Sprintf("%v-secrets", name), "-n", namespace, "--ignore-not-found=true"})
	}
}

func deleteIngressForVisibilityChange(params Params, name, namespace string) {
	if params.Visibility == "public" {
		// public uses service of type loadbalancer and doesn't need ingress
		log.Printf("Deleting ingress if it exists, which is used for visibility private or iap...\n")
		runCommand("kubectl", []string{"delete", "ingress", name, "-n", namespace, "--ignore-not-found=true"})
	}
}

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func runCommand(command string, args []string) {
	log.Printf("Running command '%v %v'...", command, strings.Join(args, " "))
	cmd := exec.Command(command, args...)
	cmd.Dir = "/estafette-work"
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	handleError(err)
}
