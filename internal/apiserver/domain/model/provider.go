// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import (
	"gorm.io/gorm"

	v1 "github.com/coding-hui/iam/pkg/api/apiserver/v1"

	metav1 "github.com/coding-hui/common/meta/v1"
	"github.com/coding-hui/common/util/idutil"
)

func init() {
	RegisterModel(&Provider{})
}

// Provider represent a provider model.
type Provider struct {
	// Standard object's metadata.
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Type          string           `json:"type"          gorm:"column:type;type:varchar(64)"`
	Category      string           `json:"category"      gorm:"column:category;type:varchar(20)"`
	Status        string           `json:"status"        gorm:"column:status;type:varchar(20);default:0"`
	Owner         string           `json:"owner"         gorm:"column:owner;type:varchar(100)"`
	DisplayName   string           `json:"displayName"   gorm:"column:display_name;type:varchar(100)"`
	Description   string           `json:"description"   gorm:"column:description;type:varchar(100)"`
	MappingMethod v1.MappingMethod `json:"mappingMethod" gorm:"column:mapping_method;type:varchar(20)"`
}

// TableName maps to mysql table name.
func (p *Provider) TableName() string {
	return TableNamePrefix + "provider"
}

// AfterCreate run after create database record.
func (p *Provider) AfterCreate(tx *gorm.DB) error {
	p.InstanceID = idutil.GetInstanceID(p.ID, "p-")

	return tx.Save(p).Error
}

func (p *Provider) ToIdentityProvider() *IdentityProvider {
	idp := &IdentityProvider{
		ObjectMeta:  metav1.ObjectMeta{},
		Type:        p.Type,
		Category:    "",
		Status:      p.Status,
		Owner:       p.Owner,
		DisplayName: p.DisplayName,
		Description: p.Description,
	}
	if clientID, ok := p.Extend.Get("clientid").(string); ok {
		idp.ClientID = clientID
	}
	if clientSecret, ok := p.Extend.Get("clientsecret").(string); ok {
		idp.ClientSecret = clientSecret
	}
	if redirectURL, ok := p.Extend.Get("redirecturl").(string); ok {
		idp.RedirectURL = redirectURL
	}
	if insecureSkipVerify, ok := p.Extend.Get("insecureskipverify").(bool); ok {
		idp.InsecureSkipVerify = insecureSkipVerify
	}
	if scopes, ok := p.Extend.Get("scopes").([]string); ok {
		idp.Scopes = scopes
	}
	return idp
}

type IdentityProvider struct {
	// Standard object's metadata.
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Type          string              `json:"type"`
	Category      v1.ProviderCategory `json:"category"`
	MappingMethod v1.MappingMethod    `json:"mappingMethod"`
	Status        string              `json:"status"`
	Owner         string              `json:"owner"`
	DisplayName   string              `json:"displayName"`
	Description   string              `json:"description"`

	// ClientID is the application's ID.
	ClientID string `json:"clientID"`
	// ClientSecret is the application's secret.
	ClientSecret string `json:"-"`
	// RedirectURL is the URL to redirect users going through
	// the OAuth flow, after the resource owner's URLs.
	RedirectURL string `json:"redirectURL"`
	// Used to turn off TLS certificate checks
	InsecureSkipVerify bool `json:"insecureSkipVerify"`
	// Scope specifies optional requested permissions.
	Scopes []string `json:"scopes"`
}
