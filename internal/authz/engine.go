// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package authz

import (
	"context"
	"encoding/json"
	"sync"
	"time"
)

// Engine is the authorization engine with zero external dependencies.
type Engine struct {
	mu       sync.RWMutex
	policies map[string]*Policy
	cache    map[string]*CachedDecision
	cacheTTL time.Duration
}

// Policy represents a cached policy for authorization.
type Policy struct {
	ID         string
	Subjects   []string
	Effect     string
	Actions    []string
	Resources  []string
	Conditions json.RawMessage
}

// CachedDecision represents a cached authorization decision.
type CachedDecision struct {
	Decision string
	Reason   string
	CachedAt time.Time
}

// AuthzRequest represents an authorization request.
type AuthzRequest struct {
	Subject  string
	Action   string
	Resource string
	Context  map[string]any
}

// AuthzResponse represents an authorization response.
type AuthzResponse struct {
	Decision string
	Reason   string
}

// NewEngine creates a new authorization engine.
func NewEngine() *Engine {
	return &Engine{
		policies: make(map[string]*Policy),
		cache:    make(map[string]*CachedDecision),
		cacheTTL: 5 * time.Minute,
	}
}

// Authorize makes an authorization decision.
func (e *Engine) Authorize(ctx context.Context, req *AuthzRequest) (*AuthzResponse, error) {
	cacheKey := e.cacheKey(req)
	if decision, ok := e.getCachedDecision(cacheKey); ok {
		return decision, nil
	}

	e.mu.RLock()
	defer e.mu.RUnlock()

	for _, p := range e.policies {
		if e.matchesPolicy(req, p) {
			decision := &AuthzResponse{
				Decision: p.Effect,
				Reason:   "matched policy",
			}
			e.setCachedDecision(cacheKey, decision)
			return decision, nil
		}
	}

	return &AuthzResponse{
		Decision: "deny",
		Reason:   "no matching policy",
	}, nil
}

// LoadPolicies loads policies into the engine.
func (e *Engine) LoadPolicies(policies []*Policy) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.policies = make(map[string]*Policy)
	for _, p := range policies {
		for _, subject := range p.Subjects {
			for _, resource := range p.Resources {
				key := subject + ":" + resource
				e.policies[key] = p
			}
		}
	}
}

func (e *Engine) matchesPolicy(req *AuthzRequest, p *Policy) bool {
	// Check subject
	matchedSubject := false
	for _, s := range p.Subjects {
		if s == req.Subject || s == "*" {
			matchedSubject = true
			break
		}
	}
	if !matchedSubject {
		return false
	}

	// Check action
	matchedAction := false
	for _, a := range p.Actions {
		if a == req.Action || a == "*" {
			matchedAction = true
			break
		}
	}
	if !matchedAction {
		return false
	}

	// Check resource
	matchedResource := false
	for _, r := range p.Resources {
		if r == req.Resource || r == "*" {
			matchedResource = true
			break
		}
	}
	return matchedResource
}

func (e *Engine) cacheKey(req *AuthzRequest) string {
	return req.Subject + ":" + req.Action + ":" + req.Resource
}

func (e *Engine) getCachedDecision(key string) (*AuthzResponse, bool) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	decision, ok := e.cache[key]
	if !ok {
		return nil, false
	}

	if time.Since(decision.CachedAt) > e.cacheTTL {
		return nil, false
	}

	return &AuthzResponse{
		Decision: decision.Decision,
		Reason:   decision.Reason,
	}, true
}

func (e *Engine) setCachedDecision(key string, decision *AuthzResponse) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.cache[key] = &CachedDecision{
		Decision: decision.Decision,
		Reason:   decision.Reason,
		CachedAt: time.Now(),
	}
}
