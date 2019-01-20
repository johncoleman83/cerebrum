package mock

import (
	"github.com/labstack/echo"

	"github.com/johncoleman83/cerebrum/pkg/utl/models"
)

// RBAC Mock
type RBAC struct {
	UserFn           func(echo.Context) *models.AuthUser
	EnforceRoleFn    func(echo.Context, models.AccessRole) error
	EnforceUserFn    func(echo.Context, uint) error
	EnforceAccountFn func(echo.Context, uint) error
	EnforceTeamFn    func(echo.Context, uint) error
	AccountCreateFn  func(echo.Context, models.AccessRole, uint, uint) error
	IsLowerRoleFn    func(echo.Context, models.AccessRole) error
}

// User mock
func (a *RBAC) User(c echo.Context) *models.AuthUser {
	return a.UserFn(c)
}

// EnforceRole mock
func (a *RBAC) EnforceRole(c echo.Context, role models.AccessRole) error {
	return a.EnforceRoleFn(c, role)
}

// EnforceUser mock
func (a *RBAC) EnforceUser(c echo.Context, id uint) error {
	return a.EnforceUserFn(c, id)
}

// EnforceAccount mock
func (a *RBAC) EnforceAccount(c echo.Context, id uint) error {
	return a.EnforceAccountFn(c, id)
}

// EnforceTeam mock
func (a *RBAC) EnforceTeam(c echo.Context, id uint) error {
	return a.EnforceTeamFn(c, id)
}

// AccountCreate mock
func (a *RBAC) AccountCreate(c echo.Context, roleID models.AccessRole, accountID, teamID uint) error {
	return a.AccountCreateFn(c, roleID, accountID, teamID)
}

// IsLowerRole mock
func (a *RBAC) IsLowerRole(c echo.Context, role models.AccessRole) error {
	return a.IsLowerRoleFn(c, role)
}
