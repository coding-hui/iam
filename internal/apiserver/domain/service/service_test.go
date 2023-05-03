// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package service

import (
	"testing"

	"gotest.tools/assert"

	"github.com/coding-hui/iam/cmd/iam-apiserver/app/options"
	"github.com/coding-hui/iam/internal/apiserver/config"
)

func TestInitServiceBean(t *testing.T) {
	cfg, err := config.CreateConfigFromOptions(options.NewOptions())
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, len(InitServiceBean(*cfg)), 4)
}
