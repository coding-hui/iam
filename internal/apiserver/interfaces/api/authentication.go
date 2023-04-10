package api

import (
	"strings"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"k8s.io/klog/v2"

	"github.com/coding-hui/iam/internal/apiserver/config"
	"github.com/coding-hui/iam/internal/apiserver/domain/service"
	"github.com/coding-hui/iam/internal/pkg/api"
	"github.com/coding-hui/iam/pkg/api/apiserver/v1alpha1"
)

type authentication struct {
	UserService           service.UserService           `inject:""`
	AuthenticationService service.AuthenticationService `inject:""`

	cfg config.Config
}

// NewAuthentication is the  of authentication
func NewAuthentication(c config.Config) Interface {
	return &authentication{cfg: c}
}

func (a *authentication) GetApiGroup() InitApiGroup {
	return InitApiGroup{
		BaseUrl: "",
		Filters: gin.HandlersChain{authCheckFilter},
		Apis: []InitApi{
			{
				Method:  POST,
				Path:    "/login",
				Handler: a.authenticate,
			},
		},
	}
}

func authCheckFilter(c *gin.Context) {
	var tokenValue string
	tokenHeader := c.Request.Header.Get("Authorization")
	if tokenHeader != "" {
		splitted := strings.Split(tokenHeader, " ")
		if len(splitted) != 2 {
			api.Fail(c)
			return
		}
		tokenValue = splitted[1]
	}

	klog.Info("determine token: %s", tokenValue)

	c.Next()
}

func (a *authentication) authenticate(c *gin.Context) {
	var login v1alpha1.AuthenticateRequest
	var err error

	// support header and body both
	if c.Request.Header.Get("Authorization") != "" {
		login, err = parseWithHeader(c)
	} else {
		login, err = parseWithBody(c)
	}
	if err != nil {
		api.FailWithErrCode(err, c)
	}

	response, err := a.AuthenticationService.Authenticate(c.Request.Context(), login)
	if err != nil {
		api.FailWithErrCode(err, c)
	}

	api.OkWithData(response, c)
}

func parseWithHeader(c *gin.Context) (v1alpha1.AuthenticateRequest, error) {
	username, password, ok := c.Request.BasicAuth()
	if !ok {
		return v1alpha1.AuthenticateRequest{}, jwt.ErrFailedAuthentication
	}

	return v1alpha1.AuthenticateRequest{
		Username: username,
		Password: password,
	}, nil
}

func parseWithBody(c *gin.Context) (v1alpha1.AuthenticateRequest, error) {
	var login v1alpha1.AuthenticateRequest
	if err := c.ShouldBindJSON(&login); err != nil {
		return v1alpha1.AuthenticateRequest{}, jwt.ErrFailedAuthentication
	}

	return login, nil
}
