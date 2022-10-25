/*
 * Copyright (c) Dolittle. All rights reserved.
 * Licensed under the MIT license. See LICENSE file in the project root for full license information.
 */

package neo4j

import (
	"context"
	"dolittle.io/fleet-observer/entities"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Configurations struct {
	session neo4j.SessionWithContext
	ctx     context.Context
}

func NewConfigurations(session neo4j.SessionWithContext, ctx context.Context) *Configurations {
	return &Configurations{
		session: session,
		ctx:     ctx,
	}
}

func (c *Configurations) SetArtifact(config entities.ArtifactConfiguration) error {
	return nil
}

func (c *Configurations) ListArtifacts() ([]entities.ArtifactConfiguration, error) {
	return nil, nil
}

func (c *Configurations) SetRuntime(config entities.RuntimeConfiguration) error {
	return nil
}

func (c *Configurations) ListRuntimes() ([]entities.RuntimeConfiguration, error) {
	return nil, nil
}
