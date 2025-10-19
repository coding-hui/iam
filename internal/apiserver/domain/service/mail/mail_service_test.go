// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package mail

import (
	"testing"

	metav1 "github.com/coding-hui/common/meta/v1"

	"github.com/coding-hui/iam/internal/apiserver/domain/model"
)

func TestMailService(t *testing.T) {
	// Test with disabled mail service
	service := NewService(nil)

	user := &model.User{
		ObjectMeta: metav1.ObjectMeta{
			Name: "testuser",
		},
		Email: "test@example.com",
	}

	// Should not fail even with disabled service
	err := service.SendWelcomeEmail(user, "initialpassword")
	if err != nil {
		t.Errorf("SendWelcomeEmail should not fail with disabled service: %v", err)
	}

	err = service.SendPasswordResetEmail(user, "reset-token")
	if err != nil {
		t.Errorf("SendPasswordResetEmail should not fail with disabled service: %v", err)
	}

	// Test with user without email
	userWithoutEmail := &model.User{
		ObjectMeta: metav1.ObjectMeta{
			Name: "testuser2",
		},
		Email: "",
	}

	err = service.SendWelcomeEmail(userWithoutEmail, "initialpassword")
	if err != nil {
		t.Errorf("SendWelcomeEmail should not fail with user without email: %v", err)
	}
}
