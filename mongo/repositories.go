package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repositories struct {
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
		Customers:      NewCustomers(database, ctx),
		Applications:   NewApplications(database, ctx),
		Environments:   NewEnvironments(database, ctx),
		Artifacts:      NewArtifacts(database, ctx),
		Runtimes:       NewRuntimes(database, ctx),
		Deployments:    NewDeployments(database, ctx),
		Configurations: NewConfigurations(database, ctx),
	}
}
