package query_test

import (
	"testing"

	"github.com/johncoleman83/cerebrum/pkg/utl/models"

	"github.com/labstack/echo"

	"github.com/johncoleman83/cerebrum/pkg/utl/query"
	"github.com/stretchr/testify/assert"
)

func TestList(t *testing.T) {
	type args struct {
		user *models.AuthUser
	}
	cases := []struct {
		name         string
		args         args
		expectedData *models.ListQuery
		expectedErr  error
	}{
		{
			name: "Super admin user",
			args: args{user: &models.AuthUser{
				AccessLevel: models.SuperAdminRole,
			}},
		},
		{
			name: "Company admin user",
			args: args{user: &models.AuthUser{
				AccessLevel: models.CompanyAdminRole,
				CompanyID:   1,
			}},
			expectedData: &models.ListQuery{
				Query: "company_id = ?",
				ID:    1},
		},
		{
			name: "Location admin user",
			args: args{user: &models.AuthUser{
				AccessLevel: models.LocationAdminRole,
				CompanyID:   1,
				LocationID:  2,
			}},
			expectedData: &models.ListQuery{
				Query: "location_id = ?",
				ID:    2},
		},
		{
			name: "Normal user",
			args: args{user: &models.AuthUser{
				AccessLevel: models.UserRole,
			}},
			expectedErr: echo.ErrForbidden,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			q, err := query.List(tt.args.user)
			assert.Equal(t, tt.expectedData, q)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}
