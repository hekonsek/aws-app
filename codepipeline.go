package awsom

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/codepipeline"
	"os"
	"strings"
)

// Service

func CodePipelineService() (*codepipeline.CodePipeline, error) {
	sess, err := CreateSession()
	if err != nil {
		return nil, err
	}
	return codepipeline.New(sess), err
}

// Constants

const codePipelineRoleName = "awsom-codepipeline"

type CodePipeline struct {
	Name   string
	GitUrl string
}

func (codePipeline *CodePipeline) CreateOrUpdate() error {
	codePipelineService, err := CodePipelineService()
	if err != nil {
		return err
	}

	existingPipelines, err := codePipelineService.ListPipelines(&codepipeline.ListPipelinesInput{})
	if err != nil {
		return err
	}
	for _, pipeline := range existingPipelines.Pipelines {
		if *pipeline.Name == codePipeline.Name {
			return nil
		}
	}

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
				{
					Name: aws.String("configure"),
					Actions: []*codepipeline.ActionDeclaration{
						{
							Name: aws.String("configure"),
							ActionTypeId: &codepipeline.ActionTypeId{
								Owner:    aws.String(codepipeline.ActionOwnerAws),
								Provider: aws.String("CodeBuild"),
								Category: aws.String(codepipeline.ActionCategoryBuild),
								Version:  aws.String("1"),
							},
							Configuration: map[string]*string{
								"ProjectName": aws.String(ConfigureStageName(codePipeline.Name)),
							},
							InputArtifacts: []*codepipeline.InputArtifact{
								{
									Name: aws.String("source"),
								},
							},
							OutputArtifacts: []*codepipeline.OutputArtifact{
								{
									Name: aws.String("configured-source"),
								},
							},
						},
					},
				},
				{
					Name: aws.String("version"),
					Actions: []*codepipeline.ActionDeclaration{
						{
							Name: aws.String("version"),
							ActionTypeId: &codepipeline.ActionTypeId{
								Owner:    aws.String(codepipeline.ActionOwnerAws),
								Provider: aws.String("CodeBuild"),
								Category: aws.String(codepipeline.ActionCategoryBuild),
								Version:  aws.String("1"),
							},
							Configuration: map[string]*string{
								"ProjectName": aws.String(VersionStageName(codePipeline.Name)),
							},
							InputArtifacts: []*codepipeline.InputArtifact{
								{
									Name: aws.String("configured-source"),
								},
							},
							OutputArtifacts: []*codepipeline.OutputArtifact{
								{
									Name: aws.String("versioned-source"),
								},
							},
						},
					},
				},
				{
					Name: aws.String("build"),
					Actions: []*codepipeline.ActionDeclaration{
						{
							Name: aws.String("build"),
							ActionTypeId: &codepipeline.ActionTypeId{
								Owner:    aws.String(codepipeline.ActionOwnerAws),
								Provider: aws.String("CodeBuild"),
								Category: aws.String(codepipeline.ActionCategoryBuild),
								Version:  aws.String("1"),
							},
							Configuration: map[string]*string{
								"ProjectName": aws.String(BuildStageName(codePipeline.Name)),
							},
							InputArtifacts: []*codepipeline.InputArtifact{
								{
									Name: aws.String("versioned-source"),
								},
							},
							OutputArtifacts: []*codepipeline.OutputArtifact{
								{
									Name: aws.String("build"),
								},
							},
						},
					},
				},
				{
					Name: aws.String("dockerize"),
					Actions: []*codepipeline.ActionDeclaration{
						{
							Name: aws.String("dockerize"),
							ActionTypeId: &codepipeline.ActionTypeId{
								Owner:    aws.String(codepipeline.ActionOwnerAws),
								Provider: aws.String("CodeBuild"),
								Category: aws.String(codepipeline.ActionCategoryBuild),
								Version:  aws.String("1"),
							},
							Configuration: map[string]*string{
								"ProjectName": aws.String(DockerizeStageName(codePipeline.Name)),
							},
							InputArtifacts: []*codepipeline.InputArtifact{
								{
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

type CodePipelineRecord struct {
	Name string
}

func ListCodePipelines() ([]CodePipelineRecord, error) {
	codePipelineService, err := CodePipelineService()
	if err != nil {
		return nil, err
	}

	pipelines, err := codePipelineService.ListPipelines(&codepipeline.ListPipelinesInput{})
	pipelinesRecords := []CodePipelineRecord{}
	for _, pipeline := range pipelines.Pipelines {
		pipelinesRecords = append(pipelinesRecords, CodePipelineRecord{
			Name: *pipeline.Name,
		})
	}
	return pipelinesRecords, nil
}

func DeleteCodePipeline(name string) error {
	codePipelineService, err := CodePipelineService()
	if err != nil {
		return err
	}

	_, err = codePipelineService.DeletePipeline(&codepipeline.DeletePipelineInput{
		Name: aws.String(name),
	})
	if err != nil {
		return err
	}

	err = DeleteS3Bucket(name)
	if err != nil {
		return err
	}

	return nil
}

func PipelineExists(name string) (exists bool, err error) {
	codePipelineService, err := CodePipelineService()
	if err != nil {
		return
	}

	pipelines, err := codePipelineService.ListPipelines(&codepipeline.ListPipelinesInput{})
	for len(pipelines.Pipelines) > 0 {
		for _, pipeline := range pipelines.Pipelines {
			if *pipeline.Name == name {
				return true, err
			}
		}
		if pipelines.NextToken != nil {
			pipelines, err = codePipelineService.ListPipelines(&codepipeline.ListPipelinesInput{
				NextToken: pipelines.NextToken,
			})
		} else {
			break
		}
	}

	return
}

func ReadPipelineSource(pipelineName string) (owner string, repo string, err error) {
	codePipelineService, err := CodePipelineService()
	if err != nil {
		return
	}

	exists, err := PipelineExists(pipelineName)
	if err != nil {
		return
	}
	if !exists {
		return
	}

	pipeline, err := codePipelineService.GetPipeline(&codepipeline.GetPipelineInput{Name: aws.String(pipelineName)})
	if err != nil {
		return
	}
	owner = *pipeline.Pipeline.Stages[0].Actions[0].Configuration["Owner"]
	repo = *pipeline.Pipeline.Stages[0].Actions[0].Configuration["Repo"]
	return
}

func ConfigureStageName(name string) string {
	return name + "-configure"
}

func VersionStageName(name string) string {
	return name + "-version"
}

func BuildStageName(name string) string {
	return name + "-build"
}

func DockerizeStageName(name string) string {
	return name + "-dockerize"
}
