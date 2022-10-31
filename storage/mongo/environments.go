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

type Environments struct {
	collection *mongo.Collection
	ctx        context.Context
}

func NewEnvironments(database *mongo.Database, ctx context.Context) *Environments {
	return &Environments{
		collection: database.Collection("environments"),
		ctx:        ctx,
	}
}

func (e *Environments) Set(environment entities.Environment) error {
	update := bson.D{{"$set", environment}}
	_, err := e.collection.UpdateByID(e.ctx, environment.UID, update, options.Update().SetUpsert(true))
	return err
}

func (e *Environments) Get(id entities.EnvironmentUID) (*entities.Environment, bool, error) {
	result := e.collection.FindOne(e.ctx, bson.D{{"_id", id}})
	err := result.Err()
	if err == mongo.ErrNoDocuments {
		return nil, false, nil
	} else if err != nil {
		return nil, true, err
	}

	environment := &entities.Environment{}
	err = result.Decode(environment)
	if err != nil {
		return nil, true, err
	}

	return environment, true, nil
}

func (e *Environments) List() ([]entities.Environment, error) {
	cursor, err := e.collection.Find(e.ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	var environments []entities.Environment
	if err := cursor.All(e.ctx, &environments); err != nil {
		return nil, err
	}

	return environments, cursor.Close(e.ctx)
}
