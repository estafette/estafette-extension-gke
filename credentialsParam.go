package main

import (
	"fmt"
)

// CredentialsParam is used to first retrieve credentials and use any defaults set there
type CredentialsParam struct {
	Credentials string `json:"credentials,omitempty"`
}

// SetDefaults fills in empty fields with convention-based defaults
func (p *CredentialsParam) SetDefaults(releaseName string) {
	// default credentials to release name prefixed with gke if no override in stage params
	if p.Credentials == "" && releaseName != "" {
		p.Credentials = fmt.Sprintf("gke-%v", releaseName)
	}
}

// ValidateRequiredProperties checks whether all needed properties are set
func (p *CredentialsParam) ValidateRequiredProperties() (bool, []error) {

	errors := []error{}

	// validate control params
	if p.Credentials == "" {
		errors = append(errors, fmt.Errorf("Credentials property is required; set it via credentials property on this stage"))
	}

	return len(errors) == 0, errors
}
