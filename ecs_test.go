package awsom_test

import (
	"github.com/hekonsek/awsom"
	"github.com/hekonsek/awsom/aws"
	"github.com/hekonsek/random-strings"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateTaskDefinition(t *testing.T) {
	t.Parallel()

	// Given
	name := randomstrings.ForHumanWithHash()
	defer func() {
		err := awsom.DeleteEcsTaskDefinition(name)
		assert.NoError(t, err)
	}()

	// When
	err := awsom.NewEcsTaskDefinitionBuilder(name, "hekonsek/http-echo").Create()

	// Then
	assert.NoError(t, err)
}

func TestCreateCluster(t *testing.T) {
	t.Parallel()

	// Given
	name := randomstrings.ForHumanWithHash()
	defer func() {
		err := aws.DeleteVpc(name)
		assert.NoError(t, err)
	}()
	err := aws.NewVpcBuilder(name).Create()
	assert.NoError(t, err)
	defer func() {
		err := awsom.DeleteEcsCluster(name)
		assert.NoError(t, err)
	}()

	// When
	err = awsom.NewEcsClusterBuilder(name).Create()

	// Then
	assert.NoError(t, err)
}

func TestCreateEcsApplication(t *testing.T) {
	t.Parallel()

	// Given
	name := randomstrings.ForHumanWithHash()
	defer func() {
		err := awsom.DeleteEcsTaskDefinition(name)
		assert.NoError(t, err)
	}()
	defer func() {
		err := awsom.DeleteEcsCluster(name)
		assert.NoError(t, err)
		err = aws.DeleteElasticLoadBalancer(name)
		assert.NoError(t, err)
		err = aws.DeleteVpc(name)
		assert.NoError(t, err)
	}()
	err := aws.NewVpcBuilder(name).Create()
	assert.NoError(t, err)
	err = aws.NewElasticLoadBalancer(name).Create()
	assert.NoError(t, err)
	err = awsom.NewEcsClusterBuilder(name).Create()
	assert.NoError(t, err)

	// When
	err = awsom.NewEcsDeploymentBuilder(name, name, "hekonsek/http-echo").Create()

	// Then
	assert.NoError(t, err)
}
