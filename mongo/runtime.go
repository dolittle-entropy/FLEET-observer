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
	id := version.VersionString()
	update := bson.D{{"$set", version}}
	_, err := r.versionsCollection.UpdateByID(r.ctx, id, update, options.Update().SetUpsert(true))
	return err
}
