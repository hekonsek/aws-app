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
		out, err := awsom.Exec{
			Command:    "git tag -l",
			WorkingDir: "cloned",
		}.Run()
		if err != nil {
			awsom.ExitOnCliError(err)
		}
		version := out[len(out)-1]
		fmt.Println(version)
	},
}
