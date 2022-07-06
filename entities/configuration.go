package entities

type CustomerConfiguration struct {
	ContentHash                string `bson:"content_hash"`
	ConfigForArtifactID        string `bson:"config_for_artifact_id"`
	DeployedInEnvironmentName  string `bson:"deployed_in_environment_name"`
	EnvironmentOfApplicationID string `bson:"environment_of_application_id"`
	OwnedByCustomerID          string `bson:"owned_by_customer_id"`
}

type RuntimeConfiguration struct {
	ContentHash                string `bson:"content_hash"`
	ConfigForArtifactID        string `bson:"config_for_artifact_id"`
	DeployedInEnvironmentName  string `bson:"deployed_in_environment_name"`
	EnvironmentOfApplicationID string `bson:"environment_of_application_id"`
	OwnedByCustomerID          string `bson:"owned_by_customer_id"`
}
