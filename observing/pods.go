package observing

import (
	"dolittle.io/fleet-observer/entities"
	"dolittle.io/fleet-observer/kubernetes"
	"dolittle.io/fleet-observer/mongo"
	"fmt"
	"github.com/rs/zerolog"
	coreV1 "k8s.io/api/core/v1"
	listersAppsV1 "k8s.io/client-go/listers/apps/v1"
	listersCoreV1 "k8s.io/client-go/listers/core/v1"
	"strings"
)

type PodsHandler struct {
	configurations *mongo.Configurations
	deployments    *mongo.Deployments
	configmaps     listersCoreV1.ConfigMapLister
	secrets        listersCoreV1.SecretLister
	replicasets    listersAppsV1.ReplicaSetLister
	logger         zerolog.Logger
}

func NewPodsHandler(configurations *mongo.Configurations, deployments *mongo.Deployments, configmaps listersCoreV1.ConfigMapLister, secrets listersCoreV1.SecretLister, replicasets listersAppsV1.ReplicaSetLister, logger zerolog.Logger) *PodsHandler {
	return &PodsHandler{
		configurations: configurations,
		deployments:    deployments,
		configmaps:     configmaps,
		secrets:        secrets,
		replicasets:    replicasets,
		logger:         logger,
	}
}

func (ph *PodsHandler) Handle(obj any) error {
	pod, ok := obj.(*coreV1.Pod)
	if !ok {
		return ReceivedWrongType(obj, "Pod")
	}

	logger := ph.logger.With().Str("namespace", pod.GetNamespace()).Str("name", pod.GetName()).Logger()

	tenantID, applicationID, environmentName, microserviceID, ok := GetMicroserviceIdentifiers(pod.ObjectMeta)
	if !ok {
		logger.Trace().Msg("Skipping pod because it is missing microservice identifiers")
		return nil
	}

	var tenantsConfName, dolittleConfName, filesConfName, envConfName, envSecName string
	var hasTenantsConf, hasDolittleConf, hasFilesConf, hasEnvConf, hasEnvSec bool

	for _, volume := range pod.Spec.Volumes {
		if volume.Name == "tenants-config" && volume.ConfigMap != nil {
			tenantsConfName = volume.ConfigMap.Name
			hasTenantsConf = true
		}
		if volume.Name == "dolittle-config" && volume.ConfigMap != nil {
			dolittleConfName = volume.ConfigMap.Name
			hasDolittleConf = true
		}
		if volume.Name == "config-files" && volume.ConfigMap != nil {
			filesConfName = volume.ConfigMap.Name
			hasFilesConf = true
		}
	}
	for _, container := range pod.Spec.Containers {
		if container.Name == "head" {
			for _, source := range container.EnvFrom {
				if source.ConfigMapRef != nil && strings.HasSuffix(source.ConfigMapRef.Name, "-env-variables") {
					envConfName = source.ConfigMapRef.Name
					hasEnvConf = true
				}
				if source.SecretRef != nil && strings.HasSuffix(source.SecretRef.Name, "-secret-env-variables") {
					envSecName = source.SecretRef.Name
					hasEnvSec = true
				}
			}
		}
	}

	if !hasTenantsConf || !hasDolittleConf || !hasFilesConf || !hasEnvConf || !hasEnvSec {
		logger.Trace().Msg("Skipping pod because it is missing configuration references")
		return nil
	}

	tenantsConfig, err := ph.configmaps.ConfigMaps(pod.GetNamespace()).Get(tenantsConfName)
	if err != nil {
		return err
	}
	dolittleConfig, err := ph.configmaps.ConfigMaps(pod.GetNamespace()).Get(dolittleConfName)
	if err != nil {
		return err
	}
	filesConfig, err := ph.configmaps.ConfigMaps(pod.GetNamespace()).Get(filesConfName)
	if err != nil {
		return err
	}
	envConfig, err := ph.configmaps.ConfigMaps(pod.GetNamespace()).Get(envConfName)
	if err != nil {
		return err
	}
	envSecret, err := ph.secrets.Secrets(pod.GetNamespace()).Get(envSecName)
	if err != nil {
		return err
	}

	runtimeConfigHasher := kubernetes.NewConfigHasher()
	runtimeConfigHasher.WriteConfigMap(tenantsConfig)
	runtimeConfigHasher.WriteConfigMap(dolittleConfig)

	customerConfigHasher := kubernetes.NewConfigHasher()
	customerConfigHasher.WriteConfigMap(filesConfig)
	customerConfigHasher.WriteConfigMap(envConfig)
	customerConfigHasher.WriteSecret(envSecret)

	owners, err := ph.replicasets.GetPodReplicaSets(pod)
	if err != nil {
		return err
	}
	if len(owners) != 1 {
		return PodOwnerNotFound
	}

	runtimeConfig := entities.RuntimeConfiguration{
		ContentHash:                runtimeConfigHasher.GetComputedHash(),
		ConfigForArtifactID:        microserviceID,
		DeployedInEnvironmentName:  environmentName,
		EnvironmentOfApplicationID: applicationID,
		OwnedByCustomerID:          tenantID,
	}
	if err := ph.configurations.SetRuntime(runtimeConfig); err != nil {
		return err
	}
	logger.Debug().Interface("config", runtimeConfig).Msg("Updated runtime configuration")

	customerConfig := entities.CustomerConfiguration{
		ContentHash:                customerConfigHasher.GetComputedHash(),
		ConfigForArtifactID:        microserviceID,
		DeployedInEnvironmentName:  environmentName,
		EnvironmentOfApplicationID: applicationID,
		OwnedByCustomerID:          tenantID,
	}
	if err := ph.configurations.SetCustomer(customerConfig); err != nil {
		return err
	}
	logger.Debug().Interface("config", customerConfig).Msg("Updated customer configuration")

	instance := entities.DeploymentInstance{
		ID:                            string(pod.GetUID()),
		Started:                       pod.GetCreationTimestamp().UTC(),
		InstanceOfDeploymentID:        fmt.Sprintf("%v", owners[0].GetGeneration()),
		DeploymentOfArtifactID:        microserviceID,
		DeployedInEnvironmentName:     environmentName,
		EnvironmentOfApplicationID:    applicationID,
		OwnedByCustomerID:             tenantID,
		UsesArtifactConfigurationHash: customerConfigHasher.GetComputedHash(),
		UsesRuntimeConfigurationHash:  runtimeConfigHasher.GetComputedHash(),
	}
	if err := ph.deployments.SetInstance(instance); err != nil {
		return err
	}
	logger.Debug().Interface("instance", instance).Msg("Updated deployment instance")

	return nil
}
