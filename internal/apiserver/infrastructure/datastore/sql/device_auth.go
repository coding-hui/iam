// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package sql

import (
	"context"
	"time"

	"gorm.io/gorm"

	"github.com/coding-hui/iam/internal/apiserver/domain/model"
	"github.com/coding-hui/iam/internal/apiserver/domain/repository"
)

// deviceAuthRepository represents the device authorization repository implementation.
type deviceAuthRepository struct {
	db *gorm.DB
}

// newDeviceAuthRepository creates a new device authorization repository.
func newDeviceAuthRepository(db *gorm.DB) repository.DeviceAuthRepository {
	return &deviceAuthRepository{db: db}
}

// Create creates a new device authorization.
func (r *deviceAuthRepository) Create(ctx context.Context, deviceAuth *model.DeviceAuthorization) error {
	return r.db.Create(deviceAuth).Error
}

// GetByDeviceCode retrieves device authorization by device code.
func (r *deviceAuthRepository) GetByDeviceCode(ctx context.Context, deviceCode string) (*model.DeviceAuthorization, error) {
	var deviceAuth model.DeviceAuthorization
	if err := r.db.Where("device_code = ?", deviceCode).First(&deviceAuth).Error; err != nil {
		return nil, err
	}
	return &deviceAuth, nil
}

// GetByUserCode retrieves device authorization by user code.
func (r *deviceAuthRepository) GetByUserCode(ctx context.Context, userCode string) (*model.DeviceAuthorization, error) {
	var deviceAuth model.DeviceAuthorization
	if err := r.db.Where("user_code = ?", userCode).First(&deviceAuth).Error; err != nil {
		return nil, err
	}
	return &deviceAuth, nil
}

// Update updates an existing device authorization.
func (r *deviceAuthRepository) Update(ctx context.Context, deviceAuth *model.DeviceAuthorization) error {
	return r.db.Save(deviceAuth).Error
}

// Delete deletes a device authorization by device code.
func (r *deviceAuthRepository) Delete(ctx context.Context, deviceCode string) error {
	return r.db.Where("device_code = ?", deviceCode).Delete(&model.DeviceAuthorization{}).Error
}

// CleanupExpired removes expired device authorizations.
func (r *deviceAuthRepository) CleanupExpired(ctx context.Context) error {
	return r.db.Where("expires_at < ?", time.Now()).Delete(&model.DeviceAuthorization{}).Error
}
