/*
 * Copyright (c) Dolittle. All rights reserved.
 * Licensed under the MIT license. See LICENSE file in the project root for full license information.
 */

package neo4j

import (
	"context"
	"github.com/knadh/koanf"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/rs/zerolog"
)

func ConnectToNeo4j(config *koanf.Koanf, logger zerolog.Logger, ctx context.Context) (neo4j.SessionWithContext, error) {
	connectionString := config.String("neo4j.connection-string")
	driver, err := neo4j.NewDriverWithContext(connectionString, neo4j.NoAuth())
	if err != nil {
		return nil, err
	}

	logger = logger.With().Str("component", "neo4j").Logger()
	logger.Debug().Str("connection-string", connectionString).Msg("Connecting to Neo4j")

	err = driver.VerifyConnectivity(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to connect to Neo4j")
		return nil, err
	}

	logger.Info().Msg("Connected to Neo4j")

	return driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite}), nil
}
