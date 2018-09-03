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
	action = kingpin.Flag("action", "Any of the following actions: build, push, tag.").Envar("ESTAFETTE_EXTENSION_ACTION").String()
)

func main() {

	// parse command line parameters
	kingpin.Parse()

	// log to stdout and hide timestamp
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	// log startup message
	log.Printf("Starting %v version %v...", app, version)

	// todo set data and render templates
	// .Name
	// .Namespace
	// .Labels
	// .AppLabelSelector
	// .Hosts
	// .ServiceType
	// .MinPods
	// .MaxPods
	// .TargetCPU

	// // set defaults
	// appLabel := os.Getenv("ESTAFETTE_LABEL_APP")
	// if *container == "" && appLabel != "" {
	// 	*container = appLabel
	// }

	// // get private container registries credentials
	// credentialsJSON := os.Getenv("ESTAFETTE_CI_REPOSITORY_CREDENTIALS_JSON")
	// var credentials []*contracts.ContainerRepositoryCredentialConfig
	// if credentialsJSON != "" {
	// 	json.Unmarshal([]byte(credentialsJSON), &credentials)
	// }

	// // validate inputs
	// validateRepositories(*repositories)

	// // split into arrays and set other variables
	// var repositoriesSlice []string
	// if *repositories != "" {
	// 	repositoriesSlice = strings.Split(*repositories, ",")
	// }
	// var tagsSlice []string
	// if *tags != "" {
	// 	tagsSlice = strings.Split(*tags, ",")
	// }
	// var copySlice []string
	// if *copy != "" {
	// 	copySlice = strings.Split(*copy, ",")
	// }
	// // var argsSlice []string
	// // if *args != "" {
	// // 	argsSlice = strings.Split(*args, ",")
	// // }
	// estafetteBuildVersion := os.Getenv("ESTAFETTE_BUILD_VERSION")

	// switch *action {
	// case "build":

	// 	// image: extensions/docker:stable
	// 	// action: build
	// 	// container: docker
	// 	// repositories:
	// 	// - extensions
	// 	// path: .
	// 	// copy:
	// 	// - Dockerfile
	// 	// - /etc/ssl/certs/ca-certificates.crt

	// 	// make build dir if it doesn't exist
	// 	log.Printf("Ensuring build directory %v exists\n", *path)
	// 	runCommand("mkdir", []string{"-p", *path})

	// 	// copy files/dirs from copySlice to build path
	// 	for _, c := range copySlice {
	// 		log.Printf("Copying %v to %v\n", c, *path)
	// 		runCommand("cp", []string{"-r", c, *path})
	// 	}

	// 	// build docker image
	// 	log.Printf("Building docker image %v/%v:%v...\n", repositoriesSlice[0], *container, estafetteBuildVersion)
	// 	args := []string{
	// 		"build",
	// 	}
	// 	for _, r := range repositoriesSlice {
	// 		args = append(args, "--tag")
	// 		args = append(args, fmt.Sprintf("%v/%v:%v", r, *container, estafetteBuildVersion))
	// 		for _, t := range tagsSlice {
	// 			args = append(args, "--tag")
	// 			args = append(args, fmt.Sprintf("%v/%v:%v", r, *container, t))
	// 		}
	// 	}

	// 	args = append(args, "--file")
	// 	args = append(args, fmt.Sprintf("%v/%v", *path, *dockerfile))
	// 	args = append(args, *path)
	// 	runCommand("docker", args)

	// case "push":

	// 	// image: extensions/docker:stable
	// 	// action: push
	// 	// container: docker
	// 	// repositories:
	// 	// - extensions
	// 	// tags:
	// 	// - dev

	// 	sourceContainerPath := fmt.Sprintf("%v/%v:%v", repositoriesSlice[0], *container, estafetteBuildVersion)

	// 	// push each repository + tag combination
	// 	for i, r := range repositoriesSlice {

	// 		targetContainerPath := fmt.Sprintf("%v/%v:%v", r, *container, estafetteBuildVersion)

	// 		if i > 0 {
	// 			// tag container with default tag (it already exists for the first repository)
	// 			log.Printf("Tagging container image %v\n", targetContainerPath)
	// 			tagArgs := []string{
	// 				"tag",
	// 				sourceContainerPath,
	// 				targetContainerPath,
	// 			}
	// 			err := exec.Command("docker", tagArgs...).Run()
	// 			handleError(err)
	// 		}

	// 		loginIfRequired(credentials, targetContainerPath)

	// 		// push container with default tag
	// 		log.Printf("Pushing container image %v\n", targetContainerPath)
	// 		pushArgs := []string{
	// 			"push",
	// 			targetContainerPath,
	// 		}
	// 		runCommand("docker", pushArgs)

	// 		// push additional tags
	// 		for _, t := range tagsSlice {

	// 			targetContainerPath := fmt.Sprintf("%v/%v:%v", r, *container, t)

	// 			// tag container with additional tag
	// 			log.Printf("Tagging container image %v\n", targetContainerPath)
	// 			tagArgs := []string{
	// 				"tag",
	// 				sourceContainerPath,
	// 				targetContainerPath,
	// 			}
	// 			runCommand("docker", tagArgs)

	// 			loginIfRequired(credentials, targetContainerPath)

	// 			log.Printf("Pushing container image %v\n", targetContainerPath)
	// 			pushArgs := []string{
	// 				"push",
	// 				targetContainerPath,
	// 			}
	// 			runCommand("docker", pushArgs)
	// 		}
	// 	}

	// case "tag":

	// 	// image: extensions/docker:stable
	// 	// action: tag
	// 	// container: docker
	// 	// repositories:
	// 	// - extensions
	// 	// tags:
	// 	// - stable
	// 	// - latest

	// 	sourceContainerPath := fmt.Sprintf("%v/%v:%v", repositoriesSlice[0], *container, estafetteBuildVersion)

	// 	loginIfRequired(credentials, sourceContainerPath)

	// 	// pull source container first
	// 	log.Printf("Pulling container image %v\n", sourceContainerPath)
	// 	pullArgs := []string{
	// 		"pull",
	// 		sourceContainerPath,
	// 	}
	// 	runCommand("docker", pullArgs)

	// 	// push each repository + tag combination
	// 	for i, r := range repositoriesSlice {

	// 		targetContainerPath := fmt.Sprintf("%v/%v:%v", r, *container, estafetteBuildVersion)

	// 		if i > 0 {
	// 			// tag container with default tag
	// 			log.Printf("Tagging container image %v\n", targetContainerPath)
	// 			tagArgs := []string{
	// 				"tag",
	// 				sourceContainerPath,
	// 				targetContainerPath,
	// 			}
	// 			runCommand("docker", tagArgs)

	// 			loginIfRequired(credentials, targetContainerPath)

	// 			// push container with default tag
	// 			log.Printf("Pushing container image %v\n", targetContainerPath)
	// 			pushArgs := []string{
	// 				"push",
	// 				targetContainerPath,
	// 			}
	// 			runCommand("docker", pushArgs)
	// 		}

	// 		// push additional tags
	// 		for _, t := range tagsSlice {

	// 			targetContainerPath := fmt.Sprintf("%v/%v:%v", r, *container, t)

	// 			// tag container with additional tag
	// 			log.Printf("Tagging container image %v\n", targetContainerPath)
	// 			tagArgs := []string{
	// 				"tag",
	// 				sourceContainerPath,
	// 				targetContainerPath,
	// 			}
	// 			runCommand("docker", tagArgs)

	// 			loginIfRequired(credentials, targetContainerPath)

	// 			log.Printf("Pushing container image %v\n", targetContainerPath)
	// 			pushArgs := []string{
	// 				"push",
	// 				targetContainerPath,
	// 			}
	// 			runCommand("docker", pushArgs)
	// 		}
	// 	}

	// default:
	// 	log.Fatal("Set `command: <command>` on this step to build, push or tag")
	// }
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
