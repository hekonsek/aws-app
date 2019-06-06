package main

import (
	"github.com/hekonsek/osexit"
	"github.com/spf13/cobra"
)

func init() {
	awsCommand.AddCommand(awsElbCommand)
}

var awsElbCommand = &cobra.Command{
	Use:   "elb",
	Short: "Low-level AWS ELB related operations.",
	Run: func(cmd *cobra.Command, args []string) {
		osexit.ExitOnError(cmd.Help())
	},
}
