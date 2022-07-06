package mongo

import (
	"context"
	"dolittle.io/fleet-observer/entities"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Configurations struct {
	customerCollection *mongo.Collection
	runtimeCollection  *mongo.Collection
	ctx                context.Context
}

func NewConfigurations(database *mongo.Database, ctx context.Context) *Configurations {
	return &Configurations{
		customerCollection: database.Collection("artifact-configurations"),
		runtimeCollection:  database.Collection("runtime-configurations"),
		ctx:                ctx,
	}
}

func (c *Configurations) SetCustomer(config entities.CustomerConfiguration) error {
	id := fmt.Sprintf("%v/%v/%v/%v/%v", config.OwnedByCustomerID, config.EnvironmentOfApplicationID, config.DeployedInEnvironmentName, config.ConfigForArtifactID, config.ContentHash)
	update := bson.D{{"$set", config}}
	_, err := c.customerCollection.UpdateByID(c.ctx, id, update, options.Update().SetUpsert(true))
	return err
}

func (c *Configurations) SetRuntime(config entities.RuntimeConfiguration) error {
	id := fmt.Sprintf("%v/%v/%v/%v/%v", config.OwnedByCustomerID, config.EnvironmentOfApplicationID, config.DeployedInEnvironmentName, config.ConfigForArtifactID, config.ContentHash)
	update := bson.D{{"$set", config}}
	_, err := c.runtimeCollection.UpdateByID(c.ctx, id, update, options.Update().SetUpsert(true))
	return err
}
