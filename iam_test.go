package awsom

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/hekonsek/random-strings"
	"strconv"
	"testing"
)
import "github.com/stretchr/testify/assert"

func TestCodeBuildRoleExists(t *testing.T) {
	t.Parallel()

	// Given
	assumeRolePolicyDocument, err := AssumeServiceRolePolicyDocument("codebuild.amazonaws.com")
	assert.NoError(t, err)
	roleName := randomstrings.ForHumanWithHash()

	// When
	arn, err := (&Role{
		Name:                     roleName,
		AssumeRolePolicyDocument: assumeRolePolicyDocument,
	}).CreateOrUpdate()
	assert.NoError(t, err)

	// Then
	assert.NotEmpty(t, arn)

	// Clean up
	err = DeleteRole(roleName)
	assert.NoError(t, err)
}

func TestRoleHasPolicy(t *testing.T) {
	t.Parallel()

	// Given
	assumeRolePolicyDocument, err := AssumeServiceRolePolicyDocument("codebuild.amazonaws.com")
	assert.NoError(t, err)
	roleName := randomstrings.ForHumanWithHash()

	// When
	_, err = (&Role{
		Name:                     roleName,
		AssumeRolePolicyDocument: assumeRolePolicyDocument,
		Polices:                  []string{PolicyCloudWatchLogsFullAccess},
	}).CreateOrUpdate()
	assert.NoError(t, err)

	// Then
	iamService, err := iamService()
	assert.NoError(t, err)
	policies, err := iamService.ListAttachedRolePolicies(&iam.ListAttachedRolePoliciesInput{RoleName: aws.String(roleName)})
	assert.NoError(t, err)
	assert.True(t, len(policies.AttachedPolicies) == 1)

	// Clean up
	err = DeleteRole(roleName)
	assert.NoError(t, err)
}

func TestReadAccountId(t *testing.T) {
	t.Parallel()

	id, err := AccountId()
	assert.NoError(t, err)
	assert.NotEmpty(t, id)
	_, err = strconv.ParseInt(id, 0, 64)
	assert.NoError(t, err)
}
