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
	return multiUpdate(
		e.session,
		e.ctx,
		map[string]any{
			"uid":                  environment.UID,
			"name":                 environment.Properties.Name,
			"link_application_uid": environment.Links.EnvironmentOfApplicationUID,
		},
		`
			MERGE (environment:Environment { _uid: $uid })
			SET environment = { _uid: $uid, name: $name }
			RETURN id(environment)
		`,
		`
			MATCH (environment:Environment { _uid: $uid })
			WITH environment
				MERGE (application:Application { _uid: $link_application_uid })
				WITH environment, application
					MERGE (environment)-[:EnvironmentOf]->(application)
					WITH environment, application
						MATCH (environment)-[r:EnvironmentOf]->(other)
						WHERE other._uid <> application._uid
						DELETE r
			RETURN id(environment)
		`)
}

func (e *Environments) Get(id entities.EnvironmentUID) (*entities.Environment, bool, error) {
	environment := &entities.Environment{}
	found, err := findSingleJson(
		e.session,
		e.ctx,
		map[string]any{
			"uid": id,
		},
		`
			MATCH (environment:Environment { _uid: $uid })-[:EnvironmentOf]->(application:Application)
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
			RETURN apoc.convert.toJson(entry) as json
		`,
		environment)
	return environment, found, err
}

func (e *Environments) List() ([]entities.Environment, error) {
	var environments []entities.Environment
	return environments, findAllJson(
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
