package awsom

import (
	"github.com/hekonsek/random-strings"
	"testing"
)
import "github.com/stretchr/testify/assert"

func TestCreateApplication(t *testing.T) {
	t.Parallel()

	// Given
	name := randomstrings.ForHumanWithHash()
	defer func() {
		err := DeleteApplication(name)
		assert.NoError(t, err)
		for _, bucket := range []string{
			ConfigureStageName(name) + "-codebuild-artifacts",
			VersionStageName(name) + "-codebuild-artifacts",
			BuildStageName(name) + "-codebuild-artifacts",
			DockerizeStageName(name) + "-codebuild-artifacts"} {
			println("Deleting " + bucket)
			err = DeleteS3Bucket(bucket)
			assert.NoError(t, err)
		}
	}()

	// When
	err := (&Application{
		Name:   name,
		GitUrl: "https://github.com/hekonsek/awsom-spring-rest",
	}).CreateOrUpdate()

	// Then
	assert.NoError(t, err)
}
