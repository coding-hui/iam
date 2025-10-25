// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package service

import (
	"context"
	"crypto/rand"
	"fmt"
	"strings"

	"github.com/coding-hui/iam/internal/apiserver/domain/model"
	"github.com/coding-hui/iam/internal/apiserver/domain/repository"
	assembler "github.com/coding-hui/iam/internal/apiserver/interfaces/api/assembler/v1"
	"github.com/coding-hui/iam/internal/pkg/request"
	v1 "github.com/coding-hui/iam/pkg/api/apiserver/v1"
	"github.com/coding-hui/iam/pkg/code"
	"github.com/coding-hui/iam/pkg/log"

	"github.com/coding-hui/common/errors"
	metav1 "github.com/coding-hui/common/meta/v1"
	"github.com/coding-hui/common/util/auth"
)

// ApiKeyService API Key manage api.
type ApiKeyService interface {
	// CreateApiKey creates a new API Key.
	CreateApiKey(ctx context.Context, req v1.CreateApiKeyRequest) (*v1.CreateApiKeyResponse, error)

	// UpdateApiKey updates an existing API Key.
	UpdateApiKey(ctx context.Context, instanceId string, req v1.UpdateApiKeyRequest) (*v1.ApiKeyBase, error)

	// DeleteApiKey deletes an API Key by instance ID.
	DeleteApiKey(ctx context.Context, instanceId string, opts metav1.DeleteOptions) error

	// BatchDeleteApiKeys deletes multiple API Keys by instance IDs.
	BatchDeleteApiKeys(ctx context.Context, instanceIds []string, opts metav1.DeleteOptions) error

	// GetApiKey retrieves an API Key by instance ID.
	GetApiKey(ctx context.Context, instanceId string, opts metav1.GetOptions) (*v1.ApiKeyBase, error)

	// GetApiKeyByKey retrieves an API Key by key value.
	GetApiKeyByKey(ctx context.Context, key string, opts metav1.GetOptions) (*v1.ApiKeyBase, error)

	// ListApiKeys retrieves a list of API Keys with pagination and filtering.
	ListApiKeys(ctx context.Context, opts v1.ListApiKeyOptions) (*v1.ApiKeyList, error)

	// EnableApiKey enables a disabled API Key.
	EnableApiKey(ctx context.Context, instanceId string) error

	// DisableApiKey disables an enabled API Key.
	DisableApiKey(ctx context.Context, instanceId string) error

	// RegenerateSecret regenerates the secret for an API Key.
	RegenerateSecret(ctx context.Context, instanceId string) (*v1.CreateApiKeyResponse, error)

	// ValidateApiKey validates an API Key and returns the associated user and API Key object.
	ValidateApiKey(ctx context.Context, key string) (*model.User, *model.ApiKey, error)

	// CleanupExpiredApiKeys deletes expired API Keys.
	CleanupExpiredApiKeys(ctx context.Context) error
}

type apiKeyServiceImpl struct {
	Store       repository.Factory `inject:"repository"`
	UserService UserService        `inject:""`
}

// NewApiKeyService new API Key service.
func NewApiKeyService() ApiKeyService {
	return &apiKeyServiceImpl{}
}

// CreateApiKey creates a new API Key.
func (s *apiKeyServiceImpl) CreateApiKey(ctx context.Context, req v1.CreateApiKeyRequest) (*v1.CreateApiKeyResponse, error) {
	// Validate request
	if req.Name == "" {
		return nil, errors.WithCode(code.ErrValidation, "API Key name is required")
	}

	// Generate API Key and Secret
	key, secret, err := s.generateApiKeyAndSecret()
	if err != nil {
		return nil, errors.WithMessage(err, "Failed to generate API Key and Secret")
	}

	// Validate generated key format
	if !strings.HasPrefix(key, "sk-") || len(key) != 35 { // sk- + 32 hex chars = 35 chars
		return nil, errors.WithCode(code.ErrUnknown, "Generated API Key format is invalid")
	}

	if len(secret) != 64 { // 64 hex characters for 256-bit secret
		return nil, errors.WithCode(code.ErrUnknown, "Generated API Secret format is invalid")
	}

	// Encrypt the secret
	encryptedSecret, err := auth.Encrypt(secret)
	if err != nil {
		return nil, errors.WithMessage(err, "Failed to encrypt API Secret")
	}

	// Get current user from context
	user, err := s.getCurrentUser(ctx)
	if err != nil {
		return nil, err
	}

	// Create API Key model
	apiKey := &model.ApiKey{
		ObjectMeta: metav1.ObjectMeta{
			Name: req.Name,
		},
		Key:    key,
		Secret: encryptedSecret,
		UserID: user.GetInstanceID(),
		Status: model.ApiKeyStatusActive,
	}

	// Set expiration time if provided (non-zero time)
	if !req.ExpiresAt.IsZero() {
		apiKey.ExpiresAt = &req.ExpiresAt
	}

	// Save to database
	createdApiKey, err := s.Store.ApiKeyRepository().Create(ctx, apiKey, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}

	// Convert to response
	base := assembler.ConvertApiKeyModelToBase(createdApiKey)

	log.Infof("Created API Key '%s' for user '%s' with key: %s", req.Name, user.Name, key)

	return &v1.CreateApiKeyResponse{
		ApiKeyBase: *base,
	}, nil
}

