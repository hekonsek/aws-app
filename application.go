package awsom

// Constants

const codeBuildRoleName = "aws-app-codebuild"

// Application type

type Application struct {
	Name string
}

func (application *Application) CreateOrUpdate() error {
	assumeRolePolicyDocument, err := AssumeServiceRolePolicyDocument("codebuild.amazonaws.com")
	if err != nil {
		return err
	}
	err = (&Role{
		Name:                     codeBuildRoleName,
		AssumeRolePolicyDocument: assumeRolePolicyDocument,
		Polices:                  []string{PolicyCloudWatchLogsFullAccess, PolicyAmazonS3FullAccess, PolicyAmazonEC2ContainerRegistryFullAccess},
	}).CreateOrUpdate()
	if err != nil {
		return err
	}

	return nil
}
