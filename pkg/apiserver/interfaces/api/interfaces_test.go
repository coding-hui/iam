package api

import (
	"testing"

	"gotest.tools/assert"
)

func TestInitAPIBean(t *testing.T) {
	assert.Equal(t, len(InitAPIBean()), 1)
}
