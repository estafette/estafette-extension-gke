package gcp

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/estafette/estafette-extension-gke/api"
	foundation "github.com/estafette/estafette-foundation"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2/google"
	containerv1 "google.golang.org/api/container/v1beta1"
	"google.golang.org/api/googleapi"
	iamv1 "google.golang.org/api/iam/v1"
	servicemanagementv1 "google.golang.org/api/servicemanagement/v1"
	"gopkg.in/yaml.v2"
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

	// ErrServiceNotFound is returned when an a cloud endpoints service cannot be found
	ErrServiceNotFound = wrapError{msg: "The cloud endpoints service is not found"}
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

	servicemanagementv1Service, err := servicemanagementv1.New(googleClient)
	if err != nil {
		return nil, err
	}

	return &client{
		containerv1Service:         containerv1Service,
		servicemanagementv1Service: servicemanagementv1Service,
	}, nil
}

type client struct {
	containerv1Service         *containerv1.Service
	servicemanagementv1Service *servicemanagementv1.APIService
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
	if params.EspOpenAPIYamlPath != "" {
		return c.deployGoogleCloudEndpointsWithOpenApi(ctx, params)
	} else {
		return c.deployGoogleCloudEndpointsWithGrpc(ctx, params)
	}
}

func (c *client) deployGoogleCloudEndpointsWithOpenApi(ctx context.Context, params api.Params) (err error) {
	log.Info().Msgf("Checking if openapi spec at path %v exists...", params.EspOpenAPIYamlPath)
	// get host from openapi file
	if !foundation.FileExists(params.EspOpenAPIYamlPath) {
		return fmt.Errorf("File at path %v does not exist. Did you forget to use `clone: true` for your release?", params.EspOpenAPIYamlPath)
	}

	openapiSpecBytes, err := ioutil.ReadFile(params.EspOpenAPIYamlPath)
	if err != nil {
		return
	}

	var openapiSpec struct {
		Host string `yaml:"host"`
	}

	log.Info().Msg("Unmarshalling openapi spec to get service name...")
	err = yaml.Unmarshal(openapiSpecBytes, &openapiSpec)
	if err != nil {
		return
	}

	if openapiSpec.Host == "" {
		return fmt.Errorf("The host field in the openapi spec at %v is empty, please set it", params.EspOpenAPIYamlPath)
	}

	serviceName := openapiSpec.Host
	log.Info().Msgf("Found service name %v in openapi", serviceName)

	if err = c.createService(ctx, serviceName, params); err != nil {
		return err
	}

	configID, err := c.submitServiceConfiguration(
		ctx,
		params,
		serviceName,
		[]*servicemanagementv1.ConfigFile{
			{
				FileContents: base64.StdEncoding.EncodeToString(openapiSpecBytes),
				FilePath:     filepath.Base(params.EspOpenAPIYamlPath),
				FileType:     "OPEN_API_YAML",
			},
		})

	if err != nil {
		return err
	}

	if err = c.rolloutService(ctx, params, serviceName, configID); err != nil {
		return err
	}

	return nil
}

func (c *client) deployGoogleCloudEndpointsWithGrpc(ctx context.Context, params api.Params) (err error) {
	log.Info().Msgf("Checking if gRPC service config file at path %v exists...", params.EspGrpcConfigYamlPath)
	if !foundation.FileExists(params.EspGrpcConfigYamlPath) {
		return fmt.Errorf("File at path %v does not exist. Did you forget to use `clone: true` for your release?", params.EspGrpcConfigYamlPath)
	}

	log.Info().Msgf("Checking if the Proto descriptor file at path %v exists...", params.EspGrpcProtoDescriptorPath)
	if !foundation.FileExists(params.EspGrpcProtoDescriptorPath) {
		return fmt.Errorf("File at path %v does not exist. Did you forget to use `clone: true` for your release?", params.EspGrpcProtoDescriptorPath)
	}

	grpcServiceConfigBytes, err := ioutil.ReadFile(params.EspGrpcConfigYamlPath)
	if err != nil {
		return
	}

	grpcProtoDescriptorBytes, err := ioutil.ReadFile(params.EspGrpcProtoDescriptorPath)
	if err != nil {
		return
	}

	var grpcServiceConfig struct {
		Name string `yaml:"name"`
	}

	log.Info().Msg("Unmarshalling the gRPC config spec to get the service name...")
	err = yaml.Unmarshal(grpcServiceConfigBytes, &grpcServiceConfig)
	if err != nil {
		return
	}

	if grpcServiceConfig.Name == "" {
		return fmt.Errorf("The name field in the gRPC config spec at %v is empty, please set it", params.EspGrpcConfigYamlPath)
	}

	serviceName := grpcServiceConfig.Name
	log.Info().Msgf("Found service name %v in the gRPC config", serviceName)

	if err = c.createService(ctx, serviceName, params); err != nil {
		return err
	}

	configID, err := c.submitServiceConfiguration(
		ctx,
		params,
		serviceName,
		[]*servicemanagementv1.ConfigFile{
			{
				FileContents: base64.StdEncoding.EncodeToString(grpcProtoDescriptorBytes),
				FilePath:     filepath.Base(params.EspGrpcProtoDescriptorPath),
				FileType:     "FILE_DESCRIPTOR_SET_PROTO",
			},
			{
				FileContents: base64.StdEncoding.EncodeToString(grpcServiceConfigBytes),
				FilePath:     filepath.Base(params.EspGrpcConfigYamlPath),
				FileType:     "SERVICE_CONFIG_YAML",
			},
		})

	if err != nil {
		return err
	}

	if err = c.rolloutService(ctx, params, serviceName, configID); err != nil {
		return err
	}

	return nil
}

