/*
 * Copyright (c) Dolittle. All rights reserved.
 * Licensed under the MIT license. See LICENSE file in the project root for full license information.
 */

package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repositories struct {
	Nodes          *Nodes
	Customers      *Customers
	Applications   *Applications
	Environments   *Environments
	Artifacts      *Artifacts
	Runtimes       *Runtimes
	Deployments    *Deployments
	Configurations *Configurations
}

func NewRepositories(database *mongo.Database, ctx context.Context) *Repositories {
	return &Repositories{
		Nodes:          NewNodes(database, ctx),
		Customers:      NewCustomers(database, ctx),
		Applications:   NewApplications(database, ctx),
		Environments:   NewEnvironments(database, ctx),
		Artifacts:      NewArtifacts(database, ctx),
		Runtimes:       NewRuntimes(database, ctx),
		Deployments:    NewDeployments(database, ctx),
		Configurations: NewConfigurations(database, ctx),
	}
}
