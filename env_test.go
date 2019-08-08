package awsom_test

import (
	"github.com/hekonsek/awsom"
	"github.com/hekonsek/awsom/aws"
	"github.com/hekonsek/random-strings"
	"testing"
)
import "github.com/stretchr/testify/assert"

func TestCreateEnvironment(t *testing.T) {
	t.Parallel()

	// Given
	name := randomstrings.ForHumanWithDashAndHash()
	domain := randomstrings.ForHuman() + ".com"
	defer func() {
		err := aws.DeleteHostedZone(domain)
		assert.NoError(t, err)
	}()
	defer func() {
		err := awsom.DeleteEnv(name)
		assert.NoError(t, err)
	}()

	// When
	err := awsom.NewEnvBuilder(name, domain).Create()

	// Then
	assert.NoError(t, err)
}
