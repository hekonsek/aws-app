package cmd

import (
	"fmt"
	"github.com/hekonsek/awsom"
	"github.com/spf13/cobra"
)

func init() {
	stepCommand.AddCommand(stepVersionCurrentCommand)
}

var stepVersionCurrentCommand = &cobra.Command{
	Use: "version-current",
	Run: func(cmd *cobra.Command, args []string) {
		version, err := awsom.CurrentVersion("cloned")
		awsom.ExitOnCliError(err)
		fmt.Println(version)
	},
}
