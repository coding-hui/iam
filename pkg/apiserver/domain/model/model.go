package model

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"
	"sigs.k8s.io/yaml"
)

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

// JSONStruct json struct, same with runtime.RawExtension
type JSONStruct map[string]interface{}

// NewJSONStruct new json struct from runtime.RawExtension
func NewJSONStruct(raw *runtime.RawExtension) (*JSONStruct, error) {
	if raw == nil || raw.Raw == nil {
		return nil, nil
	}
	var data JSONStruct
	err := json.Unmarshal(raw.Raw, &data)
	if err != nil {
		return nil, fmt.Errorf("parse raw data failure %w", err)
	}
	return &data, nil
}

// NewJSONStructByString new json struct from string
func NewJSONStructByString(source string) (*JSONStruct, error) {
	if source == "" {
		return nil, nil
	}
	var data JSONStruct
	err := json.Unmarshal([]byte(source), &data)
	if err != nil {
		return nil, fmt.Errorf("parse raw data failure %w", err)
	}
	return &data, nil
}

// NewJSONStructByStruct new json struct from struct object
func NewJSONStructByStruct(object interface{}) (*JSONStruct, error) {
	if object == nil {
		return nil, nil
	}
	var data JSONStruct
	out, err := yaml.Marshal(object)
	if err != nil {
		return nil, fmt.Errorf("marshal object data failure %w", err)
	}
	if err := yaml.Unmarshal(out, &data); err != nil {
		return nil, fmt.Errorf("unmarshal object data failure %w", err)
	}
	return &data, nil
}

// JSON Encoded as a JSON string
func (j *JSONStruct) JSON() string {
	b, err := json.Marshal(j)
	if err != nil {
		klog.Errorf("json marshal failure %s", err.Error())
	}
	return string(b)
}

// Properties return the map
func (j *JSONStruct) Properties() map[string]interface{} {
	return *j
}

// RawExtension Encoded as a RawExtension
func (j *JSONStruct) RawExtension() *runtime.RawExtension {
	yamlByte, err := yaml.Marshal(j)
	if err != nil {
		klog.Errorf("yaml marshal failure %s", err.Error())
		return nil
	}
	b, err := yaml.YAMLToJSON(yamlByte)
	if err != nil {
		klog.Errorf("yaml to json failure %s", err.Error())
		return nil
	}
	if len(b) == 0 || string(b) == "null" {
		return nil
	}
	return &runtime.RawExtension{Raw: b}
}

// BaseModel common model
type BaseModel struct {
	// CreatedAt is a timestamp representing the server time when this object was
	// created. It is not guaranteed to be set in happens-before order across separate operations.
	// Clients may not set this value. It is represented in RFC3339 form and is in UTC.
	CreatedAt time.Time `json:"created_at,omitempty" gorm:"column:created_at"`

	// UpdatedAt is a timestamp representing the server time when this object was updated.
	// Clients may not set this value. It is represented in RFC3339 form and is in UTC.
	UpdatedAt time.Time `json:"updated_at,omitempty" gorm:"column:updated_at"`
}

// SetCreateTime set create time
func (m *BaseModel) SetCreateTime(time time.Time) {
	m.CreatedAt = time
}

// SetUpdateTime set update time
func (m *BaseModel) SetUpdateTime(time time.Time) {
	m.UpdatedAt = time
}

func deepCopy(src interface{}) interface{} {
	dst := reflect.New(reflect.TypeOf(src).Elem())

	val := reflect.ValueOf(src).Elem()
	nVal := dst.Elem()
	for i := 0; i < val.NumField(); i++ {
		nvField := nVal.Field(i)
		nvField.Set(val.Field(i))
	}

	return dst.Interface()
}
