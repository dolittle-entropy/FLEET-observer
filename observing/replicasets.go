/*
 * Copyright (c) Dolittle. All rights reserved.
 * Licensed under the MIT license. See LICENSE file in the project root for full license information.
 */

package observing

import (
	"dolittle.io/fleet-observer/entities"
	"dolittle.io/fleet-observer/storage"
	"github.com/rs/zerolog"
	appsV1 "k8s.io/api/apps/v1"
	coreV1 "k8s.io/api/core/v1"
	"regexp"
	"strconv"
	"time"
)

type ReplicasetHandler struct {
	environments storage.Environments
	artifacts    storage.Artifacts
	runtimes     storage.Runtimes
	deployments  storage.Deployments
	logger       zerolog.Logger
}

func NewReplicasetHandler(environments storage.Environments, artifacts storage.Artifacts, runtimes storage.Runtimes, deployments storage.Deployments, logger zerolog.Logger) *ReplicasetHandler {
	return &ReplicasetHandler{
		environments: environments,
		artifacts:    artifacts,
		runtimes:     runtimes,
		deployments:  deployments,
		logger:       logger.With().Str("handler", "replicasets").Logger(),
	}
}

func (rh *ReplicasetHandler) Handle(obj any, _deleted bool) error {
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

	deploymentName, ok := replicaset.GetLabels()["microservice"]
	if !ok {
		logger.Trace().Msg("Skipping replicaset because it does not have a microservice label")
		return nil
	}

	revision, ok := replicaset.GetAnnotations()["deployment.kubernetes.io/revision"]
	if !ok {
		logger.Trace().Msg("Skipping replicaset because it does not have a revision annotation")
		return nil
	}

	runtimeContainer, headContainer, ok := getRuntimeAndHeadContainer(replicaset.Spec.Template.Spec)
	if !ok {
		logger.Trace().Msg("Skipping replicaset because it does not have a runtime and head container")
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
		revision,
		deploymentName,
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
