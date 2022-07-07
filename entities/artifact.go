package entities

import (
	"fmt"
	"time"
)

type ArtifactUID string

var ArtifactType = "Artifact"

type Artifact struct {
	UID  ArtifactUID `bson:"_id" json:"uid"`
	Type string      `bson:"_type" json:"type"`

	Properties struct {
		ID string `bson:"id" json:"id"`
	} `bson:"properties" json:"properties"`

	Links struct {
		DevelopedByCustomerUID CustomerUID `bson:"developed_by_customer_uid" json:"developedBy"`
	} `bson:"links" json:"links"`
}

func NewArtifactUID(customerID, artifactID string) ArtifactUID {
	return ArtifactUID(fmt.Sprintf("%v/%v", customerID, artifactID))
}

func NewArtifact(customerID, id string) Artifact {
	artifact := Artifact{}
	artifact.UID = NewArtifactUID(customerID, id)
	artifact.Type = ArtifactType
	artifact.Properties.ID = id
	artifact.Links.DevelopedByCustomerUID = NewCustomerUID(customerID)
	return artifact
}

type ArtifactVersionUID string

var ArtifactVersionType = "ArtifactVersion"

type ArtifactVersion struct {
	UID  ArtifactVersionUID `bson:"_id" json:"uid"`
	Type string             `bson:"_type" json:"type"`

	Properties struct {
		Name     string    `bson:"name" json:"name"`
		Released time.Time `bson:"released" json:"-"`
	} `bson:"properties" json:"properties"`

	Links struct {
		VersionOfArtifactUID ArtifactUID `bson:"version_of_artifact_uid" json:"versionOf"`
	} `bson:"links" json:"links"`
}

func NewArtifactVersionUID(customerID, artifactID, artifactVersionID string) ArtifactVersionUID {
	return ArtifactVersionUID(fmt.Sprintf("%v/%v", NewArtifactUID(customerID, artifactID), artifactVersionID))
}

func NewArtifactVersion(customerID, artifactID, name string, released time.Time) ArtifactVersion {
	version := ArtifactVersion{}
	version.UID = NewArtifactVersionUID(customerID, artifactID, name)
	version.Type = ArtifactVersionType
	version.Properties.Name = name
	version.Properties.Released = released
	version.Links.VersionOfArtifactUID = NewArtifactUID(customerID, artifactID)
	return version
}
