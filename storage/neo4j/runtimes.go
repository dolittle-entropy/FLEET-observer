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

type Runtimes struct {
	session neo4j.SessionWithContext
	ctx     context.Context
}

func NewRuntimes(session neo4j.SessionWithContext, ctx context.Context) *Runtimes {
	return &Runtimes{
		session: session,
		ctx:     ctx,
	}
}

func (r *Runtimes) SetVersion(version entities.RuntimeVersion) error {
	var prerelease any = nil
	if version.Properties.Prerelease != "" {
		prerelease = version.Properties.Prerelease
	}
	return multiUpdate(
		r.session,
		r.ctx,
		map[string]any{
			"uid":        version.UID,
			"major":      version.Properties.Major,
			"minor":      version.Properties.Minor,
			"patch":      version.Properties.Patch,
			"prerelease": prerelease,
		},
		`
			MERGE (version:RuntimeVersion { _uid: $uid })
			SET version = {
				_uid: $uid,
				major: $major,
				minor: $minor,
				patch: $patch,
				prerelease: $prerelease
			}
			RETURN id(version)
		`)
}

func (r *Runtimes) ListVersions() ([]entities.RuntimeVersion, error) {
	var versions []entities.RuntimeVersion
	return versions, findAllJson(
		r.session,
		r.ctx,
		`
			MATCH (version:RuntimeVersion)
			WITH {
				uid: version._uid,
				type: "RuntimeVersion",
				properties: {
					major: version.major,
					minor: version.minor,
					patch: version.patch,
					prerelease: version.prerelease
				}
			} as entry
			RETURN apoc.convert.toJson(collect(entry)) as json
		`,
		&versions)
}
