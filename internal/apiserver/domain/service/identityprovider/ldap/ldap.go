// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package ldap

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"os"
	"time"

	"github.com/go-ldap/ldap"
	"github.com/mitchellh/mapstructure"

	"github.com/coding-hui/common/errors"

	"github.com/coding-hui/iam/internal/apiserver/domain/service/identityprovider"
	"github.com/coding-hui/iam/internal/pkg/code"
	"github.com/coding-hui/iam/internal/pkg/options"
	v1 "github.com/coding-hui/iam/pkg/api/apiserver/v1"
	"github.com/coding-hui/iam/pkg/log"
)

const (
	defaultReadTimeout = 15000
)

func init() {
	identityprovider.RegisterGenericProvider(&ldapProviderFactory{})
}

type ldapProvider struct {
	// Host and optional port of the LDAP server in the form "host:port".
	// If the port is not supplied, 389 for insecure or StartTLS connections, 636
	Host string `json:"host,omitempty"                 yaml:"host"`
	// Timeout duration when reading data from remote server. Default to 15s.
	ReadTimeout int `json:"readTimeout"                    yaml:"readTimeout"`
	// If specified, connections will use the ldaps:// protocol
	StartTLS bool `json:"startTLS,omitempty"             yaml:"startTLS"`
	// Used to turn off TLS certificate checks
	InsecureSkipVerify bool `json:"insecureSkipVerify"             yaml:"insecureSkipVerify"`
	// Path to a trusted root certificate file. Default: use the host's root CA.
	RootCA string `json:"rootCA,omitempty"               yaml:"rootCA"`
	// A raw certificate file can also be provided inline. Base64 encoded PEM file
	RootCAData string `json:"rootCAData,omitempty"           yaml:"rootCAData"`
	// Username (DN) of the "manager" user identity.
	ManagerDN string `json:"managerDN,omitempty"            yaml:"managerDN"`
	// The password for the manager DN.
	ManagerPassword string `json:"managerPassword,omitempty"      yaml:"managerPassword"`
	// User search scope.
	UserSearchBase string `json:"userSearchBase,omitempty"       yaml:"userSearchBase"`
	// LDAP filter used to identify objects of type user. e.g. (objectClass=person)
	UserSearchFilter string `json:"userSearchFilter,omitempty"     yaml:"userSearchFilter"`
	// Group search scope.
	GroupSearchBase string `json:"groupSearchBase,omitempty"      yaml:"groupSearchBase"`
	// LDAP filter used to identify objects of type group. e.g. (objectclass=group)
	GroupSearchFilter string `json:"groupSearchFilter,omitempty"    yaml:"groupSearchFilter"`
	// Attribute on a user object storing the groups the user is a member of.
	UserMemberAttribute string `json:"userMemberAttribute,omitempty"  yaml:"userMemberAttribute"`
	// Attribute on a group object storing the information for primary group membership.
	GroupMemberAttribute string `json:"groupMemberAttribute,omitempty" yaml:"groupMemberAttribute"`
	// The following three fields are direct mappings of attributes on the user entry.
	// login attribute used for comparing user entries.
	LoginAttribute string `json:"loginAttribute"                 yaml:"loginAttribute"`
	MailAttribute  string `json:"mailAttribute"                  yaml:"mailAttribute"`
}

type ldapProviderFactory struct{}

func (l *ldapProviderFactory) Type() v1.IdentityProviderType {
	return v1.LDAPIdentityProvider
}

func (l *ldapProviderFactory) Create(opts options.DynamicOptions) (identityprovider.GenericProvider, error) {
	var ldapProvider ldapProvider
	if err := mapstructure.Decode(opts, &ldapProvider); err != nil {
		return nil, err
	}
	if ldapProvider.ReadTimeout <= 0 {
		ldapProvider.ReadTimeout = defaultReadTimeout
	}
	return &ldapProvider, nil
}

type ldapIdentity struct {
	Username string
	Email    string
}

func (l *ldapIdentity) GetUserID() string {
	return l.Username
}

func (l *ldapIdentity) GetUsername() string {
	return l.Username
}

func (l *ldapIdentity) GetEmail() string {
	return l.Email
}

func (l *ldapIdentity) GetAvatar() string {
	return ""
}

func (l ldapProvider) Authenticate(req v1.AuthenticateRequest) (identityprovider.Identity, error) {
	username, password := req.Username, req.Password
	if username == "" && password == "" {
		return nil, errors.WithCode(code.ErrMissingLoginValues, "")
	}
	conn, err := l.newConn()
	if err != nil {
		log.Errorf("failed to authenticate ldap: %v", err)
		return nil, err
	}

	conn.SetTimeout(time.Duration(l.ReadTimeout) * time.Millisecond)
	defer conn.Close()

	if err = conn.Bind(l.ManagerDN, l.ManagerPassword); err != nil {
		log.Errorf("failed connection ldap: %v", err)
		return nil, err
	}

	filter := fmt.Sprintf("(%s=%s)", l.LoginAttribute, ldap.EscapeFilter(username))
	if l.UserSearchFilter != "" {
		filter = fmt.Sprintf("(&%s%s)", filter, l.UserSearchFilter)
	}
	result, err := conn.Search(&ldap.SearchRequest{
		BaseDN:       l.UserSearchBase,
		Scope:        ldap.ScopeWholeSubtree,
		DerefAliases: ldap.NeverDerefAliases,
		SizeLimit:    1,
		TimeLimit:    0,
		TypesOnly:    false,
		Filter:       filter,
		Attributes:   []string{l.LoginAttribute, l.MailAttribute},
	})
	if err != nil {
		return nil, err
	}

	if len(result.Entries) == 0 {
		return nil, errors.WithCode(code.ErrUnauthorized, "ldap: no results returned for filter: %v", filter)
	}

	if len(result.Entries) > 1 {
		return nil, errors.WithCode(code.ErrUnauthorized, "ldap: filter returned multiple results: %v", filter)
	}

	// len(result.Entries) == 1
	entry := result.Entries[0]
	if err = conn.Bind(entry.DN, password); err != nil {
		if ldap.IsErrorWithCode(err, ldap.LDAPResultInvalidCredentials) {
			log.V(4).Infof("ldap: %v", err)
			return nil, errors.WithCode(code.ErrUnknown, "ldap: incorrect password")
		}
		return nil, err
	}
	email := entry.GetAttributeValue(l.MailAttribute)
	uid := entry.GetAttributeValue(l.LoginAttribute)
	return &ldapIdentity{
		Username: uid,
		Email:    email,
	}, nil
}

func (l ldapProvider) newConn() (*ldap.Conn, error) {
	if !l.StartTLS {
		return ldap.Dial("tcp", l.Host)
	}
	tlsConfig := tls.Config{}
	if l.InsecureSkipVerify {
		tlsConfig.InsecureSkipVerify = true
	}
	tlsConfig.RootCAs = x509.NewCertPool()
	var caCert []byte
	var err error
	// Load CA cert
	if l.RootCA != "" {
		if caCert, err = os.ReadFile(l.RootCA); err != nil {
			return nil, err
		}
	}
	if l.RootCAData != "" {
		if caCert, err = base64.StdEncoding.DecodeString(l.RootCAData); err != nil {
			return nil, err
		}
	}
	if caCert != nil {
		tlsConfig.RootCAs.AppendCertsFromPEM(caCert)
	}
	return ldap.DialTLS("tcp", l.Host, &tlsConfig)
}
