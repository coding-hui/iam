// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"context"

	"github.com/coding-hui/iam/internal/apiserver/domain/model"
	v1 "github.com/coding-hui/iam/pkg/api/apiserver/v1"

	metav1 "github.com/coding-hui/common/meta/v1"
)

// ApiKeyRepository defines the API Key repository interface.
type ApiKeyRepository interface {
	// Create creates a new API Key.
	Create(ctx context.Context, apiKey *model.ApiKey, opts metav1.CreateOptions) (*model.ApiKey, error)

	// Update updates an existing API Key.
	Update(ctx context.Context, apiKey *model.ApiKey, opts metav1.UpdateOptions) error

	// DeleteByInstanceId deletes an API Key by instance ID.
	DeleteByInstanceId(ctx context.Context, instanceId string, opts metav1.DeleteOptions) error

	// BatchDelete deletes multiple API Keys by instance IDs.
	BatchDelete(ctx context.Context, instanceIds []string, opts metav1.DeleteOptions) error

	// GetByInstanceId retrieves an API Key by instance ID.
	GetByInstanceId(ctx context.Context, instanceId string, opts metav1.GetOptions) (*model.ApiKey, error)

	// GetByKey retrieves an API Key by key value.
	GetByKey(ctx context.Context, key string, opts metav1.GetOptions) (*model.ApiKey, error)

	// GetByUser retrieves API Keys for a specific user.
	GetByUser(ctx context.Context, userId string, opts metav1.GetOptions) ([]model.ApiKey, error)

	// List retrieves a list of API Keys with pagination and filtering.
	List(ctx context.Context, opts v1.ListApiKeyOptions) ([]model.ApiKey, error)

	// Count returns the total count of API Keys matching the filter.
	Count(ctx context.Context, opts v1.ListApiKeyOptions) (int64, error)

	// UpdateUsage updates the usage statistics for an API Key.
	UpdateUsage(ctx context.Context, instanceId string) error

	// CleanupExpired deletes expired API Keys.
	CleanupExpired(ctx context.Context) error
}
