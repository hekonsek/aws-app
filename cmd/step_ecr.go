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
	Use:   "ecr",
	Short: "Returns ECR repository URI for a current application. Creates new repository if necessary.",
	Run: func(cmd *cobra.Command, args []string) {
		applicationName := awsom.ApplicationNameFromCurrentBuild()
		repositoryUri, err := awsom.EnsureEcrRepositoryExists(applicationName)
		awsom.ExitOnCliError(err)
		fmt.Println(repositoryUri)
	},
}
