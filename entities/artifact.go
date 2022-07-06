package entities

import "time"

type Artifact struct {
	ID                    string `bson:"id"`
	DevelopedByCustomerID string `bson:"developed_by_customer_id"`
}

type ArtifactVersion struct {
	Name                  string    `bson:"name"`
	Released              time.Time `bson:"released"`
	VersionOfArtifactID   string    `bson:"version_of_artifact_id"`
	DevelopedByCustomerID string    `bson:"developed_by_customer_id"`
}
