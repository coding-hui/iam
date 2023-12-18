// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package cache

import (
	"fmt"

	"github.com/coding-hui/iam/pkg/options"
)

type Options struct {
	Type    string                 `json:"type"    mapstructure:"type"`
	Options options.DynamicOptions `json:"options" mapstructure:"options"`
}

// NewCacheOptions returns options points to nowhere,
// because redis is not required for some components.
func NewCacheOptions() *Options {
	return &Options{
		Type:    TypeInMemoryCache,
		Options: map[string]interface{}{},
	}
}

// Validate check options.
func (r *Options) Validate() []error {
	errors := make([]error, 0)

	if r.Type == "" {
		errors = append(errors, fmt.Errorf("invalid cache type"))
	}

	return errors
}
