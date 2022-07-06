package cmd

import (
	"github.com/rs/zerolog"
	"os"
	"os/signal"
)

func StopChannelFromSignals(logger zerolog.Logger) <-chan struct{} {
	stopCh := make(chan struct{})

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	go func() {
		<-signals
		logger.Info().Msg("Caught interrupt signal, shutting down...")
		close(stopCh)
	}()

	return stopCh
}

func WaitForStop(stopCh <-chan struct{}, logger zerolog.Logger) error {
	<-stopCh
	logger.Info().Msg("Goodbye")
	return nil
}
