package main

import (
	"github.com/hekonsek/osexit"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(runtimeCommand)
}

var runtimeCommand = &cobra.Command{
	Use:   "runtime",
	Short: "Commands related to 'runtimes' (ECS, EKS, Lambda, etc).",
	Run: func(cmd *cobra.Command, args []string) {
		osexit.ExitOnError(cmd.Help())
	},
}
