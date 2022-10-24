/*
 * Copyright (c) Dolittle. All rights reserved.
 * Licensed under the MIT license. See LICENSE file in the project root for full license information.
 */

package storage

import "dolittle.io/fleet-observer/entities"

type Runtimes interface {
	SetVersion(version entities.RuntimeVersion) error
	ListVersions() ([]entities.RuntimeVersion, error)
}
