package cmd

import (
	config "dolittle.io/fleet-observer/config"
	"dolittle.io/fleet-observer/mongo"
	"github.com/spf13/cobra"
)

var drop = &cobra.Command{
	Use:   "drop",
	Short: "Drops the stored data in the database",
	RunE: func(cmd *cobra.Command, args []string) error {
		config, logger, err := config.SetupFor(cmd)
		if err != nil {
			return err
		}

		ctx := ContextFromSignals(logger)

		database, err := mongo.ConnectToMongo(config, logger, ctx)
		if err != nil {
			return err
		}

		logger.Warn().Str("database", database.Name()).Msg("WILL DROP ALL DATA FROM THE DATABASE!")
		logger.Warn().Msg("Are you sure you want to continue?")
		logger.Warn().Msg("Type 'yes' to drop the database...")

		answer, err := ReadLineFromInput(cmd, ctx)
		if err != nil || answer != "yes" {
			logger.Info().Err(err).Msg("Dropping aborted")
			return nil
		}

		logger.Info().Msg("Dropping database...")
		return database.Drop(ctx)
	},
}
