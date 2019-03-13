package cmd

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/codepipeline"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/hekonsek/awsom"
	"github.com/spf13/cobra"
	"os/exec"
	"strconv"
	"strings"
)

func init() {
	stepCommand.AddCommand(stepVersionCommand)
}

var stepVersionCommand = &cobra.Command{
	Use: "version",
	Run: func(cmd *cobra.Command, args []string) {
		applicationName := awsom.ApplicationNameFromCurrentBuild()

		sess, err := awsom.CreateSession()
		if err != nil {
			panic(err)
		}
		codePipelineService := codepipeline.New(sess)
		pipeline, err := codePipelineService.GetPipeline(&codepipeline.GetPipelineInput{Name: aws.String(applicationName)})
		if err != nil {
			panic(err)
		}
		gitHubOwner := *pipeline.Pipeline.Stages[0].Actions[0].Configuration["Owner"]
		gitHubRepo := *pipeline.Pipeline.Stages[0].Actions[0].Configuration["Repo"]

		secretsManagerService := secretsmanager.New(sess)
		token, err := secretsManagerService.GetSecretValue(&secretsmanager.GetSecretValueInput{SecretId: aws.String(applicationName)})
		if err != nil {
			panic(err)
		}
		gitHubToken := *token.SecretString

		cmdx := exec.Command("git", "tag", "-l")
		cmdx.Dir = "cloned"
		stdoutStderr, err := cmdx.CombinedOutput()
		if err != nil {
			fmt.Println(string(stdoutStderr))
			fmt.Println()
			panic(err)
		} else {
			fmt.Println(string(stdoutStderr))
		}

		version := ""
		if string(stdoutStderr) == "" {
			version = "0.0"
		} else {
			version = strings.Split(string(stdoutStderr), "\n")[len(strings.Split(string(stdoutStderr), "\n"))-2]
			versionNumber, err := strconv.ParseInt(strings.Split(version, ".")[1], 0, 64)
			if err != nil {
				panic(err)
			}
			versionNumber++
			version = fmt.Sprintf("0.%d", versionNumber)
		}

		cmdx = exec.Command("git", "tag", version)
		cmdx.Dir = "cloned"
		stdoutStderr, err = cmdx.CombinedOutput()
		if err != nil {
			fmt.Println(string(stdoutStderr))
			fmt.Println()
			fmt.Println(err.Error())
		} else {
			fmt.Println(string(stdoutStderr))
		}

		cmdx = exec.Command("git", "push", "--tags", fmt.Sprintf("https://%s@github.com/%s/%s.git", gitHubToken, gitHubOwner, gitHubRepo))
		cmdx.Dir = "cloned"
		stdoutStderr, err = cmdx.CombinedOutput()
		if err != nil {
			fmt.Println(string(stdoutStderr))
			fmt.Println()
			panic(err)
		} else {
			fmt.Println(string(stdoutStderr))
		}
	},
}
