package awsom

import (
	"bytes"
	"github.com/GeertJohan/go.rice"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	awsom_session "github.com/hekonsek/awsom-session"
	"os"
	"strings"
	"text/template"
)

// Policies

const PolicyCloudWatchLogsFullAccess = "arn:aws:iam::aws:policy/CloudWatchLogsFullAccess"

const PolicyAmazonS3FullAccess = "arn:aws:iam::aws:policy/AmazonS3FullAccess"

const PolicyAmazonEC2ContainerRegistryFullAccess = "arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryFullAccess"

const PolicyAWSCodeBuildDeveloperAccess = "arn:aws:iam::aws:policy/AWSCodeBuildDeveloperAccess"

const PolicySecretsManagerReadWrite = "arn:aws:iam::aws:policy/SecretsManagerReadWrite"

const PolicyAWSCodePipelineReadOnlyAccess = "arn:aws:iam::aws:policy/AWSCodePipelineReadOnlyAccess"

const PolicyIAMReadOnlyAccess = "arn:aws:iam::aws:policy/IAMReadOnlyAccess"

// Service

func iamService() (*iam.IAM, error) {
	sess, err := awsom_session.NewSession()
	if err != nil {
		return nil, err
	}
	return iam.New(sess), nil
}

// Account

const codeBuildEnvBuildArn = "CODEBUILD_BUILD_ARN"

func AccountId() (string, error) {
	if os.Getenv(codeBuildEnvBuildArn) == "" {
		iamService, err := iamService()
		if err != nil {
			return "", err
		}
		user, err := iamService.GetUser(&iam.GetUserInput{})
		if err != nil {
			return "", err
		}
		arnWithoutPrefix := strings.Replace(*user.User.Arn, "arn:aws:iam::", "", 1)
		return strings.Split(arnWithoutPrefix, ":")[0], nil
	} else {
		return strings.Split(os.Getenv(codeBuildEnvBuildArn), ":")[4], nil
	}
}

// Role

type Role struct {
	Name                     string
	AssumeRolePolicyDocument string
	Polices                  []string
}

func (r *Role) CreateOrUpdate() (string, error) {
	iamService, err := iamService()
	if err != nil {
		return "", err
	}
	roles, err := iamService.ListRoles(&iam.ListRolesInput{})
	if err != nil {
		return "", err
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
			return "", err
		}
		codeBuildRole = createRoleResponse.Role

		for _, policy := range r.Polices {
			_, err = iamService.AttachRolePolicy(&iam.AttachRolePolicyInput{
				RoleName:  codeBuildRole.RoleName,
				PolicyArn: aws.String(string(policy)),
			})
			if err != nil {
				return "", err
			}
		}
	}

	return *codeBuildRole.Arn, nil
}

func DeleteRole(roleName string) error {
	iamService, err := iamService()
	if err != nil {
		return err
	}

	policies, err := iamService.ListAttachedRolePolicies(&iam.ListAttachedRolePoliciesInput{RoleName: aws.String(roleName)})
	for _, policy := range policies.AttachedPolicies {
		_, err = iamService.DetachRolePolicy(&iam.DetachRolePolicyInput{
			RoleName:  aws.String(roleName),
			PolicyArn: policy.PolicyArn,
		})
		if err != nil {
			return err
		}
	}

	_, err = iamService.DeleteRole(&iam.DeleteRoleInput{
		RoleName: aws.String(roleName),
	})
	if err != nil {
		return err
	}

	return nil
}

func RoleArn(roleName string) (string, error) {
	iamService, err := iamService()
	if err != nil {
		return "", err
	}

	role, err := iamService.GetRole(&iam.GetRoleInput{
		RoleName: aws.String(roleName),
	})
	if role.Role == nil {
		return "", nil
	}
	if err != nil {
		return "", err
	}

	return *role.Role.Arn, nil
}

func AssumeServiceRolePolicyDocument(serviceName string) (string, error) {
	box, err := rice.FindBox("rice")
	if err != nil {
		return "", err
	}
	roleTemplate, err := box.String("assume_service_role_template.json")
	if err != nil {
		return "", err
	}
	templateParser, err := template.New("roleTemplate").Parse(roleTemplate)
	if err != nil {
		return "", err
	}

	var buffer bytes.Buffer
	err = templateParser.Execute(&buffer, map[string]string{
		"Service": serviceName,
	})
	if err != nil {
		return "", err
	}
	return buffer.String(), nil
}
