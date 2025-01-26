package user

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_randString(t *testing.T) {
	l := 5

	str, err := randString(l)
	assert.NoError(t, err)

	assert.Len(t, str, l)
}
