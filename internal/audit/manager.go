// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package audit

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// ManagerImpl implements audit.Manager.
type ManagerImpl struct {
	pool     Pool
	recorder Recorder
}

// NewManagerImpl creates a new audit manager.
func NewManagerImpl(pool Pool, recorder Recorder) *ManagerImpl {
	return &ManagerImpl{
		pool:     pool,
		recorder: recorder,
	}
}

// RecordEvent records an audit event.
func (m *ManagerImpl) RecordEvent(ctx context.Context, req *RecordEventRequest) error {
	event := &AuditEvent{
		ID:         uuid.New(),
		NetworkID:  req.NetworkID,
		Type:       req.Type,
		ActorID:    req.ActorID,
		ActorType:  req.ActorType,
		TargetID:   req.TargetID,
		TargetType: req.TargetType,
		Outcome:    req.Outcome,
		ClientIP:   req.ClientIP,
		UserAgent:  req.UserAgent,
		RequestID:  req.RequestID,
		Metadata:   req.Metadata,
		Timestamp:  time.Now(),
	}
	return m.recorder.Record(ctx, event)
}

// ListEvents lists audit events.
func (m *ManagerImpl) ListEvents(ctx context.Context, networkID uuid.UUID, limit, offset int) ([]*AuditEvent, int, error) {
	return m.pool.ListEvents(ctx, networkID, limit, offset)
}

// Ensure ManagerImpl implements Manager.
var _ Manager = (*ManagerImpl)(nil)
