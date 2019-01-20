package rbac_test

import (
	"testing"

	"github.com/johncoleman83/cerebrum/pkg/utl/models"

	"github.com/johncoleman83/cerebrum/pkg/utl/mock"
	"github.com/johncoleman83/cerebrum/pkg/utl/rbac"

	"github.com/labstack/echo"

	"github.com/stretchr/testify/assert"
)

func TestUser(t *testing.T) {
	ctx := mock.EchoCtxWithKeys([]string{
		"id", "company_id", "location_id", "username", "email", "role"},
		uint(9), uint(15), uint(52), "rocinante", "rocinante@gmail.com", models.SuperAdminRole)
	expectedUser := &models.AuthUser{
		ID:          uint(9),
		Username:    "rocinante",
		CompanyID:   uint(15),
		LocationID:  uint(52),
		Email:       "rocinante@gmail.com",
		AccessLevel: models.SuperAdminRole,
	}
	rbacSvc := rbac.New()
	assert.Equal(t, expectedUser, rbacSvc.User(ctx))
}

func TestEnforceRole(t *testing.T) {
	type args struct {
		ctx  echo.Context
		role models.AccessRole
	}
	cases := []struct {
		name        string
		args        args
		expectedErr bool
	}{
		{
			name:        "Not authorized",
			args:        args{ctx: mock.EchoCtxWithKeys([]string{"role"}, models.CompanyAdminRole), role: models.SuperAdminRole},
			expectedErr: true,
		},
		{
			name:        "Authorized",
			args:        args{ctx: mock.EchoCtxWithKeys([]string{"role"}, models.SuperAdminRole), role: models.CompanyAdminRole},
			expectedErr: false,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			rbacSvc := rbac.New()
			res := rbacSvc.EnforceRole(tt.args.ctx, tt.args.role)
			assert.Equal(t, tt.expectedErr, res == echo.ErrForbidden)
		})
	}
}

func TestEnforceUser(t *testing.T) {
	type args struct {
		ctx echo.Context
		id  uint
	}
	cases := []struct {
		name        string
		args        args
		expectedErr bool
	}{
		{
			name:        "Not same user, not an admin",
			args:        args{ctx: mock.EchoCtxWithKeys([]string{"id", "role"}, uint(15), models.LocationAdminRole), id: uint(122)},
			expectedErr: true,
		},
		{
			name:        "Not same user, but admin",
			args:        args{ctx: mock.EchoCtxWithKeys([]string{"id", "role"}, uint(22), models.SuperAdminRole), id: uint(44)},
			expectedErr: false,
		},
		{
			name:        "Same user",
			args:        args{ctx: mock.EchoCtxWithKeys([]string{"id", "role"}, uint(8), models.AdminRole), id: uint(8)},
			expectedErr: false,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			rbacSvc := rbac.New()
			res := rbacSvc.EnforceUser(tt.args.ctx, tt.args.id)
			assert.Equal(t, tt.expectedErr, res == echo.ErrForbidden)
		})
	}
}

func TestEnforceCompany(t *testing.T) {
	type args struct {
		ctx echo.Context
		id  uint
	}
	cases := []struct {
		name        string
		args        args
		expectedErr bool
	}{
		{
			name:        "Not same company, not an admin",
			args:        args{ctx: mock.EchoCtxWithKeys([]string{"company_id", "role"}, uint(7), models.UserRole), id: uint(9)},
			expectedErr: true,
		},
		{
			name:        "Same company, not company admin or admin",
			args:        args{ctx: mock.EchoCtxWithKeys([]string{"company_id", "role"}, uint(22), models.UserRole), id: uint(22)},
			expectedErr: true,
		},
		{
			name:        "Same company, company admin",
			args:        args{ctx: mock.EchoCtxWithKeys([]string{"company_id", "role"}, uint(5), models.CompanyAdminRole), id: uint(5)},
			expectedErr: false,
		},
		{
			name:        "Not same company but admin",
			args:        args{ctx: mock.EchoCtxWithKeys([]string{"company_id", "role"}, uint(8), models.AdminRole), id: uint(9)},
			expectedErr: false,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			rbacSvc := rbac.New()
			res := rbacSvc.EnforceCompany(tt.args.ctx, tt.args.id)
			assert.Equal(t, tt.expectedErr, res == echo.ErrForbidden)
		})
	}
}

