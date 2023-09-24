// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

var client Factory

// Factory defines the repository interface.
type Factory interface {
	Transaction

	UserRepository() UserRepository
	CasbinRepository() CasbinRepository
	ResourceRepository() ResourceRepository
	RoleRepository() RoleRepository
	PolicyRepository() PolicyRepository
	OrganizationRepository() OrganizationRepository
	Close() error
}

// Client return the repository client instance.
func Client() Factory {
	return client
}

// SetClient set the repository client.
func SetClient(factory Factory) {
	client = factory
}
