package awsom

import (
	"testing"
)
import "github.com/stretchr/testify/assert"

func TestCreateApplication(t *testing.T) {
	t.Parallel()

	err := (&Application{
		Name:   RandomName(),
		GitUrl: "https://github.com/hekonsek/awsom-spring-rest",
	}).CreateOrUpdate()

	assert.NoError(t, err)
}
