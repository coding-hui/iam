// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package api

import (
	"github.com/gin-gonic/gin"

	"github.com/coding-hui/iam/internal/apiserver/config"
)

// versionPrefix API version prefix.
var versionPrefix = "/api/v1"

// GetAPIPrefix return the prefix of the api route path.
func GetAPIPrefix() []string {
	return []string{versionPrefix, "/v1"}
}

// Interface the API should define the http route.
type Interface interface {
	RegisterApiGroup(*gin.Engine)
}

var registeredAPI []Interface

// RegisterAPI register API handler.
func RegisterAPI(apis Interface) {
	registeredAPI = append(registeredAPI, apis)
}

// GetRegisteredAPI return all API handlers.
func GetRegisteredAPI() []Interface {
	return registeredAPI
}

// InitAPIBean inits all API handlers, pass in the required parameter object.
// It can be implemented using the idea of dependency injection.
func InitAPIBean(c config.Config) []interface{} {
	// Ping
	RegisterAPI(NewPing())

	// Authentication
	RegisterAPI(NewAuthentication(c))
	RegisterAPI(NewUser())
	RegisterAPI(NewResource())
	RegisterAPI(NewRole())
	RegisterAPI(NewOrganization())
	RegisterAPI(NewDepartment())

	// policies
	RegisterAPI(NewPolicy())

	// grpc cache
	RegisterAPI(NewCacheServer())

	// providers
	RegisterAPI(NewIdentityProvider())

	// apps
	RegisterAPI(NewApplication())

	beans := make([]interface{}, 0, len(registeredAPI))
	for i := range registeredAPI {
		beans = append(beans, registeredAPI[i])
	}

	return beans
}
