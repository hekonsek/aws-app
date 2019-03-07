package aws_app

import (
	"github.com/Pallinder/sillyname-go"
	"strings"
	"testing"
)
import "github.com/stretchr/testify/assert"

func TestCodeBuildRoleExists(t *testing.T) {
	t.Parallel()

	err := (&Role{
		Name:                     strings.Replace(sillyname.GenerateStupidName(), " ", "-", -1),
		AssumeRolePolicyDocument: codeBuildAssumeRolePolicyDocument,
		Polices:                  []string{PolicyCloudWatchLogsFullAccess, PolicyAmazonS3FullAccess, PolicyAmazonEC2ContainerRegistryFullAccess},
	}).CreateOrUpdate()
	assert.NoError(t, err)
}
