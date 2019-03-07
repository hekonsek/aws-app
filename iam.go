package awsom

import (
	"bytes"
	"github.com/GeertJohan/go.rice"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"text/template"
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

func DeleteRole(roleName string) error {
	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		return err
	}
	iamService := iam.New(sess)

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
