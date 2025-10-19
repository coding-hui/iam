// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package service

import (
	"testing"

	"gotest.tools/assert"

	"github.com/coding-hui/iam/internal/apiserver/domain/model"
)

func TestApiKey_ValidateIP(t *testing.T) {
	// Test with no IP restrictions
	apiKey := &model.ApiKey{
		AllowedIPs: nil,
	}
	assert.Equal(t, apiKey.ValidateIP("192.168.1.1"), true)
	assert.Equal(t, apiKey.ValidateIP("10.0.0.1"), true)

	// Test with specific IP addresses
	apiKey.AllowedIPs = &model.AllowedIPs{
		IPs: []string{"192.168.1.1", "10.0.0.1"},
	}
	assert.Equal(t, apiKey.ValidateIP("192.168.1.1"), true)
	assert.Equal(t, apiKey.ValidateIP("10.0.0.1"), true)
	assert.Equal(t, apiKey.ValidateIP("192.168.1.2"), false)
	assert.Equal(t, apiKey.ValidateIP("invalid-ip"), false)

	// Test with CIDR blocks
	apiKey.AllowedIPs = &model.AllowedIPs{
		CIDRs: []string{"192.168.1.0/24", "10.0.0.0/8"},
	}
	assert.Equal(t, apiKey.ValidateIP("192.168.1.100"), true)
	assert.Equal(t, apiKey.ValidateIP("10.0.0.1"), true)
	assert.Equal(t, apiKey.ValidateIP("10.255.255.254"), true)
	assert.Equal(t, apiKey.ValidateIP("192.168.2.1"), false)
	assert.Equal(t, apiKey.ValidateIP("172.16.0.1"), false)

	// Test with both IPs and CIDRs
	apiKey.AllowedIPs = &model.AllowedIPs{
		IPs:   []string{"192.168.1.100"},
		CIDRs: []string{"10.0.0.0/8"},
	}
	assert.Equal(t, apiKey.ValidateIP("192.168.1.100"), true)
	assert.Equal(t, apiKey.ValidateIP("10.0.0.1"), true)
	assert.Equal(t, apiKey.ValidateIP("192.168.1.101"), false)
	assert.Equal(t, apiKey.ValidateIP("172.16.0.1"), false)

	// Test invalid CIDR format
	apiKey.AllowedIPs = &model.AllowedIPs{
		CIDRs: []string{"invalid-cidr"},
	}
	assert.Equal(t, apiKey.ValidateIP("192.168.1.1"), false)
}
