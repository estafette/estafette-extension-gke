package main

// Parameters controls the behaviour of this extension
type Parameters struct {
	Name      string `json:"name,omitempty"`
	Namespace string `json:"namespace,omitempty"`
	DryRun    bool   `json:"dryrun,omitempty"`
}
