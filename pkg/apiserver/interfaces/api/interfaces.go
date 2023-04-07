package api

import "github.com/gin-gonic/gin"

// versionPrefix API version prefix.
var versionPrefix = "/api/v1"

// viewPrefix the path prefix for view page
var viewPrefix = "/view"

// GetAPIPrefix return the prefix of the api route path
func GetAPIPrefix() []string {
	return []string{versionPrefix, viewPrefix, "/v1"}
}

// InitApiGroup the API should define the http route
type InitApiGroup struct {
	BaseUrl string
	Apis    []InitApi
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
func InitAPIBean() []interface{} {
	// Ping
	RegisterAPI(NewPing())

	// Authentication
	RegisterAPI(NewUser())

	var beans []interface{}
	for i := range registeredAPI {
		beans = append(beans, registeredAPI[i])
	}
	return beans
}
