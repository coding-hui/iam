// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package identity

import (
	"crypto/rand"
	"errors"

	"golang.org/x/crypto/argon2"
)

// ErrInvalidHash is returned when the hash is invalid.
var ErrInvalidHash = errors.New("invalid hash")

// Hasher defines the interface for password hashing.
type Hasher interface {
	Hash(password string) ([]byte, error)
	Verify(password string, hash []byte) error
}

// Argon2idHasher implements password hashing using Argon2id.
type Argon2idHasher struct {
	memory      int
	iterations  int
	parallelism int
	keyLen      int
	saltLen     int
}

// NewArgon2idHasher creates a new Argon2id hasher with recommended parameters.
func NewArgon2idHasher() *Argon2idHasher {
	return &Argon2idHasher{
		memory:      64 * 1024, // 64MB
		iterations:  3,
		parallelism: 4,
		keyLen:      32,
		saltLen:     16,
	}
}

// Hash generates a hash from the password using Argon2id.
func (h *Argon2idHasher) Hash(password string) ([]byte, error) {
	salt := make([]byte, h.saltLen)
	if _, err := rand.Read(salt); err != nil {
		return nil, err
	}

	hash := argon2.IDKey(
		[]byte(password),
		salt,
		uint32(h.iterations),
		uint32(h.memory),
		uint8(h.parallelism),
		uint32(h.keyLen),
	)

	// Prepend salt to hash for storage
	result := make([]byte, h.saltLen+len(hash))
	copy(result[:h.saltLen], salt)
	copy(result[h.saltLen:], hash)

	return result, nil
}

// Verify checks if the password matches the hash.
func (h *Argon2idHasher) Verify(password string, hash []byte) error {
	if len(hash) < h.saltLen+h.keyLen {
		return ErrInvalidHash
	}

	salt := hash[:h.saltLen]
	expectedHash := hash[h.saltLen:]

	actualHash := argon2.IDKey(
		[]byte(password),
		salt,
		uint32(h.iterations),
		uint32(h.memory),
		uint8(h.parallelism),
		uint32(h.keyLen),
	)

	if !constantTimeCompare(actualHash, expectedHash) {
		return ErrInvalidCredentials
	}

	return nil
}

// constantTimeCompare compares two byte slices in constant time.
func constantTimeCompare(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	var result byte
	for i := 0; i < len(a); i++ {
		result |= a[i] ^ b[i]
	}
	return result == 0
}
