package entities

type Application struct {
	ID                string `bson:"id"`
	Name              string `bson:"name"`
	OwnedByCustomerID string `bson:"owned_by_customer_id"`
}
