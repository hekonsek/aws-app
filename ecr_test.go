package awsom

import (
	"github.com/hekonsek/awsom/random-strings"
	"testing"
)
import "github.com/stretchr/testify/assert"

func TestRepositoryDoesntExist(t *testing.T) {
	t.Parallel()

	repositoryUri, err := EcrRepositoryExists(randomstrings.GenerateLowercaseNameWithHash())
	assert.NoError(t, err)
	assert.Empty(t, repositoryUri)
}

func TestEnsureRepositoryDoesntExist(t *testing.T) {
	t.Parallel()

	repositoryUri, err := EnsureEcrRepositoryExists(randomstrings.GenerateLowercaseNameWithHash())
	assert.NoError(t, err)
	assert.NotEmpty(t, repositoryUri)
}
