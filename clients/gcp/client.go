package gcp

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/estafette/estafette-extension-gke/api"
	foundation "github.com/estafette/estafette-foundation"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2/google"
	containerv1 "google.golang.org/api/container/v1beta1"
	"google.golang.org/api/googleapi"
	iamv1 "google.golang.org/api/iam/v1"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

var (
	// ErrAPIForbidden is returned when the api returns a 401
	ErrAPIForbidden = wrapError{msg: "The api is not allowed for the current service account"}

	// ErrAPINotEnabled is returned when an api is not enabled
	ErrAPINotEnabled = wrapError{msg: "The api is not enabled"}

	// ErrUnknownProjectID is returned when an api throws 'googleapi: Error 400: Unknown project id: 0, invalid'
	ErrUnknownProjectID = wrapError{msg: "The project id is unknown"}

	// ErrProjectNotFound is returned when an api throws 'googleapi: Error 404: The requested project was not found., notFound'
	ErrProjectNotFound = wrapError{msg: "The project is not found"}

	// ErrEntityNotFound is returned when pubsub topics return html with a 404
	ErrEntityNotFound = wrapError{msg: "Entity is not found"}

	// ErrEntityNotActive is returned when cloud sql instance is not running and its databases cannot be fetched
	ErrEntityNotActive = wrapError{msg: "Entity is not runactivening"}
)

//go:generate mockgen -package=gcp -destination ./mock.go -source=client.go
type Client interface {
	LoadGKEClusterKubeConfig(ctx context.Context, credential *api.GKECredentials) (kubeContextName string, err error)
	GetGKECluster(ctx context.Context, projectID, location, clusterID string) (cluster *containerv1.Cluster, err error)
	DeployGoogleCloudEndpoints(ctx context.Context, params api.Params) (err error)
}

// NewClient returns a new gcp.Client
func NewClient(ctx context.Context) (Client, error) {

	// use service account to authenticate against gcp apis
	googleClient, err := google.DefaultClient(ctx, iamv1.CloudPlatformScope)
	if err != nil {
		return nil, err
	}

	containerv1Service, err := containerv1.New(googleClient)
	if err != nil {
		return nil, err
	}

	return &client{
		containerv1Service: containerv1Service,
	}, nil
}

type client struct {
	containerv1Service *containerv1.Service
}

func (c *client) LoadGKEClusterKubeConfig(ctx context.Context, credential *api.GKECredentials) (kubeContextName string, err error) {
	if credential == nil {
		return kubeContextName, fmt.Errorf("LoadGKEClusterKubeConfig argument credential is nil")
	}
	if credential.AdditionalProperties.Project == "" {
		return kubeContextName, fmt.Errorf("LoadGKEClusterKubeConfig credential argument has empty Project")
	}
	if credential.AdditionalProperties.Cluster == "" {
		return kubeContextName, fmt.Errorf("LoadGKEClusterKubeConfig credential argument has empty Cluster")
	}
	if credential.AdditionalProperties.Region == "" && credential.AdditionalProperties.Zone == "" {
		return kubeContextName, fmt.Errorf("LoadGKEClusterKubeConfig credential argument has empty Region or Zone")
	}

	kubeContextName = fmt.Sprintf("%v-%v-%v", credential.AdditionalProperties.Project, credential.GetLocation(), credential.AdditionalProperties.Cluster)

	log.Info().Msgf("Generating .kube/config sections for context %v", kubeContextName)

	cluster, err := c.GetGKECluster(ctx, credential.AdditionalProperties.Project, credential.GetLocation(), credential.AdditionalProperties.Cluster)
	if err != nil {
		return
	}

	kubeConfigPath := os.Getenv("KUBECONFIG")
	if kubeConfigPath == "" {
		return kubeContextName, fmt.Errorf("Value of envvar KUBECONFIG is empty, cannot create kube config")
	}

	// check if kubeconfig exists and read if it does
	currentConfig := clientcmdapi.NewConfig()
	if foundation.FileExists(kubeConfigPath) {
		currentConfig, err = clientcmd.LoadFromFile(kubeConfigPath)
		if err != nil {
			return
		}
	}

	decodedClusterCaCertificate, err := base64.StdEncoding.DecodeString(cluster.MasterAuth.ClusterCaCertificate)
	if err != nil {
		return
	}

	// add cluster if it doesn't exist
	if _, exists := currentConfig.Clusters[kubeContextName]; !exists {
		currentConfig.Clusters[kubeContextName] = &clientcmdapi.Cluster{
			Server:                   fmt.Sprintf("https://%v", cluster.Endpoint),
			CertificateAuthorityData: decodedClusterCaCertificate,
		}
	}

	// add authinfo if it doesn't exist
	if _, exists := currentConfig.AuthInfos[kubeContextName]; !exists {
		currentConfig.AuthInfos[kubeContextName] = &clientcmdapi.AuthInfo{
			AuthProvider: &clientcmdapi.AuthProviderConfig{
				Name: "gcp",
			},
		}
	}

	// add context if it doesn't exist
	if _, exists := currentConfig.Contexts[kubeContextName]; !exists {
		currentConfig.Contexts[kubeContextName] = &clientcmdapi.Context{
			Cluster:  kubeContextName,
			AuthInfo: kubeContextName,
		}
	}

	// set apiversion if empty
	if currentConfig.APIVersion == "" {
		currentConfig.APIVersion = "v1"
	}

	// set kind if empty
	if currentConfig.Kind == "" {
		currentConfig.Kind = "Config"
	}

	currentConfig.CurrentContext = kubeContextName

	// write kube config file
	err = clientcmd.WriteToFile(*currentConfig, kubeConfigPath)
	if err != nil {
		return
	}

	return
}

