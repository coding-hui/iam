// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import (
	"fmt"
	"strings"
)

// TableNamePrefix table name prefix.
var TableNamePrefix = "iam_"

var registeredModels = map[string]Interface{}

// Interface model interface.
type Interface interface {
	TableName() string
}

// RegisterModel register model.
func RegisterModel(models ...Interface) {
	for _, model := range models {
		if _, exist := registeredModels[model.TableName()]; exist {
			panic(fmt.Errorf("model table name %s conflict", model.TableName()))
		}
		registeredModels[model.TableName()] = model
	}
}

// GetRegisterModels will return the register models.
func GetRegisterModels() map[string]Interface {
	return registeredModels
}

func GetResourceIdentifier(instanceID string) string {
	if len(instanceID) == 0 {
		return ""
	}
	prefixIdx := strings.Index(instanceID, "-")
	if prefixIdx > -1 {
		return instanceID[:prefixIdx]
	}
	return instanceID
}
