/*
 * Copyright (c) Dolittle. All rights reserved.
 * Licensed under the MIT license. See LICENSE file in the project root for full license information.
 */

package storage

import "dolittle.io/fleet-observer/entities"

type Applications interface {
	Set(application entities.Application) error
	Get(id entities.ApplicationUID) (*entities.Application, bool, error)
	List() ([]entities.Application, error)
}
