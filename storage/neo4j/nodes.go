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

type Nodes struct {
	session neo4j.SessionWithContext
	ctx     context.Context
}

func NewNodes(session neo4j.SessionWithContext, ctx context.Context) *Nodes {
	return &Nodes{
		session: session,
		ctx:     ctx,
	}
}

func (n *Nodes) Set(node entities.Node) error {
	return runUpdate(
		n.session,
		n.ctx,
		`
			MERGE (node:Node { _uid: $uid })
			SET node = { _uid: $uid, hostname: $hostname, image: $image, type: $type }
			RETURN id(node)
		`,
		map[string]any{
			"uid":      node.UID,
			"hostname": node.Properties.Hostname,
			"image":    node.Properties.Image,
			"type":     node.Properties.Type,
		})
}

func (n *Nodes) List() ([]entities.Node, error) {
	var nodes []entities.Node
	return nodes, querySingleJson(
		n.session,
		n.ctx,
		`
			MATCH (node:Node)
			WITH {
				uid: node._uid,
				type: "Node",
				properties: {
					hostname: node.hostname,
					image: node.image,
					type: node.type
				}
			} as entry
			RETURN apoc.convert.toJson(collect(entry)) as json
		`,
		&nodes)
}
