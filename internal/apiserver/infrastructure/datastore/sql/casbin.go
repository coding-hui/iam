// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package sql

import (
	"fmt"
	"sync"

	"github.com/casbin/casbin/v2"
	rediswatcher "github.com/casbin/redis-watcher/v2"
	"github.com/redis/go-redis/v9"

	"github.com/coding-hui/iam/internal/apiserver/domain/model"
	"github.com/coding-hui/iam/internal/apiserver/domain/repository"
	"github.com/coding-hui/iam/internal/pkg/channel"
	"github.com/coding-hui/iam/internal/pkg/mycasbin"
	genericoptions "github.com/coding-hui/iam/pkg/options"
)

type casbinRepositoryImpl struct {
	once      sync.Once
	client    *Client
	redisOpts *genericoptions.RedisOptions
	enforcer  *casbin.SyncedEnforcer
}

// newCasbinRepository new SyncedEnforcer Repository.
func newCasbinRepository(client *Client, redisOpts *genericoptions.RedisOptions) repository.CasbinRepository {
	return &casbinRepositoryImpl{client: client, redisOpts: redisOpts}
}

func (c *casbinRepositoryImpl) SyncedEnforcer() *casbin.SyncedEnforcer {
	c.once.Do(func() {
		c.enforcer = mycasbin.Setup(c.client.db, model.TableNamePrefix, "casbin_rules")

		// Initialize the watcher.
		// Use the Redis host as parameter.
		w, _ := rediswatcher.NewPublishWatcher(
			fmt.Sprintf("%s:%d", c.redisOpts.Host, c.redisOpts.Port),
			rediswatcher.WatcherOptions{
				Options: redis.Options{
					Password: c.redisOpts.Password,
				},
				Channel: channel.RedisPubSubChannel,
			})

		// Set the watcher for the enforcer.
		_ = c.enforcer.SetWatcher(w)
	})
	return c.enforcer
}
