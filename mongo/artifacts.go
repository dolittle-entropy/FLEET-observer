package mongo

import (
	"context"
	"dolittle.io/fleet-observer/entities"
	"fmt"
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
	id := fmt.Sprintf("%v/%v", artifact.DevelopedByCustomerID, artifact.ID)
	update := bson.D{{"$set", artifact}}
	_, err := a.collection.UpdateByID(a.ctx, id, update, options.Update().SetUpsert(true))
	return err
}

func (a *Artifacts) SetVersion(version entities.ArtifactVersion) error {
	id := fmt.Sprintf("%v/%v/%v", version.DevelopedByCustomerID, version.VersionOfArtifactID, version.Name)
	update := bson.D{{"$set", version}}
	_, err := a.versionsCollection.UpdateByID(a.ctx, id, update, options.Update().SetUpsert(true))
	return err
}
