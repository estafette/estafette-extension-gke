package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCredentialsByName(t *testing.T) {

	t.Run("ReturnsCredentialIfCredentialWithNameExists", func(t *testing.T) {

		credentials := []GKECredentials{
			GKECredentials{
				Name: "gke-production",
			},
		}

		// act
		credential := GetCredentialsByName(credentials, "gke-production")

		assert.NotNil(t, credential)
		assert.Equal(t, "gke-production", credential.Name)
	})

	t.Run("ReturnsNilIfCredentialWithNameDoesNotExist", func(t *testing.T) {

		credentials := []GKECredentials{
			GKECredentials{
				Name: "gke-production",
			},
		}

		// act
		credential := GetCredentialsByName(credentials, "gke-staging")

		assert.Nil(t, credential)
	})
}
