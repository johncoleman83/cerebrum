// Package rbac Role Based Access Control
package rbac

import (
	"github.com/labstack/echo"

	"github.com/johncoleman83/cerebrum/pkg/utl/models"
)

// New creates new RBAC service
func New() *Service {
	return &Service{}
}

// Service is RBAC application service
type Service struct{}

func checkBool(b bool) error {
	if b {
		return nil
	}
	return echo.ErrForbidden
}

// User returns user data stored in jwt token
func (s *Service) User(c echo.Context) *models.AuthUser {
	id := c.Get("id").(uint)
	accountID := c.Get("account_id").(uint)
	teamID := c.Get("team_id").(uint)
	user := c.Get("username").(string)
	email := c.Get("email").(string)
	role := c.Get("role").(models.AccessRole)
	return &models.AuthUser{
		ID:          id,
		Username:    user,
		AccountID:   accountID,
		TeamID:      teamID,
		Email:       email,
		AccessLevel: role,
	}
}

// EnforceRole authorizes request by AccessRole
func (s *Service) EnforceRole(c echo.Context, r models.AccessRole) error {
	return checkBool(!(c.Get("role").(models.AccessRole) > r))
}

// EnforceUser checks whether the request to change user data is done by the same user
func (s *Service) EnforceUser(c echo.Context, ID uint) error {
	// TODO: Implement querying db and checking the requested user's account_id/team_id
	// to allow account/team admins to view the user
	if s.isAdmin(c) {
		return nil
	}
	return checkBool(c.Get("id").(uint) == ID)
}

// EnforceAccount checks whether the request to apply change to account data
// is done by the user belonging to the that account and that the user has role AccountAdmin.
// If user has admin role, the check for account doesnt need to pass.
func (s *Service) EnforceAccount(c echo.Context, ID uint) error {
	if s.isAdmin(c) {
		return nil
	}
	if err := s.EnforceRole(c, models.AccountAdminRole); err != nil {
		return err
	}
	return checkBool(c.Get("account_id").(uint) == ID)
}

// EnforceTeam checks whether the request to change team data
// is done by the user belonging to the requested team
func (s *Service) EnforceTeam(c echo.Context, ID uint) error {
	if s.isAccountAdmin(c) {
		return nil
	}
	if err := s.EnforceRole(c, models.TeamAdminRole); err != nil {
		return err
	}
	return checkBool(c.Get("team_id").(uint) == ID)
}

func (s *Service) isAdmin(c echo.Context) bool {
	return !(c.Get("role").(models.AccessRole) > models.AdminRole)
}

func (s *Service) isAccountAdmin(c echo.Context) bool {
	// Must query account ID in database for the given user
	return !(c.Get("role").(models.AccessRole) > models.AccountAdminRole)
}

// AccountCreate performs auth check when creating a new account
// Team admin cannot create accounts, needs to be fixed on EnforceTeam function
func (s *Service) AccountCreate(c echo.Context, roleID models.AccessRole, accountID, teamID uint) error {
	if err := s.EnforceTeam(c, teamID); err != nil {
		return err
	}
	return s.IsLowerRole(c, roleID)
}

// IsLowerRole checks whether the requesting user has higher role than the user it wants to change
// Used for account creation/deletion
func (s *Service) IsLowerRole(c echo.Context, r models.AccessRole) error {
	return checkBool(c.Get("role").(models.AccessRole) < r)
}
