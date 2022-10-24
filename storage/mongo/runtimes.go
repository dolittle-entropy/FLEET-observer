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

type Runtimes struct {
	versionsCollection *mongo.Collection
	ctx                context.Context
}

func NewRuntimes(database *mongo.Database, ctx context.Context) *Runtimes {
	return &Runtimes{
		versionsCollection: database.Collection("runtime-versions"),
		ctx:                ctx,
	}
}

func (r *Runtimes) SetVersion(version entities.RuntimeVersion) error {
	update := bson.D{{"$set", version}}
	_, err := r.versionsCollection.UpdateByID(r.ctx, version.UID, update, options.Update().SetUpsert(true))
	return err
}

func (r *Runtimes) ListVersions() ([]entities.RuntimeVersion, error) {
	cursor, err := r.versionsCollection.Find(r.ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	var versions []entities.RuntimeVersion
	if err := cursor.All(r.ctx, &versions); err != nil {
		return nil, err
	}

	return versions, cursor.Close(r.ctx)
}
