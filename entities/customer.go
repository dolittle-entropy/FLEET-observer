package entities

type Customer struct {
	ID   string `bson:"id"`
	Name string `bson:"name"`
}
