// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package webhook

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
)

// Event represents a webhook event.
type Event struct {
	ID        uuid.UUID       `json:"id"`
	Type      string          `json:"type"`
	TargetURL string          `json:"target_url"`
	Payload   json.RawMessage `json:"payload"`
	Retries   int             `json:"retries"`
	CreatedAt time.Time       `json:"created_at"`
}

// Manager defines the interface for webhook operations.
type Manager interface {
	RegisterHook(ctx context.Context, eventType, targetURL string) error
	UnregisterHook(ctx context.Context, eventType, targetURL string) error
	SendEvent(ctx context.Context, eventType string, payload any) error
}

// managerImpl implements webhook.Manager.
type managerImpl struct {
	mu    sync.Mutex
	hooks map[string][]string // eventType -> targetURLs
}

// NewManager creates a new webhook manager.
func NewManager() *managerImpl {
	return &managerImpl{
		hooks: make(map[string][]string),
	}
}

// RegisterHook registers a webhook for an event type.
func (m *managerImpl) RegisterHook(ctx context.Context, eventType, targetURL string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.hooks[eventType] = append(m.hooks[eventType], targetURL)
	return nil
}

// UnregisterHook unregisters a webhook.
func (m *managerImpl) UnregisterHook(ctx context.Context, eventType, targetURL string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	hooks := m.hooks[eventType]
	for i, url := range hooks {
		if url == targetURL {
			m.hooks[eventType] = append(hooks[:i], hooks[i+1:]...)
			break
		}
	}
	return nil
}

// SendEvent sends an event to all registered webhooks.
func (m *managerImpl) SendEvent(ctx context.Context, eventType string, payload any) error {
	m.mu.Lock()
	targetURLs := m.hooks[eventType]
	m.mu.Unlock()

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	for _, targetURL := range targetURLs {
		go func(url string) {
			// Fire and forget with timeout
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			req, _ := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(payloadBytes))
			req.Header.Set("Content-Type", "application/json")

			client := &http.Client{}
			_, _ = client.Do(req)
		}(targetURL)
	}

	return nil
}

// Ensure managerImpl implements Manager.
var _ Manager = (*managerImpl)(nil)
