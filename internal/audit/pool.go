// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package audit

import (
	"context"

	"github.com/google/uuid"

	"github.com/coding-hui/iam/internal/persistence"
)

// auditPool implements Pool using persistence.AuditEventPersister.
type auditPool struct {
	persister auditPersister
}

// auditPersister is the persistence interface for audit operations.
type auditPersister interface {
	CreateAuditEvent(ctx context.Context, event *persistence.AuditEvent) error
	ListAuditEvents(ctx context.Context, networkID string, limit, offset int, filter *persistence.AuditFilter) ([]*persistence.AuditEvent, int, error)
}

// NewPool creates a new audit pool.
func NewPool(p auditPersister) Pool {
	return &auditPool{persister: p}
}

// ListEvents lists audit events with pagination.
func (p *auditPool) ListEvents(ctx context.Context, networkID uuid.UUID, limit, offset int) ([]*AuditEvent, int, error) {
	ms, total, err := p.persister.ListAuditEvents(ctx, networkID.String(), limit, offset, nil)
	if err != nil {
		return nil, 0, err
	}
	events := make([]*AuditEvent, len(ms))
	for i := range ms {
		events[i] = p.modelToDomain(ms[i])
	}
	return events, total, nil
}

func (p *auditPool) modelToDomain(m *persistence.AuditEvent) *AuditEvent {
	if m == nil {
		return nil
	}
	return &AuditEvent{
		ID:         parseUUID(m.ID),
		NetworkID:  parseUUID(m.NetworkID),
		Type:       m.Type,
		ActorID:    parseUUID(m.ActorID),
		ActorType:  m.ActorType,
		TargetID:   parseUUID(m.TargetID),
		TargetType: m.TargetType,
		Outcome:    m.Outcome,
		ClientIP:   m.ClientIP,
		UserAgent:  m.UserAgent,
		RequestID:  m.RequestID,
		Metadata:   m.Metadata,
		Timestamp:  m.Timestamp,
	}
}

func parseUUID(s string) uuid.UUID {
	if s == "" {
		return uuid.Nil
	}
	id, err := uuid.Parse(s)
	if err != nil {
		return uuid.Nil
	}
	return id
}

// Ensure auditPool implements Pool.
var _ Pool = (*auditPool)(nil)

// recorder implements Recorder.
type recorder struct {
	persister auditPersister
}

// NewRecorder creates a new audit recorder.
func NewRecorder(p auditPersister) Recorder {
	return &recorder{persister: p}
}

// Record records an audit event.
func (r *recorder) Record(ctx context.Context, event *AuditEvent) error {
	m := r.domainToModel(event)
	return r.persister.CreateAuditEvent(ctx, m)
}

func (r *recorder) domainToModel(e *AuditEvent) *persistence.AuditEvent {
	return &persistence.AuditEvent{
		ID:         e.ID.String(),
		NetworkID:  e.NetworkID.String(),
		Type:       e.Type,
		ActorID:    e.ActorID.String(),
		ActorType:  e.ActorType,
		TargetID:   e.TargetID.String(),
		TargetType: e.TargetType,
		Outcome:    e.Outcome,
		ClientIP:   e.ClientIP,
		UserAgent:  e.UserAgent,
		RequestID:  e.RequestID,
		Metadata:   e.Metadata,
		Timestamp:  e.Timestamp,
	}
}

// Ensure recorder implements Recorder.
var _ Recorder = (*recorder)(nil)
