// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package api

import (
	"github.com/gin-gonic/gin"

	"github.com/coding-hui/iam/internal/apiserver/config"
)

const (
	GET    string = "GET"
	POST   string = "POST"
	PUT    string = "PUT"
	DELETE string = "DELETE"
)

// versionPrefix API version prefix.
var versionPrefix = "/api/v1"

// GetAPIPrefix return the prefix of the api route path
func GetAPIPrefix() []string {
	return []string{versionPrefix, "/v1"}
}

// InitApiGroup the API should define the http route
type InitApiGroup struct {
	BaseUrl string
	Apis    []InitApi
	Filters gin.HandlersChain
}

// InitApi the API should define the http route
type InitApi struct {
	Method  string
	Path    string
	Handler gin.HandlerFunc
}

// Interface the API should define the http route
type Interface interface {
	GetApiGroup() InitApiGroup
}

var registeredAPI []Interface

// RegisterAPI register API handler
func RegisterAPI(apis Interface) {
	registeredAPI = append(registeredAPI, apis)
}

// GetRegisteredAPI return all API handlers
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

	var beans []interface{}
	for i := range registeredAPI {
		beans = append(beans, registeredAPI[i])
	}

	return beans
}
