// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import (
	"net"
	"time"

	"gorm.io/gorm"

	metav1 "github.com/coding-hui/common/meta/v1"
	"github.com/coding-hui/common/util/idutil"
)

func init() {
	RegisterModel(&ApiKey{})
}

// ApiKeyStatus represents the status of an API Key.
type ApiKeyStatus int

const (
	// ApiKeyStatusActive indicates the API Key is active and can be used.
	ApiKeyStatusActive ApiKeyStatus = 1
	// ApiKeyStatusInactive indicates the API Key is inactive and cannot be used.
	ApiKeyStatusInactive ApiKeyStatus = 0
	// ApiKeyStatusExpired indicates the API Key has expired.
	ApiKeyStatusExpired ApiKeyStatus = 2
)

// ApiKey represents an API Key for programmatic access.
type ApiKey struct {
	// Standard object's metadata.
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Name is the human-readable name of the API Key.
	Name string `json:"name" gorm:"column:name;type:varchar(128);not null"`

	// Key is the API Key identifier (public part).
	Key string `json:"key" gorm:"column:key;type:varchar(256);uniqueIndex;not null"`

	// Secret is the API Secret (private part, encrypted storage).
	Secret string `json:"-" gorm:"column:secret;type:varchar(512);not null"`

	// UserID is the ID of the user who owns this API Key.
	UserID string `json:"userId" gorm:"column:user_id;type:varchar(64);not null;index"`

	// ExpiresAt is the expiration time of the API Key.
	ExpiresAt *time.Time `json:"expiresAt,omitempty" gorm:"column:expires_at"`

	// Status indicates the current status of the API Key.
	Status ApiKeyStatus `json:"status" gorm:"column:status;type:tinyint;default:1"`

	// Permissions defines the scope of permissions for this API Key.
	Permissions *PermissionSpec `json:"permissions,omitempty" gorm:"column:permissions;type:json"`

	// LastUsedAt is the timestamp when this API Key was last used.
	LastUsedAt *time.Time `json:"lastUsedAt,omitempty" gorm:"column:last_used_at"`

	// UsageCount tracks how many times this API Key has been used.
	UsageCount int64 `json:"usageCount" gorm:"column:usage_count;default:0"`

	// AllowedIPs specifies the allowed IP addresses for this API Key.
	AllowedIPs *AllowedIPs `json:"allowedIps,omitempty" gorm:"column:allowed_ips;type:json"`

	// Description provides additional information about the API Key.
	Description string `json:"description,omitempty" gorm:"column:description;type:text"`
}

// PermissionSpec defines the permission scope for an API Key.
type PermissionSpec struct {
	// Roles contains the role names that this API Key can assume.
	Roles []string `json:"roles,omitempty"`

	// Resources defines the resources this API Key can access.
	Resources []ResourcePermission `json:"resources,omitempty"`

	// Actions defines the actions this API Key can perform.
	Actions []string `json:"actions,omitempty"`

	// Scopes defines OAuth2-like scopes for the API Key.
	Scopes []string `json:"scopes,omitempty"`
}

// ResourcePermission defines permission for a specific resource.
type ResourcePermission struct {
	// ResourceType is the type of resource (e.g., "user", "policy").
	ResourceType string `json:"resourceType"`

	// ResourceIDs are the specific resource IDs, empty means all resources.
	ResourceIDs []string `json:"resourceIds,omitempty"`

	// Actions are the allowed actions on the resource.
	Actions []string `json:"actions"`
}

// AllowedIPs defines IP address restrictions for an API Key.
type AllowedIPs struct {
	// IPs contains the allowed IP addresses in CIDR notation.
	IPs []string `json:"ips,omitempty"`

	// CIDRs contains the allowed CIDR blocks.
	CIDRs []string `json:"cidrs,omitempty"`
}

// TableName maps to mysql table name.
func (a *ApiKey) TableName() string {
	return TableNamePrefix + "apikey"
}

// AfterCreate run after create database record.
func (a *ApiKey) AfterCreate(tx *gorm.DB) error {
	a.InstanceID = idutil.GetInstanceID(a.ID, "apikey-")
	return tx.Save(a).Error
}

// IsActive returns true if the API Key is active and not expired.
func (a *ApiKey) IsActive() bool {
	if a.Status != ApiKeyStatusActive {
		return false
	}

	if a.ExpiresAt != nil && a.ExpiresAt.Before(time.Now()) {
		return false
	}

	return true
}

// MarkUsed updates the usage statistics.
func (a *ApiKey) MarkUsed() {
	now := time.Now()
	a.LastUsedAt = &now
	a.UsageCount++
}

// ValidateIP checks if the given IP address is allowed.
func (a *ApiKey) ValidateIP(ip string) bool {
	if a.AllowedIPs == nil {
		return true
	}

	// Parse the input IP
	inputIP := net.ParseIP(ip)
	if inputIP == nil {
		return false // Invalid IP format
	}

	// Check against allowed IPs
	for _, allowedIP := range a.AllowedIPs.IPs {
		// Handle single IP addresses
		if allowedIP == ip {
			return true
		}

		// Handle CIDR notation
		_, cidr, err := net.ParseCIDR(allowedIP)
		if err == nil && cidr.Contains(inputIP) {
			return true
		}
	}

	// Check against allowed CIDR blocks
	for _, cidrStr := range a.AllowedIPs.CIDRs {
		_, cidr, err := net.ParseCIDR(cidrStr)
		if err != nil {
			continue // Skip invalid CIDR
		}
		if cidr.Contains(inputIP) {
			return true
		}
	}

	return false
}
