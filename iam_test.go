package awsom

import (
	"testing"
)
import "github.com/stretchr/testify/assert"

func TestCodeBuildRoleExists(t *testing.T) {
	t.Parallel()

	// Given
	assumeRolePolicyDocument, err := AssumeServiceRolePolicyDocument("codebuild.amazonaws.com")
	assert.NoError(t, err)
	roleName := RandomName()

	// When
	_, err = (&Role{
		Name:                     roleName,
		AssumeRolePolicyDocument: assumeRolePolicyDocument,
		Polices:                  []string{PolicyCloudWatchLogsFullAccess, PolicyAmazonS3FullAccess, PolicyAmazonEC2ContainerRegistryFullAccess},
	}).CreateOrUpdate()

	// Then
	assert.NoError(t, err)

	// Clean up
	err = DeleteRole(roleName)
	assert.NoError(t, err)
}
