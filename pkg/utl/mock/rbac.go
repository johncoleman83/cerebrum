package mock

import (
	"github.com/labstack/echo"

	cerebrum "github.com/johncoleman83/cerebrum/pkg/utl/model"
)

// RBAC Mock
type RBAC struct {
	UserFn            func(echo.Context) *cerebrum.AuthUser
	EnforceRoleFn     func(echo.Context, cerebrum.AccessRole) error
	EnforceUserFn     func(echo.Context, uint) error
	EnforceCompanyFn  func(echo.Context, uint) error
	EnforceLocationFn func(echo.Context, uint) error
	AccountCreateFn   func(echo.Context, cerebrum.AccessRole, uint, uint) error
	IsLowerRoleFn     func(echo.Context, cerebrum.AccessRole) error
}

// User mock
func (a *RBAC) User(c echo.Context) *cerebrum.AuthUser {
	return a.UserFn(c)
}

// EnforceRole mock
func (a *RBAC) EnforceRole(c echo.Context, role cerebrum.AccessRole) error {
	return a.EnforceRoleFn(c, role)
}

// EnforceUser mock
func (a *RBAC) EnforceUser(c echo.Context, id uint) error {
	return a.EnforceUserFn(c, id)
}

// EnforceCompany mock
func (a *RBAC) EnforceCompany(c echo.Context, id uint) error {
	return a.EnforceCompanyFn(c, id)
}

// EnforceLocation mock
func (a *RBAC) EnforceLocation(c echo.Context, id uint) error {
	return a.EnforceLocationFn(c, id)
}

// AccountCreate mock
func (a *RBAC) AccountCreate(c echo.Context, roleID cerebrum.AccessRole, companyID, locationID uint) error {
	return a.AccountCreateFn(c, roleID, companyID, locationID)
}

// IsLowerRole mock
func (a *RBAC) IsLowerRole(c echo.Context, role cerebrum.AccessRole) error {
	return a.IsLowerRoleFn(c, role)
}
