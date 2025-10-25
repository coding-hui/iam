// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"context"

	"github.com/coding-hui/iam/internal/apiserver/domain/model"
)

// DeviceAuthRepository defines the device authorization repository interface.
type DeviceAuthRepository interface {
	// Create creates a new device authorization.
	Create(ctx context.Context, deviceAuth *model.DeviceAuthorization) error
	// GetByDeviceCode retrieves device authorization by device code.
	GetByDeviceCode(ctx context.Context, deviceCode string) (*model.DeviceAuthorization, error)
	// GetByUserCode retrieves device authorization by user code.
	GetByUserCode(ctx context.Context, userCode string) (*model.DeviceAuthorization, error)
	// Update updates an existing device authorization.
	Update(ctx context.Context, deviceAuth *model.DeviceAuthorization) error
	// Delete deletes a device authorization by device code.
	Delete(ctx context.Context, deviceCode string) error
	// CleanupExpired removes expired device authorizations.
	CleanupExpired(ctx context.Context) error
}
