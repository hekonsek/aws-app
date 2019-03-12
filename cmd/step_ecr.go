package cmd

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/hekonsek/awsom"
	"github.com/spf13/cobra"
	"os"
)

func init() {
	stepCommand.AddCommand(stepEcrCommand)
}

var stepEcrCommand = &cobra.Command{
	Use: "ecr",
	Run: func(cmd *cobra.Command, args []string) {
		sess, err := awsom.CreateSession()
		ecrService := ecr.New(sess)

		applicationName := awsom.ApplicationNameFromCurrentBuild()
		accountId, err := awsom.AccountId()
		if err != nil {
			panic(err)
		}
		repositoryName := fmt.Sprintf("%s.dkr.ecr.%s.amazonaws.com/%s", accountId, os.Getenv("AWS_REGION"), applicationName)

		repositories, err := ecrService.DescribeRepositories(&ecr.DescribeRepositoriesInput{})
		if err != nil {
			panic(err)
		}
		repositoryExists := false
		for _, repository := range repositories.Repositories {
			if *repository.RepositoryName == applicationName {
				repositoryExists = true
				break
			}
		}

		if !repositoryExists {
			_, err = ecrService.CreateRepository(&ecr.CreateRepositoryInput{
				RepositoryName: aws.String(applicationName),
			})
			if err != nil {
				panic(err)
			}
		}

		fmt.Println(repositoryName)
	},
}
