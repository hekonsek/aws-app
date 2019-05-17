package main

import (
	"github.com/hekonsek/awsom"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(stepCommand)
}

var stepCommand = &cobra.Command{
	Use:   "step",
	Short: "Commands related to 'steps' (commands executed by CI/CD builder process).",
	Run: func(cmd *cobra.Command, args []string) {
		awsom.ExitOnCliError(cmd.Help())
	},
}
