package mongo

import (
	"context"
	"dolittle.io/fleet-observer/entities"
	"fmt"
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
	id := fmt.Sprintf("%v/%v", application.OwnedByCustomerID, application.ID)
	update := bson.D{{"$set", application}}
	_, err := a.collection.UpdateByID(a.ctx, id, update, options.Update().SetUpsert(true))
	return err
}
