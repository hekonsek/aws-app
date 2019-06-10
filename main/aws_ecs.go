package main

import (
	"github.com/hekonsek/osexit"
	"github.com/spf13/cobra"
)

func init() {
	awsCommand.AddCommand(awsEcsCommand)
}

var awsEcsCommand = &cobra.Command{
	Use:   "ecs",
	Short: "Low-level AWS ECS related operations.",
	Run: func(cmd *cobra.Command, args []string) {
		osexit.ExitOnError(cmd.Help())
	},
}
