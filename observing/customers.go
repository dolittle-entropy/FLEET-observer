package observing

import (
	"dolittle.io/fleet-observer/entities"
	"dolittle.io/fleet-observer/mongo"
	"github.com/rs/zerolog"
	coreV1 "k8s.io/api/core/v1"
)

type CustomersHandler struct {
	customers *mongo.Customers
	logger    zerolog.Logger
}

func NewCustomersHandler(customers *mongo.Customers, logger zerolog.Logger) *CustomersHandler {
	return &CustomersHandler{
		customers: customers,
		logger:    logger.With().Str("handler", "customers").Logger(),
	}
}

func (c *CustomersHandler) Handle(obj any) error {
	namespace, ok := obj.(*coreV1.Namespace)
	if !ok {
		return ReceivedWrongType(obj, "Namespace")
	}

	logger := c.logger.With().Str("namespace", namespace.GetName()).Logger()

	tenantID, ok := namespace.GetAnnotations()["dolittle.io/tenant-id"]
	if !ok {
		logger.Trace().Msg("Skipping namespace because it does not have a tenant-id annotation")
		return nil
	}

	tenant := namespace.GetLabels()["tenant"]

	customer := entities.Customer{
		ID:   tenantID,
		Name: tenant,
	}
	if err := c.customers.Set(customer); err != nil {
		return err
	}

	logger.Debug().Interface("customer", customer).Msg("Updated customer")
	return nil
}
