// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package schema

import (
	"encoding/json"
	"errors"
)

// ErrInvalidJSON is returned when traits is not valid JSON.
var ErrInvalidJSON = errors.New("traits must be valid JSON")

// ErrValidation is returned when traits validation fails.
var ErrValidation = errors.New("validation failed")

// Validator validates identity traits against a schema.
type Validator struct {
	registry *SchemaRegistry
}

// NewValidator creates a new schema validator.
func NewValidator(registry *SchemaRegistry) *Validator {
	return &Validator{registry: registry}
}

// Validate validates the given traits against the schema identified by schemaID.
func (v *Validator) Validate(schemaID string, traits json.RawMessage) error {
	schema, err := v.registry.Get(schemaID)
	if err != nil {
		return err
	}

	return schema.Validate(traits)
}

// Validate validates the given traits against this schema.
func (s *Schema) Validate(traits json.RawMessage) error {
	if s.JSONSchema == nil {
		return nil // No schema defined, skip validation
	}

	// Parse the JSON Schema
	var js schemaDefinition
	if err := json.Unmarshal(s.JSONSchema, &js); err != nil {
		return err
	}

	// Parse the traits
	var traitsData map[string]any
	if err := json.Unmarshal(traits, &traitsData); err != nil {
		return ErrInvalidJSON
	}

	// Check required fields
	for _, required := range js.Required {
		if _, ok := traitsData[required]; !ok {
			return errors.New("required field missing: " + required)
		}
	}

	// Basic type validation
	for field, fieldSchema := range js.Properties {
		if value, ok := traitsData[field]; ok {
			expectedType := fieldSchema.Type
			if err := validateType(value, expectedType); err != nil {
				return err
			}
		}
	}

	return nil
}

// schemaDefinition represents a subset of JSON Schema for validation.
type schemaDefinition struct {
	Type       string                     `json:"type"`
	Required   []string                   `json:"required"`
	Properties map[string]fieldDefinition `json:"properties"`
}

// fieldDefinition defines a field in the schema.
type fieldDefinition struct {
	Type   string `json:"type"`
	Format string `json:"format"`
}

func validateType(value any, expectedType string) error {
	switch expectedType {
	case "string":
		if _, ok := value.(string); !ok {
			return errors.New("expected string")
		}
	case "number", "integer":
		switch value.(type) {
		case float64, int, int64:
			// JSON numbers are float64 by default
		default:
			return errors.New("expected number")
		}
	case "boolean":
		if _, ok := value.(bool); !ok {
			return errors.New("expected boolean")
		}
	case "array":
		if _, ok := value.([]any); !ok {
			return errors.New("expected array")
		}
	case "object":
		if _, ok := value.(map[string]any); !ok {
			return errors.New("expected object")
		}
	}
	return nil
}
