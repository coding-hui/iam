package v1

import "time"

// UserBase is the base info of user
type UserBase struct {
	CreateTime    time.Time `json:"create_time"`
	LastLoginTime time.Time `json:"last_login_time"`
	Name          string    `json:"name"`
	Email         string    `json:"email"`
	Type          string    `json:"type"`
	State         string    `json:"state"`
	Alias         string    `json:"alias,omitempty"`
	Disabled      bool      `json:"disabled"`
}

// ListUserOptions list user options
type ListUserOptions struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Alias string `json:"alias"`
}

// ListUserResponse list user response
type ListUserResponse struct {
	Users []*DetailUserResponse `json:"users"`
	Total int64                 `json:"total"`
}

// CreateUserRequest create user request
type CreateUserRequest struct {
	Name     string   `json:"name" validate:"checkname"`
	Alias    string   `json:"alias,omitempty" validate:"checkalias" optional:"true"`
	Email    string   `json:"email" validate:"checkemail"`
	Password string   `json:"password" validate:"checkpassword"`
	Roles    []string `json:"roles"`
}

// UpdateUserRequest update user request
type UpdateUserRequest struct {
	Alias    string    `json:"alias,omitempty" optional:"true"`
	Password string    `json:"password,omitempty" validate:"checkpassword" optional:"true"`
	Email    string    `json:"email,omitempty" validate:"checkemail" optional:"true"`
	Roles    *[]string `json:"roles"`
}

// DetailUserResponse is the response of user detail
type DetailUserResponse struct {
	UserBase
}
