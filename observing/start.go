/*
 * Copyright (c) Dolittle. All rights reserved.
 * Licensed under the MIT license. See LICENSE file in the project root for full license information.
 */

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

	namespacesHandler := NewNamespacesHandler(
		repositories.Customers,
		repositories.Applications,
		logger,
	)
	namespaces := kubernetes.NewObserver("namespaces", factory.Core().V1().Namespaces().Informer(), logger)
	namespaces.Start(namespacesHandler, stop)

	replicasetsHandler := NewReplicasetHandler(
		repositories.Environments,
		repositories.Artifacts,
		repositories.Runtimes,
		repositories.Deployments,
		logger,
	)
	replicasets := kubernetes.NewObserver("replicasets", factory.Apps().V1().ReplicaSets().Informer(), logger)
	replicasets.Start(replicasetsHandler, stop)

	podsHandler := NewPodsHandler(
		repositories.Configurations,
		repositories.Deployments,
		factory.Core().V1().ConfigMaps().Lister(),
		factory.Core().V1().Secrets().Lister(),
		factory.Apps().V1().ReplicaSets().Lister(),
		logger,
	)
	pods := kubernetes.NewObserver("pods", factory.Core().V1().Pods().Informer(), logger)
	pods.Start(podsHandler, stop)
}
