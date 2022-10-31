/*
 * Copyright (c) Dolittle. All rights reserved.
 * Licensed under the MIT license. See LICENSE file in the project root for full license information.
 */

package neo4j

import (
	"context"
	"dolittle.io/fleet-observer/entities"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"time"
)

type Deployments struct {
	session neo4j.SessionWithContext
	ctx     context.Context
}

func NewDeployments(session neo4j.SessionWithContext, ctx context.Context) *Deployments {
	return &Deployments{
		session: session,
		ctx:     ctx,
	}
}

func (d *Deployments) Set(deployment entities.Deployment) error {
	return multiUpdate(
		d.session,
		d.ctx,
		map[string]any{
			"uid":                       deployment.UID,
			"id":                        deployment.Properties.ID,
			"name":                      deployment.Properties.Name,
			"created":                   deployment.Properties.Created.Format(time.RFC3339),
			"link_environment_uid":      deployment.Links.DeployedInEnvironmentUID,
			"link_artifact_version_uid": deployment.Links.UsesArtifactVersionUID,
			"link_runtime_version_uid":  deployment.Links.UsesRuntimeVersionUID,
		},
		`
			MERGE (deployment:Deployment { _uid: $uid })
			SET deployment = { _uid: $uid, id: $id, name: $name, created: datetime($created) }
			RETURN id(deployment)
		`,
		`
			MATCH (deployment:Deployment { _uid: $uid })
			WITH deployment
				MERGE (environment:Environment { _uid: $link_environment_uid})
				WITH deployment, environment
					MERGE (deployment)-[:DeployedIn]->(environment)
					WITH deployment, environment
						MATCH (deployment)-[r:DeployedIn]->(other)
						WHERE other._uid <> environment._uid
						DELETE r
			RETURN id(deployment)
		`,
		`
			MATCH (deployment:Deployment { _uid: $uid })
			WITH deployment
				MERGE (version:ArtifactVersion { _uid: $link_artifact_version_uid})
				WITH deployment, version
					MERGE (deployment)-[:UsesArtifact]->(version)
					WITH deployment, version
						MATCH (deployment)-[r:UsesArtifact]->(other)
						WHERE other._uid <> version._uid
						DELETE r
			RETURN id(deployment)
		`,
		`
			MATCH (deployment:Deployment { _uid: $uid })
			WITH deployment
				MERGE (version:RuntimeVersion { _uid: $link_runtime_version_uid})
				WITH deployment, version
					MERGE (deployment)-[:UsesRuntime]->(version)
					WITH deployment, version
						MATCH (deployment)-[r:UsesRuntime]->(other)
						WHERE other._uid <> version._uid
						DELETE r
			RETURN id(deployment)
		`)
}

func (d *Deployments) List() ([]entities.Deployment, error) {
	var deployments []entities.Deployment
	return deployments, findAllJson(
		d.session,
		d.ctx,
		`
			MATCH (deployment:Deployment)-[:DeployedIn]->(environment:Environment)
			WITH deployment, environment
				MATCH (deployment)-[:UsesArtifact]->(artifact:ArtifactVersion)
			WITH deployment, environment, artifact
				MATCH (deployment)-[:UsesRuntime]->(runtime:RuntimeVersion)
			WITH {
				uid: deployment._uid,
				type: "Deployment",
				properties: {
					id: deployment.id,
					created: toString(deployment.created)
				},
				links: {
					deployedIn: environment._uid,
					usesArtifact: artifact._uid,
					usesRuntime: runtime._uid
				}
			} as entry
			RETURN apoc.convert.toJson(collect(entry)) as json
		`,
		&deployments)
}

