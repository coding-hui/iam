package service

import (
	"context"
	"fmt"

	"github.com/wecoding/iam/pkg/apiserver/config"
)

// needInitData register the service that need to init data
var needInitData []DataInit

// InitServiceBean init all service instance
func InitServiceBean(c config.Config) []interface{} {
	authenticationService := NewAuthenticationService(c)
	userService := NewUserService()

	needInitData = []DataInit{userService}

	return []interface{}{userService, authenticationService}
}

// DataInit the service set that needs init data
type DataInit interface {
	Init(ctx context.Context) error
}

// InitData init data
func InitData(ctx context.Context) error {
	for _, init := range needInitData {
		if err := init.Init(ctx); err != nil {
			return fmt.Errorf("database init failure %w", err)
		}
	}

	return nil
}
