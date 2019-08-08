package aws_test

import (
	"github.com/hekonsek/awsom/aws"
	randomstrings "github.com/hekonsek/random-strings"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateHostedZone(t *testing.T) {
	t.Parallel()

	// Given
	domain := randomstrings.ForHuman() + ".com"
	defer func() {
		err := aws.DeleteHostedZone(domain)
		assert.NoError(t, err)
	}()

	// When
	err := aws.NewHostedZoneBuilder(domain).Create()
	assert.NoError(t, err)

	// Then
	zoneId, err := aws.HostedZoneIdByDomain(domain)
	assert.NoError(t, err)
	assert.NotEmpty(t, zoneId)
}

func TestTagHostedZone(t *testing.T) {
	t.Parallel()

	// Given
	domain := randomstrings.ForHuman() + ".com"
	defer func() {
		err := aws.DeleteHostedZone(domain)
		assert.NoError(t, err)
	}()
	err := aws.NewHostedZoneBuilder(domain).Create()
	assert.NoError(t, err)

	// When
	err = aws.TagHostedZone(domain, "foo", "bar")
	assert.NoError(t, err)

	// Then
	tags, err := aws.HostedZoneTags(domain)
	assert.NoError(t, err)
	assert.Equal(t, "bar", tags["foo"])
}
