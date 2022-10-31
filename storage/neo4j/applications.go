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

type Applications struct {
	session neo4j.SessionWithContext
	ctx     context.Context
}

func NewApplications(session neo4j.SessionWithContext, ctx context.Context) *Applications {
	return &Applications{
		session: session,
		ctx:     ctx,
	}
}

func (a *Applications) Set(application entities.Application) error {
	return multiUpdate(
		a.session,
		a.ctx,
		map[string]any{
			"uid":               application.UID,
			"id":                application.Properties.ID,
			"name":              application.Properties.Name,
			"link_customer_uid": application.Links.OwnedByCustomerUID,
		},
		`
			MERGE (application:Application { _uid: $uid })
			SET application = { _uid: $uid, id: $id, name: $name }
			RETURN id(application)
		`, `
			MATCH (application:Application { _uid: $uid })
			WITH application
				MERGE (customer:Customer { _uid: $link_customer_uid})
				WITH application, customer
					MERGE (application)-[:OwnedBy]->(customer)
					WITH application, customer
						MATCH (application)-[r:OwnedBy]->(other)
						WHERE other._uid <> customer._uid
						DELETE r
			RETURN id(application)
		`)
}

func (a *Applications) Get(id entities.ApplicationUID) (*entities.Application, bool, error) {
	application := &entities.Application{}
	found, err := findSingleJson(
		a.session,
		a.ctx,
		map[string]any{
			"uid": id,
		},
		`
			MATCH (application:Application { _uid: $uid })-[:OwnedBy]->(customer:Customer)
			WITH {
				uid: application._uid,
				type: "Application",
				properties: {
					id: application.id,
					name: application.name
				},
				links: {
					ownedBy: customer._uid
				}
			} as entry
			RETURN apoc.convert.toJson(entry) as json
		`,
		application)
	return application, found, err
}

func (a *Applications) List() ([]entities.Application, error) {
	var applications []entities.Application
	return applications, findAllJson(
		a.session,
		a.ctx,
		`
			MATCH (application:Application)-[:OwnedBy]->(customer:Customer)
			WITH {
				uid: application._uid,
				type: "Application",
				properties: {
					id: application.id,
					name: application.name
				},
				links: {
					ownedBy: customer._uid
				}
			} as entry
			RETURN apoc.convert.toJson(collect(entry)) as json
		`,
		&applications)
}
