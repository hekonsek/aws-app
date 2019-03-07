package awsom

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
)

// Constants

const PolicyCloudWatchLogsFullAccess = "arn:aws:iam::aws:policy/CloudWatchLogsFullAccess"

const PolicyAmazonS3FullAccess = "arn:aws:iam::aws:policy/AmazonS3FullAccess"

const PolicyAmazonEC2ContainerRegistryFullAccess = "arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryFullAccess"

// Roles

type Role struct {
	Name                     string
	AssumeRolePolicyDocument string
	Polices                  []string
}

func (r *Role) CreateOrUpdate() error {
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
		if *role.RoleName == r.Name {
			codeBuildRole = role
			break
		}
	}

	if codeBuildRole == nil {
		createRoleResponse, err := iamService.CreateRole(&iam.CreateRoleInput{
			RoleName:                 aws.String(r.Name),
			AssumeRolePolicyDocument: aws.String(r.AssumeRolePolicyDocument),
		})
		if err != nil {
			return err
		}
		codeBuildRole = createRoleResponse.Role

		for _, policy := range r.Polices {
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
