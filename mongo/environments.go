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
