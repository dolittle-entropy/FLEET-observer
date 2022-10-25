/*
 * Copyright (c) Dolittle. All rights reserved.
 * Licensed under the MIT license. See LICENSE file in the project root for full license information.
 */

package neo4j

import (
	"context"
	"dolittle.io/fleet-observer/entities"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Deployments struct {
	session neo4j.SessionWithContext
	ctx     context.Context
}

func NewDeployments(session neo4j.SessionWithContext, ctx context.Context) *Deployments {
	return &Deployments{
		session: session,
		ctx:     ctx,
	}
}

func (d *Deployments) Set(deployment entities.Deployment) error {
	return nil
}

func (d *Deployments) List() ([]entities.Deployment, error) {
	return nil, nil
}

func (d *Deployments) SetInstance(instance entities.DeploymentInstance) error {
	return nil
}

func (d *Deployments) ListInstances() ([]entities.DeploymentInstance, error) {
	return nil, nil
}
