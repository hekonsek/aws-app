package awsom

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/codepipeline"
	"os"
	"strings"
)

// Constants

const codePipelineRoleName = "awsom-codepipeline"

type CodePipeline struct {
	Name   string
	GitUrl string
}

func (codePipeline *CodePipeline) CreateOrUpdate() error {
	roleArn, err := RoleArn(codePipelineRoleName)
	if err != nil {
		return err
	}
	if roleArn == "" {
		assumeRolePolicyDocument, err := AssumeServiceRolePolicyDocument("codepipeline.amazonaws.com")
		if err != nil {
			return err
		}
		roleArn, err = (&Role{
			Name:                     codePipelineRoleName,
			AssumeRolePolicyDocument: assumeRolePolicyDocument,
			Polices:                  []string{PolicyAmazonS3FullAccess, PolicyAWSCodeBuildDeveloperAccess},
		}).CreateOrUpdate()
		if err != nil {
			return err
		}
	}

	err = (&S3Bucket{
		Name: codePipeline.Name,
	}).CreateOrUpdate()
	if err != nil {
		return err
	}

	sess, err := CreateSession()
	if err != nil {
		return err
	}
	codePipelineService := codepipeline.New(sess)

	gitProjectWithGitSuffix := strings.Replace(codePipeline.GitUrl, "https://github.com/", "", 1)
	gitProjectInlined := strings.Replace(gitProjectWithGitSuffix, ".git", "", -1)
	gitProject := strings.Split(gitProjectInlined, "/")

	_, err = codePipelineService.CreatePipeline(&codepipeline.CreatePipelineInput{
		Pipeline: &codepipeline.PipelineDeclaration{
			Name:    aws.String(codePipeline.Name),
			RoleArn: aws.String(roleArn),
			Stages: []*codepipeline.StageDeclaration{
				{
					Name: aws.String("source"),
					Actions: []*codepipeline.ActionDeclaration{
						{
							Name: aws.String("source"),
							ActionTypeId: &codepipeline.ActionTypeId{
								Owner:    aws.String(codepipeline.ActionOwnerThirdParty),
								Provider: aws.String("GitHub"),
								Category: aws.String(codepipeline.ActionCategorySource),
								Version:  aws.String("1"),
							},
							Configuration: map[string]*string{
								"Owner":      aws.String(gitProject[0]),
								"Repo":       aws.String(gitProject[1]),
								"Branch":     aws.String("master"),
								"OAuthToken": aws.String(os.Getenv("GITHUB_TOKEN")),
							},
							OutputArtifacts: []*codepipeline.OutputArtifact{
								{
									Name: aws.String("source"),
								},
							},
						},
					},
				},
				&codepipeline.StageDeclaration{
					Name: aws.String("build"),
					Actions: []*codepipeline.ActionDeclaration{
						&codepipeline.ActionDeclaration{
							Name: aws.String("build"),
							ActionTypeId: &codepipeline.ActionTypeId{
								Owner:    aws.String(codepipeline.ActionOwnerAws),
								Provider: aws.String("CodeBuild"),
								Category: aws.String(codepipeline.ActionCategoryBuild),
								Version:  aws.String("1"),
							},
							Configuration: map[string]*string{
								"ProjectName": aws.String(codePipeline.Name),
							},
							InputArtifacts: []*codepipeline.InputArtifact{
								&codepipeline.InputArtifact{
									Name: aws.String("source"),
								},
							},
							OutputArtifacts: []*codepipeline.OutputArtifact{
								&codepipeline.OutputArtifact{
									Name: aws.String("build"),
								},
							},
						},
					},
				},
				&codepipeline.StageDeclaration{
					Name: aws.String("dockerize"),
					Actions: []*codepipeline.ActionDeclaration{
						&codepipeline.ActionDeclaration{
							Name: aws.String("dockerize"),
							ActionTypeId: &codepipeline.ActionTypeId{
								Owner:    aws.String(codepipeline.ActionOwnerAws),
								Provider: aws.String("CodeBuild"),
								Category: aws.String(codepipeline.ActionCategoryBuild),
								Version:  aws.String("1"),
							},
							Configuration: map[string]*string{
								"ProjectName": aws.String(codePipeline.Name + "-dockerize"),
							},
							InputArtifacts: []*codepipeline.InputArtifact{
								&codepipeline.InputArtifact{
									Name: aws.String("build"),
								},
							},
						},
					},
				},
			},
			ArtifactStore: &codepipeline.ArtifactStore{
				Type:     aws.String(codepipeline.ArtifactStoreTypeS3),
				Location: aws.String(codePipeline.Name),
			},
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func DeleteCodePipeline(name string) error {
	sess, err := CreateSession()
	if err != nil {
		return err
	}
	codePipelineService := codepipeline.New(sess)

	_, err = codePipelineService.DeletePipeline(&codepipeline.DeletePipelineInput{
		Name: aws.String(name),
	})
	if err != nil {
		return err
	}

	return nil
}
