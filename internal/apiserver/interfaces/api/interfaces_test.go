package api

import (
	"gotest.tools/assert"
	"testing"

	"github.com/coding-hui/iam/internal/apiserver/config"
)

func TestInitAPIBean(t *testing.T) {
	assert.Equal(t, len(InitAPIBean(config.Config{})), 3)
}
