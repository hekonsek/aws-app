package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(appCommand)
}

var appCommand = &cobra.Command{
	Use: "app",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}
