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

type Artifacts struct {
	collection         *mongo.Collection
	versionsCollection *mongo.Collection
	ctx                context.Context
}

func NewArtifacts(database *mongo.Database, ctx context.Context) *Artifacts {
	return &Artifacts{
		collection:         database.Collection("artifacts"),
		versionsCollection: database.Collection("artifact-versions"),
		ctx:                ctx,
	}
}

func (a *Artifacts) Set(artifact entities.Artifact) error {
	update := bson.D{{"$set", artifact}}
	_, err := a.collection.UpdateByID(a.ctx, artifact.UID, update, options.Update().SetUpsert(true))
	return err
}

func (a *Artifacts) List() ([]entities.Artifact, error) {
	cursor, err := a.collection.Find(a.ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	var artifacts []entities.Artifact
	if err := cursor.All(a.ctx, &artifacts); err != nil {
		return nil, err
	}

	return artifacts, cursor.Close(a.ctx)
}

func (a *Artifacts) SetVersion(version entities.ArtifactVersion) error {
	update := bson.D{{"$set", version}}
	_, err := a.versionsCollection.UpdateByID(a.ctx, version.UID, update, options.Update().SetUpsert(true))
	return err
}

func (a *Artifacts) ListVersions() ([]entities.ArtifactVersion, error) {
	cursor, err := a.versionsCollection.Find(a.ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	var versions []entities.ArtifactVersion
	if err := cursor.All(a.ctx, &versions); err != nil {
		return nil, err
	}

	return versions, cursor.Close(a.ctx)
}
