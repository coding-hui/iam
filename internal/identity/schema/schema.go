// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package schema

import (
	"encoding/json"
	"errors"
)

// ErrSchemaNotFound is returned when a schema is not found.
var ErrSchemaNotFound = errors.New("schema not found")

// Schema represents an identity schema (Ory Kratos style).
type Schema struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	// JSONSchema is the JSON Schema definition for identity traits.
	JSONSchema json.RawMessage `json:"json_schema"`
}

// DefaultSchema is the default identity schema.
var DefaultSchema = &Schema{
	ID:   "default",
	Type: "person",
	JSONSchema: json.RawMessage(`{
		"$id": "https://schemas.iam.com/schemas/identity.json",
		"$schema": "http://json-schema.org/draft-07/schema#",
		"type": "object",
		"properties": {
			"email": {
				"type": "string",
				"format": "email"
			},
			"phone": {
				"type": "string"
			},
			"name": {
				"type": "string"
			}
		},
		"required": ["email"],
		"additionalProperties": true
	}`),
}

// OAuthSchema is the OAuth/SSO identity schema.
var OAuthSchema = &Schema{
	ID:   "oauth",
	Type: "oauth",
	JSONSchema: json.RawMessage(`{
		"$id": "https://schemas.iam.com/schemas/oauth-identity.json",
		"$schema": "http://json-schema.org/draft-07/schema#",
		"type": "object",
		"properties": {
			"email": {
				"type": "string",
				"format": "email"
			},
			"name": {
				"type": "string"
			},
			"picture": {
				"type": "string",
				"format": "uri"
			}
		},
		"required": ["email"],
		"additionalProperties": true
	}`),
}

// SchemaRegistry manages identity schemas.
type SchemaRegistry struct {
	schemas map[string]*Schema
}

// NewSchemaRegistry creates a new schema registry.
func NewSchemaRegistry() *SchemaRegistry {
	return &SchemaRegistry{
		schemas: map[string]*Schema{
			"default": DefaultSchema,
			"oauth":   OAuthSchema,
		},
	}
}

// Register registers a new schema.
func (r *SchemaRegistry) Register(schema *Schema) error {
	if schema.ID == "" {
		return errors.New("schema ID is required")
	}
	r.schemas[schema.ID] = schema
	return nil
}

// Get returns a schema by ID.
func (r *SchemaRegistry) Get(id string) (*Schema, error) {
	schema, ok := r.schemas[id]
	if !ok {
		return nil, ErrSchemaNotFound
	}
	return schema, nil
}
