package aws_test

import (
	"github.com/hekonsek/awsom"
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
		awsom.Warn(aws.DeleteHostedZone(domain))
	}()

	// When
	err := aws.NewHostedZoneBuilder(domain).Create()
	assert.NoError(t, err)

	// Then
	exists, err := aws.HostedZoneExists(domain)
	assert.NoError(t, err)
	assert.True(t, exists)
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