func (d *Deployments) SetInstance(instance entities.DeploymentInstance) error {
	var stopped any = nil
	if instance.Properties.Stopped != nil {
		stopped = instance.Properties.Stopped.Format(time.RFC3339)
	}
	return multiUpdate(
		d.session,
		d.ctx,
		map[string]any{
			"uid":                      instance.UID,
			"id":                       instance.Properties.ID,
			"started":                  instance.Properties.Started.Format(time.RFC3339),
			"stopped":                  stopped,
			"link_deployment_uid":      instance.Links.InstanceOfDeploymentUID,
			"link_artifact_config_uid": instance.Links.UsesArtifactConfigurationUID,
			"link_runtime_config_uid":  instance.Links.UsesRuntimeConfigurationUID,
			"link_node_uid":            instance.Links.ScheduledOnNodeUID,
		},
		`
			MERGE (instance:DeploymentInstance { _uid: $uid })
			SET instance = { _uid: $uid, id: $id, started: datetime($started), stopped: datetime($stopped) }
			RETURN id(instance)
		`,
		`
			MATCH (instance:DeploymentInstance { _uid: $uid })
			WITH instance
				MERGE (deployment:Deployment { _uid: $link_deployment_uid })
				WITH instance, deployment
					MERGE (instance)-[:InstanceOf]->(deployment)
					WITH instance, deployment
						MATCH (instance)-[r:InstanceOf]->(other)
						WHERE other._uid <> deployment._uid
						DELETE r
			RETURN id(instance)
		`,
		`
			MATCH (instance:DeploymentInstance { _uid: $uid })
			WITH instance
				MERGE (config:ArtifactConfiguration { _uid: $link_artifact_config_uid})
				WITH instance, config
					MERGE (instance)-[:UsesArtifactConfiguration]->(config)
					WITH instance, config
						MATCH (instance)-[r:UsesArtifactConfiguration]->(other)
						WHERE other._uid <> config._uid
						DELETE r
			RETURN id(instance)
		`,
		`
			MATCH (instance:DeploymentInstance { _uid: $uid })
			WITH instance
				MERGE (config:RuntimeConfiguration { _uid: $link_runtime_config_uid})
				WITH instance, config
					MERGE (instance)-[:UsesRuntimeConfiguration]->(config)
					WITH instance, config
						MATCH (instance)-[r:UsesRuntimeConfiguration]->(other)
						WHERE other._uid <> config._uid
						DELETE r
			RETURN id(instance)
		`,
		`
			MATCH (instance:DeploymentInstance { _uid: $uid })
			WITH instance
				MERGE (node:Node { _uid: $link_node_uid})
				WITH instance, node
					MERGE (instance)-[:ScheduledOn]->(node)
					WITH instance, node
						MATCH (instance)-[r:ScheduledOn]->(other)
						WHERE other._uid <> node._uid
						DELETE r
			RETURN id(instance)
		`)
}

func (d *Deployments) ListInstances() ([]entities.DeploymentInstance, error) {
	var instances []entities.DeploymentInstance
	return instances, findAllJson(
		d.session,
		d.ctx,
		`
			MATCH (instance:DeploymentInstance)-[:InstanceOf]->(deployment:Deployment)
			WITH instance, deployment
				MATCH (instance)-[:UsesArtifactConfiguration]->(artifact:ArtifactConfiguration)
			WITH instance, deployment, artifact
				MATCH (instance)-[:UsesRuntimeConfiguration]->(runtime:RuntimeConfiguration)
			WITH instance, deployment, artifact, runtime
				MATCH (instance)-[:ScheduledOn]->(node:Node)
			WITH {
				uid: instance._uid,
				type: "DeploymentInstance",
				properties: {
					id: instance.id,
					started: toString(instance.started),
					stopped: toString(instance.stopped)
				},
				links: {
					instanceOf: deployment._uid,
					usesArtifactConfiguration: artifact._uid,
					usesRuntimeConfiguration: runtime._uid,
					scheduledOn: node._uid
				}
			} as entry
			RETURN apoc.convert.toJson(collect(entry)) as json
		`,
		&instances)
}
