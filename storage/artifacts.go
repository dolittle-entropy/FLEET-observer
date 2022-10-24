/*
 * Copyright (c) Dolittle. All rights reserved.
 * Licensed under the MIT license. See LICENSE file in the project root for full license information.
 */

package storage

import "dolittle.io/fleet-observer/entities"

type Artifacts interface {
	Set(artifact entities.Artifact) error
	List() ([]entities.Artifact, error)
	SetVersion(version entities.ArtifactVersion) error
	ListVersions() ([]entities.ArtifactVersion, error)
}
