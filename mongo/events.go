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

type Events struct {
	collection *mongo.Collection
	ctx        context.Context
}

func NewEvents(database *mongo.Database, ctx context.Context) *Events {
	return &Events{
		collection: database.Collection("events"),
		ctx:        ctx,
	}
}

func (e *Events) Set(event entities.Event) error {
	update := bson.D{{"$set", event}}
	_, err := e.collection.UpdateByID(e.ctx, event.UID, update, options.Update().SetUpsert(true))
	return err
}

func (e *Events) List() ([]entities.Event, error) {
	cursor, err := e.collection.Find(e.ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	var events []entities.Event
	if err := cursor.All(e.ctx, &events); err != nil {
		return nil, err
	}

	return events, cursor.Close(e.ctx)
}
