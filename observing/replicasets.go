package observing

import (
	"dolittle.io/fleet-observer/entities"
	"dolittle.io/fleet-observer/mongo"
	"fmt"
	"github.com/rs/zerolog"
	appsV1 "k8s.io/api/apps/v1"
	coreV1 "k8s.io/api/core/v1"
	"regexp"
	"strconv"
	"time"
)

type ReplicasetHandler struct {
	environments *mongo.Environments
	artifacts    *mongo.Artifacts
	runtimes     *mongo.Runtimes
	deployments  *mongo.Deployments
	logger       zerolog.Logger
}

func NewReplicasetHandler(environments *mongo.Environments, artifacts *mongo.Artifacts, runtimes *mongo.Runtimes, deployments *mongo.Deployments, logger zerolog.Logger) *ReplicasetHandler {
	return &ReplicasetHandler{
		environments: environments,
		artifacts:    artifacts,
		runtimes:     runtimes,
		deployments:  deployments,
		logger:       logger.With().Str("handler", "replicasets").Logger(),
	}
}

func (rh *ReplicasetHandler) Handle(obj any) error {
	replicaset, ok := obj.(*appsV1.ReplicaSet)
	if !ok {
		return ReceivedWrongType(obj, "ReplicaSet")
	}

	logger := rh.logger.With().Str("namespace", replicaset.GetNamespace()).Str("name", replicaset.GetName()).Logger()

	// -- Get all the data --
	tenantID, applicationID, environmentName, microserviceID, ok := GetMicroserviceIdentifiers(replicaset.ObjectMeta)
	if !ok {
		logger.Trace().Msg("Skipping replicaset because it is missing microservice identifiers")
		return nil
	}

	var headContainer, runtimeContainer coreV1.Container
	var hasHeadContainer, hasRuntimeContainer = false, false
	for _, container := range replicaset.Spec.Template.Spec.Containers {
		if container.Name == "head" {
			headContainer = container
			hasHeadContainer = true
		}
		if container.Name == "runtime" {
			runtimeContainer = container
			hasRuntimeContainer = true
		}
	}
	if !hasHeadContainer {
		logger.Trace().Msg("Skipping replicaset because it does not have a head container")
		return nil
	}
	if !hasRuntimeContainer {
		logger.Trace().Msg("Skipping replicaset because it does not have a runtime container")
		return nil
	}
	artifactVersionName := getArtifactVersionName(headContainer)
	runtimeVersion, err := parseRuntimeVersion(runtimeContainer)
	if err != nil {
		return err
	}

	// -- Set all the entities --
	environment := entities.NewEnvironment(tenantID, applicationID, environmentName)
	if err := rh.environments.Set(environment); err != nil {
		return err
	}
	logger.Debug().Interface("environment", environment).Msg("Updated environment")

	artifact := entities.NewArtifact(tenantID, microserviceID)
	if err := rh.artifacts.Set(artifact); err != nil {
		return err
	}
	logger.Debug().Interface("artifact", artifact).Msg("Updated artifact")

	artifactVersion := entities.NewArtifactVersion(tenantID, microserviceID, artifactVersionName, time.Time{})
	if err := rh.artifacts.SetVersion(artifactVersion); err != nil {
		return err
	}
	logger.Debug().Interface("version", artifactVersion).Msg("Updated artifact version")

	if err := rh.runtimes.SetVersion(runtimeVersion); err != nil {
		return err
	}
	logger.Debug().Interface("version", runtimeVersion).Msg("Updated runtime version")

	deployment := entities.NewDeployment(
		tenantID,
		applicationID,
		environmentName,
		fmt.Sprintf("%v", replicaset.GetGeneration()),
		replicaset.GetCreationTimestamp().UTC(),
		artifactVersion,
		runtimeVersion,
	)
	if err := rh.deployments.Set(deployment); err != nil {
		return err
	}
	logger.Debug().Interface("deployment", deployment).Msg("Updated deployment")

	return nil
}

var containerNameExpression = regexp.MustCompile(`^([A-Za-z0-9]+\.azurecr\.io/)?(.+)$`)

func getArtifactVersionName(headContainer coreV1.Container) string {
	name := headContainer.Image
	matches := containerNameExpression.FindStringSubmatch(headContainer.Image)
	if len(matches) > 2 && len(matches[2]) > 0 {
		name = matches[2]
	}
	return name
}

var runtimeVersionExpression = regexp.MustCompile(`^dolittle/runtime:(\d+)\.(\d+)\.(\d+)(-(.+))?$`)

func parseRuntimeVersion(runtimeContainer coreV1.Container) (entities.RuntimeVersion, error) {
	matches := runtimeVersionExpression.FindStringSubmatch(runtimeContainer.Image)
	if len(matches) < 6 {
		return entities.RuntimeVersion{}, FailedToParseRuntimeVersion(runtimeContainer.Image)
	}
	return entities.NewRuntimeVersion(
		mustParseInt(matches[1]),
		mustParseInt(matches[2]),
		mustParseInt(matches[3]),
		matches[5],
		time.Time{},
	), nil
}

func mustParseInt(value string) int {
	parsed, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		panic(err)
	}
	return int(parsed)
}
