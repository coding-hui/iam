// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package persistence

import (
	"context"
	"time"
)

// AuditEvent represents an audit event in the system.
// Domain model with no persistence-specific tags (Ory style).
type AuditEvent struct {
	ID         string
	NetworkID  string
	Type       string
	ActorID    string
	ActorType  string
	TargetID   string
	TargetType string
	Outcome    string
	ClientIP   string
	UserAgent  string
	RequestID  string
	Metadata   []byte
	Timestamp  time.Time
}

// AuditFilter holds filter criteria for audit event queries.
type AuditFilter struct {
	Type      string
	ActorID   string
	TargetID  string
	Outcome   string
	StartTime *time.Time
	EndTime   *time.Time
}

// AuditEventPersister defines the interface for audit event persistence operations.
type AuditEventPersister interface {
	CreateAuditEvent(ctx context.Context, event *AuditEvent) error
	ListAuditEvents(ctx context.Context, networkID string, limit, offset int, filter *AuditFilter) ([]*AuditEvent, int, error)
}
