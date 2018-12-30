package query_test

import (
	"testing"

	cerebrum "github.com/johncoleman83/cerebrum/pkg/utl/model"

	"github.com/labstack/echo"

	"github.com/johncoleman83/cerebrum/pkg/utl/query"
	"github.com/stretchr/testify/assert"
)

func TestList(t *testing.T) {
	type args struct {
		user *cerebrum.AuthUser
	}
	cases := []struct {
		name         string
		args         args
		expectedData *cerebrum.ListQuery
		expectedErr  error
	}{
		{
			name: "Super admin user",
			args: args{user: &cerebrum.AuthUser{
				Role: cerebrum.SuperAdminRole,
			}},
		},
		{
			name: "Company admin user",
			args: args{user: &cerebrum.AuthUser{
				Role:      cerebrum.CompanyAdminRole,
				CompanyID: 1,
			}},
			expectedData: &cerebrum.ListQuery{
				Query: "company_id = ?",
				ID:    1},
		},
		{
			name: "Location admin user",
			args: args{user: &cerebrum.AuthUser{
				Role:       cerebrum.LocationAdminRole,
				CompanyID:  1,
				LocationID: 2,
			}},
			expectedData: &cerebrum.ListQuery{
				Query: "location_id = ?",
				ID:    2},
		},
		{
			name: "Normal user",
			args: args{user: &cerebrum.AuthUser{
				Role: cerebrum.UserRole,
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
