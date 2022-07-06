package observing

import (
	"dolittle.io/fleet-observer/entities"
	"dolittle.io/fleet-observer/mongo"
	"fmt"
	"github.com/rs/zerolog"
	coreV1 "k8s.io/api/core/v1"
	listersAppsV1 "k8s.io/client-go/listers/apps/v1"
	listersCoreV1 "k8s.io/client-go/listers/core/v1"
)

type PodsHandler struct {
	deployments *mongo.Deployments
	configmaps  listersCoreV1.ConfigMapLister
	secrets     listersCoreV1.SecretLister
	replicasets listersAppsV1.ReplicaSetLister
	logger      zerolog.Logger
}

func NewPodsHandler(deployments *mongo.Deployments, configmaps listersCoreV1.ConfigMapLister, secrets listersCoreV1.SecretLister, replicasets listersAppsV1.ReplicaSetLister, logger zerolog.Logger) *PodsHandler {
	return &PodsHandler{
		deployments: deployments,
		configmaps:  configmaps,
		secrets:     secrets,
		replicasets: replicasets,
		logger:      logger,
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

	customerConfigHash, err := ComputeCustomerConfigHashFor(pod.ObjectMeta, ph.configmaps, ph.secrets)
	if err != nil {
		return err
	}

	runtimeConfigHash, err := ComputeRuntimeConfigHashFor(pod.ObjectMeta, ph.configmaps)
	if err != nil {
		return err
	}

	owners, err := ph.replicasets.GetPodReplicaSets(pod)
	if err != nil {
		return err
	}
	if len(owners) != 1 {
		return PodOwnerNotFound
	}

	instance := entities.DeploymentInstance{
		ID:                            string(pod.GetUID()),
		Started:                       pod.GetCreationTimestamp().UTC(),
		InstanceOfDeploymentID:        fmt.Sprintf("%v", owners[0].GetGeneration()),
		DeploymentOfArtifactID:        microserviceID,
		DeployedInEnvironmentName:     environmentName,
		EnvironmentOfApplicationID:    applicationID,
		OwnedByCustomerID:             tenantID,
		UsesArtifactConfigurationHash: customerConfigHash,
		UsesRuntimeConfigurationHash:  runtimeConfigHash,
	}
	if err := ph.deployments.SetInstance(instance); err != nil {
		return err
	}
	logger.Debug().Interface("instance", instance).Msg("Updated deployment instance")

	return nil
}
