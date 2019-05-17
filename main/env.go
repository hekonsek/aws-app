package main

import (
	"github.com/hekonsek/osexit"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(envCommand)
}

var envCommand = &cobra.Command{
	Use:   "env",
	Short: "Commands related to environments.",
	Run: func(cmd *cobra.Command, args []string) {
		osexit.ExitOnError(cmd.Help())
	},
}
