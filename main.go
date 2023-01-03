package main

import (
	"context"
	"runtime"

	"github.com/alecthomas/kingpin"
	"github.com/estafette/estafette-extension-gke/api"
	"github.com/estafette/estafette-extension-gke/clients/credentials"
	"github.com/estafette/estafette-extension-gke/clients/gcp"
	"github.com/estafette/estafette-extension-gke/clients/parameters"
	"github.com/estafette/estafette-extension-gke/services/builder"
	"github.com/estafette/estafette-extension-gke/services/extension"
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
	gitSource        = kingpin.Flag("git-source", "Repository source.").Envar("ESTAFETTE_GIT_SOURCE").String()
	gitOwner         = kingpin.Flag("git-owner", "Repository owner.").Envar("ESTAFETTE_GIT_OWNER").String()
	gitName          = kingpin.Flag("git-name", "Repository name, used as application name if not passed explicitly and app label not being set.").Envar("ESTAFETTE_GIT_NAME").String()
	gitBranch        = kingpin.Flag("git-branch", "Repository commit branch.").Envar("ESTAFETTE_GIT_BRANCH").String()
	gitRevision      = kingpin.Flag("git-revision", "Repository commit revisition.").Envar("ESTAFETTE_GIT_REVISION").String()
	appLabel         = kingpin.Flag("app-name", "App label, used as application name if not passed explicitly.").Envar("ESTAFETTE_LABEL_APP").String()
	buildVersion     = kingpin.Flag("build-version", "Version number, used if not passed explicitly.").Envar("ESTAFETTE_BUILD_VERSION").String()
	releaseName      = kingpin.Flag("release-name", "Name of the release section, which is used by convention to resolve the credentials.").Envar("ESTAFETTE_RELEASE_NAME").String()
	releaseAction    = kingpin.Flag("release-action", "Name of the release action, to control the type of release.").Envar("ESTAFETTE_RELEASE_ACTION").String()
	releaseID        = kingpin.Flag("release-id", "ID of the release, to use as a label.").Envar("ESTAFETTE_RELEASE_ID").String()
	triggeredBy      = kingpin.Flag("triggered-by", "The user id of the person triggering the release.").Envar("ESTAFETTE_TRIGGER_MANUAL_USER_ID").String()
	builderImageSHA  = kingpin.Flag("builder-image-sha", "The SHA of the image that is running the stage").Envar("ESTAFETTE_STAGE_IMAGE_SHA").String()
	builderImageDate = kingpin.Flag("builder-image-date", "The creation date of the image that is running the stage").Envar("ESTAFETTE_STAGE_IMAGE_CREATED_DATE").String()

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

	extensionService, err := extension.NewService(ctx, credentialsClient, parametersClient, gcpClient, builderService, generatorService)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed creating extension.Service")
	}

	err = extensionService.Run(ctx, credential, *releaseName, *paramsYAML, *gitSource, *gitOwner, *gitName, *appLabel, *buildVersion, *releaseAction, *releaseID, *gitBranch, *gitRevision, *builderImageSHA, *builderImageDate, *triggeredBy)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed running extension.Service")
	}
}
