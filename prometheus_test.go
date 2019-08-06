package awsom_test

import (
	"github.com/hekonsek/awsom"
	"github.com/hekonsek/awsom/aws"
	randomstrings "github.com/hekonsek/random-strings"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMonitoringEnvironmentCreated(t *testing.T) {
	t.Parallel()

	vpcName := randomstrings.ForHumanWithHash()[5:]
	// Given
	defer func() {
		err := aws.DeleteVpc(vpcName)
		assert.NoError(t, err)
	}()
	defer func() {
		err := aws.DeleteElasticLoadBalancer(vpcName)
		assert.NoError(t, err)
	}()
	defer func() {
		err := aws.DeleteEcsCluster(vpcName)
		assert.NoError(t, err)
	}()
	defer func() {
		err := aws.DeleteEcsApplication(vpcName, "prometheus")
		assert.NoError(t, err)
	}()

	// When
	err := awsom.NewPrometheusBuilder().WithVPc(vpcName).Create()
	assert.NoError(t, err)

	// Then
	vpcExists, err := aws.VpcExistsByName(vpcName)
	assert.True(t, vpcExists)
	clusterExists, err := aws.EcsClusterExistsByName(vpcName)
	assert.True(t, clusterExists)
}
