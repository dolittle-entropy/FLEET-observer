/*
 * Copyright (c) Dolittle. All rights reserved.
 * Licensed under the MIT license. See LICENSE file in the project root for full license information.
 */

package storage

import "dolittle.io/fleet-observer/entities"

type Nodes interface {
	Set(node entities.Node) error
	List() ([]entities.Node, error)
}
