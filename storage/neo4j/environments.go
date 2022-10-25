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

type Environments struct {
	session neo4j.SessionWithContext
	ctx     context.Context
}

func NewEnvironments(session neo4j.SessionWithContext, ctx context.Context) *Environments {
	return &Environments{
		session: session,
		ctx:     ctx,
	}
}

func (e *Environments) Set(environment entities.Environment) error {
	return runUpdate(
		e.session,
		e.ctx,
		`
			MERGE (environment:Environment { _uid: $uid })
			SET environment = { _uid: $uid, name: $name }
			WITH environment
				MERGE (application:Application { _uid: $link_application_uid })
				WITH environment, application
					MERGE (environment)-[:EnvironmentOf]->(application)
					WITH environment, application
						MATCH (environment)-[r:EnvironmentOf]->(other)
						WHERE other._uid <> application._uid
						DELETE r
			RETURN id(environment)
		`,
		map[string]any{
			"uid":                  environment.UID,
			"name":                 environment.Properties.Name,
			"link_application_uid": environment.Links.EnvironmentOfApplicationUID,
		})
}

func (e *Environments) List() ([]entities.Environment, error) {
	var environments []entities.Environment
	return environments, querySingleJson(
		e.session,
		e.ctx,
		`
			MATCH (environment:Environment)-[:EnvironmentOf]->(application:Application)
			WITH {
				uid: environment._uid,
				type: "Environment",
				properties: {
					name: environment.name
				},
				links: {
					environmentOf: application._uid
				}
			} as entry
			RETURN apoc.convert.toJson(collect(entry)) as json
		`,
		&environments)
}
