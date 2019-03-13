package awsom

import (
	"testing"
)
import "github.com/stretchr/testify/assert"

func TestCreateApplication(t *testing.T) {
	// Given
	t.Parallel()
	name := GenerateLowercaseName()

	// When
	err := (&Application{
		Name:   name,
		GitUrl: "https://github.com/hekonsek/awsom-spring-rest",
	}).CreateOrUpdate()

	// Then
	assert.NoError(t, err)

	// Clean up
	err = DeleteApplication(name)
	assert.NoError(t, err)
}
