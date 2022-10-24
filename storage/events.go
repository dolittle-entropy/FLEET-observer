/*
 * Copyright (c) Dolittle. All rights reserved.
 * Licensed under the MIT license. See LICENSE file in the project root for full license information.
 */

package storage

import "dolittle.io/fleet-observer/entities"

type Events interface {
	Set(event entities.Event) error
	Get(id entities.EventUID) (*entities.Event, bool, error)
	List() ([]entities.Event, error)
}
