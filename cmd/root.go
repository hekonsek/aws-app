package cmd

import (
	"github.com/hekonsek/awsom"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "awsom",
	Short: "Awsom - toolkit making AWS application delivery simple.",

	Run: func(cmd *cobra.Command, args []string) {
		awsom.ExitOnCliError(cmd.Help())
	},
}
