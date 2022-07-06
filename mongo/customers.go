package mongo

import (
	"context"
	"dolittle.io/fleet-observer/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Customers struct {
	collection *mongo.Collection
	ctx        context.Context
}

func NewCustomers(database *mongo.Database, ctx context.Context) *Customers {
	return &Customers{
		collection: database.Collection("customers"),
		ctx:        ctx,
	}
}

func (c *Customers) Set(customer entities.Customer) error {
	_, err := c.collection.ReplaceOne(c.ctx, bson.D{{"_id", customer.ID}}, customer, options.Replace().SetUpsert(true))
	return err
}

func (c *Customers) List() ([]entities.Customer, error) {
	cursor, err := c.collection.Find(c.ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	var customers []entities.Customer
	if err := cursor.All(c.ctx, &customers); err != nil {
		return nil, err
	}

	return customers, cursor.Close(c.ctx)
}
