package main

import (
	"github.com/hekonsek/osexit"
	"github.com/spf13/cobra"
)

func init() {
	awsCommand.AddCommand(awsVpcCommand)
}

var awsVpcCommand = &cobra.Command{
	Use:   "vpc",
	Short: "Low-level AWS VPCs related operations.",
	Run: func(cmd *cobra.Command, args []string) {
		osexit.ExitOnError(cmd.Help())
	},
}
