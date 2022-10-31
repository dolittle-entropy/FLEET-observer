/*
 * Copyright (c) Dolittle. All rights reserved.
 * Licensed under the MIT license. See LICENSE file in the project root for full license information.
 */

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

func (d *Deployments) Get(id entities.DeploymentUID) (*entities.Deployment, bool, error) {
	result := d.collection.FindOne(d.ctx, bson.D{{"_id", id}})
	err := result.Err()
	if err == mongo.ErrNoDocuments {
		return nil, false, nil
	} else if err != nil {
		return nil, true, err
	}

	deployment := &entities.Deployment{}
	err = result.Decode(deployment)
	if err != nil {
		return nil, true, err
	}

	return deployment, true, nil
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

func (d *Deployments) GetInstance(id entities.DeploymentInstanceUID) (*entities.DeploymentInstance, bool, error) {
	result := d.instancesCollection.FindOne(d.ctx, bson.D{{"_id", id}})
	err := result.Err()
	if err == mongo.ErrNoDocuments {
		return nil, false, nil
	} else if err != nil {
		return nil, true, err
	}

	instance := &entities.DeploymentInstance{}
	err = result.Decode(instance)
	if err != nil {
		return nil, true, err
	}

	return instance, true, nil
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

func (d *Deployments) ListRunningInstances() ([]entities.DeploymentInstance, error) {
	cursor, err := d.instancesCollection.Find(d.ctx, bson.D{
		{"$or", bson.A{
			bson.D{{"properties.stopped", bson.D{{"$exists", false}}}},
			bson.D{{"properties.stopped", nil}},
		}},
	})
	if err != nil {
		return nil, err
	}

	var instances []entities.DeploymentInstance
	if err := cursor.All(d.ctx, &instances); err != nil {
		return nil, err
	}

	return instances, cursor.Close(d.ctx)
}
