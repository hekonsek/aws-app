package awsom

import (
	"testing"
)
import "github.com/stretchr/testify/assert"

func TestCreateApplication(t *testing.T) {
	t.Parallel()

	err := (&Application{
		Name: RandomName(),
	}).CreateOrUpdate()

	assert.NoError(t, err)
}
