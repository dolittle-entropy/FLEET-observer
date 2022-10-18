/*
 * Copyright (c) Dolittle. All rights reserved.
 * Licensed under the MIT license. See LICENSE file in the project root for full license information.
 */

package mongo

import (
	"context"
	"dolittle.io/fleet-observer/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Nodes struct {
	collection *mongo.Collection
	ctx        context.Context
}

func NewNodes(database *mongo.Database, ctx context.Context) *Nodes {
	return &Nodes{
		collection: database.Collection("nodes"),
		ctx:        ctx,
	}
}

func (n *Nodes) Set(node entities.Node) error {
	update := bson.D{{"$set", node}}
	_, err := n.collection.UpdateByID(n.ctx, node.UID, update, options.Update().SetUpsert(true))
	return err
}

func (n *Nodes) List() ([]entities.Node, error) {
	cursor, err := n.collection.Find(n.ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	var nodes []entities.Node
	if err := cursor.All(n.ctx, &nodes); err != nil {
		return nil, err
	}

	return nodes, cursor.Close(n.ctx)
}
