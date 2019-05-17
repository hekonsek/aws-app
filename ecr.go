package awsom

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecr"
	awsom_session "github.com/hekonsek/awsom-session"
)

func EnsureEcrRepositoryExists(name string) (string, error) {
	sess, err := awsom_session.NewSession()
	ecrService := ecr.New(sess)

	repositoryUri, err := EcrRepositoryExists(name)
	if err != nil {
		return "", err
	}

	if repositoryUri != "" {
		return repositoryUri, nil
	}

	createRepositoryResponse, err := ecrService.CreateRepository(&ecr.CreateRepositoryInput{
		RepositoryName: aws.String(name),
	})
	if err != nil {
		return "", err
	}

	return *createRepositoryResponse.Repository.RepositoryUri, nil
}

func EcrRepositoryExists(name string) (string, error) {
	sess, err := awsom_session.NewSession()
	ecrService := ecr.New(sess)

	repositories, err := ecrService.DescribeRepositories(&ecr.DescribeRepositoriesInput{})
	if err != nil {
		return "", err
	}
	for _, repository := range repositories.Repositories {
		if *repository.RepositoryName == name {
			return *repository.RepositoryUri, nil
		}
	}

	return "", nil
}
