/*
 * Copyright (c) Dolittle. All rights reserved.
 * Licensed under the MIT license. See LICENSE file in the project root for full license information.
 */

package cmd

import (
	"dolittle.io/fleet-observer/config"
	"dolittle.io/fleet-observer/exporting"
	"dolittle.io/fleet-observer/storage"
	"github.com/spf13/cobra"
)

var export = &cobra.Command{
	Use:   "export",
	Short: "Exports the stored data in the database as NDJSON",
	RunE: func(cmd *cobra.Command, args []string) error {
		config, logger, err := config.SetupFor(cmd)
		if err != nil {
			return err
		}

		ctx := ContextFromSignals(logger)

		repositories, err := storage.Connect(config, logger, ctx)
		if err != nil {
			return err
		}

		exporter := exporting.NewExporter(repositories, logger, ctx)
		return exporter.ExportToFile(config.String("output"))
	},
}

func init() {
	export.Flags().String("output", "./export.ndjson", "The output file to export to")
}
