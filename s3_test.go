package awsom

import (
	"testing"
)
import "github.com/stretchr/testify/assert"

func TestCreateBucket(t *testing.T) {
	// Given
	t.Parallel()
	name := RandomName()

	// When
	err := (&S3Bucket{
		Name: name,
	}).CreateOrUpdate()
	assert.NoError(t, err)

	// Then
	exists, err := S3BucketExists(name)
	assert.NoError(t, err)
	assert.True(t, exists)

	// Clean up
	err = DeleteS3Bucket(name)
	assert.NoError(t, err)
}
