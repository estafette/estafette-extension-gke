package credentials

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"

	"github.com/estafette/estafette-extension-gke/api"
	foundation "github.com/estafette/estafette-foundation"
	"github.com/rs/zerolog/log"
)

//go:generate mockgen -package=credentials -destination ./mock.go -source=client.go
type Client interface {
	Init(ctx context.Context, paramsJSON, releaseName, credentialsPath string) (credential *api.GKECredentials, err error)
	GetCredentialsByName(c []api.GKECredentials, credentialName string) *api.GKECredentials
}

// NewClient returns a new gcp.Client
func NewClient(ctx context.Context) (Client, error) {
	return &client{}, nil
}

type client struct {
}

func (c *client) Init(ctx context.Context, paramsJSON, releaseName, credentialsPath string) (credential *api.GKECredentials, err error) {
	log.Info().Msg("Unmarshalling credentials parameter...")
	var credentialsParam api.CredentialsParam
	err = json.Unmarshal([]byte(paramsJSON), &credentialsParam)
	if err != nil {
		return

	}

	log.Info().Msg("Setting default for credential parameter...")
	credentialsParam.SetDefaults(releaseName)

	log.Info().Msg("Validating required credential parameter...")
	valid, errors := credentialsParam.ValidateRequiredProperties()
	if !valid {
		return nil, fmt.Errorf("Not all valid fields are set: %v", errors)
	}

	log.Info().Msg("Unmarshalling injected credentials...")
	var credentials []api.GKECredentials

	// use mounted credential file if present instead of relying on an envvar
	if runtime.GOOS == "windows" {
		credentialsPath = "C:" + credentialsPath
	}
	if foundation.FileExists(credentialsPath) {
		log.Info().Msgf("Reading credentials from file at path %v...", credentialsPath)
		credentialsFileContent, err := ioutil.ReadFile(credentialsPath)
		if err != nil {
			return nil, fmt.Errorf("Failed reading credential file at path %v.", credentialsPath)
		}
		err = json.Unmarshal(credentialsFileContent, &credentials)
		if err != nil {
			return nil, fmt.Errorf("Failed unmarshalling injected credentials: %w", err)
		}
		if len(credentials) == 0 {
			log.Warn().Str("data", string(credentialsFileContent)).Msgf("Found 0 credentials in file %v", credentialsPath)
		}
		log.Debug().Msgf("Read %v credentials", len(credentials))
	}

	log.Info().Msgf("Checking if credential %v exists...", credentialsParam.Credentials)
	credential = c.GetCredentialsByName(credentials, credentialsParam.Credentials)
	if credential == nil {
		return nil, fmt.Errorf("Credential with name %v does not exist", credentialsParam.Credentials)
	}

	log.Info().Msgf("Storing gke credential %v on disk at path %v...", credentialsParam.Credentials, os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))
	err = ioutil.WriteFile(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"), []byte(credential.AdditionalProperties.ServiceAccountKeyfile), 0666)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed writing service account keyfile")
	}

	return
}

func (c *client) GetCredentialsByName(creds []api.GKECredentials, credentialName string) *api.GKECredentials {
	for _, cred := range creds {
		if cred.Name == credentialName {
			return &cred
		}
	}

	return nil
}
