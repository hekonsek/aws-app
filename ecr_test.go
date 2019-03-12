package awsom

import (
	"testing"
)
import "github.com/stretchr/testify/assert"

func TestRepositoryDoesntExist(t *testing.T) {
	t.Parallel()

	repositoryUri, err := EcrRepositoryExists(RandomName())
	assert.NoError(t, err)
	assert.Empty(t, repositoryUri)
}

func TestEnsureRepositoryDoesntExist(t *testing.T) {
	t.Parallel()

	repositoryUri, err := EnsureEcrRepositoryExists(RandomName())
	assert.NoError(t, err)
	assert.NotEmpty(t, repositoryUri)
}
