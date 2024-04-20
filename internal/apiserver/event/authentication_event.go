// Copyright (c) 2024 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package event

import (
	"context"

	"github.com/coding-hui/iam/internal/apiserver/domain/repository"
	"github.com/coding-hui/iam/pkg/log"
)

type AuthenticationEvent struct {
	BasicEvent
	Success        bool
	FailMessage    string
	Username       string
	UserInstanceID string
}

func (e *AuthenticationEvent) Name() string {
	return AuthenticationEventType
}

type authenticationSuccessListener struct {
	Store repository.Factory `inject:"repository"`
}

func NewAuthenticationEventListener() Listener {
	return &authenticationSuccessListener{}
}

func (l *authenticationSuccessListener) Handle(raw Event) error {
	e := raw.(*AuthenticationEvent)
	err := l.Store.UserRepository().FlushLastLoginTime(context.Background(), e.Username)
	if err != nil {
		log.Errorf("Failed to flush user [%s] last login time: %v", e.Username, err)
		return err
	}
	return nil
}
