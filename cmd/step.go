package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(stepCommand)
}

var stepCommand = &cobra.Command{
	Use: "step",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}
