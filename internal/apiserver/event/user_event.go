// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package event

import "github.com/coding-hui/iam/internal/apiserver/domain/model"

// UserCreatedEvent user created event
type UserCreatedEvent struct {
	BasicEvent
}

// Name get event name
func (e *UserCreatedEvent) Name() string {
	return UserCreatedEventType
}

// NewUserCreatedEvent create a new user created event
func NewUserCreatedEvent(user *model.User) *UserCreatedEvent {
	return &UserCreatedEvent{
		BasicEvent: *NewBasic(UserCreatedEventType, map[string]any{
			"user": user,
		}),
	}
}
