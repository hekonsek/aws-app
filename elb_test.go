package awsom

import (
	randomstrings "github.com/hekonsek/random-strings"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateLoadBalancer(t *testing.T) {
	t.Parallel()

	// Given
	name := randomstrings.ForHumanWithHash()
	err := NewVpcBuilder(name).Create()
	assert.NoError(t, err)
	defer func() {
		err := DeleteLoadBalancer(name)
		assert.NoError(t, err)
	}()
	defer func() {
		err = DeleteVpc(name)
		assert.NoError(t, err)
	}()

	// When
	err = (&ApplicationLoadBalancerBuilder{Name: name}).Create()

	// Then
	assert.NoError(t, err)
}
