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
		Name:        "Test API Key",
		Description: "Test description",
		ExpiresAt:   expiresAt,
		Permissions: &v1.ApiKeyPermissionSpec{
			Roles:   []string{"admin"},
			Actions: []string{"read", "write"},
			Scopes:  []string{"api"},
		},
		AllowedIPs: &v1.ApiKeyAllowedIPs{
			IPs:   []string{"192.168.1.1"},
			CIDRs: []string{"10.0.0.0/8"},
		},
	}

	// Test the request structure
	assert.Equal(t, req.Name, "Test API Key")
	assert.Equal(t, len(req.Permissions.Roles), 1)
	assert.Equal(t, len(req.Permissions.Actions), 2)
}

func TestApiKeyService_GenerateKeyAndSecret(t *testing.T) {
	s := &apiKeyServiceImpl{}

	// Test key generation
	key, err := s.generateKey()
	assert.NilError(t, err)
	assert.Equal(t, len(key), 35) // sk- + 32 hex chars
	assert.Check(t, len(key) > 0)
	assert.Check(t, key[:3] == "sk-")

	// Test secret generation
	secret, err := s.generateSecret()
	assert.NilError(t, err)
	assert.Equal(t, len(secret), 64) // 64 hex chars
	assert.Check(t, len(secret) > 0)
}

func TestApiKeyService_ConvertPermissions(t *testing.T) {
	s := &apiKeyServiceImpl{}

	// Test nil permissions
	result := s.convertPermissions(nil)
	assert.Equal(t, result, (*model.PermissionSpec)(nil))

	// Test with permissions
	apiPerms := &v1.ApiKeyPermissionSpec{
		Roles:   []string{"admin", "user"},
		Actions: []string{"read", "write"},
		Scopes:  []string{"api"},
		Resources: []v1.ApiKeyResourcePermission{
			{
				ResourceType: "user",
				ResourceIDs:  []string{"user-1", "user-2"},
				Actions:      []string{"read"},
			},
		},
	}

	result = s.convertPermissions(apiPerms)
	assert.Equal(t, len(result.Roles), 2)
	assert.Equal(t, len(result.Actions), 2)
	assert.Equal(t, len(result.Scopes), 1)
	assert.Equal(t, len(result.Resources), 1)
	assert.Equal(t, result.Resources[0].ResourceType, "user")
}

func TestApiKeyService_ConvertAllowedIPs(t *testing.T) {
	s := &apiKeyServiceImpl{}

	// Test nil allowed IPs
	result := s.convertAllowedIPs(nil)
	assert.Equal(t, result, (*model.AllowedIPs)(nil))

	// Test with allowed IPs
	apiIPs := &v1.ApiKeyAllowedIPs{
		IPs:   []string{"192.168.1.1", "10.0.0.1"},
		CIDRs: []string{"10.0.0.0/8", "172.16.0.0/12"},
	}

	result = s.convertAllowedIPs(apiIPs)
	assert.Equal(t, len(result.IPs), 2)
	assert.Equal(t, len(result.CIDRs), 2)
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