// UpdateApiKey updates an existing API Key.
func (s *apiKeyServiceImpl) UpdateApiKey(ctx context.Context, instanceId string, req v1.UpdateApiKeyRequest) (*v1.ApiKeyBase, error) {
	apiKey, err := s.Store.ApiKeyRepository().GetByInstanceId(ctx, instanceId, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	// Verify ownership
	user, err := s.getCurrentUser(ctx)
	if err != nil {
		return nil, err
	}
	if apiKey.UserID != user.GetInstanceID() {
		return nil, errors.WithCode(code.ErrPermissionDenied, "Permission denied")
	}

	// Update fields
	apiKey.Name = req.Name

	// Update expiration time if provided (non-zero time)
	if !req.ExpiresAt.IsZero() {
		apiKey.ExpiresAt = &req.ExpiresAt
	} else {
		// If ExpiresAt is zero, set it to nil to clear expiration
		apiKey.ExpiresAt = nil
	}

	if req.Status != 0 {
		apiKey.Status = model.ApiKeyStatus(req.Status)
	}

	if err := s.Store.ApiKeyRepository().Update(ctx, apiKey, metav1.UpdateOptions{}); err != nil {
		return nil, err
	}

	return assembler.ConvertApiKeyModelToBase(apiKey), nil
}

// DeleteApiKey deletes an API Key by instance ID.
func (s *apiKeyServiceImpl) DeleteApiKey(ctx context.Context, instanceId string, opts metav1.DeleteOptions) error {
	apiKey, err := s.Store.ApiKeyRepository().GetByInstanceId(ctx, instanceId, metav1.GetOptions{})
	if err != nil {
		return err
	}

	// Verify ownership
	user, err := s.getCurrentUser(ctx)
	if err != nil {
		return err
	}
	if apiKey.UserID != user.GetInstanceID() {
		return errors.WithCode(code.ErrPermissionDenied, "Permission denied")
	}

	return s.Store.ApiKeyRepository().DeleteByInstanceId(ctx, instanceId, opts)
}

// BatchDeleteApiKeys deletes multiple API Keys by instance IDs.
func (s *apiKeyServiceImpl) BatchDeleteApiKeys(ctx context.Context, instanceIds []string, opts metav1.DeleteOptions) error {
	return s.Store.ApiKeyRepository().BatchDelete(ctx, instanceIds, opts)
}

// GetApiKey retrieves an API Key by instance ID.
func (s *apiKeyServiceImpl) GetApiKey(ctx context.Context, instanceId string, opts metav1.GetOptions) (*v1.ApiKeyBase, error) {
	apiKey, err := s.Store.ApiKeyRepository().GetByInstanceId(ctx, instanceId, opts)
	if err != nil {
		return nil, err
	}

	// Verify ownership
	user, err := s.getCurrentUser(ctx)
	if err != nil {
		return nil, err
	}
	if apiKey.UserID != user.GetInstanceID() {
		return nil, errors.WithCode(code.ErrPermissionDenied, "Permission denied")
	}

	// Check if API Key has expired and update status if needed
	if !apiKey.IsActive() && apiKey.Status == model.ApiKeyStatusActive {
		// API Key has expired but status is still active, update to expired status
		apiKey.Status = model.ApiKeyStatusExpired
		if err := s.Store.ApiKeyRepository().Update(ctx, apiKey, metav1.UpdateOptions{}); err != nil {
			log.Warnf("Failed to update expired API Key status: %v", err)
			// Continue with the original status if update fails
		}
	}

	return assembler.ConvertApiKeyModelToBase(apiKey), nil
}

// GetApiKeyByKey retrieves an API Key by key value.
func (s *apiKeyServiceImpl) GetApiKeyByKey(ctx context.Context, key string, opts metav1.GetOptions) (*v1.ApiKeyBase, error) {
	apiKey, err := s.Store.ApiKeyRepository().GetByKey(ctx, key, opts)
	if err != nil {
		return nil, err
	}

	// Check if API Key has expired and update status if needed
	if !apiKey.IsActive() && apiKey.Status == model.ApiKeyStatusActive {
		// API Key has expired but status is still active, update to expired status
		apiKey.Status = model.ApiKeyStatusExpired
		if err := s.Store.ApiKeyRepository().Update(ctx, apiKey, metav1.UpdateOptions{}); err != nil {
			log.Warnf("Failed to update expired API Key status: %v", err)
			// Continue with the original status if update fails
		}
	}

	return assembler.ConvertApiKeyModelToBase(apiKey), nil
}

// ListApiKeys retrieves a list of API Keys with pagination and filtering.
func (s *apiKeyServiceImpl) ListApiKeys(ctx context.Context, opts v1.ListApiKeyOptions) (*v1.ApiKeyList, error) {
	// Set user ID filter to current user's keys
	user, err := s.getCurrentUser(ctx)
	if err != nil {
		return nil, err
	}
	opts.UserID = user.GetInstanceID()

	apiKeys, err := s.Store.ApiKeyRepository().List(ctx, opts)
	if err != nil {
		return nil, err
	}

	// Check and update expired API Keys status
	for i := range apiKeys {
		apiKey := &apiKeys[i]
		if !apiKey.IsActive() && apiKey.Status == model.ApiKeyStatusActive {
			// API Key has expired but status is still active, update to expired status
			apiKey.Status = model.ApiKeyStatusExpired
			if err := s.Store.ApiKeyRepository().Update(ctx, apiKey, metav1.UpdateOptions{}); err != nil {
				log.Warnf("Failed to update expired API Key status: %v", err)
				// Continue with the original status if update fails
			}
		}
	}

	count, err := s.Store.ApiKeyRepository().Count(ctx, opts)
	if err != nil {
		return nil, err
	}

	items := make([]*v1.ApiKeyBase, len(apiKeys))
	for i, apiKey := range apiKeys {
		items[i] = assembler.ConvertApiKeyModelToBase(&apiKey)
	}

	return &v1.ApiKeyList{
		Items: items,
		ListMeta: metav1.ListMeta{
			TotalCount: count,
		},
	}, nil
}

// EnableApiKey enables a disabled API Key.
func (s *apiKeyServiceImpl) EnableApiKey(ctx context.Context, instanceId string) error {
	apiKey, err := s.Store.ApiKeyRepository().GetByInstanceId(ctx, instanceId, metav1.GetOptions{})
	if err != nil {
		return err
	}

	// Verify ownership
	user, err := s.getCurrentUser(ctx)
	if err != nil {
		return err
	}
	if apiKey.UserID != user.GetInstanceID() {
		return errors.WithCode(code.ErrPermissionDenied, "Permission denied")
	}

	if apiKey.Status == model.ApiKeyStatusActive {
		return errors.WithCode(code.ErrApiKeyAlreadyEnabled, "API Key is already enabled")
	}

	apiKey.Status = model.ApiKeyStatusActive
	return s.Store.ApiKeyRepository().Update(ctx, apiKey, metav1.UpdateOptions{})
}

// DisableApiKey disables an enabled API Key.
func (s *apiKeyServiceImpl) DisableApiKey(ctx context.Context, instanceId string) error {
	apiKey, err := s.Store.ApiKeyRepository().GetByInstanceId(ctx, instanceId, metav1.GetOptions{})
	if err != nil {
		return err
	}

	// Verify ownership
	user, err := s.getCurrentUser(ctx)
	if err != nil {
		return err
	}
	if apiKey.UserID != user.GetInstanceID() {
		return errors.WithCode(code.ErrPermissionDenied, "Permission denied")
	}

	if apiKey.Status == model.ApiKeyStatusInactive {
		return errors.WithCode(code.ErrApiKeyAlreadyDisabled, "API Key is already disabled")
	}

	apiKey.Status = model.ApiKeyStatusInactive
	return s.Store.ApiKeyRepository().Update(ctx, apiKey, metav1.UpdateOptions{})
}

// RegenerateSecret regenerates the secret for an API Key.
func (s *apiKeyServiceImpl) RegenerateSecret(ctx context.Context, instanceId string) (*v1.CreateApiKeyResponse, error) {
	apiKey, err := s.Store.ApiKeyRepository().GetByInstanceId(ctx, instanceId, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	// Verify ownership
	user, err := s.getCurrentUser(ctx)
	if err != nil {
		return nil, err
	}
	if apiKey.UserID != user.GetInstanceID() {
		return nil, errors.WithCode(code.ErrPermissionDenied, "Permission denied")
	}

	// Generate new secret
	secret, err := s.generateSecret()
	if err != nil {
		return nil, errors.WithMessage(err, "Failed to generate new secret")
	}

	// Encrypt the new secret
	encryptedSecret, err := auth.Encrypt(secret)
	if err != nil {
		return nil, errors.WithMessage(err, "Failed to encrypt new secret")
	}

	// Update API Key
	apiKey.Secret = encryptedSecret
	if err := s.Store.ApiKeyRepository().Update(ctx, apiKey, metav1.UpdateOptions{}); err != nil {
		return nil, err
	}

	base := assembler.ConvertApiKeyModelToBase(apiKey)

	return &v1.CreateApiKeyResponse{
		ApiKeyBase: *base,
	}, nil
}

// ValidateApiKey validates an API Key and returns the associated user and API Key object.
func (s *apiKeyServiceImpl) ValidateApiKey(ctx context.Context, key string) (*model.User, *model.ApiKey, error) {
	apiKey, err := s.Store.ApiKeyRepository().GetByKey(ctx, key, metav1.GetOptions{})
	if err != nil {
		return nil, nil, errors.WithCode(code.ErrApiKeyInvalid, "Invalid API Key")
	}

	if apiKey.IsExpired() {
		return nil, nil, errors.WithCode(code.ErrApiKeyExpired, "API Key is expired")
	}

	if !apiKey.IsActive() {
		return nil, nil, errors.WithCode(code.ErrApiKeyInactive, "API Key is not active")
	}

	// Update usage statistics
	if err := s.Store.ApiKeyRepository().UpdateUsage(ctx, apiKey.GetInstanceID()); err != nil {
		log.Warnf("Failed to update API Key usage: %v", err)
	}

	// Get user
	user, err := s.UserService.GetUserByInstanceId(ctx, apiKey.UserID, metav1.GetOptions{})
	if err != nil {
		return nil, nil, errors.WithMessage(err, "Failed to get API Key owner")
	}

	return user, apiKey, nil
}

// CleanupExpiredApiKeys deletes expired API Keys.
func (s *apiKeyServiceImpl) CleanupExpiredApiKeys(ctx context.Context) error {
	return s.Store.ApiKeyRepository().CleanupExpired(ctx)
}

// Helper functions

func (s *apiKeyServiceImpl) generateApiKeyAndSecret() (string, string, error) {
	// Generate key with retry logic for collision detection
	var key string
	for i := 0; i < 3; i++ { // Retry up to 3 times in case of collision
		generatedKey, err := s.generateKey()
		if err != nil {
			return "", "", err
		}

		// Check if key already exists
		ctx := context.Background()
		_, err = s.Store.ApiKeyRepository().GetByKey(ctx, generatedKey, metav1.GetOptions{})
		if err != nil {
			// Key doesn't exist, we can use it
			key = generatedKey
			break
		}
		// Key exists, try again
		if i == 2 {
			// Last attempt failed
			return "", "", errors.WithCode(code.ErrApiKeyGenerationFailed, "Failed to generate unique API Key after multiple attempts")
		}
	}

	secret, err := s.generateSecret()
	if err != nil {
		return "", "", err
	}

	return key, secret, nil
}

func (s *apiKeyServiceImpl) generateSecret() (string, error) {
	// Generate 32 bytes (256 bits) for high security
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	// Format: 64 hex characters (lowercase)
	return fmt.Sprintf("%064x", bytes), nil
}

func (s *apiKeyServiceImpl) generateKey() (string, error) {
	// Generate 16 bytes (128 bits) for hex encoding
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	// Format: sk- + 32 hex characters (e.g., sk-716273ec41584ab7ab09ca112367c1e5)
	return "sk-" + fmt.Sprintf("%032x", bytes), nil
}

func (s *apiKeyServiceImpl) getCurrentUser(ctx context.Context) (*model.User, error) {
	currentUser, ok := request.UserFrom(ctx)
	if !ok {
		return nil, errors.WithCode(code.ErrPermissionDenied, "Failed to obtain the current user")
	}

	user, err := s.UserService.GetUserByInstanceId(ctx, currentUser.InstanceID, metav1.GetOptions{})
	if err != nil {
		return nil, errors.WithMessage(err, "Failed to get current user details")
	}

	return user, nil
}
