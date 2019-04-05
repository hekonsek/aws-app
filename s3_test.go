package awsom

import (
	"github.com/hekonsek/random-strings"
	"testing"
)
import "github.com/stretchr/testify/assert"

func TestCreateBucket(t *testing.T) {
	t.Parallel()

	// Given
	name := randomstrings.ForHumanWithHash()
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
	exists, err := S3BucketExists(randomstrings.ForHumanWithHash())

	// Then
	assert.NoError(t, err)
	assert.False(t, exists)
}
