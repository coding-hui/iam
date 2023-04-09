package model

import (
	"fmt"
)

// TableNamePrefix table name prefix
var TableNamePrefix = "iam_"

var registeredModels = map[string]Interface{}

// Interface model interface
type Interface interface {
	TableName() string
	ShortTableName() string
}

// RegisterModel register model
func RegisterModel(models ...Interface) {
	for _, model := range models {
		if _, exist := registeredModels[model.TableName()]; exist {
			panic(fmt.Errorf("model table name %s conflict", model.TableName()))
		}
		registeredModels[model.TableName()] = model
	}
}

// GetRegisterModels will return the register models
func GetRegisterModels() map[string]Interface {
	return registeredModels
}
