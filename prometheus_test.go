package awsom_test

import (
	"github.com/hekonsek/awsom"
	"github.com/hekonsek/awsom/aws"
	"testing"
)
import "github.com/stretchr/testify/assert"

func TestMonitoringEnvironmentCreated(t *testing.T) {
	t.Parallel()

	// Given
	defer func() {
		err := awsom.DeleteEcsApplication("monitoring", "prometheus")
		assert.NoError(t, err)
	}()
	defer func() {
		err := awsom.DeleteEcsCluster("monitoring")
		assert.NoError(t, err)
	}()
	defer func() {
		err := aws.DeleteLoadBalancer("monitoring")
		assert.NoError(t, err)
	}()
	defer func() {
		err := aws.DeleteVpc("monitoring")
		assert.NoError(t, err)
	}()

	// When
	err := awsom.NewPrometheusBuilder().Create()
	assert.NoError(t, err)

	// Then
	vpcExists, err := aws.VpcExistsByName("monitoring")
	assert.True(t, vpcExists)
	clusterExists, err := awsom.EcsClusterExistsByName("monitoring")
	assert.True(t, clusterExists)
}
