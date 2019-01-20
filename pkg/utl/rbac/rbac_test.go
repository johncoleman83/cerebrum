package rbac_test

import (
	"testing"

	cerebrum "github.com/johncoleman83/cerebrum/pkg/utl/model"

	"github.com/johncoleman83/cerebrum/pkg/utl/mock"
	"github.com/johncoleman83/cerebrum/pkg/utl/rbac"

	"github.com/labstack/echo"

	"github.com/stretchr/testify/assert"
)

func TestUser(t *testing.T) {
	ctx := mock.EchoCtxWithKeys([]string{
		"id", "company_id", "location_id", "username", "email", "role"},
		uint(9), uint(15), uint(52), "rocinante", "rocinante@gmail.com", cerebrum.SuperAdminRole)
	expectedUser := &cerebrum.AuthUser{
		ID:          uint(9),
		Username:    "rocinante",
		CompanyID:   uint(15),
		LocationID:  uint(52),
		Email:       "rocinante@gmail.com",
		AccessLevel: cerebrum.SuperAdminRole,
	}
	rbacSvc := rbac.New()
	assert.Equal(t, expectedUser, rbacSvc.User(ctx))
}

func TestEnforceRole(t *testing.T) {
	type args struct {
		ctx  echo.Context
		role cerebrum.AccessRole
	}
	cases := []struct {
		name        string
		args        args
		expectedErr bool
	}{
		{
			name:        "Not authorized",
			args:        args{ctx: mock.EchoCtxWithKeys([]string{"role"}, cerebrum.CompanyAdminRole), role: cerebrum.SuperAdminRole},
			expectedErr: true,
		},
		{
			name:        "Authorized",
			args:        args{ctx: mock.EchoCtxWithKeys([]string{"role"}, cerebrum.SuperAdminRole), role: cerebrum.CompanyAdminRole},
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
			args:        args{ctx: mock.EchoCtxWithKeys([]string{"id", "role"}, uint(15), cerebrum.LocationAdminRole), id: uint(122)},
			expectedErr: true,
		},
		{
			name:        "Not same user, but admin",
			args:        args{ctx: mock.EchoCtxWithKeys([]string{"id", "role"}, uint(22), cerebrum.SuperAdminRole), id: uint(44)},
			expectedErr: false,
		},
		{
			name:        "Same user",
			args:        args{ctx: mock.EchoCtxWithKeys([]string{"id", "role"}, uint(8), cerebrum.AdminRole), id: uint(8)},
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
			args:        args{ctx: mock.EchoCtxWithKeys([]string{"company_id", "role"}, uint(7), cerebrum.UserRole), id: uint(9)},
			expectedErr: true,
		},
		{
			name:        "Same company, not company admin or admin",
			args:        args{ctx: mock.EchoCtxWithKeys([]string{"company_id", "role"}, uint(22), cerebrum.UserRole), id: uint(22)},
			expectedErr: true,
		},
		{
			name:        "Same company, company admin",
			args:        args{ctx: mock.EchoCtxWithKeys([]string{"company_id", "role"}, uint(5), cerebrum.CompanyAdminRole), id: uint(5)},
			expectedErr: false,
		},
		{
			name:        "Not same company but admin",
			args:        args{ctx: mock.EchoCtxWithKeys([]string{"company_id", "role"}, uint(8), cerebrum.AdminRole), id: uint(9)},
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
			args:        args{ctx: mock.EchoCtxWithKeys([]string{"location_id", "role"}, uint(7), cerebrum.UserRole), id: uint(9)},
			expectedErr: true,
		},
		{
			name:        "Same location, not company admin or admin",
			args:        args{ctx: mock.EchoCtxWithKeys([]string{"location_id", "role"}, uint(22), cerebrum.UserRole), id: uint(22)},
			expectedErr: true,
		},
		{
			name:        "Same location, company admin",
			args:        args{ctx: mock.EchoCtxWithKeys([]string{"location_id", "role"}, uint(5), cerebrum.CompanyAdminRole), id: uint(5)},
			expectedErr: false,
		},
		{
			name:        "Location admin",
			args:        args{ctx: mock.EchoCtxWithKeys([]string{"location_id", "role"}, uint(5), cerebrum.LocationAdminRole), id: uint(5)},
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
		roleID      cerebrum.AccessRole
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
			args:        args{ctx: mock.EchoCtxWithKeys([]string{"company_id", "location_id", "role"}, uint(2), uint(3), cerebrum.UserRole), roleID: cerebrum.AccessRole(500), company_id: uint(7), location_id: uint(8)},
			expectedErr: true,
		},
		{
			name:        "Same location, not company, creating user role, not an admin",
			args:        args{ctx: mock.EchoCtxWithKeys([]string{"company_id", "location_id", "role"}, uint(2), uint(3), cerebrum.UserRole), roleID: cerebrum.AccessRole(500), company_id: uint(2), location_id: uint(8)},
			expectedErr: true,
		},
		{
			name:        "Different location, company, creating user role, not an admin",
			args:        args{ctx: mock.EchoCtxWithKeys([]string{"company_id", "location_id", "role"}, uint(2), uint(3), cerebrum.CompanyAdminRole), roleID: cerebrum.AccessRole(400), company_id: uint(2), location_id: uint(4)},
			expectedErr: false,
		},
		{
			name:        "Same location, company, creating user role, not an admin",
			args:        args{ctx: mock.EchoCtxWithKeys([]string{"company_id", "location_id", "role"}, uint(2), uint(3), cerebrum.CompanyAdminRole), roleID: cerebrum.AccessRole(500), company_id: uint(2), location_id: uint(3)},
			expectedErr: false,
		},
		{
			name:        "Same location, company, creating user role, admin",
			args:        args{ctx: mock.EchoCtxWithKeys([]string{"company_id", "location_id", "role"}, uint(2), uint(3), cerebrum.CompanyAdminRole), roleID: cerebrum.AccessRole(500), company_id: uint(2), location_id: uint(3)},
			expectedErr: false,
		},
		{
			name:        "Different everything, admin",
			args:        args{ctx: mock.EchoCtxWithKeys([]string{"company_id", "location_id", "role"}, uint(2), uint(3), cerebrum.AdminRole), roleID: cerebrum.AccessRole(200), company_id: uint(7), location_id: uint(4)},
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
	ctx := mock.EchoCtxWithKeys([]string{"role"}, cerebrum.CompanyAdminRole)
	rbacSvc := rbac.New()
	if rbacSvc.IsLowerRole(ctx, cerebrum.LocationAdminRole) != nil {
		t.Error("The requested user is higher role than the user requesting it")
	}
	if rbacSvc.IsLowerRole(ctx, cerebrum.AdminRole) == nil {
		t.Error("The requested user is lower role than the user requesting it")
	}
}
