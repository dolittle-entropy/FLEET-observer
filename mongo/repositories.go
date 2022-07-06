package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repositories struct {
	Customers *Customers
}

func NewRepositories(database *mongo.Database, ctx context.Context) *Repositories {
	return &Repositories{
		Customers: NewCustomers(database, ctx),
	}
}
