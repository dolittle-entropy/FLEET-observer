/*
 * Copyright (c) Dolittle. All rights reserved.
 * Licensed under the MIT license. See LICENSE file in the project root for full license information.
 */

package cleanup

import (
	"context"
	"dolittle.io/fleet-observer/storage"
	"github.com/rs/zerolog"
	coreV1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	listerCoreV1 "k8s.io/client-go/listers/core/v1"
	"time"
)

type Instances struct {
	deployments  storage.Deployments
	environments storage.Environments
	applications storage.Applications
	pods         listerCoreV1.PodLister
	logger       zerolog.Logger
}

func (i *Instances) Cleanup(ctx context.Context) error {
	runningInstances, err := i.deployments.ListRunningInstances()
	if err != nil {
		return err
	}

	existingPods, err := i.pods.List(labels.Everything())
	if err != nil {
		return err
	}

	runningPodsByUID := map[string]bool{}
	for _, pod := range existingPods {
		if pod.Status.Phase == coreV1.PodFailed || pod.Status.Phase == coreV1.PodSucceeded {
			continue
		}
		runningPodsByUID[string(pod.GetUID())] = true
	}

	if err := ctx.Err(); err != nil {
		return err
	}

	for _, instance := range runningInstances {
		if err := ctx.Err(); err != nil {
			return err
		}

		if instance.Properties.Stopped != nil {
			i.logger.Warn().
				Str("uid", string(instance.UID)).
				Msg("DeploymentInstance already stopped")
			continue
		}

		if podStillRunning := runningPodsByUID[instance.Properties.ID]; podStillRunning {
			continue
		}

		i.logger.Info().
			Str("uid", string(instance.UID)).
			Msg("Marking DeploymentInstance as stopped since Pod doesn't exist anymore")

		now := time.Now().UTC()
		instance.Properties.Stopped = &now
		err = i.deployments.SetInstance(instance)
		if err != nil {
			return err
		}
	}

	return nil
}