func (c *client) createService(ctx context.Context, serviceName string, params api.Params) (err error) {
	log.Info().Msgf("Checking if service %v exists...", serviceName)
	// GET https://servicemanagement.googleapis.com/v1/services/<servicename>
	var service *servicemanagementv1.ManagedService
	err = c.substituteErrorsWithPredefinedErrors(foundation.Retry(func() error {
		// https://cloud.google.com/service-infrastructure/docs/service-management/reference/rest/v1/services/get
		service, err = c.servicemanagementv1Service.Services.Get(serviceName).Context(ctx).Do()
		if err != nil {
			return err
		}
		return nil
	}, c.getRetryOptions()...))
	if err != nil && !errors.Is(err, ErrServiceNotFound) {
		return fmt.Errorf("Can't get service %v for project %v: %w", serviceName, params.EspEndpointsProjectID, err)
	}

	if service == nil {
		log.Info().Msgf("Creating service %v in project %v...", serviceName, params.EspEndpointsProjectID)
		// create the service
		// POST https://servicemanagement.googleapis.com/v1/services/<servicename>
		service = &servicemanagementv1.ManagedService{
			ProducerProjectId: params.EspEndpointsProjectID,
			ServiceName:       serviceName,
		}

		var operation *servicemanagementv1.Operation
		err = c.substituteErrorsWithPredefinedErrors(foundation.Retry(func() error {
			// https://cloud.google.com/service-infrastructure/docs/service-management/reference/rest/v1/services/create
			operation, err = c.servicemanagementv1Service.Services.Create(service).Context(ctx).Do()
			if err != nil {
				return err
			}
			return nil
		}, c.getRetryOptions()...))
		if err != nil {
			return fmt.Errorf("Can't get service %v for project %v: %w", serviceName, params.EspEndpointsProjectID, err)
		}

		err = c.waitForServiceManagementV1Operation(ctx, params.EspEndpointsProjectID, operation)
		if err != nil {
			return
		}
	}

	return
}

func (c *client) submitServiceConfiguration(ctx context.Context, params api.Params, serviceName string, configFiles []*servicemanagementv1.ConfigFile) (configID string, err error) {
	log.Info().Msgf("Submitting config for service %v in project %v...", serviceName, params.EspEndpointsProjectID)
	// POST https://servicemanagement.googleapis.com/v1/services/<servicename>/configs:submit
	var operation *servicemanagementv1.Operation
	err = c.substituteErrorsWithPredefinedErrors(foundation.Retry(func() error {
		// https://cloud.google.com/service-infrastructure/docs/service-management/reference/rest/v1/services.configs/submit
		operation, err = c.servicemanagementv1Service.Services.Configs.Submit(serviceName, &servicemanagementv1.SubmitConfigSourceRequest{
			ConfigSource: &servicemanagementv1.ConfigSource{
				Files: configFiles,
			},
			ValidateOnly: false,
		}).Context(ctx).Do()
		if err != nil {
			return err
		}
		return nil
	}, c.getRetryOptions()...))
	if err != nil {
		return "", fmt.Errorf("Can't submit config for service %v for project %v: %w", serviceName, params.EspEndpointsProjectID, err)
	}

	// GET https://servicemanagement.googleapis.com/v1/operations/serviceConfigs.<servicename>%3<config id>
	err = c.waitForServiceManagementV1Operation(ctx, params.EspEndpointsProjectID, operation)
	if err != nil {
		return "", err
	}

	var response servicemanagementv1.SubmitConfigSourceResponse
	err = json.Unmarshal(operation.Response, &response)
	if err != nil {
		return "", err
	}

	configID = response.ServiceConfig.Id
	log.Info().Msgf("Submitted config with id %v for service %v in project %v", configID, serviceName, params.EspEndpointsProjectID)

	return configID, nil
}

