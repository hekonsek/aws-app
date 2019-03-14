package awsom

import (
	"github.com/hekonsek/awsom/random-strings"
	"testing"
)
import "github.com/stretchr/testify/assert"

func TestCreateBucket(t *testing.T) {
	t.Parallel()

	// Given
	name := randomstrings.GenerateLowercaseNameWithHash()
	defer func() {
		err := DeleteS3Bucket(name)
		assert.NoError(t, err)
	}()

	// When
	err := (&S3Bucket{
		Name: name,
	}).CreateOrUpdate()
	assert.NoError(t, err)

	// Then
	exists, err := S3BucketExists(name)
	assert.NoError(t, err)
	assert.True(t, exists)
}

func TestDetectingNonExistingBucket(t *testing.T) {
	t.Parallel()

	// When
	exists, err := S3BucketExists(randomstrings.GenerateLowercaseNameWithHash())

	// Then
	assert.NoError(t, err)
	assert.False(t, exists)
}
