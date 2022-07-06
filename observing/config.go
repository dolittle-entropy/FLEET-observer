package observing

import (
	"dolittle.io/fleet-observer/entities"
	"dolittle.io/fleet-observer/kubernetes"
	"dolittle.io/fleet-observer/mongo"
	"github.com/rs/zerolog"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/listers/core/v1"
	"strings"
)

type ConfigHandler struct {
	configurations *mongo.Configurations
	configmaps     v1.ConfigMapLister
	secrets        v1.SecretLister
	logger         zerolog.Logger
}

func NewConfigHandler(configurations *mongo.Configurations, configmaps v1.ConfigMapLister, secrets v1.SecretLister, logger zerolog.Logger) *ConfigHandler {
	return &ConfigHandler{
		configurations: configurations,
		configmaps:     configmaps,
		secrets:        secrets,
		logger:         logger,
	}
}

func (ch *ConfigHandler) Handle(obj any) error {
	configmap, ok := obj.(*coreV1.ConfigMap)
	if ok {
		if strings.HasSuffix(configmap.GetName(), "-dolittle") {
			return ch.HandleDolittleConfig(configmap)
		} else {
			return ch.HandleCustomerConfig(configmap.ObjectMeta)
		}
	}

	secret, ok := obj.(*coreV1.Secret)
	if ok {
		return ch.HandleCustomerConfig(secret.ObjectMeta)
	}

	return ReceivedWrongType(obj, "ConfigMap or Secret")
}

func (ch *ConfigHandler) HandleCustomerConfig(meta metaV1.ObjectMeta) error {
	logger := ch.logger.With().Str("namespace", meta.GetNamespace()).Str("name", meta.GetName()).Str("type", "customer").Logger()

	tenantID, applicationID, environmentName, microserviceID, ok := GetMicroserviceIdentifiers(meta)
	if !ok {
		logger.Trace().Msg("Skipping customer config because it is missing microservice identifiers")
		return nil
	}

	hash, err := ComputeCustomerConfigHashFor(meta, ch.configmaps, ch.secrets)
	if err != nil {
		return err
	}

	config := entities.CustomerConfiguration{
		ContentHash:                hash,
		ConfigForArtifactID:        microserviceID,
		DeployedInEnvironmentName:  environmentName,
		EnvironmentOfApplicationID: applicationID,
		OwnedByCustomerID:          tenantID,
	}
	if err := ch.configurations.SetCustomer(config); err != nil {
		return err
	}
	logger.Debug().Interface("config", config).Msg("Updated configuration")

	return nil
}

func (ch *ConfigHandler) HandleDolittleConfig(configMap *coreV1.ConfigMap) error {
	logger := ch.logger.With().Str("namespace", configMap.GetNamespace()).Str("name", configMap.GetName()).Str("type", "dolittle").Logger()

	tenantID, applicationID, environmentName, microserviceID, ok := GetMicroserviceIdentifiers(configMap.ObjectMeta)
	if !ok {
		logger.Trace().Msg("Skipping customer config because it is missing microservice identifiers")
		return nil
	}

	hash, err := ComputeRuntimeConfigHashFor(configMap.ObjectMeta, ch.configmaps)
	if err != nil {
		return err
	}

	config := entities.RuntimeConfiguration{
		ContentHash:                hash,
		ConfigForArtifactID:        microserviceID,
		DeployedInEnvironmentName:  environmentName,
		EnvironmentOfApplicationID: applicationID,
		OwnedByCustomerID:          tenantID,
	}
	if err := ch.configurations.SetRuntime(config); err != nil {
		return err
	}
	logger.Debug().Interface("config", config).Msg("Updated configuration")

	return nil
}

func ComputeCustomerConfigHashFor(object metaV1.ObjectMeta, configmaps v1.ConfigMapLister, secrets v1.SecretLister) (string, error) {
	selector, err := createMicroserviceLabelSelectorFor(object)
	if err != nil {
		return "", err
	}

	microserviceConfigMaps, err := configmaps.ConfigMaps(object.GetNamespace()).List(selector)
	if err != nil {
		return "", err
	}

	microserviceSecrets, err := secrets.Secrets(object.GetNamespace()).List(selector)
	if err != nil {
		return "", err
	}

	hasher := kubernetes.NewConfigHasher()

	for _, configMap := range microserviceConfigMaps {
		if strings.HasSuffix(configMap.GetName(), "-config-files") && microserviceEquals(object, configMap.ObjectMeta) {
			hasher.WriteConfigMap(configMap)
			goto FoundConfigFiles
		}
	}
	return "", CouldNotFindConfiguration("-config-files")

FoundConfigFiles:
	for _, configMap := range microserviceConfigMaps {
		if strings.HasSuffix(configMap.GetName(), "-env-variables") && microserviceEquals(object, configMap.ObjectMeta) {
			hasher.WriteConfigMap(configMap)
			goto FoundEnvVariables
		}
	}
	return "", CouldNotFindConfiguration("-env-variables")

FoundEnvVariables:
	for _, secret := range microserviceSecrets {
		if strings.HasSuffix(secret.GetName(), "-secret-env-variables") && microserviceEquals(object, secret.ObjectMeta) {
			hasher.WriteSecret(secret)
			goto FoundSecretEnvVariables
		}
	}
	return "", CouldNotFindConfiguration("-secret-env-variables")

FoundSecretEnvVariables:
	return hasher.GetComputedHash(), nil
}

func ComputeRuntimeConfigHashFor(object metaV1.ObjectMeta, configmaps v1.ConfigMapLister) (string, error) {
	environmentSelector, err := createEnvironmentLabelSelectorFor(object)
	if err != nil {
		return "", err
	}

	microserviceSelector, err := createMicroserviceLabelSelectorFor(object)
	if err != nil {
		return "", err
	}

	environmentConfigMaps, err := configmaps.ConfigMaps(object.GetNamespace()).List(environmentSelector)
	if err != nil {
		return "", err
	}

	microserviceConfigMaps, err := configmaps.ConfigMaps(object.GetNamespace()).List(microserviceSelector)
	if err != nil {
		return "", err
	}

	hasher := kubernetes.NewConfigHasher()

	for _, configMap := range environmentConfigMaps {
		if strings.HasSuffix(configMap.GetName(), "-tenants") && environmentEquals(object, configMap.ObjectMeta) {
			hasher.WriteConfigMap(configMap)
			goto FoundTenants
		}
	}
	return "", CouldNotFindConfiguration("-tenants")

FoundTenants:
	for _, configMap := range microserviceConfigMaps {
		if strings.HasSuffix(configMap.GetName(), "-dolittle") && microserviceEquals(object, configMap.ObjectMeta) {
			hasher.WriteConfigMap(configMap)
			goto FoundDolittle
		}
	}
	return "", CouldNotFindConfiguration("-dolittle")

FoundDolittle:
	return hasher.GetComputedHash(), nil
}
