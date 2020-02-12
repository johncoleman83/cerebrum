package rbac_test

import (
	"testing"

	"github.com/johncoleman83/cerebrum/pkg/utl/models"

	"github.com/johncoleman83/cerebrum/pkg/utl/mock"
	rbacService "github.com/johncoleman83/cerebrum/pkg/utl/rbac"

	"github.com/labstack/echo"

	"github.com/stretchr/testify/assert"
)

func TestUser(t *testing.T) {
	ctx := mock.EchoCtxWithKeys([]string{
		"id", "account_id", "team_id", "username", "email", "role"},
		uint(9), uint(15), uint(52), "rocinante", "rocinante@gmail.com", models.SuperAdminRole)
	expectedUser := &models.AuthUser{
		ID:          uint(9),
		Username:    "rocinante",
		AccountID:   uint(15),
		TeamID:      uint(52),
		Email:       "rocinante@gmail.com",
		AccessLevel: models.SuperAdminRole,
	}
	rbac := rbacService.New()
	assert.Equal(t, expectedUser, rbac.User(ctx))
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
			args:        args{ctx: mock.EchoCtxWithKeys([]string{"role"}, models.AccountAdminRole), role: models.SuperAdminRole},
			expectedErr: true,
		},
		{
			name:        "Authorized",
			args:        args{ctx: mock.EchoCtxWithKeys([]string{"role"}, models.SuperAdminRole), role: models.AccountAdminRole},
			expectedErr: false,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			rbac := rbacService.New()
			res := rbac.EnforceRole(tt.args.ctx, tt.args.role)
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
			args:        args{ctx: mock.EchoCtxWithKeys([]string{"id", "role"}, uint(15), models.TeamAdminRole), id: uint(122)},
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
			rbac := rbacService.New()
			res := rbac.EnforceUser(tt.args.ctx, tt.args.id)
			assert.Equal(t, tt.expectedErr, res == echo.ErrForbidden)
		})
	}
}

func TestEnforceAccount(t *testing.T) {
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
			name:        "Not same account, not an admin",
			args:        args{ctx: mock.EchoCtxWithKeys([]string{"account_id", "role"}, uint(7), models.UserRole), id: uint(9)},
			expectedErr: true,
		},
		{
			name:        "Same account, not account admin or admin",
			args:        args{ctx: mock.EchoCtxWithKeys([]string{"account_id", "role"}, uint(22), models.UserRole), id: uint(22)},
			expectedErr: true,
		},
		{
			name:        "Same account, account admin",
			args:        args{ctx: mock.EchoCtxWithKeys([]string{"account_id", "role"}, uint(5), models.AccountAdminRole), id: uint(5)},
			expectedErr: false,
		},
		{
			name:        "Not same account but admin",
			args:        args{ctx: mock.EchoCtxWithKeys([]string{"account_id", "role"}, uint(8), models.AdminRole), id: uint(9)},
			expectedErr: false,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			rbac := rbacService.New()
			res := rbac.EnforceAccount(tt.args.ctx, tt.args.id)
			assert.Equal(t, tt.expectedErr, res == echo.ErrForbidden)
		})
	}
}

func TestEnforceTeam(t *testing.T) {
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
			name:        "Not same team, not an admin",
			args:        args{ctx: mock.EchoCtxWithKeys([]string{"team_id", "role"}, uint(7), models.UserRole), id: uint(9)},
			expectedErr: true,
		},
		{
			name:        "Same team, not account admin or admin",
			args:        args{ctx: mock.EchoCtxWithKeys([]string{"team_id", "role"}, uint(22), models.UserRole), id: uint(22)},
			expectedErr: true,
		},
		{
			name:        "Same team, account admin",
			args:        args{ctx: mock.EchoCtxWithKeys([]string{"team_id", "role"}, uint(5), models.AccountAdminRole), id: uint(5)},
			expectedErr: false,
		},
		{
			name:        "Team admin",
			args:        args{ctx: mock.EchoCtxWithKeys([]string{"team_id", "role"}, uint(5), models.TeamAdminRole), id: uint(5)},
			expectedErr: false,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			rbac := rbacService.New()
			res := rbac.EnforceTeam(tt.args.ctx, tt.args.id)
			assert.Equal(t, tt.expectedErr, res == echo.ErrForbidden)
		})
	}
}

func TestAccountCreate(t *testing.T) {
	type args struct {
		ctx        echo.Context
		roleID     models.AccessRole
		account_id uint
		team_id    uint
	}
	cases := []struct {
		name        string
		args        args
		expectedErr bool
	}{
		{
			name:        "Different team, account, creating user role, not an admin",
			args:        args{ctx: mock.EchoCtxWithKeys([]string{"account_id", "team_id", "role"}, uint(2), uint(3), models.UserRole), roleID: models.AccessRole(500), account_id: uint(7), team_id: uint(8)},
			expectedErr: true,
		},
		{
			name:        "Same team, not account, creating user role, not an admin",
			args:        args{ctx: mock.EchoCtxWithKeys([]string{"account_id", "team_id", "role"}, uint(2), uint(3), models.UserRole), roleID: models.AccessRole(500), account_id: uint(2), team_id: uint(8)},
			expectedErr: true,
		},
		{
			name:        "Different team, account, creating user role, not an admin",
			args:        args{ctx: mock.EchoCtxWithKeys([]string{"account_id", "team_id", "role"}, uint(2), uint(3), models.AccountAdminRole), roleID: models.AccessRole(400), account_id: uint(2), team_id: uint(4)},
			expectedErr: false,
		},
		{
			name:        "Same team, account, creating user role, not an admin",
			args:        args{ctx: mock.EchoCtxWithKeys([]string{"account_id", "team_id", "role"}, uint(2), uint(3), models.AccountAdminRole), roleID: models.AccessRole(500), account_id: uint(2), team_id: uint(3)},
			expectedErr: false,
		},
		{
			name:        "Same team, account, creating user role, admin",
			args:        args{ctx: mock.EchoCtxWithKeys([]string{"account_id", "team_id", "role"}, uint(2), uint(3), models.AccountAdminRole), roleID: models.AccessRole(500), account_id: uint(2), team_id: uint(3)},
			expectedErr: false,
		},
		{
			name:        "Different everything, admin",
			args:        args{ctx: mock.EchoCtxWithKeys([]string{"account_id", "team_id", "role"}, uint(2), uint(3), models.AdminRole), roleID: models.AccessRole(200), account_id: uint(7), team_id: uint(4)},
			expectedErr: false,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			rbac := rbacService.New()
			res := rbac.AccountCreate(tt.args.ctx, tt.args.roleID, tt.args.account_id, tt.args.team_id)
			assert.Equal(t, tt.expectedErr, res == echo.ErrForbidden)
		})
	}
}

func TestIsLowerRole(t *testing.T) {
	ctx := mock.EchoCtxWithKeys([]string{"role"}, models.AccountAdminRole)
	rbac := rbacService.New()
	if rbac.IsLowerRole(ctx, models.TeamAdminRole) != nil {
		t.Error("The requested user is higher role than the user requesting it")
	}
	if rbac.IsLowerRole(ctx, models.AdminRole) == nil {
		t.Error("The requested user is lower role than the user requesting it")
	}
}
