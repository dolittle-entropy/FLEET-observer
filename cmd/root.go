/*
 * Copyright (c) Dolittle. All rights reserved.
 * Licensed under the MIT license. See LICENSE file in the project root for full license information.
 */

package cmd

import "github.com/spf13/cobra"

var root = &cobra.Command{
	Use:   "fleet-observer",
	Short: "fleet-observer observes Dolittle microservices in a Kubernetes cluster",
	RunE: func(cmd *cobra.Command, _ []string) error {
		return cmd.Usage()
	},
}

// Execute starts the cobra.Command execution
func Execute() {
	cobra.CheckErr(root.Execute())
}

func init() {
	root.PersistentFlags().StringSlice("config", nil, "A configuration file to load, can be specified multiple times")
	root.PersistentFlags().String("logger.format", "console", "The logging format to use, 'json' or 'console'")
	root.PersistentFlags().String("logger.level", "info", "The logging minimum log level to output")
	root.PersistentFlags().String("mongodb.connection-string", "mongodb://localhost:27017/observer", "The connection string to MongoDB")
	root.PersistentFlags().String("neo4j.connection-string", "", "The connection string string to Neo4j")
	root.PersistentFlags().String("neo4j.username", "neo4j", "The username to use for authenticating with Neo4j")
	root.PersistentFlags().String("neo4j.password", "", "The password to use for authenticating with Neo4j. If not set, authentication will not be performed.")

	root.AddCommand(observe)
	root.AddCommand(drop)
	root.AddCommand(export)
}
