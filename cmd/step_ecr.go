package cmd

import (
	"fmt"
	"github.com/hekonsek/awsom"
	"github.com/spf13/cobra"
)

func init() {
	stepCommand.AddCommand(stepEcrCommand)
}

var stepEcrCommand = &cobra.Command{
	Use: "ecr",
	Run: func(cmd *cobra.Command, args []string) {
		applicationName := awsom.ApplicationNameFromCurrentBuild()
		repositoryUri, err := awsom.EnsureEcrRepositoryExists(applicationName)
		if err != nil {
			panic(err)
		}
		fmt.Println(repositoryUri)
	},
}
