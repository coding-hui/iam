package model

import (
	"time"
)

func init() {
	RegisterModel(&User{})
}

type User struct {
	BaseModel
	Name          string    `json:"name"`
	Email         string    `json:"email"`
	Alias         string    `json:"alias,omitempty"`
	Password      string    `json:"password,omitempty"`
	Disabled      bool      `json:"disabled"`
	LastLoginTime time.Time `json:"last_login_time,omitempty"`
	// UserRoles binding the platform level roles
	UserRoles []string `json:"user_roles"`
}

// TableName return custom table name
func (u *User) TableName() string {
	return tableNamePrefix + "user"
}

// ShortTableName return custom table name
func (u *User) ShortTableName() string {
	return "usr"
}

// PrimaryKey return custom primary key
func (u *User) PrimaryKey() string {
	return u.Name
}

// Index return custom index
func (u *User) Index() map[string]interface{} {
	index := make(map[string]interface{})
	if u.Name != "" {
		index["name"] = u.Name
	}
	if u.Email != "" {
		index["email"] = u.Email
	}
	return index
}
