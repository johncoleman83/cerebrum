// Package query contains support functions for making db queries
package query

import (
	"github.com/labstack/echo"

	"github.com/johncoleman83/cerebrum/pkg/utl/models"
)

// List prepares data for list queries
func List(u *models.AuthUser) (*models.ListQuery, error) {
	switch true {
	case u.AccessLevel <= models.AdminRole: // user is SuperAdmin or Admin
		return nil, nil
	case u.AccessLevel == models.AccountAdminRole:
		return &models.ListQuery{Query: "account_id = ?", ID: u.AccountID}, nil
	case u.AccessLevel == models.TeamAdminRole:
		return &models.ListQuery{Query: "primary_team_id = ?", ID: u.PrimaryTeamID}, nil
	default:
		return nil, echo.ErrForbidden
	}
}
