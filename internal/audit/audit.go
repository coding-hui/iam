// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package audit

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// AuditEvent represents an audit event in the system.
type AuditEvent struct {
	ID         uuid.UUID       `json:"id"`
	NetworkID  uuid.UUID       `json:"network_id"`
	Type       string          `json:"type"`
	ActorID    uuid.UUID       `json:"actor_id"`
	ActorType  string          `json:"actor_type"`
	TargetID   uuid.UUID       `json:"target_id"`
	TargetType string          `json:"target_type"`
	Outcome    string          `json:"outcome"`
	ClientIP   string          `json:"client_ip"`
	UserAgent  string          `json:"user_agent"`
	RequestID  string          `json:"request_id"`
	Metadata   json.RawMessage `json:"metadata,omitempty"`
	Timestamp  time.Time       `json:"timestamp"`
}

// Pool defines the interface for reading audit data.
type Pool interface {
	ListEvents(ctx context.Context, networkID uuid.UUID, limit, offset int) ([]*AuditEvent, int, error)
}

// Recorder defines the interface for writing audit events.
type Recorder interface {
	Record(ctx context.Context, event *AuditEvent) error
}

// Manager defines the interface for audit business logic.
type Manager interface {
	RecordEvent(ctx context.Context, req *RecordEventRequest) error
	ListEvents(ctx context.Context, networkID uuid.UUID, limit, offset int) ([]*AuditEvent, int, error)
}

// RecordEventRequest holds data for recording an audit event.
type RecordEventRequest struct {
	NetworkID  uuid.UUID       `json:"network_id"`
	Type       string          `json:"type"`
	ActorID    uuid.UUID       `json:"actor_id"`
	ActorType  string          `json:"actor_type"`
	TargetID   uuid.UUID       `json:"target_id"`
	TargetType string          `json:"target_type"`
	Outcome    string          `json:"outcome"`
	ClientIP   string          `json:"client_ip"`
	UserAgent  string          `json:"user_agent"`
	RequestID  string          `json:"request_id"`
	Metadata   json.RawMessage `json:"metadata,omitempty"`
}