func TestEnforceLocation(t *testing.T) {
	type args struct {
		ctx echo.Context
		id  uint
	}
	cases := []struct {
		name        string
		args        args
		expectedErr bool
	}{
		{
			name:        "Not same location, not an admin",
			args:        args{ctx: mock.EchoCtxWithKeys([]string{"location_id", "role"}, uint(7), models.UserRole), id: uint(9)},
			expectedErr: true,
		},
		{
			name:        "Same location, not company admin or admin",
			args:        args{ctx: mock.EchoCtxWithKeys([]string{"location_id", "role"}, uint(22), models.UserRole), id: uint(22)},
			expectedErr: true,
		},
		{
			name:        "Same location, company admin",
			args:        args{ctx: mock.EchoCtxWithKeys([]string{"location_id", "role"}, uint(5), models.CompanyAdminRole), id: uint(5)},
			expectedErr: false,
		},
		{
			name:        "Location admin",
			args:        args{ctx: mock.EchoCtxWithKeys([]string{"location_id", "role"}, uint(5), models.LocationAdminRole), id: uint(5)},
			expectedErr: false,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			rbacSvc := rbac.New()
			res := rbacSvc.EnforceLocation(tt.args.ctx, tt.args.id)
			assert.Equal(t, tt.expectedErr, res == echo.ErrForbidden)
		})
	}
}

func TestAccountCreate(t *testing.T) {
	type args struct {
		ctx         echo.Context
		roleID      models.AccessRole
		company_id  uint
		location_id uint
	}
	cases := []struct {
		name        string
		args        args
		expectedErr bool
	}{
		{
			name:        "Different location, company, creating user role, not an admin",
			args:        args{ctx: mock.EchoCtxWithKeys([]string{"company_id", "location_id", "role"}, uint(2), uint(3), models.UserRole), roleID: models.AccessRole(500), company_id: uint(7), location_id: uint(8)},
			expectedErr: true,
		},
		{
			name:        "Same location, not company, creating user role, not an admin",
			args:        args{ctx: mock.EchoCtxWithKeys([]string{"company_id", "location_id", "role"}, uint(2), uint(3), models.UserRole), roleID: models.AccessRole(500), company_id: uint(2), location_id: uint(8)},
			expectedErr: true,
		},
		{
			name:        "Different location, company, creating user role, not an admin",
			args:        args{ctx: mock.EchoCtxWithKeys([]string{"company_id", "location_id", "role"}, uint(2), uint(3), models.CompanyAdminRole), roleID: models.AccessRole(400), company_id: uint(2), location_id: uint(4)},
			expectedErr: false,
		},
		{
			name:        "Same location, company, creating user role, not an admin",
			args:        args{ctx: mock.EchoCtxWithKeys([]string{"company_id", "location_id", "role"}, uint(2), uint(3), models.CompanyAdminRole), roleID: models.AccessRole(500), company_id: uint(2), location_id: uint(3)},
			expectedErr: false,
		},
		{
			name:        "Same location, company, creating user role, admin",
			args:        args{ctx: mock.EchoCtxWithKeys([]string{"company_id", "location_id", "role"}, uint(2), uint(3), models.CompanyAdminRole), roleID: models.AccessRole(500), company_id: uint(2), location_id: uint(3)},
			expectedErr: false,
		},
		{
			name:        "Different everything, admin",
			args:        args{ctx: mock.EchoCtxWithKeys([]string{"company_id", "location_id", "role"}, uint(2), uint(3), models.AdminRole), roleID: models.AccessRole(200), company_id: uint(7), location_id: uint(4)},
			expectedErr: false,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			rbacSvc := rbac.New()
			res := rbacSvc.AccountCreate(tt.args.ctx, tt.args.roleID, tt.args.company_id, tt.args.location_id)
			assert.Equal(t, tt.expectedErr, res == echo.ErrForbidden)
		})
	}
}

func TestIsLowerRole(t *testing.T) {
	ctx := mock.EchoCtxWithKeys([]string{"role"}, models.CompanyAdminRole)
	rbacSvc := rbac.New()
	if rbacSvc.IsLowerRole(ctx, models.LocationAdminRole) != nil {
		t.Error("The requested user is higher role than the user requesting it")
	}
	if rbacSvc.IsLowerRole(ctx, models.AdminRole) == nil {
		t.Error("The requested user is lower role than the user requesting it")
	}
}
