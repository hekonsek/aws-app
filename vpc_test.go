package awsom_test

import (
	"github.com/hekonsek/awsom"
	"github.com/hekonsek/random-strings"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateVpc(t *testing.T) {
	t.Parallel()

	// Given
	name := randomstrings.ForHumanWithHash()
	defer func() {
		err := awsom.DeleteVpc(name)
		assert.NoError(t, err)
	}()

	// When
	err := awsom.NewVpcBuilder(name).Create()
	assert.NoError(t, err)

	// Then
	exists, err := awsom.VpcExistsByName(name)
	assert.NoError(t, err)
	assert.True(t, exists)
}

func TestCreateVpcWithCustomCdirPrefix(t *testing.T) {
	t.Parallel()

	// Given
	name := randomstrings.ForHumanWithHash()
	defer func() {
		err := awsom.DeleteVpc(name)
		assert.NoError(t, err)
	}()

	// When
	err := awsom.NewVpcBuilder(name).WithCidrBlockPrefix("15.10").Create()
	assert.NoError(t, err)

	// Then
	exists, err := awsom.VpcExistsByName(name)
	assert.NoError(t, err)
	assert.True(t, exists)
}

func TestCreateVpcWithThreeSubnets(t *testing.T) {
	t.Parallel()

	// Given
	name := randomstrings.ForHumanWithHash()
	defer func() {
		err := awsom.DeleteVpc(name)
		assert.NoError(t, err)
	}()

	// When
	err := awsom.NewVpcBuilder(name).Create()
	assert.NoError(t, err)

	// Then
	subnets, err := awsom.VpcSubnetsByName(name)
	assert.NoError(t, err)
	assert.Len(t, subnets, 3)
}

func TestCreateVpcTwice(t *testing.T) {
	t.Parallel()

	// Given
	name := randomstrings.ForHumanWithHash()
	defer func() {
		err := awsom.DeleteVpc(name)
		assert.NoError(t, err)
	}()
	err := awsom.NewVpcBuilder(name).Create()
	assert.NoError(t, err)

	// When
	err = awsom.NewVpcBuilder(name).Create()

	// Then
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "already exists")
}
