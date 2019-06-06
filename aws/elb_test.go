package aws_test

import (
	"github.com/hekonsek/awsom/aws"
	randomstrings "github.com/hekonsek/random-strings"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateLoadBalancer(t *testing.T) {
	t.Parallel()

	// Given
	name := randomstrings.ForHumanWithHash()
	err := aws.NewVpcBuilder(name).Create()
	assert.NoError(t, err)
	defer func() {
		err := aws.DeleteElasticLoadBalancer(name)
		assert.NoError(t, err)
		err = aws.DeleteVpc(name)
		assert.NoError(t, err)
	}()

	// When
	err = aws.NewElasticLoadBalancer(name).Create()

	// Then
	assert.NoError(t, err)
}
