package awsom

// Application type

type Application struct {
	Name   string
	GitUrl string
}

func (application *Application) CreateOrUpdate() error {
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

	return nil
}
