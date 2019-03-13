package awsom

import (
	"github.com/stretchr/testify/assert"
	"testing"
)
import "github.com/go-errors/errors"

func TestGenerateLowercaseNameWithHash(t *testing.T) {
	t.Parallel()

	name := GenerateLowercaseNameWithHash()
	assert.NotEmpty(t, name)
}

func TestCliCapture(t *testing.T) {
	output, err := CliCapture(func() error {
		CliError(errors.New("someError"))
		return nil
	})
	assert.NoError(t, err)
	assert.Equal(t, "Something went wrong: someError", output)
}

func TestExec(t *testing.T) {
	t.Parallel()

	output, err := Exec{Command: "echo foo\nbar"}.Run()
	assert.NoError(t, err)
	assert.Equal(t, []string{"foo", "bar"}, output)
}
