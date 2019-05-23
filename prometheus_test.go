package awsom_test

import (
	"github.com/hekonsek/awsom"
	"testing"
)
import "github.com/stretchr/testify/assert"

func TestMonitoringEnvironmentCreated(t *testing.T) {
	t.Parallel()

	// Given
	defer func() {
		err := awsom.DeleteEcsCluster("monitoring")
		assert.NoError(t, err)
	}()
	defer func() {
		err := awsom.DeleteVpc("monitoring")
		assert.NoError(t, err)
	}()

	// When
	err := awsom.NewPrometheusBuilder().Create()
	assert.NoError(t, err)

	// Then
	vpcExists, err := awsom.VpcExistsByName("monitoring")
	assert.True(t, vpcExists)
	clusterExists, err := awsom.EcsClusterExistsByName("monitoring")
	assert.True(t, clusterExists)
}
