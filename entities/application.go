package entities

import "fmt"

type ApplicationUID string

var ApplicationType = "Application"

type Application struct {
	UID  ApplicationUID `bson:"_id" json:"uid"`
	Type string         `bson:"_type" json:"type"`

	Properties struct {
		ID   string `bson:"id" json:"id"`
		Name string `bson:"name" json:"name"`
	} `bson:"properties" json:"properties"`

	Links struct {
		OwnedByCustomerUID CustomerUID `bson:"owned_by_customer_uid" json:"ownedBy"`
	} `bson:"links" json:"links"`
}

func NewApplicationUID(customerID, applicationID string) ApplicationUID {
	return ApplicationUID(fmt.Sprintf("%v/%v", NewCustomerUID(customerID), applicationID))
}

func NewApplication(customerID, id, name string) Application {
	application := Application{}
	application.UID = NewApplicationUID(customerID, id)
	application.Type = ApplicationType
	application.Properties.ID = id
	application.Properties.Name = name
	application.Links.OwnedByCustomerUID = NewCustomerUID(customerID)
	return application
}
