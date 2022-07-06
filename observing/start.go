package observing

import (
	"context"
	"dolittle.io/fleet-observer/kubernetes"
	"dolittle.io/fleet-observer/mongo"
	"github.com/rs/zerolog"
	"k8s.io/client-go/informers"
)

func StartAllObservers(factory informers.SharedInformerFactory, repositories *mongo.Repositories, logger zerolog.Logger, ctx context.Context) {
	stop := ctx.Done()
	kubernetes.NewObserver("customers", factory.Core().V1().Namespaces().Informer(), logger).Start(NewCustomersHandler(repositories.Customers, logger), stop)
}
