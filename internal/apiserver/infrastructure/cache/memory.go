// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package cache

import (
	"regexp"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
	"k8s.io/apimachinery/pkg/util/wait"

	"github.com/coding-hui/common/errors"

	"github.com/coding-hui/iam/pkg/options"
)

var ErrNoSuchKey = errors.New("no such key")

const (
	TypeInMemoryCache    = "InMemoryCache"
	defaultCleanupPeriod = 2 * time.Hour
)

type simpleObject struct {
	value       string
	neverExpire bool
	expiredAt   time.Time
}

func (so *simpleObject) IsExpired() bool {
	if so.neverExpire {
		return false
	}
	if time.Now().After(so.expiredAt) {
		return true
	}
	return false
}

// InMemoryCacheOptions used to create inMemoryCache in memory.
// CleanupPeriod specifies cleans up expired token every period.
// Note the SimpleCache cannot be used in multi-replicas apiserver,
// which will lead to data inconsistency.
type InMemoryCacheOptions struct {
	CleanupPeriod time.Duration `json:"cleanupPeriod" yaml:"cleanupPeriod" mapstructure:"cleanupperiod"`
}

// imMemoryCache implements cache.Interface use memory objects, it should be used only for testing.
type inMemoryCache struct {
	store map[string]simpleObject
}

func NewInMemoryCache(options *InMemoryCacheOptions, stopCh <-chan struct{}) (Interface, error) {
	var cleanupPeriod time.Duration
	cache := &inMemoryCache{
		store: make(map[string]simpleObject),
	}

	if options == nil || options.CleanupPeriod == 0 {
		cleanupPeriod = defaultCleanupPeriod
	} else {
		cleanupPeriod = options.CleanupPeriod
	}
	go wait.Until(cache.cleanInvalidToken, cleanupPeriod, stopCh)

	return cache, nil
}

func (s *inMemoryCache) cleanInvalidToken() {
	for k, v := range s.store {
		if v.IsExpired() {
			delete(s.store, k)
		}
	}
}

func (s *inMemoryCache) Keys(pattern string) ([]string, error) {
	// There is a little difference between go regexp and redis key pattern
	// In redis, * means any character, while in go . means match everything.
	pattern = strings.ReplaceAll(pattern, "*", ".*")

	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	var keys []string
	for k := range s.store {
		if re.MatchString(k) {
			keys = append(keys, k)
		}
	}

	return keys, nil
}

func (s *inMemoryCache) Set(key string, value string, duration time.Duration) error {
	sobject := simpleObject{
		value:       value,
		neverExpire: false,
		expiredAt:   time.Now().Add(duration),
	}

	if duration == NeverExpire {
		sobject.neverExpire = true
	}

	s.store[key] = sobject
	return nil
}

func (s *inMemoryCache) Del(keys ...string) error {
	for _, key := range keys {
		delete(s.store, key)
	}
	return nil
}

func (s *inMemoryCache) Get(key string) (string, error) {
	if sobject, ok := s.store[key]; ok {
		if sobject.neverExpire || time.Now().Before(sobject.expiredAt) {
			return sobject.value, nil
		}
	}

	return "", ErrNoSuchKey
}

func (s *inMemoryCache) Exists(keys ...string) (bool, error) {
	for _, key := range keys {
		if _, ok := s.store[key]; !ok {
			return false, nil
		}
	}

	return true, nil
}

func (s *inMemoryCache) Expire(key string, duration time.Duration) error {
	value, err := s.Get(key)
	if err != nil {
		return err
	}

	sobject := simpleObject{
		value:       value,
		neverExpire: false,
		expiredAt:   time.Now().Add(duration),
	}

	if duration == NeverExpire {
		sobject.neverExpire = true
	}

	s.store[key] = sobject
	return nil
}

type inMemoryCacheFactory struct{}

func (sf *inMemoryCacheFactory) Type() string {
	return TypeInMemoryCache
}

func (sf *inMemoryCacheFactory) Create(options options.DynamicOptions, stopCh <-chan struct{}) (Interface, error) {
	var sOptions InMemoryCacheOptions
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		DecodeHook:       mapstructure.StringToTimeDurationHookFunc(),
		WeaklyTypedInput: true,
		Result:           &sOptions,
	})
	if err != nil {
		return nil, err
	}
	if err := decoder.Decode(options); err != nil {
		return nil, err
	}

	return NewInMemoryCache(&sOptions, stopCh)
}

func init() {
	RegisterCacheFactory(&inMemoryCacheFactory{})
}
