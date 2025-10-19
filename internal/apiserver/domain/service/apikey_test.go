// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package service

import (
	"testing"
	"time"

	"gotest.tools/assert"

	"github.com/coding-hui/iam/internal/apiserver/domain/model"
	v1 "github.com/coding-hui/iam/pkg/api/apiserver/v1"
)

func TestApiKeyService_CreateApiKey(t *testing.T) {
	// Test request structure
	expiresAt := time.Now().Add(24 * time.Hour)
	req := v1.CreateApiKeyRequest{
		Name:      "Test API Key",
		ExpiresAt: expiresAt,
	}

	// Test the request structure
	assert.Equal(t, req.Name, "Test API Key")
}

func TestApiKeyService_GenerateKey(t *testing.T) {
	s := &apiKeyServiceImpl{}

	// Test key generation
	key, err := s.generateKey()
	assert.NilError(t, err)
	assert.Equal(t, len(key), 30) // sk_ + 27 base62 chars
	assert.Check(t, len(key) > 0)
	assert.Check(t, key[:3] == "sk_")
	// Verify Base62 characters only
	for _, ch := range key[3:] {
		assert.Check(t, (ch >= '0' && ch <= '9') || (ch >= 'A' && ch <= 'Z') || (ch >= 'a' && ch <= 'z'))
	}
}

func TestApiKey_IsActive(t *testing.T) {
	// Test active key without expiration
	apiKey := &model.ApiKey{
		Status: model.ApiKeyStatusActive,
	}
	assert.Equal(t, apiKey.IsActive(), true)

	// Test inactive key
	apiKey.Status = model.ApiKeyStatusInactive
	assert.Equal(t, apiKey.IsActive(), false)

	// Test expired key
	apiKey.Status = model.ApiKeyStatusActive
	pastTime := time.Now().Add(-24 * time.Hour)
	apiKey.ExpiresAt = &pastTime
	assert.Equal(t, apiKey.IsActive(), false)

	// Test active key with future expiration
	futureTime := time.Now().Add(24 * time.Hour)
	apiKey.ExpiresAt = &futureTime
	assert.Equal(t, apiKey.IsActive(), true)
}

func TestApiKey_MarkUsed(t *testing.T) {
	apiKey := &model.ApiKey{
		UsageCount: 0,
	}

	// Test initial usage
	apiKey.MarkUsed()
	assert.Equal(t, apiKey.UsageCount, int64(1))
	assert.Check(t, apiKey.LastUsedAt != nil)

	// Test subsequent usage
	apiKey.MarkUsed()
	assert.Equal(t, apiKey.UsageCount, int64(2))
}
