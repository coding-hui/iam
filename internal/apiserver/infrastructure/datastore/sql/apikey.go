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
	"github.com/coding-hui/iam/internal/apiserver/infrastructure/datastore"
	v1 "github.com/coding-hui/iam/pkg/api/apiserver/v1"
	"github.com/coding-hui/iam/pkg/code"

	"github.com/coding-hui/common/errors"
	metav1 "github.com/coding-hui/common/meta/v1"
)

type apiKeyRepository struct {
	client *Client
}

func newApiKeyRepository(client *Client) repository.ApiKeyRepository {
	return &apiKeyRepository{client: client}
}

// Create creates a new API Key.
func (r *apiKeyRepository) Create(ctx context.Context, apiKey *model.ApiKey, opts metav1.CreateOptions) (*model.ApiKey, error) {
	if err := r.client.WithCtx(ctx).Create(apiKey).Error; err != nil {
		return nil, errors.WithCode(code.ErrDatabaseCreate, "Failed to create API Key: %s", err.Error())
	}

	return apiKey, nil
}

// Update updates an existing API Key.
func (r *apiKeyRepository) Update(ctx context.Context, apiKey *model.ApiKey, opts metav1.UpdateOptions) error {
	if err := r.client.WithCtx(ctx).Save(apiKey).Error; err != nil {
		return errors.WithCode(code.ErrDatabaseUpdate, "Failed to update API Key: %s", err.Error())
	}

	return nil
}

// DeleteByInstanceId deletes an API Key by instance ID.
func (r *apiKeyRepository) DeleteByInstanceId(ctx context.Context, instanceId string, opts metav1.DeleteOptions) error {
	if err := r.client.WithCtx(ctx).Where("instance_id = ?", instanceId).Delete(&model.ApiKey{}).Error; err != nil {
		return errors.WithCode(code.ErrDatabaseDelete, "Failed to delete API Key: %s", err.Error())
	}

	return nil
}

// BatchDelete deletes multiple API Keys by instance IDs.
func (r *apiKeyRepository) BatchDelete(ctx context.Context, instanceIds []string, opts metav1.DeleteOptions) error {
	if err := r.client.WithCtx(ctx).Where("instance_id IN ?", instanceIds).Delete(&model.ApiKey{}).Error; err != nil {
		return errors.WithCode(code.ErrDatabaseDelete, "Failed to delete API Keys: %s", err.Error())
	}

	return nil
}

// GetByInstanceId retrieves an API Key by instance ID.
func (r *apiKeyRepository) GetByInstanceId(ctx context.Context, instanceId string, opts metav1.GetOptions) (*model.ApiKey, error) {
	var apiKey model.ApiKey
	if err := r.client.WithCtx(ctx).Where("instance_id = ?", instanceId).First(&apiKey).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithCode(code.ErrApiKeyNotFound, "API Key not found")
		}
		return nil, errors.WithCode(code.ErrDatabaseQuery, "Failed to query API Key: %s", err.Error())
	}

	return &apiKey, nil
}

// GetByKey retrieves an API Key by key value.
func (r *apiKeyRepository) GetByKey(ctx context.Context, key string, opts metav1.GetOptions) (*model.ApiKey, error) {
	var apiKey model.ApiKey
	if err := r.client.WithCtx(ctx).Where("`key` = ?", key).First(&apiKey).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithCode(code.ErrApiKeyNotFound, "API Key not found")
		}
		return nil, errors.WithCode(code.ErrDatabaseQuery, "Failed to query API Key: %s", err.Error())
	}

	return &apiKey, nil
}

// GetByUser retrieves API Keys for a specific user.
func (r *apiKeyRepository) GetByUser(ctx context.Context, userId string, opts metav1.GetOptions) ([]model.ApiKey, error) {
	var apiKeys []model.ApiKey
	if err := r.client.WithCtx(ctx).Where("user_id = ?", userId).Find(&apiKeys).Error; err != nil {
		return nil, errors.WithCode(code.ErrDatabaseQuery, "Failed to query user API Keys: %s", err.Error())
	}

	return apiKeys, nil
}

// List retrieves a list of API Keys with pagination and filtering.
func (r *apiKeyRepository) List(ctx context.Context, opts v1.ListApiKeyOptions) ([]model.ApiKey, error) {
	var apiKeys []model.ApiKey

	db := r.client.WithCtx(ctx)

	// Apply filters
	if opts.UserID != "" {
		db = db.Where("user_id = ?", opts.UserID)
	}
	if opts.Status != 0 {
		db = db.Where("status = ?", opts.Status)
	}

	// Apply field selector and pagination
	db = db.Scopes(makeCondition(opts.ListOptions), paginate(opts.ListOptions))

	if err := db.Order("created_at DESC").Find(&apiKeys).Error; err != nil {
		return nil, errors.WithCode(code.ErrDatabaseQuery, "Failed to list API Keys: %s", err.Error())
	}

	return apiKeys, nil
}

// Count returns the total count of API Keys matching the filter.
func (r *apiKeyRepository) Count(ctx context.Context, opts v1.ListApiKeyOptions) (int64, error) {
	var count int64

	db := r.client.WithCtx(ctx).Model(&model.ApiKey{})

	// Apply filters
	if opts.UserID != "" {
		db = db.Where("user_id = ?", opts.UserID)
	}
	if opts.Status != 0 {
		db = db.Where("status = ?", opts.Status)
	}

	// Apply field selector
	db = db.Scopes(makeCondition(opts.ListOptions))

	if err := db.Count(&count).Error; err != nil {
		return 0, datastore.NewDBError(err, "count API Keys")
	}

	return count, nil
}

// UpdateUsage updates the usage statistics for an API Key.
func (r *apiKeyRepository) UpdateUsage(ctx context.Context, instanceId string) error {
	now := time.Now()
	return r.client.WithCtx(ctx).Model(&model.ApiKey{}).
		Where("instance_id = ?", instanceId).
		Updates(map[string]interface{}{
			"last_used_at": now,
			"usage_count":  gorm.Expr("usage_count + 1"),
		}).Error
}

// CleanupExpired deletes expired API Keys.
func (r *apiKeyRepository) CleanupExpired(ctx context.Context) error {
	return r.client.WithCtx(ctx).
		Where("expires_at IS NOT NULL AND expires_at < ?", time.Now()).
		Delete(&model.ApiKey{}).Error
}
