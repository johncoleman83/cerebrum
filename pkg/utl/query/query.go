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
	case u.AccessLevel == models.CompanyAdminRole:
		return &models.ListQuery{Query: "company_id = ?", ID: u.CompanyID}, nil
	case u.AccessLevel == models.LocationAdminRole:
		return &models.ListQuery{Query: "location_id = ?", ID: u.LocationID}, nil
	default:
		return nil, echo.ErrForbidden
	}
}
