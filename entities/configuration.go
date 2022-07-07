/*
 * Copyright (c) Dolittle. All rights reserved.
 * Licensed under the MIT license. See LICENSE file in the project root for full license information.
 */

package entities

import "fmt"

type ArtifactConfigurationUID string

var ArtifactConfigurationType = "ArtifactConfiguration"

type ArtifactConfiguration struct {
	UID  ArtifactConfigurationUID `bson:"_id" json:"uid"`
	Type string                   `bson:"_type" json:"type"`

	Properties struct {
		ContentHash string `bson:"content_hash" json:"hash"`
	} `bson:"properties" json:"properties"`

	Links struct {
	} `bson:"links" json:"-"`
}

func NewArtifactConfigurationUID(customerID, applicationID, environment, artifactID, contentHash string) ArtifactConfigurationUID {
	return ArtifactConfigurationUID(configurationUID(customerID, applicationID, environment, artifactID, contentHash))
}

func NewArtifactConfiguration(customerID, applicationID, environment, artifactID, contentHash string) ArtifactConfiguration {
	configuration := ArtifactConfiguration{}
	configuration.UID = NewArtifactConfigurationUID(customerID, applicationID, environment, artifactID, contentHash)
	configuration.Type = ArtifactConfigurationType
	configuration.Properties.ContentHash = contentHash
	return configuration
}

type RuntimeConfigurationUID string

var RuntimeConfigurationType = "RuntimeConfiguration"

type RuntimeConfiguration struct {
	UID  RuntimeConfigurationUID `bson:"_id" json:"uid"`
	Type string                  `bson:"_type" json:"type"`

	Properties struct {
		ContentHash string `bson:"content_hash" json:"hash"`
	} `bson:"properties" json:"properties"`

	Links struct {
	} `bson:"links" json:"-"`
}

func NewRuntimeConfigurationUID(customerID, applicationID, environment, artifactID, contentHash string) RuntimeConfigurationUID {
	return RuntimeConfigurationUID(configurationUID(customerID, applicationID, environment, artifactID, contentHash))
}

func NewRuntimeConfiguration(customerID, applicationID, environment, artifactID, contentHash string) RuntimeConfiguration {
	configuration := RuntimeConfiguration{}
	configuration.UID = NewRuntimeConfigurationUID(customerID, applicationID, environment, artifactID, contentHash)
	configuration.Type = RuntimeConfigurationType
	configuration.Properties.ContentHash = contentHash
	return configuration
}

func configurationUID(customerID, applicationID, environment, artifactID, contentHash string) string {
	return fmt.Sprintf("%v/%v/%v", NewEnvironmentUID(customerID, applicationID, environment), artifactID, contentHash)
}
