package cmd

import (
	"context"
	"github.com/rs/zerolog"
	"os"
	"os/signal"
)

func ContextFromSignals(logger zerolog.Logger) context.Context {
	ctx, cancel := context.WithCancel(context.Background())

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	go func() {
		<-signals
		logger.Info().Msg("Caught interrupt signal, shutting down...")
		cancel()
	}()

	return ctx
}

func WaitForStop(logger zerolog.Logger, ctx context.Context) error {
	<-ctx.Done()
	logger.Info().Msg("Goodbye")
	return nil
}