func (c *client) rolloutService(ctx context.Context, params api.Params, serviceName string, configID string) (err error) {
	log.Info().Msgf("Creating rollout for config with id %v for service %v in project %v...", configID, serviceName, params.EspEndpointsProjectID)
	// POST https://servicemanagement.googleapis.com/v1/services/<servicename>/rollouts
	var operation *servicemanagementv1.Operation
	err = c.substituteErrorsWithPredefinedErrors(foundation.Retry(func() error {
		// https://cloud.google.com/service-infrastructure/docs/service-management/reference/rest/v1/services.rollouts/create
		operation, err = c.servicemanagementv1Service.Services.Rollouts.Create(serviceName, &servicemanagementv1.Rollout{
			ServiceName: serviceName,
			TrafficPercentStrategy: &servicemanagementv1.TrafficPercentStrategy{
				Percentages: map[string]float64{
					configID: 100,
				},
			},
		}).Context(ctx).Do()
		if err != nil {
			return err
		}
		return nil
	}, c.getRetryOptions()...))
	if err != nil {
		return fmt.Errorf("Can't create rollout for service %v for project %v: %w", serviceName, params.EspEndpointsProjectID, err)
	}

	// GET https://servicemanagement.googleapis.com/v1/operations/rollouts.<servicename>%3A9b4bc80c-94a5-49e8-8984-28631648a1d1
	err = c.waitForServiceManagementV1Operation(ctx, params.EspEndpointsProjectID, operation)
	if err != nil {
		return
	}

	log.Info().Msgf("Retrieving service %v in project %v after rollout has finished...", serviceName, params.EspEndpointsProjectID)
	// GET https://servicemanagement.googleapis.com/v1/services/<servicename>
	err = c.substituteErrorsWithPredefinedErrors(foundation.Retry(func() error {
		// https://cloud.google.com/service-infrastructure/docs/service-management/reference/rest/v1/services/get
		_, err = c.servicemanagementv1Service.Services.Get(serviceName).Context(ctx).Do()
		if err != nil {
			return err
		}
		return nil
	}, c.getRetryOptions()...))
	if err != nil {
		return fmt.Errorf("Can't get service %v for project %v: %w", serviceName, params.EspEndpointsProjectID, err)
	}

	return
}

func (c *client) substituteErrorsWithPredefinedErrors(err error) error {
	if err == nil {
		return nil
	}

	if googleapiErr, ok := err.(*googleapi.Error); ok && googleapiErr.Code == http.StatusForbidden && strings.Contains(err.Error(), "Service") && strings.Contains(err.Error(), "not found or permission denied") {
		return ErrServiceNotFound.wrap(err)
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

func (c *client) waitForServiceManagementV1Operation(ctx context.Context, projectID string, operation *servicemanagementv1.Operation) (err error) {

	if operation != nil {
		startTime := time.Now().UTC()
		timeTakenSeconds := time.Now().UTC().Sub(startTime).Seconds()
		operationName := operation.Name
		for operation != nil && operationName != "" && !operation.Done {
			log.Info().Msgf("Waiting for operation %v in project %v to finish (%vs)...", operationName, projectID, math.Round(timeTakenSeconds))

			// sleep first, otherwise the check immediately follows the action
			c.sleepForWait(startTime)

			err = c.substituteErrorsWithPredefinedErrors(foundation.Retry(func() error {
				// https://cloud.google.com/resource-manager/reference/rest/v1/operations/get
				operation, err = c.servicemanagementv1Service.Operations.Get(operationName).Context(ctx).Do()
				if err != nil {
					return err
				}

				return nil
			}, c.getRetryOptions()...))
			if err != nil {
				return
			}
			timeTakenSeconds = time.Now().UTC().Sub(startTime).Seconds()
		}

		// check if the operation ended with an error
		if operation != nil && operation.Error != nil {
			return fmt.Errorf("Operation %v finished with an error: %v", operationName, *operation.Error)
		}

		log.Info().Msgf("Done waiting for operation %v in project %v to finish (%vs)", operationName, projectID, math.Round(timeTakenSeconds))

	} else {
		log.Warn().Interface("operation", operation).Msgf("Cannot wait for compute operation in project %v, it's nil", projectID)
	}

	return nil
}

func (c *client) sleepForWait(startTime time.Time) {

	sleepTimeSeconds := 5
	secondsSinceStart := time.Now().UTC().Sub(startTime).Seconds()
	if secondsSinceStart > 300 {
		sleepTimeSeconds = 30
	} else if secondsSinceStart > 120 {
		sleepTimeSeconds = 15
	} else if secondsSinceStart > 30 {
		sleepTimeSeconds = 10
	}

	sleepTime := foundation.ApplyJitter(sleepTimeSeconds)
	time.Sleep(time.Duration(sleepTime) * time.Second)
}
