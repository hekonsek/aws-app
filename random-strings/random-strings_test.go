package randomstrings

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerateLowercaseNameWithHash(t *testing.T) {
	t.Parallel()

	name := GenerateLowercaseNameWithHash()
	assert.NotEmpty(t, name)
}

func TestGenerateLowercaseNameLongerThanEightCharacter(t *testing.T) {
	t.Parallel()

	name := GenerateLowercaseName()
	assert.True(t, len(name) >= 8)
}
