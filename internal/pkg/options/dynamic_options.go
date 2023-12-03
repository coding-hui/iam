// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package options

import (
	"encoding/json"
	"strings"

	"github.com/mitchellh/mapstructure"
)

// DynamicOptions accept dynamic configuration, the type of key MUST be string.
type DynamicOptions map[string]interface{}

func (o DynamicOptions) MarshalJSON() ([]byte, error) {
	data, err := json.Marshal(desensitize(o))
	return data, err
}

func (o DynamicOptions) To(target interface{}) error {
	err := mapstructure.Decode(o, target)
	if err != nil {
		return err
	}
	return nil
}

var sensitiveKeys = [...]string{"password", "secret"}

// isSensitiveData returns whether the input string contains sensitive information.
func isSensitiveData(key string) bool {
	for _, v := range sensitiveKeys {
		if strings.Contains(strings.ToLower(key), v) {
			return true
		}
	}
	return false
}

// desensitize returns the desensitized data.
func desensitize(data map[string]interface{}) map[string]interface{} {
	output := make(map[string]interface{})
	for k, v := range data {
		if isSensitiveData(k) {
			continue
		}
		switch v := v.(type) {
		case map[interface{}]interface{}:
			output[k] = desensitize(convert(v))
		default:
			output[k] = v
		}
	}
	return output
}

// convert returns formatted data.
func convert(m map[interface{}]interface{}) map[string]interface{} {
	output := make(map[string]interface{})
	for k, v := range m {
		switch k := k.(type) {
		case string:
			output[k] = v
		}
	}
	return output
}
