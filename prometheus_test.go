package awsom_test

import (
	"github.com/hekonsek/awsom"
	"github.com/hekonsek/awsom/aws"
	randomstrings "github.com/hekonsek/random-strings"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMonitoringClusterCreated(t *testing.T) {
	t.Parallel()

	environment := randomstrings.ForHumanWithHash()
	// Given
	defer func() {
		awsom.Warn(awsom.DeleteEnv(environment))
		awsom.Warn(aws.DeleteHostedZone(environment + ".com"))
	}()
	err := awsom.NewEnvBuilder(environment, environment+".com").Create()
	assert.NoError(t, err)

	// When
	err = awsom.NewPrometheusBuilder().WithVPc(environment).Create()
	assert.NoError(t, err)

	// Then
	clusterExists, err := aws.EcsClusterExistsByName(environment)
	assert.True(t, clusterExists)
}
