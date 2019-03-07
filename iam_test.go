package awsom

import (
	"testing"
)
import "github.com/stretchr/testify/assert"

func TestCodeBuildRoleExists(t *testing.T) {
	t.Parallel()
	assumeRolePolicyDocument, err := AssumeServiceRolePolicyDocument("codebuild.amazonaws.com")
	assert.NoError(t, err)

	err = (&Role{
		Name:                     RandomName(),
		AssumeRolePolicyDocument: assumeRolePolicyDocument,
		Polices:                  []string{PolicyCloudWatchLogsFullAccess, PolicyAmazonS3FullAccess, PolicyAmazonEC2ContainerRegistryFullAccess},
	}).CreateOrUpdate()
	assert.NoError(t, err)
}
