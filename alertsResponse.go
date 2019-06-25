package main

import (
	"time"
)

type alertsResponse struct {
	Status string `json:"status"`
	Data   struct {
		Alerts []struct {
			Labels label `json:"labels,omitempty"`
			Annotations struct {
				Description string `json:"description"`
				Summary     string `json:"summary"`
			} `json:"annotations"`
			State    string    `json:"state"`
			ActiveAt time.Time `json:"activeAt"`
			Value    float64   `json:"value"`
		} `json:"alerts"`
	} `json:"data"`
}

type label struct {
	Alertname           string `json:"alertname"`
	Deployment          string `json:"deployment"`
	GcloudProject       string `json:"gcloud_project"`
	GcpProject          string `json:"gcp_project"`
	GkeCluster          string `json:"gke_cluster"`
	Instance            string `json:"instance"`
	Job                 string `json:"job"`
	K8SApp              string `json:"k8s_app"`
	KubernetesCluster   string `json:"kubernetes_cluster"`
	KubernetesName      string `json:"kubernetes_name"`
	KubernetesNamespace string `json:"kubernetes_namespace"`
	LabelTeam           string `json:"label_team"`
	Namespace           string `json:"namespace"`
	TravixEnvironment   string `json:"travix_environment"`
}