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

type Applications struct {
	collection *mongo.Collection
	ctx        context.Context
}

func NewApplications(database *mongo.Database, ctx context.Context) *Applications {
	return &Applications{
		collection: database.Collection("applications"),
		ctx:        ctx,
	}
}

func (a *Applications) Set(application entities.Application) error {
	update := bson.D{{"$set", application}}
	_, err := a.collection.UpdateByID(a.ctx, application.UID, update, options.Update().SetUpsert(true))
	return err
}

func (a *Applications) Get(id entities.ApplicationUID) (*entities.Application, bool, error) {
	result := a.collection.FindOne(a.ctx, bson.D{{"_id", id}})
	err := result.Err()
	if err == mongo.ErrNoDocuments {
		return nil, false, nil
	} else if err != nil {
		return nil, true, err
	}

	application := &entities.Application{}
	err = result.Decode(application)
	if err != nil {
		return nil, true, err
	}

	return application, true, nil
}

func (a *Applications) List() ([]entities.Application, error) {
	cursor, err := a.collection.Find(a.ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	var applications []entities.Application
	if err := cursor.All(a.ctx, &applications); err != nil {
		return nil, err
	}

	return applications, cursor.Close(a.ctx)
}
