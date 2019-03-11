package awsom

import "errors"

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

	err := ApplyCodeBuildDefaults(CodeBuild{
		Name:   application.Name,
		GitUrl: application.GitUrl,
	}).CreateOrUpdate()
	if err != nil {
		return err
	}

	err = ApplyCodeBuildDefaults(CodeBuild{
		Name:       application.Name + "-dockerize",
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
	err = DeleteCodeBuild(name + "-dockerize")
	if err != nil {
		return err
	}
	err = DeleteCodePipeline(name)
	if err != nil {
		return err
	}

	return nil
}
