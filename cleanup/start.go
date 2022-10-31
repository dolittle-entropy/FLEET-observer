/*
 * Copyright (c) Dolittle. All rights reserved.
 * Licensed under the MIT license. See LICENSE file in the project root for full license information.
 */

package cleanup

import (
	"context"
	"dolittle.io/fleet-observer/storage"
	"github.com/rs/zerolog"
	"k8s.io/client-go/informers"
	"time"
)

func StartAllCleanup(period time.Duration, factory informers.SharedInformerFactory, repositories *storage.Repositories, logger zerolog.Logger, ctx context.Context) {
	instancesLogger := logger.With().Str("cleanup", "instances").Logger()
	instances := &Instances{
		deployments:  repositories.Deployments,
		environments: repositories.Environments,
		applications: repositories.Applications,
		pods:         factory.Core().V1().Pods().Lister(),
		logger:       instancesLogger,
	}
	go RunCleaner(instances, period, factory, instancesLogger, ctx)
}

type Cleaner interface {
	Cleanup(ctx context.Context) error
}

func RunCleaner(cleaner Cleaner, period time.Duration, factory informers.SharedInformerFactory, logger zerolog.Logger, ctx context.Context) {
	timer := time.NewTimer(1 * time.Second)

	for {
		factory.WaitForCacheSync(ctx.Done())

		select {
		case <-ctx.Done():
			logger.Debug().Msg("Stopping cleanup")
			return
		case <-timer.C:
		}

		logger.Debug().Msg("Running cleanup")

		err := cleaner.Cleanup(ctx)
		if err == context.Canceled || err == context.DeadlineExceeded {
			logger.Debug().Msg("Stopping cleanup")
			return
		}

		if err != nil {
			logger.Error().Err(err).Msg("Cleanup failed")
		}

		timer.Reset(period)
	}
}
