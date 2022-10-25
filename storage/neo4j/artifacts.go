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

type Artifacts struct {
	session neo4j.SessionWithContext
	ctx     context.Context
}

func NewArtifacts(session neo4j.SessionWithContext, ctx context.Context) *Artifacts {
	return &Artifacts{
		session: session,
		ctx:     ctx,
	}
}

func (a *Artifacts) Set(artifact entities.Artifact) error {
	return runUpdate(
		a.session,
		a.ctx,
		`
			MERGE (artifact:Artifact { _uid: $uid })
			SET artifact = { _uid: $uid, id: $id }
			WITH artifact
				MERGE (customer:Customer { _uid: $link_customer_uid})
				WITH artifact, customer
					MERGE (artifact)-[:DevelopedBy]->(customer)
					WITH artifact, customer
						MATCH (artifact)-[r:DevelopedBy]->(other)
						WHERE other._uid <> customer._uid
						DELETE r
			RETURN id(artifact)
		`,
		map[string]any{
			"uid":               artifact.UID,
			"id":                artifact.Properties.ID,
			"link_customer_uid": artifact.Links.DevelopedByCustomerUID,
		})
}

func (a *Artifacts) List() ([]entities.Artifact, error) {
	var artifacts []entities.Artifact
	return artifacts, querySingleJson(
		a.session,
		a.ctx,
		`
			MATCH (artifact:Artifact)-[:DevelopedBy]->(customer:Customer)
			WITH {
				uid: artifact._uid,
				type: "Artifact",
				properties: {
					id: artifact.id
				},
				links: {
					developedBy: customer._uid
				}
			} as entry
			RETURN apoc.convert.toJson(collect(entry)) as json
		`,
		&artifacts)
}

func (a *Artifacts) SetVersion(version entities.ArtifactVersion) error {
	return runUpdate(
		a.session,
		a.ctx,
		`
			MERGE (version:ArtifactVersion { _uid: $uid })
			SET version = { _uid: $uid, name: $name }
			WITH version
				MERGE (artifact:Artifact { _uid: $link_artifact_uid})
				WITH version, artifact
					MERGE (version)-[:VersionOf]->(artifact)
					WITH version, artifact
						MATCH (version)-[r:VersionOf]->(other)
						WHERE other._uid <> artifact._uid
						DELETE r
			RETURN id(version)
		`,
		map[string]any{
			"uid":               version.UID,
			"name":              version.Properties.Name,
			"link_artifact_uid": version.Links.VersionOfArtifactUID,
		})
}

func (a *Artifacts) ListVersions() ([]entities.ArtifactVersion, error) {
	var versions []entities.ArtifactVersion
	return versions, querySingleJson(
		a.session,
		a.ctx,
		`
			MATCH (version:ArtifactVersion)-[:VersionOf]->(artifact:Artifact)
			WITH {
				uid: version._uid,
				type: "ArtifactVersion",
				properties: {
					name: version.name
				},
				links: {
					versionOf: artifact._uid
				}
			} as entry
			RETURN apoc.convert.toJson(collect(entry)) as json
		`,
		&versions)
}
