package main

import (
	"github.com/hekonsek/awsom"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(prometheusCommand)
}

var prometheusCommand = &cobra.Command{
	Use:   "prometheus",
	Short: "Commands related to Prometheus monitoring service.",
	Run: func(cmd *cobra.Command, args []string) {
		awsom.ExitOnCliError(cmd.Help())
	},
}
