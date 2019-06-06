package awsom_test

import (
	"github.com/hekonsek/awsom"
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
		err := awsom.DeleteVpc(name)
		assert.NoError(t, err)
	}()
	err := awsom.NewVpcBuilder(name).Create()
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
	}()
	defer func() {
		err := awsom.DeleteLoadBalancer(name)
		assert.NoError(t, err)
	}()
	defer func() {
		err := awsom.DeleteVpc(name)
		assert.NoError(t, err)
	}()
	err := awsom.NewVpcBuilder(name).Create()
	assert.NoError(t, err)
	err = (&awsom.ApplicationLoadBalancerBuilder{Name: name}).Create()
	assert.NoError(t, err)
	err = awsom.NewEcsClusterBuilder(name).Create()
	assert.NoError(t, err)

	// When
	err = awsom.NewEcsDeploymentBuilder(name, name, "hekonsek/http-echo").Create()

	// Then
	assert.NoError(t, err)
}
