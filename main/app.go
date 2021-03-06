package main

import (
	"github.com/hekonsek/awsom"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(appCommand)
}

var appCommand = &cobra.Command{
	Use:   "app",
	Short: "Commands related to 'applications' (CI/CD pipelines).",
	Run: func(cmd *cobra.Command, args []string) {
		awsom.ExitOnCliError(cmd.Help())
	},
}
