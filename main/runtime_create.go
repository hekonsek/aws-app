package main

import (
	"github.com/hekonsek/osexit"
	"github.com/spf13/cobra"
)

func init() {
	runtimeCommand.AddCommand(runtimeCreateCommand)
}

var runtimeCreateCommand = &cobra.Command{
	Use:   "create",
	Short: "Commands related to creating 'runtimes' (ECS, EKS, Lambda, etc).",
	Run: func(cmd *cobra.Command, args []string) {
		osexit.ExitOnError(cmd.Help())
	},
}
