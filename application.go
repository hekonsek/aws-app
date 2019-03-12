package awsom

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"os"
	"strings"
)

const ErrorApplicationNameTooShort = "ERR_TO_SHORT"

// Application type

type Application struct {
	Name   string
	GitUrl string
}

func (application *Application) CreateOrUpdate() error {
	if len(application.Name) < 3 {
		return errors.New(ErrorApplicationNameTooShort)
	}

	sess, err := CreateSession()
	if err != nil {
		panic(err)
	}
	secretsManagerService := secretsmanager.New(sess)
	secrets, err := secretsManagerService.ListSecrets(&secretsmanager.ListSecretsInput{})
	if err != nil {
		panic(err)
	}
	secretExists := false
	for _, secret := range secrets.SecretList {
		if *secret.Name == application.Name {
			secretExists = true
			break
		}
	}
	if !secretExists {
		_, err := secretsManagerService.CreateSecret(&secretsmanager.CreateSecretInput{
			Name:         aws.String(application.Name),
			SecretString: aws.String(os.Getenv("GITHUB_TOKEN")),
		})
		if err != nil {
			panic(err)
		}
	}

	err = ApplyCodeBuildDefaults(CodeBuild{
		Name:   application.Name,
		GitUrl: application.GitUrl,
	}).CreateOrUpdate()
	if err != nil {
		return err
	}

	err = ApplyCodeBuildDefaults(CodeBuild{
		Name:       VersionStageName(application.Name),
		GitUrl:     application.GitUrl,
		BuildSpec:  "buildspec-version.yml",
		BuildImage: "hekonsek/awsom",
	}).CreateOrUpdate()
	if err != nil {
		return err
	}

	err = ApplyCodeBuildDefaults(CodeBuild{
		Name:       DockerizeStageName(application.Name),
		GitUrl:     application.GitUrl,
		BuildSpec:  "buildspec-docker.yml",
		BuildImage: "aws/codebuild/docker:18.09.0",
	}).CreateOrUpdate()
	if err != nil {
		return err
	}

	err = (&CodePipeline{
		Name:   application.Name,
		GitUrl: application.GitUrl,
	}).CreateOrUpdate()
	if err != nil {
		return err
	}

	return nil
}

type ApplicationRecord struct {
	Name string
}

func ListApplications() ([]ApplicationRecord, error) {
	pipelines, err := ListCodePipelines()
	if err != nil {
		return nil, err
	}
	applications := []ApplicationRecord{}
	for _, pipeline := range pipelines {
		applications = append(applications, ApplicationRecord{
			Name: pipeline.Name,
		})
	}
	return applications, nil
}

func DeleteApplication(name string) error {
	err := DeleteCodeBuild(name)
	if err != nil {
		return err
	}
	err = DeleteCodeBuild(VersionStageName(name))
	if err != nil {
		return err
	}
	err = DeleteCodeBuild(DockerizeStageName(name))
	if err != nil {
		return err
	}
	err = DeleteCodePipeline(name)
	if err != nil {
		return err
	}

	return nil
}

func ApplicationNameFromCurrentBuild() string {
	buildName := strings.Split(os.Getenv("CODEBUILD_BUILD_ID"), ":")[0]
	return strings.Split(buildName, "-")[0]
}
