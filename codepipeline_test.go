package awsom

import (
	"github.com/hekonsek/random-strings"
	"testing"
)
import "github.com/stretchr/testify/assert"

func TestReadNoSourceIfPipelinesDoesNotExist(t *testing.T) {
	t.Parallel()

	owner, repo, err := ReadPipelineSource(randomstrings.ForHumanWithHash())
	assert.NoError(t, err)
	assert.Empty(t, owner)
	assert.Empty(t, repo)
}

func TestReadPipelineSource(t *testing.T) {
	t.Parallel()

	name := randomstrings.ForHumanWithHash()
	defer func() {
		err := DeleteCodePipeline(name)
		assert.NoError(t, err)

	}()
	err := (&CodePipeline{
		Name:   name,
		GitUrl: "https://github.com/hekonsek/awsom-spring-rest",
	}).CreateOrUpdate()
	assert.NoError(t, err)

	owner, repo, err := ReadPipelineSource(name)
	assert.NoError(t, err)
	assert.NotEmpty(t, owner)
	assert.NotEmpty(t, repo)
}
