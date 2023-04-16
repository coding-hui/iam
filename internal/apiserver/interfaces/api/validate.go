// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package api

import (
	"regexp"
	"unicode"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

var (
	nameRegexp  = regexp.MustCompile(`^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$`)
	emailRegexp = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
)

const (
	minPageSize = 5
	maxPageSize = 100
)

func init() {
	if err := validate.RegisterValidation("name", ValidateName); err != nil {
		panic(err)
	}
	if err := validate.RegisterValidation("alias", ValidateAlias); err != nil {
		panic(err)
	}
	if err := validate.RegisterValidation("email", ValidateEmail); err != nil {
		panic(err)
	}
	if err := validate.RegisterValidation("password", ValidatePassword); err != nil {
		panic(err)
	}
}

// ValidateName custom check name field
func ValidateName(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	if len(value) > 31 || len(value) < 2 {
		return false
	}
	return nameRegexp.MatchString(value)
}

// ValidateAlias custom check alias field
func ValidateAlias(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	if value != "" && (len(value) > 64 || len(value) < 2) {
		return false
	}
	return true
}

// ValidateEmail custom check email field
func ValidateEmail(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	if value == "" {
		return true
	}
	return emailRegexp.MatchString(value)
}

// ValidatePassword custom check password field
func ValidatePassword(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	if value == "" {
		return true
	}
	if len(value) < 8 || len(value) > 16 {
		return false
	}
	// go's regex doesn't support backtracking so check the password with a loop
	letter := false
	num := false
	for _, c := range value {
		switch {
		case unicode.IsNumber(c):
			num = true
		case unicode.IsLetter(c):
			letter = true
		}
	}
	return letter && num
}
