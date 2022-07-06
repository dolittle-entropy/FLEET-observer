package entities

type Customer struct {
	ID   string `bson:"_id"`
	Name string `bson:"name"`
}
