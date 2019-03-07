package awsom

import (
	"testing"
)
import "github.com/stretchr/testify/assert"

func TestCodeBuildRoleExists(t *testing.T) {
	t.Parallel()

	err := (&Role{
		Name:                     RandomName(),
		AssumeRolePolicyDocument: codeBuildAssumeRolePolicyDocument,
		Polices:                  []string{PolicyCloudWatchLogsFullAccess, PolicyAmazonS3FullAccess, PolicyAmazonEC2ContainerRegistryFullAccess},
	}).CreateOrUpdate()
	assert.NoError(t, err)
}
