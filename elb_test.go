package awsom

import (
	"github.com/hekonsek/random-strings"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateLoadBalancer(t *testing.T) {
	t.Parallel()

	// Given
	name := randomstrings.ForHumanWithHash()
	err := DefaultVpc(name).CreateOrUpdate()
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
	err = (&ApplicationLoadBalancer{Name: name}).CreateOrUpdate()

	// Then
	assert.NoError(t, err)
}
