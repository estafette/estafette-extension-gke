package api

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

func (c *GKECredentials) GetLocation() string {
	if c.AdditionalProperties.Region != "" {
		return c.AdditionalProperties.Region
	}

	return c.AdditionalProperties.Zone
}
