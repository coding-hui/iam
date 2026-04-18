// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package sql

import (
	"context"
	"time"

	"github.com/coding-hui/iam/internal/persistence"
)

// AuditEventModel represents an audit event in the database.
type AuditEventModel struct {
	ID         string    `gorm:"primaryKey;column:id"    json:"id"`
	NetworkID  string    `gorm:"column:nid;index"        json:"network_id"`
	Type       string    `gorm:"column:type;index"       json:"type"`
	ActorID    string    `gorm:"column:actor_id;index"   json:"actor_id"`
	ActorType  string    `gorm:"column:actor_type"       json:"actor_type"`
	TargetID   string    `gorm:"column:target_id"        json:"target_id"`
	TargetType string    `gorm:"column:target_type"      json:"target_type"`
	Outcome    string    `gorm:"column:outcome;index"    json:"outcome"`
	ClientIP   string    `gorm:"column:client_ip"        json:"client_ip"`
	UserAgent  string    `gorm:"column:user_agent"       json:"user_agent"`
	RequestID  string    `gorm:"column:request_id;index" json:"request_id"`
	Metadata   []byte    `gorm:"column:metadata"         json:"metadata"`
	Timestamp  time.Time `gorm:"column:timestamp;index"  json:"timestamp"`
}

// TableName returns the table name for AuditEventModel.
func (AuditEventModel) TableName() string {
	return "iam_audit_events"
}

// AuditPool implements persistence.AuditEventPersister using GORM.
type AuditPool struct {
	db *Persister
}

// NewAuditPool creates a new audit pool.
func NewAuditPool(db *Persister) *AuditPool {
	return &AuditPool{db: db}
}

// CreateAuditEvent creates a new audit event.
func (p *AuditPool) CreateAuditEvent(ctx context.Context, event *persistence.AuditEvent) error {
	m := p.domainToModel(event)
	return p.db.Connection(ctx).Create(m).Error
}

// ListAuditEvents lists audit events with pagination.
func (p *AuditPool) ListAuditEvents(ctx context.Context, networkID string, limit, offset int, filter *persistence.AuditFilter) ([]*persistence.AuditEvent, int, error) {
	var ms []AuditEventModel
	var total int64

	query := p.db.Connection(ctx)
	if networkID != "" {
		query = query.Where("nid = ?", networkID)
	}
	if filter != nil {
		if filter.Type != "" {
			query = query.Where("type = ?", filter.Type)
		}
		if filter.ActorID != "" {
			query = query.Where("actor_id = ?", filter.ActorID)
		}
		if filter.Outcome != "" {
			query = query.Where("outcome = ?", filter.Outcome)
		}
	}

	if err := query.Model(&AuditEventModel{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Order("timestamp DESC").Limit(limit).Offset(offset).Find(&ms).Error; err != nil {
		return nil, 0, err
	}

	events := make([]*persistence.AuditEvent, len(ms))
	for i := range ms {
		events[i] = p.modelToDomain(&ms[i])
	}
	return events, int(total), nil
}

func (p *AuditPool) modelToDomain(m *AuditEventModel) *persistence.AuditEvent {
	return &persistence.AuditEvent{
		ID:         m.ID,
		NetworkID:  m.NetworkID,
		Type:       m.Type,
		ActorID:    m.ActorID,
		ActorType:  m.ActorType,
		TargetID:   m.TargetID,
		TargetType: m.TargetType,
		Outcome:    m.Outcome,
		ClientIP:   m.ClientIP,
		UserAgent:  m.UserAgent,
		RequestID:  m.RequestID,
		Metadata:   m.Metadata,
		Timestamp:  m.Timestamp,
	}
}

func (p *AuditPool) domainToModel(e *persistence.AuditEvent) *AuditEventModel {
	return &AuditEventModel{
		ID:         e.ID,
		NetworkID:  e.NetworkID,
		Type:       e.Type,
		ActorID:    e.ActorID,
		ActorType:  e.ActorType,
		TargetID:   e.TargetID,
		TargetType: e.TargetType,
		Outcome:    e.Outcome,
		ClientIP:   e.ClientIP,
		UserAgent:  e.UserAgent,
		RequestID:  e.RequestID,
		Metadata:   e.Metadata,
		Timestamp:  e.Timestamp,
	}
}

// Ensure AuditPool implements persistence.AuditEventPersister.
var _ persistence.AuditEventPersister = (*AuditPool)(nil)
