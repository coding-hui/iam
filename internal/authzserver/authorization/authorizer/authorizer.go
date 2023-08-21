// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package authorizer

import (
	"fmt"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	rediswatcher "github.com/casbin/redis-watcher/v2"
	"github.com/redis/go-redis/v9"

	"github.com/coding-hui/iam/internal/authzserver/adapter"
	"github.com/coding-hui/iam/internal/authzserver/config"
	"github.com/coding-hui/iam/internal/pkg/channel"
	"github.com/coding-hui/iam/internal/pkg/mycasbin"
	authzv1 "github.com/coding-hui/iam/pkg/api/authzserver/v1"
)

// Authorizer implement the authorize interface that use local repository to
// authorize the subject access review.
type Authorizer struct {
	enforcer *casbin.CachedEnforcer
	cfg      config.Config
}

// NewAuthorizer creates a local repository authorizer and returns it.
func NewAuthorizer(c config.Config, a *adapter.Adapter) (*Authorizer, error) {
	m, err := model.NewModelFromString(mycasbin.ModelText)
	if err != nil {
		return nil, err
	}
	enforcer, err := casbin.NewCachedEnforcer(m, a)
	if err != nil {
		return nil, err
	}

	enforcer.ClearPolicy()
	enforcer.EnableLog(true)
	enforcer.SetExpireTime(time.Minute)

	err = enforcer.LoadPolicy()
	if err != nil {
		return nil, err
	}

	// Initialize the watcher.
	// Use the Redis host as parameter.
	w, err := rediswatcher.NewWatcher(
		fmt.Sprintf("%s:%d", c.RedisOptions.Host, c.RedisOptions.Port),
		rediswatcher.WatcherOptions{
			Options: redis.Options{
				Password: c.RedisOptions.Password,
			},
			Channel: channel.RedisPubSubChannel,
		})
	if err != nil {
		return nil, err
	}

	// Set callback
	err = w.SetUpdateCallback(rediswatcher.DefaultUpdateCallback(enforcer))
	if err != nil {
		return nil, err
	}

	// Set the watcher for the enforcer.
	err = enforcer.SetWatcher(w)
	if err != nil {
		return nil, err
	}

	return &Authorizer{
		cfg:      c,
		enforcer: enforcer,
	}, err
}

func (a *Authorizer) Authorize(r *authzv1.Request) *authzv1.Response {
	allowed, err := a.enforcer.Enforce(r.Subject, r.Resource, r.Action)
	if err != nil {
		return &authzv1.Response{
			Denied: true,
			Reason: err.Error(),
		}
	}

	if !allowed {
		return &authzv1.Response{
			Denied: true,
		}
	}

	return &authzv1.Response{
		Allowed: true,
	}
}
