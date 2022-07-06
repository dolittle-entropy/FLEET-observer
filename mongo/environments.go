package mongo

import (
	"context"
	"dolittle.io/fleet-observer/entities"
	"fmt"
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
	id := fmt.Sprintf("%v/%v/%v", environment.OwnedByCustomerID, environment.EnvironmentOfApplicationID, environment.Name)
	update := bson.D{{"$set", environment}}
	_, err := e.collection.UpdateByID(e.ctx, id, update, options.Update().SetUpsert(true))
	return err
}
