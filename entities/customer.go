/*
 * Copyright (c) Dolittle. All rights reserved.
 * Licensed under the MIT license. See LICENSE file in the project root for full license information.
 */

package entities

import "fmt"

type CustomerUID string

var CustomerType = "Customer"

type Customer struct {
	UID  CustomerUID `bson:"_id" json:"uid"`
	Type string      `bson:"_type" json:"type"`

	Properties struct {
		ID   string `bson:"id" json:"id"`
		Name string `bson:"name" json:"name"`
	} `bson:"properties" json:"properties"`

	Links struct {
	} `bson:"links" json:"-"`
}

func NewCustomerUID(customerID string) CustomerUID {
	return CustomerUID(fmt.Sprintf("%v", customerID))
}

func NewCustomer(id, name string) Customer {
	customer := Customer{}
	customer.UID = NewCustomerUID(id)
	customer.Type = CustomerType
	customer.Properties.ID = id
	customer.Properties.Name = name
	return customer
}
