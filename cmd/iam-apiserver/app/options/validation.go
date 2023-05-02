// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package options

import utilerrors "k8s.io/apimachinery/pkg/util/errors"

// Validate validates iam-apiserver run options, to find options' misconfiguration.
func (s *ServerRunOptions) Validate() error {
	var errors []error

	errors = append(errors, s.GenericServerRunOptions.Validate()...)

	return utilerrors.NewAggregate(errors)
}
