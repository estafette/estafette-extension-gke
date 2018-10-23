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

	log.Printf("Unmarshalling parameters / custom properties...")
	var params Params
	err := json.Unmarshal([]byte(*paramsJSON), &params)
	if err != nil {
		log.Fatal("Failed unmarshalling parameters.", err)
	}

	log.Printf("Setting defaults for parameters that are not set in the manifest...")
	params.SetDefaults(*appLabel, *buildVersion, *releaseName)

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
		// "serviceaccount.yaml",
		// "certificate-secret.yaml",
		// "poddisruptionbudget.yaml",
		// "horizontalpodautoscaler.yaml",
		// "deployment.yaml",
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

	data := TemplateData{
		Name:      params.App,
		Namespace: params.Namespace,
	}

	// render templates
	log.Printf("Rendering merged templates...")
	var renderedTemplate bytes.Buffer
	err = tmpl.Execute(&renderedTemplate, data)

	log.Printf("Template after rendering:\n\n")
	log.Println(renderedTemplate.String())
	log.Println("")

	//log.Fatal("Extension is not finished yet")

	// templates/namespace.yaml
	// templates/service.yaml
	// templates/ingress.yaml
	// templates/serviceaccount.yaml
	// templates/certificate-secret.yaml
	// templates/configmap.yaml
	// templates/poddisruptionbudget.yaml
	// templates/horizontalpodautoscaler.yaml
	// templates/deployment.yaml

	// tmpl, err := template.ParseFiles("layout.html")
	// tmpl.Execute(w, data)

	// - echo "${GCLOUD_KEY_FILE}" | base64 -d > key-file.json
	// - gcloud auth activate-service-account ${GCLOUD_SA_NAME} --key-file ./key-file.json
	// - gcloud config set account ${GCLOUD_SA_NAME}
	// - gcloud config set project ${GCLOUD_PROJECT_NAME}
	// - gcloud container clusters get-credentials ${GCLOUD_GKE_CLUSTER_NAME} --zone ${GCLOUD_GKE_ZONE}
	// - cat kubernetes.yaml | envsubst | kubectl apply -f -
	// - kubectl rollout status deploy/estafette-ci-web -n estafette

	// log.Printf("Tagging container image %v\n", targetContainerPath)
	// tagArgs := []string{
	// 	"tag",
	// 	sourceContainerPath,
	// 	targetContainerPath,
	// }
	// runCommand("kubectl", tagArgs)

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
