package awsom

import (
	"github.com/hekonsek/awsom/random-strings"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreatePublicVpc(t *testing.T) {
	t.Parallel()

	// Given
	name := randomstrings.GenerateLowercaseNameWithHash()
	defer func() {
		err := DeleteVpc(name)
		assert.NoError(t, err)
	}()

	// When
	err := DefaultVpc(name).CreateOrUpdate()

	// Then
	assert.NoError(t, err)
}
