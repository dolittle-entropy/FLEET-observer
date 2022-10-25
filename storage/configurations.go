/*
 * Copyright (c) Dolittle. All rights reserved.
 * Licensed under the MIT license. See LICENSE file in the project root for full license information.
 */

package storage

import "dolittle.io/fleet-observer/entities"

type Configurations interface {
	SetArtifact(config entities.ArtifactConfiguration) error
	ListArtifacts() ([]entities.ArtifactConfiguration, error)
	SetRuntime(config entities.RuntimeConfiguration) error
	ListRuntimes() ([]entities.RuntimeConfiguration, error)
}
