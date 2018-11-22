package main

// GKECredentials represents the credentials of type kubernetes-engine as defined in the server config and passed to this trusted image
type GKECredentials struct {
	Name                 string                            `json:"name,omitempty"`
	Type                 string                            `json:"type,omitempty"`
	AdditionalProperties GKECredentialAdditionalProperties `json:"additionalProperties,omitempty"`
}

// GKECredentialAdditionalProperties contains the non standard fields for this type of credentials
type GKECredentialAdditionalProperties struct {
	Project               string  `json:"project,omitempty"`
	Cluster               string  `json:"cluster,omitempty"`
	Region                string  `json:"region,omitempty"`
	Zone                  string  `json:"zone,omitempty"`
	ServiceAccountKeyfile string  `json:"serviceAccountKeyfile,omitempty"`
	Defaults              *Params `json:"defaults,omitempty"`
}

// GetCredentialsByName returns a credential if the name exists
func GetCredentialsByName(c []GKECredentials, credentialName string) *GKECredentials {

	for _, cred := range c {
		if cred.Name == credentialName {
			return &cred
		}
	}

	return nil
}
