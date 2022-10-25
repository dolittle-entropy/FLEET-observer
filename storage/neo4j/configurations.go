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
	return multiUpdate(
		c.session,
		c.ctx,
		map[string]any{
			"uid":  config.UID,
			"hash": config.Properties.ContentHash,
		},
		`
			MERGE (config:ArtifactConfiguration { _uid: $uid })
			SET config = { _uid: $uid, hash: $hash }
			RETURN id(config)
		`)
}

func (c *Configurations) ListArtifacts() ([]entities.ArtifactConfiguration, error) {
	var configs []entities.ArtifactConfiguration
	return configs, findAllJson(
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
	return multiUpdate(
		c.session,
		c.ctx,
		map[string]any{
			"uid":  config.UID,
			"hash": config.Properties.ContentHash,
		},
		`
			MERGE (config:RuntimeConfiguration { _uid: $uid })
			SET config = { _uid: $uid, hash: $hash }
			RETURN id(config)
		`)
}

func (c *Configurations) ListRuntimes() ([]entities.RuntimeConfiguration, error) {
	var configs []entities.RuntimeConfiguration
	return configs, findAllJson(
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
