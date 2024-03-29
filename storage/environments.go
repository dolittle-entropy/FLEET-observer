/*
 * Copyright (c) Dolittle. All rights reserved.
 * Licensed under the MIT license. See LICENSE file in the project root for full license information.
 */

package storage

import "dolittle.io/fleet-observer/entities"

type Environments interface {
	Set(environment entities.Environment) error
	Get(id entities.EnvironmentUID) (*entities.Environment, bool, error)
	List() ([]entities.Environment, error)
}
