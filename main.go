package main

import (
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
	name      = kingpin.Flag("name", "Application name.").Envar("ESTAFETTE_EXTENSION_NAME").String()
	namespace = kingpin.Flag("namespace", "Application namespace.").Envar("ESTAFETTE_EXTENSION_NAMESPACE").String()

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

	// get some estafette envvars
	appLabel := os.Getenv("ESTAFETTE_LABEL_APP")
	estafetteBuildVersion := os.Getenv("ESTAFETTE_BUILD_VERSION")

	// validate required values are set
	if *name == "" && appLabel == "" {
		log.Fatal("Application name is required; either define an app label or use appName property.")
	}
	if *namespace == "" {
		log.Fatal("Namespace is required; use namespace property.")
	}

	// set data with defaults or overrides
	if *name == "" && appLabel != "" {
		*name = appLabel
	}

	// render templates

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
