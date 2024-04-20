// Copyright (c) 2024 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package event

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/coding-hui/iam/internal/apiserver/config"

	"github.com/coding-hui/iam/pkg/log"
)

const (
	defaultChannelSize = 100
	defaultConsumerNum = 3
)

// Bus type defines the bus interface structure.
type Bus interface {
	Publish(e Event) error
	AsyncPublish(e Event)
	AddEventListener(name string, listener Listener)

	Close() error
	CloseWait() error
	Wait() error
}

// InProcBus defines the in-process bus structure.
type InProcBus struct {
	sync.Mutex

	wg  sync.WaitGroup
	ch  chan Event
	oc  sync.Once
	err error

	cfg config.Config

	listeners map[string]*ListenerQueue

	// ChannelSize for fire events by goroutine
	ChannelSize int
	ConsumerNum int
}

// InitEvent init event bus.
func InitEvent(c config.Config) (bus Bus, listeners []interface{}) {
	bus = NewEventBus(c)

	authenticationSuccessListener := NewAuthenticationEventListener()

	bus.AddEventListener(AuthenticationEventType, authenticationSuccessListener)

	return bus, []interface{}{
		authenticationSuccessListener,
	}
}

func NewEventBus(c config.Config) *InProcBus {
	return &InProcBus{
		cfg:       c,
		listeners: make(map[string]*ListenerQueue),
	}
}

// Publish function publish a message to the bus listener.
func (b *InProcBus) Publish(e Event) error {
	b.Lock()
	defer b.Unlock()

	var msgName = e.Name()
	if lq, ok := b.listeners[msgName]; ok {
		for _, li := range lq.Sort().Items() {
			err := li.Listener.Handle(e)
			if err != nil {
				if rawErr, ok := e.(error); ok {
					return rawErr
				}
				return fmt.Errorf("expected listener to return an error, got '%T'", e)
			}
		}
	}
	return nil
}

func (b *InProcBus) AsyncPublish(e Event) {
	b.oc.Do(func() {
		b.makeConsumers()
	})

	// dispatch event
	b.ch <- e
}

func (b *InProcBus) makeConsumers() {
	if b.ConsumerNum <= 0 {
		b.ConsumerNum = defaultConsumerNum
	}
	if b.ChannelSize <= 0 {
		b.ChannelSize = defaultChannelSize
	}

	b.ch = make(chan Event, b.ChannelSize)

	for i := 0; i < b.ConsumerNum; i++ {
		b.wg.Add(1)

		go func() {
			defer func() {
				if err := recover(); err != nil {
					b.err = fmt.Errorf("async consum event error: %v", err)
				}
				b.wg.Done()
			}()

			for e := range b.ch {
				_ = b.Publish(e)
			}
		}()
	}
}

func (b *InProcBus) AddEventListener(name string, listener Listener) {
	b.addListenerItem(name, &ListenerItem{0, listener})
}

func (b *InProcBus) addListenerItem(name string, li *ListenerItem) {
	if li.Listener == nil {
		log.Warnf("Failed to add listener: the event %q listener cannot be empty", name)
		return
	}
	if reflect.ValueOf(li.Listener).Kind() == reflect.Struct {
		log.Warnf("Failed to add listener: %q - struct listener must be pointer", name)
		return
	}

	if lq, ok := b.listeners[name]; ok {
		lq.Push(li)
	} else {
		b.listeners[name] = (&ListenerQueue{}).Push(li)
	}
}

// CloseWait close channel and wait all async event done.
func (b *InProcBus) CloseWait() error {
	if err := b.Close(); err != nil {
		return err
	}
	return b.Wait()
}

// Wait wait all async event done.
func (b *InProcBus) Wait() error {
	b.wg.Wait()
	return b.err
}

// Close event channel, deny to fire new event.
func (b *InProcBus) Close() error {
	if b.ch != nil {
		close(b.ch)
	}
	return nil
}
