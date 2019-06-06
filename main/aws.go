package main

import (
	"github.com/hekonsek/osexit"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(awsCommand)
}

var awsCommand = &cobra.Command{
	Use:   "aws",
	Short: "Low-level AWS operations.",
	Run: func(cmd *cobra.Command, args []string) {
		osexit.ExitOnError(cmd.Help())
	},
}
