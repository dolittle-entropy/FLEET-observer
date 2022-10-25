/*
 * Copyright (c) Dolittle. All rights reserved.
 * Licensed under the MIT license. See LICENSE file in the project root for full license information.
 */

package exporting

import (
	"context"
	"dolittle.io/fleet-observer/entities"
	"dolittle.io/fleet-observer/storage"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog"
	"os"
)

type Exporter struct {
	repositories *storage.Repositories
	logger       zerolog.Logger
	ctx          context.Context
}

func NewExporter(repositories *storage.Repositories, logger zerolog.Logger, ctx context.Context) *Exporter {
	return &Exporter{
		repositories: repositories,
		logger:       logger,
		ctx:          ctx,
	}
}

func (e *Exporter) ExportToFile(path string) error {
	output, err := os.Create(path)
	if err != nil {
		e.logger.Error().Str("output", path).Err(err).Msg("Could not open output file")
		return err
	}
	defer output.Close()

	e.logger.Info().Str("output", path).Msg("Starting export to file")

	var data []any

	nodes, err := e.repositories.Nodes.List()
	if err != nil {
		e.logger.Error().Err(err).Msg("Failed to get nodes")
		return err
	}
	for _, node := range nodes {
		node.UID = entities.NodeUID(fmt.Sprintf("%v:%v", entities.NodeType, node.UID))
		data = append(data, node)
	}

	customers, err := e.repositories.Customers.List()
	if err != nil {
		e.logger.Error().Err(err).Msg("Failed to get customers")
		return err
	}
	for _, customer := range customers {
		customer.UID = entities.CustomerUID(fmt.Sprintf("%v:%v", entities.CustomerType, customer.UID))
		data = append(data, customer)
	}

	applications, err := e.repositories.Applications.List()
	if err != nil {
		e.logger.Error().Err(err).Msg("Failed to get applications")
		return err
	}
	for _, application := range applications {
		application.UID = entities.ApplicationUID(fmt.Sprintf("%v:%v", entities.ApplicationType, application.UID))
		application.Links.OwnedByCustomerUID = entities.CustomerUID(fmt.Sprintf("%v:%v", entities.CustomerType, application.Links.OwnedByCustomerUID))
		data = append(data, application)
	}

	environments, err := e.repositories.Environments.List()
	if err != nil {
		e.logger.Error().Err(err).Msg("Failed to get environments")
		return err
	}
	for _, environment := range environments {
		environment.UID = entities.EnvironmentUID(fmt.Sprintf("%v:%v", entities.EnvironmentType, environment.UID))
		environment.Links.EnvironmentOfApplicationUID = entities.ApplicationUID(fmt.Sprintf("%v:%v", entities.ApplicationType, environment.Links.EnvironmentOfApplicationUID))
		data = append(data, environment)
	}

	artifacts, err := e.repositories.Artifacts.List()
	if err != nil {
		e.logger.Error().Err(err).Msg("Failed to get artifacts")
		return err
	}
	for _, artifact := range artifacts {
		artifact.UID = entities.ArtifactUID(fmt.Sprintf("%v:%v", entities.ArtifactType, artifact.UID))
		artifact.Links.DevelopedByCustomerUID = entities.CustomerUID(fmt.Sprintf("%v:%v", entities.CustomerType, artifact.Links.DevelopedByCustomerUID))
		data = append(data, artifact)
	}

	artifactVersions, err := e.repositories.Artifacts.ListVersions()
	if err != nil {
		e.logger.Error().Err(err).Msg("Failed to get artifact versions")
		return err
	}
	for _, version := range artifactVersions {
		version.UID = entities.ArtifactVersionUID(fmt.Sprintf("%v:%v", entities.ArtifactVersionType, version.UID))
		version.Links.VersionOfArtifactUID = entities.ArtifactUID(fmt.Sprintf("%v:%v", entities.ArtifactType, version.Links.VersionOfArtifactUID))
		data = append(data, version)
	}

	runtimeVersions, err := e.repositories.Runtimes.ListVersions()
	if err != nil {
		e.logger.Error().Err(err).Msg("Failed to get runtime versions")
		return err
	}
	for _, version := range runtimeVersions {
		version.UID = entities.RuntimeVersionUID(fmt.Sprintf("%v:%v", entities.RuntimeVersionType, version.UID))
		data = append(data, version)
	}

	deployments, err := e.repositories.Deployments.List()
	if err != nil {
		e.logger.Error().Err(err).Msg("Failed to get deployments")
		return err
	}
	for _, deployment := range deployments {
		deployment.UID = entities.DeploymentUID(fmt.Sprintf("%v:%v", entities.DeploymentType, deployment.UID))
		deployment.Links.DeployedInEnvironmentUID = entities.EnvironmentUID(fmt.Sprintf("%v:%v", entities.EnvironmentType, deployment.Links.DeployedInEnvironmentUID))
		deployment.Links.UsesArtifactVersionUID = entities.ArtifactVersionUID(fmt.Sprintf("%v:%v", entities.ArtifactVersionType, deployment.Links.UsesArtifactVersionUID))
		deployment.Links.UsesRuntimeVersionUID = entities.RuntimeVersionUID(fmt.Sprintf("%v:%v", entities.RuntimeVersionType, deployment.Links.UsesRuntimeVersionUID))
		data = append(data, deployment)
	}

	artifactConfigs, err := e.repositories.Configurations.ListArtifacts()
	if err != nil {
		e.logger.Error().Err(err).Msg("Failed to get artifact configurations")
		return err
	}
	for _, config := range artifactConfigs {
		config.UID = entities.ArtifactConfigurationUID(fmt.Sprintf("%v:%v", entities.ArtifactConfigurationType, config.UID))
		data = append(data, config)
	}

	runtimeConfigs, err := e.repositories.Configurations.ListRuntimes()
	if err != nil {
		e.logger.Error().Err(err).Msg("Failed to get runtime configurations")
		return err
	}
	for _, config := range runtimeConfigs {
		config.UID = entities.RuntimeConfigurationUID(fmt.Sprintf("%v:%v", entities.RuntimeConfigurationType, config.UID))
		data = append(data, config)
	}

	instances, err := e.repositories.Deployments.ListInstances()
	if err != nil {
		e.logger.Error().Err(err).Msg("Failed to get deployment instances")
		return err
	}
	for _, instance := range instances {
		instance.UID = entities.DeploymentInstanceUID(fmt.Sprintf("%v:%v", entities.DeploymentInstanceType, instance.UID))
		instance.Links.InstanceOfDeploymentUID = entities.DeploymentUID(fmt.Sprintf("%v:%v", entities.DeploymentType, instance.Links.InstanceOfDeploymentUID))
		instance.Links.UsesArtifactConfigurationUID = entities.ArtifactConfigurationUID(fmt.Sprintf("%v:%v", entities.ArtifactConfigurationType, instance.Links.UsesArtifactConfigurationUID))
		instance.Links.UsesRuntimeConfigurationUID = entities.RuntimeConfigurationUID(fmt.Sprintf("%v:%v", entities.RuntimeConfigurationType, instance.Links.UsesRuntimeConfigurationUID))
		instance.Links.ScheduledOnNodeUID = entities.NodeUID(fmt.Sprintf("%v:%v", entities.NodeType, instance.Links.ScheduledOnNodeUID))
		data = append(data, instance)
	}

	events, err := e.repositories.Events.List()
	if err != nil {
		e.logger.Error().Err(err).Msg("Failed to get events")
		return err
	}
	for _, event := range events {
		event.UID = entities.EventUID(fmt.Sprintf("%v:%v", event.Type, event.UID))
		event.Links.HappenedToDeploymentInstanceUID = entities.DeploymentInstanceUID(fmt.Sprintf("%v:%v", entities.DeploymentInstanceType, event.Links.HappenedToDeploymentInstanceUID))
		data = append(data, event)
	}

	e.logger.Info().Msg("Writing to file...")
	for _, entry := range data {
		encoded, err := json.Marshal(entry)
		if err != nil {
			e.logger.Error().Err(err).Msg("Failed to convert entry to JSON")
			return err
		}

		_, err = fmt.Fprintln(output, string(encoded))
		if err != nil {
			e.logger.Error().Err(err).Msg("Could not write to file")
			return err
		}
	}

	e.logger.Info().Int("entries", len(data)).Msg("Done exporting entries!")
	return nil
}
