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

type Configurations struct {
	artifactCollection *mongo.Collection
	runtimeCollection  *mongo.Collection
	ctx                context.Context
}

func NewConfigurations(database *mongo.Database, ctx context.Context) *Configurations {
	return &Configurations{
		artifactCollection: database.Collection("artifact-configurations"),
		runtimeCollection:  database.Collection("runtime-configurations"),
		ctx:                ctx,
	}
}

func (c *Configurations) SetArtifact(config entities.ArtifactConfiguration) error {
	update := bson.D{{"$set", config}}
	_, err := c.artifactCollection.UpdateByID(c.ctx, config.UID, update, options.Update().SetUpsert(true))
	return err
}

func (c *Configurations) ListArtifacts() ([]entities.ArtifactConfiguration, error) {
	cursor, err := c.artifactCollection.Find(c.ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	var configurations []entities.ArtifactConfiguration
	if err := cursor.All(c.ctx, &configurations); err != nil {
		return nil, err
	}

	return configurations, cursor.Close(c.ctx)
}

func (c *Configurations) SetRuntime(config entities.RuntimeConfiguration) error {
	update := bson.D{{"$set", config}}
	_, err := c.runtimeCollection.UpdateByID(c.ctx, config.UID, update, options.Update().SetUpsert(true))
	return err
}

func (c *Configurations) ListRuntimes() ([]entities.RuntimeConfiguration, error) {
	cursor, err := c.runtimeCollection.Find(c.ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	var configurations []entities.RuntimeConfiguration
	if err := cursor.All(c.ctx, &configurations); err != nil {
		return nil, err
	}

	return configurations, cursor.Close(c.ctx)
}
