/*
 * Copyright (c) Dolittle. All rights reserved.
 * Licensed under the MIT license. See LICENSE file in the project root for full license information.
 */

package storage

import "dolittle.io/fleet-observer/entities"

type Deployments interface {
	Set(deployment entities.Deployment) error
	Get(id entities.DeploymentUID) (*entities.Deployment, bool, error)
	List() ([]entities.Deployment, error)
	SetInstance(instance entities.DeploymentInstance) error
	GetInstance(id entities.DeploymentInstanceUID) (*entities.DeploymentInstance, bool, error)
	ListInstances() ([]entities.DeploymentInstance, error)
	ListRunningInstances() ([]entities.DeploymentInstance, error)
}
