package entities

import "time"

type Deployment struct {
	ID                         string    `bson:"id"`
	Created                    time.Time `bson:"created"`
	DeploymentOfArtifactID     string    `bson:"deployment_of_artifact_id"`
	DeployedInEnvironmentName  string    `bson:"deployed_in_environment_name"`
	EnvironmentOfApplicationID string    `bson:"environment_of_application_id"`
	OwnedByCustomerID          string    `bson:"owned_by_customer_id"`
	UsesArtifactVersion        string    `bson:"uses_artifact_version"`
	UsesRuntimeVersion         string    `bson:"uses_runtime_version"`
}

type DeploymentInstance struct {
	ID                            string    `bson:"id"`
	Started                       time.Time `bson:"started"`
	InstanceOfDeploymentID        string    `bson:"instance_of_deployment_id"`
	DeploymentOfArtifactID        string    `bson:"deployment_of_artifact_id"`
	DeployedInEnvironmentName     string    `bson:"deployed_in_environment_name"`
	EnvironmentOfApplicationID    string    `bson:"environment_of_application_id"`
	OwnedByCustomerID             string    `bson:"owned_by_customer_id"`
	UsesArtifactConfigurationHash string    `bson:"uses_artifact_configuration_hash"`
	UsesRuntimeConfigurationHash  string    `bson:"uses_runtime_configuration_hash"`
}
