package models

import (
	"github.com/labstack/echo"
)

// AuthToken holds authentication token details with refresh token
type AuthToken struct {
	Token        string `json:"token"`
	Expires      string `json:"expires"`
	RefreshToken string `json:"refresh_token"`
}

// RefreshToken holds authentication token details
type RefreshToken struct {
	Token   string `json:"token"`
	Expires string `json:"expires"`
}

// RBACService represents role-based access control service interface
type RBACService interface {
	User(echo.Context) *AuthUser
	EnforceRole(echo.Context, AccessRole) error
	EnforceUser(echo.Context, uint) error
	EnforceAccount(echo.Context, uint) error
	EnforceTeam(echo.Context, uint) error
	AccountCreate(echo.Context, AccessRole, uint, uint) error
	IsLowerRole(echo.Context, AccessRole) error
}
