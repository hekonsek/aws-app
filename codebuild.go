package awsom

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/codebuild"
)

// Constants

const codeBuildRoleName = "aws-app-codebuild"

type CodeBuild struct {
	Name       string
	GitUrl     string
	BuildSpec  string
	BuildImage string
}

func ApplyCodeBuildDefaults(codeBuild CodeBuild) *CodeBuild {
	if codeBuild.BuildSpec == "" {
		codeBuild.BuildSpec = "buildspec.yml"
	}
	if codeBuild.BuildImage == "" {
		codeBuild.BuildImage = "aws/codebuild/java:openjdk-11"
	}
	return &codeBuild
}

func (codeBuild *CodeBuild) CreateOrUpdate() error {
	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		return err
	}

	roleArn, err := RoleArn(codeBuildRoleName)
	if err != nil {
		return err
	}
	if roleArn == "" {
		assumeRolePolicyDocument, err := AssumeServiceRolePolicyDocument("codebuild.amazonaws.com")
		if err != nil {
			return err
		}
		roleArn, err = (&Role{
			Name:                     codeBuildRoleName,
			AssumeRolePolicyDocument: assumeRolePolicyDocument,
			Polices:                  []string{PolicyCloudWatchLogsFullAccess, PolicyAmazonS3FullAccess, PolicyAmazonEC2ContainerRegistryFullAccess},
		}).CreateOrUpdate()
		if err != nil {
			return err
		}
	}

	codeBuildService := codebuild.New(sess)
	_, err = codeBuildService.CreateProject(&codebuild.CreateProjectInput{
		Name:        aws.String(codeBuild.Name),
		ServiceRole: aws.String(roleArn),
		Environment: &codebuild.ProjectEnvironment{
			Type:        aws.String(codebuild.EnvironmentTypeLinuxContainer),
			Image:       aws.String(codeBuild.BuildImage),
			ComputeType: aws.String(codebuild.ComputeTypeBuildGeneral1Small),
		},
		Source: &codebuild.ProjectSource{
			Type:      aws.String(codebuild.SourceTypeGithub),
			Location:  aws.String(codeBuild.GitUrl),
			Buildspec: aws.String(codeBuild.BuildSpec),
		},
		Artifacts: &codebuild.ProjectArtifacts{
			Type:     aws.String(codebuild.ArtifactsTypeS3),
			Location: aws.String("capsilon-hekonsek"),
		},
	})
	if err != nil {
		return err
	}

	return nil
}
