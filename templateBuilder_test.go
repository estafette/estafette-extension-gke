package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTemplates(t *testing.T) {

	t.Run("IncludesIngressIfVisibilityIsPrivate", func(t *testing.T) {

		params := Params{
			Visibility: "private",
		}

		// act
		templates := getTemplates(params)

		assert.True(t, stringArrayContains(templates, "ingress.yaml"))
	})
}

func stringArrayContains(array []string, search string) bool {
	for _, v := range array {
		if v == search {
			return true
		}
	}
	return false
}
