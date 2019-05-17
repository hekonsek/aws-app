package main

import (
	"github.com/hekonsek/awsom"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "awsom",
	Short: "Awsom - AWS applications made easy",

	Run: func(cmd *cobra.Command, args []string) {
		awsom.ExitOnCliError(cmd.Help())
	},
}

func main() {
	awsom.ExitOnCliError(RootCmd.Execute())
}
