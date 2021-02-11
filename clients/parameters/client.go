package parameters

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/estafette/estafette-extension-gke/api"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
)

//go:generate mockgen -package=parameters -destination ./mock.go -source=client.go
type Client interface {
	Init(ctx context.Context, paramsYAML string, credential *api.GKECredentials, gitSource, gitOwner, gitName, appLabel, buildVersion, releaseName, releaseAction, releaseID string) (parameters api.Params, err error)
}

// NewClient returns a new gcp.Client
func NewClient(ctx context.Context) (Client, error) {
	return &client{}, nil
}

type client struct {
}

func (c *client) Init(ctx context.Context, paramsYAML string, credential *api.GKECredentials, gitSource, gitOwner, gitName, appLabel, buildVersion, releaseName, releaseAction, releaseID string) (parameters api.Params, err error) {

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

	var params api.Params
	if credential.AdditionalProperties.Defaults != nil {
		log.Info().Msgf("Using defaults from credential %v...", credential.Name)
		// todo log just the specified defaults, not the entire parms object
		// defaultsAsYAML, err := yaml.Marshal(credential.AdditionalProperties.Defaults)
		// if err == nil {
		// 	log.Printf(string(defaultsAsYAML))
		// }
		params = *credential.AdditionalProperties.Defaults
	}

	log.Info().Msg("Unmarshalling parameters / custom properties...")
	err = yaml.Unmarshal([]byte(paramsYAML), &params)
	if err != nil {
		return params, fmt.Errorf("Failed unmarshalling parameters: %w", err)
	}

	log.Info().Msg("Setting defaults for parameters that are not set in the manifest...")
	params.SetDefaults(gitSource, gitOwner, gitName, appLabel, buildVersion, releaseName, api.ActionType(releaseAction), releaseID, estafetteLabels)

	log.Info().Msg("Validating required parameters...")
	valid, errors, warnings := params.ValidateRequiredProperties()
	if !valid {
		return params, fmt.Errorf("Not all valid fields are set: %v", errors)
	}

	for _, warning := range warnings {
		log.Printf("Warning: %s", warning)
	}

	// replacing sidecar image tags with digest
	params.ReplaceSidecarTagsWithDigest()

	return
}
