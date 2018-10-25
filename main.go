package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"text/template"

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
	appLabel     = kingpin.Flag("app-name", "App label, used as application name if not passed explicitly.").Envar("ESTAFETTE_LABEL_APP").String()
	buildVersion = kingpin.Flag("build-version", "Version number, used if not passed explicitly.").Envar("ESTAFETTE_BUILD_VERSION").String()
	releaseName  = kingpin.Flag("release-name", "Name of the release section, which is used by convention to resolve the credentials.").Envar("ESTAFETTE_RELEASE_NAME").String()
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
				estafetteLabels[envvarName] = envvarValue
			}
		}
	}

	log.Printf("Unmarshalling parameters / custom properties...")
	var params Params
	err := json.Unmarshal([]byte(*paramsJSON), &params)
	if err != nil {
		log.Fatal("Failed unmarshalling parameters.", err)
	}

	log.Printf("Setting defaults for parameters that are not set in the manifest...")
	params.SetDefaults(*appLabel, *buildVersion, *releaseName, estafetteLabels)

	log.Printf("Unmarshalling credentials...")
	var credentials []GKECredentials
	err = json.Unmarshal([]byte(*credentialsJSON), &credentials)
	if err != nil {
		log.Fatal("Failed unmarshalling credentials.", err)
	}

	log.Printf("Checking if credential %v exists...", params.Credentials)
	credential := GetCredentialsByName(credentials, params.Credentials)
	if credential == nil {
		log.Fatalf("Credential with name %v does not exist.", params.Credentials)
	}

	log.Printf("Setting default namespace from credentials in case the parameter is not set in the manifest...")
	params.SetDefaultNamespace(credential.DefaultNamespace)

	log.Printf("Validating required parameters...")
	valid, errors := params.ValidateRequiredProperties()
	if !valid {
		log.Fatal("Not all valid fields are set.", errors)
	}

	// merge templates
	templatesToMerge := []string{
		"namespace.yaml",
		"service.yaml",
		"serviceaccount.yaml",
		"certificate-secret.yaml",
		"poddisruptionbudget.yaml",
		"horizontalpodautoscaler.yaml",
		"deployment.yaml",
	}

	log.Printf("Merging templates %v...", strings.Join(templatesToMerge, ", "))

	templateStrings := []string{}
	for _, t := range templatesToMerge {
		filePath := fmt.Sprintf("/templates/%v", t)
		data, err := ioutil.ReadFile(filePath)
		if err != nil {
			log.Fatal(fmt.Sprintf("Failed reading file %v", filePath), err)
		}

		// log.Printf("Template %v:\n\n", filePath)
		// log.Println(string(data))
		// log.Println("")

		templateStrings = append(templateStrings, string(data))
	}
	templateString := strings.Join(templateStrings, "\n---\n")
	log.Printf("Template before rendering:\n\n")
	log.Println(templateString)
	log.Println("")

	// parse templates
	log.Printf("Parsing merged templates...")
	tmpl, err := template.New("kubernetes.yaml").Parse(templateString)
	if err != nil {
		log.Fatal("Failed parsing templates", err)
	}

	templateData := generateTemplateData(params)

	// render templates
	log.Printf("Rendering merged templates...")
	var renderedTemplate bytes.Buffer
	err = tmpl.Execute(&renderedTemplate, templateData)

	log.Printf("Template after rendering:\n\n")
	log.Println(renderedTemplate.String())
	log.Println("")

	log.Printf("Storing rendered manifest on disk...\n")
	err = ioutil.WriteFile("/kubernetes.yaml", renderedTemplate.Bytes(), 0600)
	if err != nil {
		log.Fatal("Failed writing manifest", err)
	}

	log.Printf("Retrieving service account email from credentials...\n")
	var keyFileMap map[string]interface{}
	err = json.Unmarshal([]byte(credential.ServiceAccountKeyfile), &keyFileMap)
	if err != nil {
		log.Fatal("Failed unmarshalling service account keyfile", err)
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
	err = ioutil.WriteFile("/key-file.json", []byte(credential.ServiceAccountKeyfile), 0600)
	if err != nil {
		log.Fatal("Failed writing service account keyfile", err)
	}

	log.Printf("Authenticating to google cloud\n")
	runCommand("gcloud", []string{"auth", "activate-service-account", saClientEmail, "--key-file", "/key-file.json"})

	log.Printf("Setting gcloud account\n")
	runCommand("gcloud", []string{"config", "set", "account", saClientEmail})

	log.Printf("Setting gcloud project\n")
	runCommand("gcloud", []string{"config", "set", "project", credential.Project})

	log.Printf("Getting gke credentials for cluster %v\n", credential.Cluster)
	clustersGetCredentialsArsgs := []string{"container", "clusters", "get-credentials", credential.Cluster}
	if credential.Zone != "" {
		clustersGetCredentialsArsgs = append(clustersGetCredentialsArsgs, "--zone", credential.Zone)
	} else if credential.Region != "" {
		clustersGetCredentialsArsgs = append(clustersGetCredentialsArsgs, "--region", credential.Region)
	} else {
		log.Fatal("Credentials have no zone or region; at least one of them has to be defined")
	}
	runCommand("gcloud", clustersGetCredentialsArsgs)

	if params.DryRun {
		runCommand("kubectl", []string{"apply", "-f", "/kubernetes.yaml", "--dry-run", "-n", params.Namespace})
	} else {
		log.Fatal("Not implemented applying manifest yet")
	}

	// - kubectl rollout status deploy/estafette-ci-web -n estafette
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
