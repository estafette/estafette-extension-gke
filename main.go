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
	customProperties = kingpin.Flag("custom-properties", "All custom properties for the stage as a json object.").Envar("ESTAFETTE_EXTENSION_CUSTOM_PROPERTIES").String()

	// Name                string
	// Namespace           string
	// Labels              map[string]string
	// AppLabelSelector    string
	// Hosts               []string
	// HostsJoined         string
	// IngressPath         string
	// UseNginxIngress     bool
	// UseGCEIngress       bool
	// ServiceType         string
	// MinReplicas         int
	// MaxReplicas         int
	// TargetCPUPercentage int
	// PreferPreemptibles  bool
	// Container           ContainerData
)

func main() {

	// parse command line parameters
	kingpin.Parse()

	// log to stdout and hide timestamp
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	// log startup message
	log.Printf("Starting %v version %v...", app, version)

	// unmarshal custom properties to parameters
	var params Parameters
	err := json.Unmarshal([]byte(*customProperties), &params)
	if err != nil {
		log.Fatal("Custom properties can't unmarshal to parameters.", err)
	}

	// get some estafette envvars
	appLabel := os.Getenv("ESTAFETTE_LABEL_APP")
	//estafetteBuildVersion := os.Getenv("ESTAFETTE_BUILD_VERSION")

	// validate required values are set
	if params.Name == "" && appLabel == "" {
		log.Fatal("Application name is required; either define an app label or use appName property.")
	}
	if params.Namespace == "" {
		log.Fatal("Namespace is required; use namespace property.")
	}

	// set data with defaults or overrides
	if params.Name == "" && appLabel != "" {
		params.Name = appLabel
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

	templateStrings := []string{}
	for _, t := range templatesToMerge {
		filePath := fmt.Sprintf("/templates/%v", t)
		data, err := ioutil.ReadFile(filePath)
		if err != nil {
			log.Fatal(fmt.Sprintf("Failed reading file %v", filePath), err)
		}

		log.Printf("Template %v:\n\n", filePath)
		log.Println(string(data))
		log.Println("")

		templateStrings = append(templateStrings, string(data))
	}
	templateString := strings.Join(templateStrings, "\n---\n")
	log.Printf("Template before rendering:\n\n")
	log.Println(templateString)
	log.Println("")

	// parse templates
	tmpl, err := template.New("kubernetes.yaml").Parse(templateString)
	if err != nil {
		log.Fatal("Failed parsing templates", err)
	}

	data := TemplateData{
		Name:      params.Name,
		Namespace: params.Namespace,
	}

	// render templates
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