func (c *client) GetGKECluster(ctx context.Context, projectID, location, clusterID string) (cluster *containerv1.Cluster, err error) {
	if projectID == "" {
		return nil, fmt.Errorf("GetGKECluster argument projectID is empty")
	}
	if location == "" {
		return nil, fmt.Errorf("GetGKECluster argument location is empty")
	}
	if clusterID == "" {
		return nil, fmt.Errorf("GetGKECluster argument clusterID is empty")
	}

	log.Debug().Msgf("Retrieving GKE cluster %v in location %v project %v...", clusterID, location, projectID)

	err = c.substituteErrorsWithPredefinedErrors(foundation.Retry(func() error {
		// https://cloud.google.com/kubernetes-engine/docs/reference/rest/v1/projects.locations.clusters/get
		cluster, err = c.containerv1Service.Projects.Locations.Clusters.Get("projects/" + projectID + "/locations/" + location + "/clusters/" + clusterID).Context(ctx).Do()
		if err != nil {
			return err
		}
		return nil
	}, c.getRetryOptions()...))
	if err != nil {
		return cluster, fmt.Errorf("Can't get gke cluster %v in location %v for project %v: %w", clusterID, location, projectID, err)
	}

	log.Debug().Msgf("Retrieved GKE cluster %v in location %v for project %v", clusterID, location, projectID)

	return
}

func (c *client) DeployGoogleCloudEndpoints(ctx context.Context, params api.Params) (err error) {
	return foundation.RunCommandWithArgsExtended(ctx, "gcloud", []string{"endpoints", "--project", params.EspEndpointsProjectID, "services", "deploy", params.EspOpenAPIYamlPath})
}

func (c *client) substituteErrorsWithPredefinedErrors(err error) error {
	if err == nil {
		return nil
	}

	if googleapiErr, ok := err.(*googleapi.Error); ok && googleapiErr.Code == http.StatusForbidden {
		return ErrAPIForbidden.wrap(err)
	}
	if googleapiErr, ok := err.(*googleapi.Error); ok && googleapiErr.Code == http.StatusBadRequest && err.Error() == "googleapi: Error 400: Unknown project id: 0, invalid" {
		return ErrUnknownProjectID.wrap(err)
	}
	if googleapiErr, ok := err.(*googleapi.Error); ok && googleapiErr.Code == http.StatusBadRequest && strings.HasSuffix(err.Error(), "has not enabled BigQuery., invalid") {
		return ErrAPINotEnabled.wrap(err)
	}
	if googleapiErr, ok := err.(*googleapi.Error); ok && googleapiErr.Code == http.StatusBadRequest && strings.HasSuffix(err.Error(), "Invalid request: Invalid request since instance is not running., invalid") {
		return ErrEntityNotActive.wrap(err)
	}
	if googleapiErr, ok := err.(*googleapi.Error); ok && googleapiErr.Code == http.StatusNotFound && err.Error() == "googleapi: Error 404: The requested project was not found., notFound" {
		return ErrProjectNotFound.wrap(err)
	}
	if googleapiErr, ok := err.(*googleapi.Error); ok && googleapiErr.Code == http.StatusNotFound {
		return ErrEntityNotFound.wrap(err)
	}
	if googleapiErr, ok := err.(*googleapi.Error); ok && googleapiErr.Code == http.StatusNoContent {
		return ErrEntityNotFound.wrap(err)
	}

	return err
}

func (c *client) getRetryOptions(extraRetryableStatuses ...int) []foundation.RetryOption {
	return []foundation.RetryOption{
		c.isRetryableErrorCustomOption(extraRetryableStatuses...),
		foundation.LastErrorOnly(true),
		foundation.Attempts(5),
		foundation.DelayMillisecond(1000),
	}
}

func (c *client) isRetryableErrorCustomOption(extraRetryableStatuses ...int) foundation.RetryOption {

	return func(c *foundation.RetryConfig) {
		c.IsRetryableError = func(err error) bool {
			switch e := err.(type) {
			case *googleapi.Error:
				// Retry on 429 and 5xx, according to
				// https://cloud.google.com/storage/docs/exponential-backoff.
				return e.Code == http.StatusTooManyRequests || (e.Code >= 500 && e.Code < 600) || e.Code == http.StatusForbidden || foundation.IntArrayContains(extraRetryableStatuses, e.Code)
			case *url.Error:
				// Retry socket-level errors ECONNREFUSED and ENETUNREACH (from syscall).
				// Unfortunately the error type is unexported, so we resort to string
				// matching.
				retriable := []string{"connection refused", "connection reset"}
				for _, s := range retriable {
					if strings.Contains(e.Error(), s) {
						return true
					}
				}
				return false
			case interface{ Temporary() bool }:
				return e.Temporary()
			default:
				return false
			}
		}
	}
}
