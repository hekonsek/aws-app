package aws_app

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
)

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

func (application *Application) Create() error {
	err := ensureCodeBuildRoleExists()
	if err != nil {
		return err
	}

	return nil
}

// Internal helpers

func ensureCodeBuildRoleExists() error {
	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		return err
	}
	iamService := iam.New(sess)
	roles, err := iamService.ListRoles(&iam.ListRolesInput{})
	if err != nil {
		return err
	}
	var codeBuildRole *iam.Role
	for _, role := range roles.Roles {
		if *role.RoleName == codeBuildRoleName {
			codeBuildRole = role
			break
		}
	}

	if codeBuildRole == nil {
		createRoleRespomse, err := iamService.CreateRole(&iam.CreateRoleInput{
			RoleName:                 aws.String(codeBuildRoleName),
			AssumeRolePolicyDocument: aws.String(codeBuildAssumeRolePolicyDocument),
		})
		if err != nil {
			return err
		}
		codeBuildRole = createRoleRespomse.Role

		for _, policy := range []string{"arn:aws:iam::aws:policy/CloudWatchLogsFullAccess", "arn:aws:iam::aws:policy/AmazonS3FullAccess", "arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryFullAccess"} {
			_, err = iamService.AttachRolePolicy(&iam.AttachRolePolicyInput{
				RoleName:  codeBuildRole.RoleName,
				PolicyArn: aws.String(policy),
			})
			if err != nil {
				return err
			}
		}
	}

	return nil
}
