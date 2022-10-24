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
		_, err := neo4j.ConnectToNeo4j(config, logger, ctx)
		if err != nil {
			return nil, err
		}

		return &Repositories{}, nil
	}

	if config.String("mongodb.connection-string") != "" {
		logger.Info().Msg("Using MongoDB for storage")
		_, err := mongo.ConnectToMongo(config, logger, ctx)
		if err != nil {
			return nil, err
		}

		return &Repositories{}, nil
	}

	return nil, ErrNoStorageConfigured
}
