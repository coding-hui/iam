// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package lockout

import (
	"context"
	"sync"
	"time"
)

// Manager defines the interface for lockout operations.
type Manager interface {
	RecordFailure(ctx context.Context, identifier string) error
	RecordSuccess(ctx context.Context, identifier string) error
	IsLocked(ctx context.Context, identifier string) (bool, time.Duration, error)
	Unlock(ctx context.Context, identifier string) error
}

// managerImpl implements account lockout functionality.
type managerImpl struct {
	mu           sync.Mutex
	failures     map[string]*failureInfo
	maxAttempts  int
	lockDuration time.Duration
}

// failureInfo tracks failed login attempts.
type failureInfo struct {
	attempts    int
	lastFail    time.Time
	lockedUntil time.Time
}

// NewManager creates a new lockout manager.
func NewManager(maxAttempts int, lockDuration time.Duration) Manager {
	return &managerImpl{
		failures:     make(map[string]*failureInfo),
		maxAttempts:  maxAttempts,
		lockDuration: lockDuration,
	}
}

// RecordFailure records a failed login attempt.
func (m *managerImpl) RecordFailure(ctx context.Context, identifier string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	info, exists := m.failures[identifier]
	if !exists {
		info = &failureInfo{}
		m.failures[identifier] = info
	}

	info.attempts++
	info.lastFail = time.Now()

	if info.attempts >= m.maxAttempts {
		info.lockedUntil = time.Now().Add(m.lockDuration)
	}

	return nil
}

// RecordSuccess records a successful login.
func (m *managerImpl) RecordSuccess(ctx context.Context, identifier string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.failures, identifier)
	return nil
}

// IsLocked checks if an identifier is locked.
func (m *managerImpl) IsLocked(ctx context.Context, identifier string) (bool, time.Duration, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	info, exists := m.failures[identifier]
	if !exists {
		return false, 0, nil
	}

	if info.lockedUntil.IsZero() {
		return false, 0, nil
	}

	if time.Now().Before(info.lockedUntil) {
		remaining := time.Until(info.lockedUntil)
		return true, remaining, nil
	}

	// Lock expired, reset
	delete(m.failures, identifier)
	return false, 0, nil
}

// Unlock unlocks an identifier.
func (m *managerImpl) Unlock(ctx context.Context, identifier string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.failures, identifier)
	return nil
}

// Ensure managerImpl implements Manager.
var _ Manager = (*managerImpl)(nil)
