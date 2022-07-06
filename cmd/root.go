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

	root.AddCommand(observe)
}
