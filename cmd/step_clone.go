package cmd

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/codepipeline"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/hekonsek/awsom"
	"github.com/spf13/cobra"
	"os/exec"
)

func init() {
	stepCommand.AddCommand(stepCloneCommand)
}

var stepCloneCommand = &cobra.Command{
	Use: "clone",
	Run: func(cmd *cobra.Command, args []string) {
		applicationName := awsom.ApplicationNameFromCurrentBuild()

		codePipelineService, err := awsom.CodePipelineService()
		if err != nil {
			panic(err)
		}
		pipeline, err := codePipelineService.GetPipeline(&codepipeline.GetPipelineInput{Name: aws.String(applicationName)})
		if err != nil {
			panic(err)
		}
		gitHubOwner := *pipeline.Pipeline.Stages[0].Actions[0].Configuration["Owner"]
		gitHubRepo := *pipeline.Pipeline.Stages[0].Actions[0].Configuration["Repo"]

		secretsManagerService, err := awsom.SecretsManagerService()
		if err != nil {
			panic(err)
		}
		token, err := secretsManagerService.GetSecretValue(&secretsmanager.GetSecretValueInput{SecretId: aws.String(applicationName)})
		if err != nil {
			panic(err)
		}
		gitHubToken := *token.SecretString

		cmdx := exec.Command("mkdir", "cloned")
		stdoutStderr, err := cmdx.CombinedOutput()
		if err != nil {
			fmt.Println(string(stdoutStderr))
			fmt.Println()
			fmt.Println(err.Error())
		} else {
			fmt.Println(string(stdoutStderr))
		}

		cmdx = exec.Command("git", "init")
		cmdx.Dir = "cloned"
		stdoutStderr, err = cmdx.CombinedOutput()
		if err != nil {
			fmt.Println(string(stdoutStderr))
			fmt.Println()
			fmt.Println(err.Error())
		} else {
			fmt.Println(string(stdoutStderr))
		}

		cmdx = exec.Command("git", "pull", fmt.Sprintf("https://%s@github.com/%s/%s.git", gitHubToken, gitHubOwner, gitHubRepo))
		cmdx.Dir = "cloned"
		stdoutStderr, err = cmdx.CombinedOutput()
		if err != nil {
			fmt.Println(string(stdoutStderr))
			fmt.Println()
			fmt.Println(err.Error())
		} else {
			fmt.Println(string(stdoutStderr))
		}

		cmdx = exec.Command("git", "fetch", "--tags", fmt.Sprintf("https://%s@github.com/%s/%s.git", gitHubToken, gitHubOwner, gitHubRepo))
		cmdx.Dir = "cloned"
		stdoutStderr, err = cmdx.CombinedOutput()
		if err != nil {
			fmt.Println(string(stdoutStderr))
			fmt.Println()
			fmt.Println(err.Error())
		} else {
			fmt.Println(string(stdoutStderr))
		}
	},
}
