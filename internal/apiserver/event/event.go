// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package event

// Event interface
type Event interface {
	Name() string
	Get(key string) any
	Set(key string, val any)
	Add(key string, val any)
	Data() any
	SetData(val any) Event
}

// Cloneable interface. event can be cloned.
type Cloneable interface {
	Event
	Clone() Event
}

// FactoryFunc for create event instance.
type FactoryFunc func() Event

// BasicEvent a built-in implements Event interface
type BasicEvent struct {
	// event name
	name string
	// user data.
	data map[string]any
}

// New create an event instance
func New(name string, data map[string]any) *BasicEvent {
	return NewBasic(name, data)
}

// NewBasic new a basic event instance
func NewBasic(name string, data map[string]any) *BasicEvent {
	if data == nil {
		data = make(map[string]any)
	}

	return &BasicEvent{
		name: name,
		data: data,
	}
}

// Get data by index
func (e *BasicEvent) Get(key string) any {
	if v, ok := e.data[key]; ok {
		return v
	}

	return nil
}

// Add value by key
func (e *BasicEvent) Add(key string, val any) {
	if _, ok := e.data[key]; !ok {
		e.Set(key, val)
	}
}

// Set value by key
func (e *BasicEvent) Set(key string, val any) {
	if e.data == nil {
		e.data = make(map[string]any)
	}

	e.data[key] = val
}

// Name get event name
func (e *BasicEvent) Name() string {
	return e.name
}

// Data get all data
func (e *BasicEvent) Data() any {
	return e.data
}

// Clone new instance
func (e *BasicEvent) Clone() Event {
	var cp = *e
	return &cp
}

// SetName set event name
func (e *BasicEvent) SetName(name string) *BasicEvent {
	e.name = name
	return e
}

// SetData set data to the event
func (e *BasicEvent) SetData(data any) Event {
	if data == nil {
		return e
	}
	if v, ok := data.(map[string]any); ok {
		e.data = v
	}
	return e
}
