package mongo

import (
	"context"
	"dolittle.io/fleet-observer/entities"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Deployments struct {
	collection          *mongo.Collection
	instancesCollection *mongo.Collection
	ctx                 context.Context
}

func NewDeployments(database *mongo.Database, ctx context.Context) *Deployments {
	return &Deployments{
		collection:          database.Collection("deployments"),
		instancesCollection: database.Collection("deployment-instances"),
		ctx:                 ctx,
	}
}

func (d *Deployments) Set(deployment entities.Deployment) error {
	id := fmt.Sprintf("%v/%v/%v/%v/%v", deployment.OwnedByCustomerID, deployment.EnvironmentOfApplicationID, deployment.DeployedInEnvironmentName, deployment.DeploymentOfArtifactID, deployment.ID)
	update := bson.D{{"$set", deployment}}
	_, err := d.collection.UpdateByID(d.ctx, id, update, options.Update().SetUpsert(true))
	return err
}

func (d *Deployments) SetInstance(instance entities.DeploymentInstance) error {
	id := fmt.Sprintf("%v/%v/%v/%v/%v/%v", instance.OwnedByCustomerID, instance.EnvironmentOfApplicationID, instance.DeployedInEnvironmentName, instance.DeploymentOfArtifactID, instance.InstanceOfDeploymentID, instance.ID)
	update := bson.D{{"$set", instance}}
	_, err := d.instancesCollection.UpdateByID(d.ctx, id, update, options.Update().SetUpsert(true))
	return err
}
