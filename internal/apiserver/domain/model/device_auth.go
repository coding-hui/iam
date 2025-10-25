package model

import (
	"time"

	metav1 "github.com/coding-hui/common/meta/v1"
)

// DeviceAuthStatus represents the status of device authorization.
type DeviceAuthStatus string

const (
	// DeviceAuthPending indicates the device authorization is pending user approval.
	DeviceAuthPending DeviceAuthStatus = "pending"
	// DeviceAuthApproved indicates the device authorization has been approved.
	DeviceAuthApproved DeviceAuthStatus = "approved"
	// DeviceAuthDenied indicates the device authorization has been denied.
	DeviceAuthDenied DeviceAuthStatus = "denied"
	// DeviceAuthExpired indicates the device authorization has expired.
	DeviceAuthExpired DeviceAuthStatus = "expired"
)

// DeviceAuthorization represents a device authorization request.
type DeviceAuthorization struct {
	// Standard object's metadata.
	metav1.ObjectMeta `json:"metadata,omitempty"`
	// DeviceCode is the unique device verification code.
	DeviceCode string `json:"deviceCode"         gorm:"column:device_code;type:varchar(64);uniqueIndex"`
	// UserCode is the user-friendly verification code.
	UserCode string `json:"userCode"           gorm:"column:user_code;type:varchar(8);uniqueIndex"`
	// ClientID is the OAuth client identifier.
	ClientID string `json:"clientID"           gorm:"column:client_id;type:varchar(128)"`
	// Scope is the OAuth scope requested.
	Scope string `json:"scope"              gorm:"column:scope;type:varchar(256)"`
	// UserID is the ID of the user who approved the authorization.
	UserID string `json:"userID"             gorm:"column:user_id;type:varchar(64)"`
	// Status indicates the current status of the authorization.
	Status DeviceAuthStatus `json:"status"             gorm:"column:status;type:varchar(32)"`
	// ExpiresAt is the expiration time of the device code.
	ExpiresAt time.Time `json:"expiresAt"          gorm:"column:expires_at"`
	// ApprovedAt is the time when the authorization was approved.
	ApprovedAt *time.Time `json:"approvedAt"         gorm:"column:approved_at"`
}

// TableName returns the table name for DeviceAuthorization.
func (d *DeviceAuthorization) TableName() string {
	return "device_authorizations"
}

func init() {
	RegisterModel(&DeviceAuthorization{})
}
