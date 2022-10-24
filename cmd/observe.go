/*
 * Copyright (c) Dolittle. All rights reserved.
 * Licensed under the MIT license. See LICENSE file in the project root for full license information.
 */

package cmd

import (
	"dolittle.io/fleet-observer/config"
	"dolittle.io/fleet-observer/kubernetes"
	"dolittle.io/fleet-observer/observing"
	"dolittle.io/fleet-observer/storage"
	"github.com/spf13/cobra"
	"k8s.io/client-go/informers"
)

var observe = &cobra.Command{
	Use:   "observe",
	Short: "Starts the observer",
	RunE: func(cmd *cobra.Command, _ []string) error {
		config, logger, err := config.SetupFor(cmd)
		if err != nil {
			return err
		}

		client, err := kubernetes.NewClientUsing(config)
		if err != nil {
			return err
		}

		factory := informers.NewSharedInformerFactory(client, config.Duration("kubernetes.sync-interval"))

		ctx := ContextFromSignals(logger)

		repositories, err := storage.Connect(config, logger, ctx)
		if err != nil {
			return err
		}

		observing.StartAllObservers(factory, repositories, logger, ctx)
		go factory.Start(ctx.Done())

		return WaitForStop(logger, ctx)
	},
}

func init() {
	observe.Flags().String("kubernetes.sync-interval", "1m", "The Kubernetes informer sync interval")
}
