/*
 * Copyright (c) Dolittle. All rights reserved.
 * Licensed under the MIT license. See LICENSE file in the project root for full license information.
 */

package entities

import (
	"fmt"
	"time"
)

type DeploymentUID string

var DeploymentType = "Deployment"

type Deployment struct {
	UID  DeploymentUID `bson:"_id" json:"uid"`
	Type string        `bson:"_type" json:"type"`

	Properties struct {
		ID      string    `bson:"id" json:"id"`
		Name    string    `bson:"name" json:"name"`
		Created time.Time `bson:"created" json:"created"`
	} `bson:"properties" json:"properties"`

	Links struct {
		DeployedInEnvironmentUID EnvironmentUID     `bson:"deployed_in_environment_uid" json:"deployedIn"`
		UsesArtifactVersionUID   ArtifactVersionUID `bson:"uses_artifact_version_uid" json:"usesArtifact"`
		UsesRuntimeVersionUID    RuntimeVersionUID  `bson:"uses_runtime_version_uid" json:"usesRuntime"`
	} `bson:"links" json:"links"`
}

func NewDeploymentUID(customerID, applicationID, environment, deploymentID string) DeploymentUID {
	return DeploymentUID(fmt.Sprintf("%v/%v", NewEnvironmentUID(customerID, applicationID, environment), deploymentID))
}

func NewDeployment(customerID, applicationID, environment, id, name string, created time.Time, artifact ArtifactVersion, runtime RuntimeVersion) Deployment {
	deployment := Deployment{}
	deployment.UID = NewDeploymentUID(customerID, applicationID, environment, id)
	deployment.Type = DeploymentType
	deployment.Properties.ID = id
	deployment.Properties.Name = name
	deployment.Properties.Created = created
	deployment.Links.DeployedInEnvironmentUID = NewEnvironmentUID(customerID, applicationID, environment)
	deployment.Links.UsesArtifactVersionUID = artifact.UID
	deployment.Links.UsesRuntimeVersionUID = runtime.UID
	return deployment
}

type DeploymentInstanceUID string

var DeploymentInstanceType = "DeploymentInstance"

type DeploymentInstance struct {
	UID  DeploymentInstanceUID `bson:"_id" json:"uid"`
	Type string                `bson:"_type" json:"type"`

	Properties struct {
		ID      string     `bson:"id" json:"id"`
		Started time.Time  `bson:"started" json:"started"`
		Stopped *time.Time `bson:"stopped" json:"stopped,omitempty"`
	} `bson:"properties" json:"properties"`

	Links struct {
		InstanceOfDeploymentUID      DeploymentUID            `bson:"instance_of_deployment_uid" json:"instanceOf"`
		UsesArtifactConfigurationUID ArtifactConfigurationUID `bson:"uses_artifact_configuration_uid" json:"usesArtifactConfiguration"`
		UsesRuntimeConfigurationUID  RuntimeConfigurationUID  `bson:"uses_runtime_configuration_uid" json:"usesRuntimeConfiguration"`
		ScheduledOnNodeUID           NodeUID                  `bson:"scheduled_on_node_uid" json:"scheduledOn"`
	} `bson:"links" json:"links"`
}

func NewDeploymentInstanceUID(customerID, applicationID, environment, deploymentID, deploymentInstanceID string) DeploymentInstanceUID {
	return DeploymentInstanceUID(fmt.Sprintf("%v/%v", NewDeploymentUID(customerID, applicationID, environment, deploymentID), deploymentInstanceID))
}

func NewDeploymentInstance(customerID, applicationID, environment, deploymentID, id string, started time.Time, stopped *time.Time, artifact ArtifactConfiguration, runtime RuntimeConfiguration, nodeName string) DeploymentInstance {
	instance := DeploymentInstance{}
	instance.UID = NewDeploymentInstanceUID(customerID, applicationID, environment, deploymentID, id)
	instance.Type = DeploymentInstanceType
	instance.Properties.ID = id
	instance.Properties.Started = started
	instance.Properties.Stopped = stopped
	instance.Links.InstanceOfDeploymentUID = NewDeploymentUID(customerID, applicationID, environment, deploymentID)
	instance.Links.UsesArtifactConfigurationUID = artifact.UID
	instance.Links.UsesRuntimeConfigurationUID = runtime.UID
	instance.Links.ScheduledOnNodeUID = NewNodeUID(nodeName)
	return instance
}
