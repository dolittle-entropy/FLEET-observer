package mongo

import (
	"context"
	"dolittle.io/fleet-observer/entities"
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
	update := bson.D{{"$set", deployment}}
	_, err := d.collection.UpdateByID(d.ctx, deployment.UID, update, options.Update().SetUpsert(true))
	return err
}

func (d *Deployments) List() ([]entities.Deployment, error) {
	cursor, err := d.collection.Find(d.ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	var deployments []entities.Deployment
	if err := cursor.All(d.ctx, &deployments); err != nil {
		return nil, err
	}

	return deployments, cursor.Close(d.ctx)
}

func (d *Deployments) SetInstance(instance entities.DeploymentInstance) error {
	update := bson.D{{"$set", instance}}
	_, err := d.instancesCollection.UpdateByID(d.ctx, instance.UID, update, options.Update().SetUpsert(true))
	return err
}

func (d *Deployments) ListInstances() ([]entities.DeploymentInstance, error) {
	cursor, err := d.instancesCollection.Find(d.ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	var instances []entities.DeploymentInstance
	if err := cursor.All(d.ctx, &instances); err != nil {
		return nil, err
	}

	return instances, cursor.Close(d.ctx)
}
