/*
 * Copyright (c) Dolittle. All rights reserved.
 * Licensed under the MIT license. See LICENSE file in the project root for full license information.
 */

package storage

import (
	"context"
	"dolittle.io/fleet-observer/storage/mongo"
	"dolittle.io/fleet-observer/storage/neo4j"
	"errors"
	"github.com/knadh/koanf"
	"github.com/rs/zerolog"
)

var (
	ErrNoStorageConfigured = errors.New("no storage configured")
)

func Connect(config *koanf.Koanf, logger zerolog.Logger, ctx context.Context) (*Repositories, error) {
	logger = logger.With().Str("component", "storage").Logger()
	if config.String("neo4j.connection-string") != "" {
		logger.Info().Msg("Using Neo4j for storage")
		session, err := neo4j.ConnectToNeo4j(config, logger, ctx)
		if err != nil {
			return nil, err
		}

		return &Repositories{
			Nodes:          neo4j.NewNodes(session, ctx),
			Customers:      neo4j.NewCustomers(session, ctx),
			Applications:   neo4j.NewApplications(session, ctx),
			Environments:   neo4j.NewEnvironments(session, ctx),
			Artifacts:      neo4j.NewArtifacts(session, ctx),
			Runtimes:       neo4j.NewRuntimes(session, ctx),
			Deployments:    neo4j.NewDeployments(session, ctx),
			Configurations: neo4j.NewConfigurations(session, ctx),
			Events:         neo4j.NewEvents(session, ctx),
		}, nil
	}

	if config.String("mongodb.connection-string") != "" {
		logger.Info().Msg("Using MongoDB for storage")
		database, err := mongo.ConnectToMongo(config, logger, ctx)
		if err != nil {
			return nil, err
		}

		return &Repositories{
			Nodes:          mongo.NewNodes(database, ctx),
			Customers:      mongo.NewCustomers(database, ctx),
			Applications:   mongo.NewApplications(database, ctx),
			Environments:   mongo.NewEnvironments(database, ctx),
			Artifacts:      mongo.NewArtifacts(database, ctx),
			Runtimes:       mongo.NewRuntimes(database, ctx),
			Deployments:    mongo.NewDeployments(database, ctx),
			Configurations: mongo.NewConfigurations(database, ctx),
			Events:         mongo.NewEvents(database, ctx),
		}, nil
	}

	return nil, ErrNoStorageConfigured
}
