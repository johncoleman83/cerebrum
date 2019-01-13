// Package query contains support functions for making db queries
package query

import (
	"github.com/labstack/echo"

	cerebrum "github.com/johncoleman83/cerebrum/pkg/utl/model"
)

// List prepares data for list queries
func List(u *cerebrum.AuthUser) (*cerebrum.ListQuery, error) {
	switch true {
	case u.Role <= cerebrum.AdminRole: // user is SuperAdmin or Admin
		return nil, nil
	case u.Role == cerebrum.CompanyAdminRole:
		return &cerebrum.ListQuery{Query: "company_id = ?", ID: u.CompanyID}, nil
	case u.Role == cerebrum.LocationAdminRole:
		return &cerebrum.ListQuery{Query: "location_id = ?", ID: u.LocationID}, nil
	default:
		return nil, echo.ErrForbidden
	}
}
