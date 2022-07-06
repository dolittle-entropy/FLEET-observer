package cmd

import (
	"dolittle.io/fleet-observer/config"
	"dolittle.io/fleet-observer/kubernetes"
	"github.com/spf13/cobra"
	"k8s.io/client-go/informers"
	"time"
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

		factory := informers.NewSharedInformerFactory(client, 1*time.Minute)

		observer := kubernetes.NewObserver("namespaces", factory.Core().V1().Namespaces().Informer(), logger)

		stop := StopChannelFromSignals(logger)

		observer.Start(kubernetes.ObserverHandlerFuncs{
			HandleFunc: func(obj any) error {
				logger.Info().Interface("obj", obj).Msg("Handling")
				return nil
			},
		}, stop)

		go factory.Start(stop)

		return WaitForStop(stop, logger)
	},
}
