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

type Customers struct {
	session neo4j.SessionWithContext
	ctx     context.Context
}

func NewCustomers(session neo4j.SessionWithContext, ctx context.Context) *Customers {
	return &Customers{
		session: session,
		ctx:     ctx,
	}
}

func (c *Customers) Set(customer entities.Customer) error {
	return runUpdate(
		c.session,
		c.ctx,
		`
			MERGE (customer:Customer { _uid: $uid })
			SET customer = { _uid: $uid, id: $id, name: $name }
			RETURN id(customer)
		`,
		map[string]any{
			"uid":  customer.UID,
			"id":   customer.Properties.ID,
			"name": customer.Properties.Name,
		})
}

func (c *Customers) List() ([]entities.Customer, error) {
	var customers []entities.Customer
	return customers, querySingleJson(
		c.session,
		c.ctx,
		`
			MATCH (customer:Customer)
			WITH {
				uid: customer._uid,
				type: "Customer",
				properties: {
					id: customer.id,
					name: customer.name
				}
			} as entry
			RETURN apoc.convert.toJson(collect(entry)) as json
		`,
		&customers)
}
