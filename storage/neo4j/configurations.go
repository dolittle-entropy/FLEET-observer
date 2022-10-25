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
	return runUpdate(
		c.session,
		c.ctx,
		`
			MERGE (config:ArtifactConfiguration { _uid: $uid })
			SET config = { _uid: $uid, hash: $hash }
			RETURN id(config)
		`,
		map[string]any{
			"uid":  config.UID,
			"hash": config.Properties.ContentHash,
		})
}

func (c *Configurations) ListArtifacts() ([]entities.ArtifactConfiguration, error) {
	var configs []entities.ArtifactConfiguration
	return configs, querySingleJson(
		c.session,
		c.ctx,
		`
			MATCH (config:ArtifactConfiguration)
			WITH {
				uid: config._uid,
				type: "ArtifactConfiguration",
				properties: {
					hash: config.hash
				}
			} as entry
			RETURN apoc.convert.toJson(collect(entry)) as json
		`,
		&configs)
}

func (c *Configurations) SetRuntime(config entities.RuntimeConfiguration) error {
	return runUpdate(
		c.session,
		c.ctx,
		`
			MERGE (config:RuntimeConfiguration { _uid: $uid })
			SET config = { _uid: $uid, hash: $hash }
			RETURN id(config)
		`,
		map[string]any{
			"uid":  config.UID,
			"hash": config.Properties.ContentHash,
		})
	return nil
}

func (c *Configurations) ListRuntimes() ([]entities.RuntimeConfiguration, error) {
	var configs []entities.RuntimeConfiguration
	return configs, querySingleJson(
		c.session,
		c.ctx,
		`
			MATCH (config:RuntimeConfiguration)
			WITH {
				uid: config._uid,
				type: "RuntimeConfiguration",
				properties: {
					hash: config.hash
				}
			} as entry
			RETURN apoc.convert.toJson(collect(entry)) as json
		`,
		&configs)
}
