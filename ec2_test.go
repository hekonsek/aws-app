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
	err := (&Vpc{
		Name:      name,
		CidrBlock: "10.0.0.0/16",
		Subnets: []Subnet{
			{Cidr: "10.0.0.0/18"},
		},
	}).CreateOrUpdate()

	// Then
	assert.NoError(t, err)
}
