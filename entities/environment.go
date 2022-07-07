package entities

import "fmt"

type EnvironmentUID string

var EnvironmentType = "Environment"

type Environment struct {
	UID  EnvironmentUID `bson:"_id" json:"uid"`
	Type string         `bson:"_type" json:"type"`

	Properties struct {
		Name string `bson:"name" json:"name"`
	} `bson:"properties" json:"properties"`

	Links struct {
		EnvironmentOfApplicationUID ApplicationUID `bson:"environment_of_application_uid" json:"environmentOf"`
	} `bson:"links" json:"links"`
}

func NewEnvironmentUID(customerID, applicationID, environment string) EnvironmentUID {
	return EnvironmentUID(fmt.Sprintf("%v/%v", NewApplicationUID(customerID, applicationID), environment))
}

func NewEnvironment(customerID, applicationID, name string) Environment {
	environment := Environment{}
	environment.UID = NewEnvironmentUID(customerID, applicationID, name)
	environment.Type = EnvironmentType
	environment.Properties.Name = name
	environment.Links.EnvironmentOfApplicationUID = NewApplicationUID(customerID, applicationID)
	return environment
}
