package awsom

import (
	"github.com/hekonsek/random-strings"
	"testing"
)
import "github.com/stretchr/testify/assert"

func TestRepositoryDoesntExist(t *testing.T) {
	t.Parallel()

	repositoryUri, err := EcrRepositoryExists(randomstrings.ForHumanWithHash())
	assert.NoError(t, err)
	assert.Empty(t, repositoryUri)
}

func TestEnsureRepositoryDoesntExist(t *testing.T) {
	t.Parallel()

	repositoryUri, err := EnsureEcrRepositoryExists(randomstrings.ForHumanWithHash())
	assert.NoError(t, err)
	assert.NotEmpty(t, repositoryUri)
}
