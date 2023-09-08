// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package code

//go:generate codegen -type=int
//go:generate codegen -type=int -doc -output ../../docs/guide/zh-CN/api/error_code_generated.md

// Common: basic errors.
// Code must start with 1xxxxx.
const (
	// ErrSuccess - 200: OK.
	ErrSuccess int = iota + 100001

	// ErrUnknown - 500: Internal server error.
	ErrUnknown

	// ErrBind - 400: Error occurred while binding the request body to the struct.
	ErrBind

	// ErrValidation - 400: Validation failed.
	ErrValidation

	// ErrTokenInvalid - 401: Token invalid.
	ErrTokenInvalid

	// ErrPageNotFound - 404: Page not found.
	ErrPageNotFound

	// ErrInvalidRequest - 404: Invalid request.
	ErrInvalidRequest
)

// common: database errors.
const (
	// ErrDatabase - 500: Database error.
	ErrDatabase int = iota + 100101

	// ErrPrimaryEmpty - 500: Primary key is empty.
	ErrPrimaryEmpty

	// ErrNilEntity - 500: Entity is nil.
	ErrNilEntity

	// ErrRecordExist - 500: Entity primary key is exist.
	ErrRecordExist

	// ErrRecordNotExist - 500: Entity primary key is not exist.
	ErrRecordNotExist

	// ErrIndexInvalid - 500: Entity index is invalid.
	ErrIndexInvalid

	// ErrEntityInvalid - 500: Entity is invalid.
	ErrEntityInvalid
)

// common: authorization and authentication errors.
const (
	// ErrEncrypt - 401: Error occurred while encrypting the user password.
	ErrEncrypt int = iota + 100201

	// ErrSignatureInvalid - 401: Signature is invalid.
	ErrSignatureInvalid

	// ErrTokenMalformed - 401: Token is malformed.
	ErrTokenMalformed

	// ErrTokenNotValidYet - 401: Token is not valid yet.
	ErrTokenNotValidYet

	// ErrExpired - 401: Token expired.
	ErrExpired

	// ErrMissingLoginValues - 401: Missing Username or Password.
	ErrMissingLoginValues

	// ErrInvalidAuthHeader - 401: Invalid authorization header.
	ErrInvalidAuthHeader

	// ErrMissingHeader - 401: The `Authorization` header was empty.
	ErrMissingHeader

	// ErrPasswordIncorrect - 401: Invalid Username or Password.
	ErrPasswordIncorrect

	// ErrInvalidRefreshToken - 401: Refresh token format is incorrect, please check.
	ErrInvalidRefreshToken

	// ErrUnauthorized - 403: Unauthorized.
	ErrUnauthorized

	// ErrPermissionDenied - 403: Permission denied.
	ErrPermissionDenied

	// ErrIdentityProviderNotFound - 401: Identity provider not found.
	ErrIdentityProviderNotFound
)

// common: encode/decode errors.
const (
	// ErrEncodingFailed - 500: Encoding failed due to an error with the data.
	ErrEncodingFailed int = iota + 100301

	// ErrDecodingFailed - 500: Decoding failed due to an error with the data.
	ErrDecodingFailed

	// ErrInvalidJSON - 500: Data is not valid JSON.
	ErrInvalidJSON

	// ErrEncodingJSON - 500: JSON data could not be encoded.
	ErrEncodingJSON

	// ErrDecodingJSON - 500: JSON data could not be decoded.
	ErrDecodingJSON

	// ErrInvalidYaml - 500: Data is not valid Yaml.
	ErrInvalidYaml

	// ErrEncodingYaml - 500: Yaml data could not be encoded.
	ErrEncodingYaml

	// ErrDecodingYaml - 500: Yaml data could not be decoded.
	ErrDecodingYaml
)
