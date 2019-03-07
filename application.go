package awsom

// Constants

const codeBuildRoleName = "aws-app-codebuild"

const codeBuildAssumeRolePolicyDocument = `{
   "Version":"2012-10-17",
   "Statement":[
      {
         "Effect":"Allow",
         "Principal":{
            "Service":"codebuild.amazonaws.com"
         },
         "Action":"sts:AssumeRole"
      }
   ]
}`

// Application type

type Application struct {
	Name string
}

func (application *Application) CreateOrUpdate() error {
	err := (&Role{
		Name:                     codeBuildRoleName,
		AssumeRolePolicyDocument: codeBuildAssumeRolePolicyDocument,
		Polices:                  []string{PolicyCloudWatchLogsFullAccess, PolicyAmazonS3FullAccess, PolicyAmazonEC2ContainerRegistryFullAccess},
	}).CreateOrUpdate()
	if err != nil {
		return err
	}

	return nil
}
