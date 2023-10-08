// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package options

import (
	_ "github.com/coding-hui/iam/internal/apiserver/domain/service/identityprovider/gitee"
	_ "github.com/coding-hui/iam/internal/apiserver/domain/service/identityprovider/github"
	_ "github.com/coding-hui/iam/internal/apiserver/domain/service/identityprovider/ldap"
	_ "github.com/coding-hui/iam/internal/apiserver/domain/service/identityprovider/wechatmini"

	"github.com/coding-hui/iam/internal/apiserver/domain/service/identityprovider"
)

// Validate checks Options and return a slice of found errs.
func (o *Options) Validate() []error {
	var errs []error

	errs = append(errs, o.GenericServerRunOptions.Validate()...)
	errs = append(errs, o.GRPCOptions.Validate()...)
	errs = append(errs, o.InsecureServing.Validate()...)
	errs = append(errs, o.SecureServing.Validate()...)
	errs = append(errs, o.MySQLOptions.Validate()...)
	errs = append(errs, o.RedisOptions.Validate()...)
	errs = append(errs, o.JwtOptions.Validate()...)
	errs = append(errs, o.FeatureOptions.Validate()...)
	errs = append(errs, o.AuthenticationOptions.Validate()...)

	if err := identityprovider.SetupWithOptions(o.AuthenticationOptions.OAuthOptions.IdentityProviders); err != nil {
		errs = append(errs, err)
	}

	return errs
}
