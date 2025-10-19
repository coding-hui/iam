// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package api

import (
	"testing"

	"gotest.tools/assert"

	"github.com/coding-hui/iam/internal/apiserver/config"
)

func TestInitAPIBean(t *testing.T) {
	assert.Equal(t, len(InitAPIBean(config.Config{})), 12)
}
