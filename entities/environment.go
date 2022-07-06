package entities

type Environment struct {
	Name                       string `bson:"name"`
	OwnedByCustomerID          string `bson:"owned_by_customer_id"`
	EnvironmentOfApplicationID string `bson:"environment_of_application_id"`
}
