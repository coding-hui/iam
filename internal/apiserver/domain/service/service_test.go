// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package service

import (
	"gotest.tools/assert"
	"testing"

	"github.com/coding-hui/iam/internal/apiserver/config"
)

func TestInitServiceBean(t *testing.T) {
	assert.Equal(t, len(InitServiceBean(config.Config{})), 3)
}
