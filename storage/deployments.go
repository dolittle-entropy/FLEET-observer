/*
 * Copyright (c) Dolittle. All rights reserved.
 * Licensed under the MIT license. See LICENSE file in the project root for full license information.
 */

package storage

import "dolittle.io/fleet-observer/entities"

type Deployments interface {
	Set(deployment entities.Deployment) error
	List() ([]entities.Deployment, error)
	SetInstance(instance entities.DeploymentInstance) error
	ListInstances() ([]entities.DeploymentInstance, error)
}
