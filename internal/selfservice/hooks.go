// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package selfservice

import (
	"context"

	"github.com/coding-hui/iam/internal/selfservice/flows"
)

// FlowHook is a hook executed during a selfservice flow.
type FlowHook interface {
	// Execute executes the hook.
	Execute(ctx context.Context, flow flows.FlowState) error
}

// FlowHooks is a list of FlowHooks.
type FlowHooks []FlowHook

// Execute executes all hooks in order.
func (h FlowHooks) Execute(ctx context.Context, flow flows.FlowState) error {
	for _, hook := range h {
		if err := hook.Execute(ctx, flow); err != nil {
			return err
		}
	}
	return nil
}

// RegistrationHooks are hooks executed during registration.
type RegistrationHooks []FlowHook

// LoginHooks are hooks executed during login.
type LoginHooks []FlowHook

// RecoveryHooks are hooks executed during password recovery.
type RecoveryHooks []FlowHook

// SettingsHooks are hooks executed during settings changes.
type SettingsHooks []FlowHook
