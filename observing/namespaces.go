/*
 * Copyright (c) Dolittle. All rights reserved.
 * Licensed under the MIT license. See LICENSE file in the project root for full license information.
 */

package observing

import (
	"dolittle.io/fleet-observer/entities"
	"dolittle.io/fleet-observer/mongo"
	"github.com/rs/zerolog"
	coreV1 "k8s.io/api/core/v1"
)

type NamespacesHandler struct {
	customers    *mongo.Customers
	applications *mongo.Applications
	logger       zerolog.Logger
}

func NewNamespacesHandler(customers *mongo.Customers, applications *mongo.Applications, logger zerolog.Logger) *NamespacesHandler {
	return &NamespacesHandler{
		customers:    customers,
		applications: applications,
		logger:       logger.With().Str("handler", "namespaces").Logger(),
	}
}

func (nh *NamespacesHandler) Handle(obj any) error {
	namespace, ok := obj.(*coreV1.Namespace)
	if !ok {
		return ReceivedWrongType(obj, "Namespace")
	}

	logger := nh.logger.With().Str("namespace", namespace.GetName()).Logger()

	tenantID, ok := namespace.GetAnnotations()["dolittle.io/tenant-id"]
	if !ok {
		logger.Trace().Msg("Skipping namespace because it does not have a tenantID annotation")
		return nil
	}
	applicationID, ok := namespace.GetAnnotations()["dolittle.io/application-id"]
	if !ok {
		logger.Trace().Msg("Skipping namespace because it does not have an applicationID annotation")
		return nil
	}

	tenantName := namespace.GetLabels()["tenant"]
	applicationName := namespace.GetLabels()["application"]

	customer := entities.NewCustomer(tenantID, tenantName)
	if err := nh.customers.Set(customer); err != nil {
		return err
	}
	logger.Debug().Interface("customer", customer).Msg("Updated customer")

	application := entities.NewApplication(tenantID, applicationID, applicationName)
	if err := nh.applications.Set(application); err != nil {
		return err
	}
	logger.Debug().Interface("application", application).Msg("Updated application")

	return nil
}
