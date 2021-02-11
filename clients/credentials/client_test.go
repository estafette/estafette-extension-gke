package credentials

import (
	"context"
	"testing"

	"github.com/estafette/estafette-extension-gke/api"
	"github.com/stretchr/testify/assert"
)

func TestGetCredentialsByName(t *testing.T) {

	t.Run("ReturnsCredentialIfCredentialWithNameExists", func(t *testing.T) {

		client, err := NewClient(context.Background())
		assert.Nil(t, err)

		credentials := []api.GKECredentials{
			{
				Name: "gke-production",
			},
		}

		// act
		credential := client.GetCredentialsByName(credentials, "gke-production")

		assert.NotNil(t, credential)
		assert.Equal(t, "gke-production", credential.Name)
	})

	t.Run("ReturnsNilIfCredentialWithNameDoesNotExist", func(t *testing.T) {

		client, err := NewClient(context.Background())
		assert.Nil(t, err)

		credentials := []api.GKECredentials{
			{
				Name: "gke-production",
			},
		}

		// act
		credential := client.GetCredentialsByName(credentials, "gke-staging")

		assert.Nil(t, credential)
	})
}
