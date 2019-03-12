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

var stepCloneApp string

func init() {
	stepCloneCommand.Flags().StringVarP(&stepCloneApp, "app", "", "", "")

	stepCommand.AddCommand(stepCloneCommand)
}

var stepCloneCommand = &cobra.Command{
	Use: "clone",
	Run: func(cmd *cobra.Command, args []string) {
		sess, err := awsom.CreateSession()
		if err != nil {
			panic(err)
		}
		codePipelineService := codepipeline.New(sess)
		pipeline, err := codePipelineService.GetPipeline(&codepipeline.GetPipelineInput{Name: aws.String(stepCloneApp)})
		if err != nil {
			panic(err)
		}
		gitHubOwner := *pipeline.Pipeline.Stages[0].Actions[0].Configuration["Owner"]
		gitHubRepo := *pipeline.Pipeline.Stages[0].Actions[0].Configuration["Repo"]

		secretsManagerService := secretsmanager.New(sess)
		token, err := secretsManagerService.GetSecretValue(&secretsmanager.GetSecretValueInput{SecretId: aws.String(stepCloneApp)})
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
