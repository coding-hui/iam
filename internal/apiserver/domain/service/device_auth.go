// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/coding-hui/common/errors"
	metav1 "github.com/coding-hui/common/meta/v1"

	"github.com/coding-hui/iam/internal/apiserver/config"
	"github.com/coding-hui/iam/internal/apiserver/domain/model"
	"github.com/coding-hui/iam/internal/apiserver/domain/repository"
	"github.com/coding-hui/iam/internal/pkg/token"
	v1 "github.com/coding-hui/iam/pkg/api/apiserver/v1"
	"github.com/coding-hui/iam/pkg/code"
)

// DeviceAuthService defines the interface for device authorization service.
type DeviceAuthService interface {
	// CreateDeviceAuthorization creates a device authorization request.
	CreateDeviceAuthorization(ctx context.Context, req *v1.DeviceAuthorizationRequest) (*v1.DeviceAuthorizationResponse, error)
	// GetDeviceToken retrieves access token using device code.
	GetDeviceToken(ctx context.Context, req *v1.DeviceTokenRequest) (*v1.DeviceTokenResponse, error)
	// VerifyUserAuthorization verifies user authorization for device flow.
	VerifyUserAuthorization(ctx context.Context, req *v1.VerifyDeviceRequest) error
}

type deviceAuthServiceImpl struct {
	Store    repository.Factory
	cfg      config.Config
	TokenSvc TokenService `inject:"tokenService"`
}

// NewDeviceAuthService creates a new device authorization service.
func NewDeviceAuthService(c config.Config) DeviceAuthService {
	return &deviceAuthServiceImpl{
		cfg: c,
	}
}

// CreateDeviceAuthorization implements DeviceAuthService.
func (d *deviceAuthServiceImpl) CreateDeviceAuthorization(ctx context.Context, req *v1.DeviceAuthorizationRequest) (*v1.DeviceAuthorizationResponse, error) {
	// Validate client
	_, err := d.cfg.AuthenticationOptions.OAuthOptions.OAuthClient(req.ClientID)
	if err != nil {
		return nil, errors.WithCode(code.ErrClientNotFound, "client not found")
	}

	// Generate device code and user code
	deviceCode, err := generateRandomString(32)
	if err != nil {
		return nil, errors.WithCode(code.ErrUnknown, "failed to generate device code")
	}

	userCode, err := generateUserCode()
	if err != nil {
		return nil, errors.WithCode(code.ErrUnknown, "failed to generate user code")
	}

	// Save device authorization
	deviceAuth := &model.DeviceAuthorization{
		DeviceCode: deviceCode,
		UserCode:   userCode,
		ClientID:   req.ClientID,
		Scope:      req.Scope,
		Status:     model.DeviceAuthPending,
		ExpiresAt:  time.Now().Add(10 * time.Minute), // 10 minutes expiration
	}

	if err := d.Store.DeviceAuthRepository().Create(ctx, deviceAuth); err != nil {
		return nil, errors.WithCode(code.ErrDatabase, "failed to create device authorization")
	}

	// Build verification URI using the insecure serving address
	var baseURL string
	if d.cfg.InsecureServing.BindPort > 0 {
		bindAddress := d.cfg.InsecureServing.BindAddress
		if bindAddress == "0.0.0.0" {
			bindAddress = "localhost"
		}
		baseURL = fmt.Sprintf("http://%s:%d", bindAddress, d.cfg.InsecureServing.BindPort)
	} else if d.cfg.SecureServing.BindPort > 0 {
		bindAddress := d.cfg.SecureServing.BindAddress
		if bindAddress == "0.0.0.0" {
			bindAddress = "localhost"
		}
		baseURL = fmt.Sprintf("https://%s:%d", bindAddress, d.cfg.SecureServing.BindPort)
	} else {
		// Fallback to relative path if no server address available
		baseURL = ""
	}

	verificationURI := baseURL + "/api/v1/device/authorize"

	return &v1.DeviceAuthorizationResponse{
		DeviceCode:              deviceCode,
		UserCode:                userCode,
		VerificationURI:         verificationURI,
		VerificationURIComplete: fmt.Sprintf("%s?user_code=%s", verificationURI, userCode),
		ExpiresIn:               600, // 10 minutes in seconds
		Interval:                5,   // 5 seconds polling interval
	}, nil
}

// GetDeviceToken implements DeviceAuthService.
func (d *deviceAuthServiceImpl) GetDeviceToken(ctx context.Context, req *v1.DeviceTokenRequest) (*v1.DeviceTokenResponse, error) {
	// Validate device code
	deviceAuth, err := d.Store.DeviceAuthRepository().GetByDeviceCode(ctx, req.DeviceCode)
	if err != nil {
		return nil, errors.WithCode(code.ErrDeviceCodeInvalid, "invalid device code")
	}

	// Check status
	if deviceAuth.Status != model.DeviceAuthApproved {
		return nil, errors.WithCode(code.ErrAuthorizationPending, "authorization pending")
	}

	if deviceAuth.ExpiresAt.Before(time.Now()) {
		return nil, errors.WithCode(code.ErrDeviceCodeExpired, "device code expired")
	}

	// Generate access token using TokenService
	tokenReq := &token.IssueRequest{
		User: v1.UserBase{
			ObjectMeta: metav1.ObjectMeta{
				Name: deviceAuth.UserID,
			},
		},
		ExpiresIn: 3600 * time.Second, // 1 hour
	}

	accessToken, err := d.TokenSvc.IssueTo(tokenReq)
	if err != nil {
		return nil, errors.WithCode(code.ErrUnknown, "failed to issue token")
	}

	// Clean up device authorization record
	_ = d.Store.DeviceAuthRepository().Delete(ctx, deviceAuth.DeviceCode)

	return &v1.DeviceTokenResponse{
		AccessToken:  accessToken,
		TokenType:    "Bearer",
		ExpiresIn:    3600, // 1 hour
		RefreshToken: "",   // Device flow doesn't provide refresh tokens
		Scope:        deviceAuth.Scope,
	}, nil
}

// VerifyUserAuthorization implements DeviceAuthService.
func (d *deviceAuthServiceImpl) VerifyUserAuthorization(ctx context.Context, req *v1.VerifyDeviceRequest) error {
	if !req.Approved {
		return errors.WithCode(code.ErrAuthorizationDenied, "authorization denied by user")
	}

	deviceAuth, err := d.Store.DeviceAuthRepository().GetByUserCode(ctx, req.UserCode)
	if err != nil {
		return errors.WithCode(code.ErrDeviceCodeInvalid, "invalid user code")
	}

	// Update authorization status
	deviceAuth.Status = model.DeviceAuthApproved
	deviceAuth.UserID = req.UserID
	now := time.Now()
	deviceAuth.ApprovedAt = &now

	return d.Store.DeviceAuthRepository().Update(ctx, deviceAuth)
}

// generateRandomString generates a random string of specified length.
func generateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// generateUserCode generates a user-friendly verification code.
func generateUserCode() (string, error) {
	// Generate a 8-character code with numbers and uppercase letters
	const charset = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	const length = 8

	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	for i := 0; i < length; i++ {
		bytes[i] = charset[int(bytes[i])%len(charset)]
	}

	// Format as ABCD-EFGH for better readability
	return fmt.Sprintf("%s-%s", string(bytes[:4]), string(bytes[4:])), nil
}
