// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package cache

import "github.com/coding-hui/iam/internal/pkg/options"

type Factory interface {
	// Type unique type of the cache
	Type() string
	// Create relevant caches by type
	Create(options options.DynamicOptions, stopCh <-chan struct{}) (Interface, error)
}
