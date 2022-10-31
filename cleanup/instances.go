/*
 * Copyright (c) Dolittle. All rights reserved.
 * Licensed under the MIT license. See LICENSE file in the project root for full license information.
 */

package cleanup

import (
	"context"
	"dolittle.io/fleet-observer/storage"
	"fmt"
	"github.com/rs/zerolog"
	"k8s.io/apimachinery/pkg/labels"
	v1 "k8s.io/client-go/listers/core/v1"
	"time"
)

type Instances struct {
	deployments  storage.Deployments
	environments storage.Environments
	applications storage.Applications
	pods         v1.PodLister
	logger       zerolog.Logger
}

func (i *Instances) Cleanup(ctx context.Context) error {
	running, err := i.deployments.ListRunningInstances()
	if err != nil {
		return err
	}

	if err := ctx.Err(); err != nil {
		return err
	}

	for _, instance := range running {
		if err := ctx.Err(); err != nil {
			return err
		}

		if instance.Properties.Stopped != nil {
			i.logger.Warn().
				Str("uid", string(instance.UID)).
				Msg("DeploymentInstance already stopped")
			continue
		}

		deployment, found, err := i.deployments.Get(instance.Links.InstanceOfDeploymentUID)
		if err != nil {
			return err
		}
		if !found {
			i.logger.Warn().
				Str("uid", string(instance.Links.InstanceOfDeploymentUID)).
				Msg("Deployment for DeploymentInstance not found in storage")
			continue
		}

		environment, found, err := i.environments.Get(deployment.Links.DeployedInEnvironmentUID)
		if err != nil {
			return err
		}
		if !found {
			i.logger.Warn().
				Str("uid", string(deployment.Links.DeployedInEnvironmentUID)).
				Msg("Environment for Deployment not found in storage")
			continue
		}

		application, found, err := i.applications.Get(environment.Links.EnvironmentOfApplicationUID)
		if err != nil {
			return err
		}
		if !found {
			i.logger.Warn().
				Str("uid", string(environment.Links.EnvironmentOfApplicationUID)).
				Msg("Application for Environment not found in storage")
			continue
		}

		namespace := fmt.Sprintf("application-%s", application.Properties.ID)
		selector, err := labels.ValidatedSelectorFromSet(labels.Set{
			"application":  application.Properties.Name,
			"environment":  environment.Properties.Name,
			"microservice": deployment.Properties.Name,
		})
		if err != nil {
			i.logger.Warn().
				Err(err).
				Msg("Failed to create pod selector")
			continue
		}

		pods, err := i.pods.Pods(namespace).List(selector)
		if err != nil {
			return err
		}

		podStillRunning := false
		for _, pod := range pods {
			if string(pod.GetUID()) == instance.Properties.ID {
				podStillRunning = true
				break
			}
		}

		if podStillRunning {
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
