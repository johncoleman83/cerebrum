package password_test

import (
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"

	"github.com/johncoleman83/cerebrum/pkg/api/password"
	"github.com/johncoleman83/cerebrum/pkg/utl/mock"
	"github.com/johncoleman83/cerebrum/pkg/utl/mock/mockstore"
	"github.com/johncoleman83/cerebrum/pkg/utl/models"
)

func TestChange(t *testing.T) {
	type args struct {
		oldpass string
		newpass string
		id      uint
	}
	cases := []struct {
		name        string
		args        args
		expectedErr bool
		udb         *mockstore.User
		rbac        *mock.RBAC
		sec         *mock.Secure
	}{
		{
			name: "Fail on EnforceUser",
			args: args{id: 1},
			rbac: &mock.RBAC{
				EnforceUserFn: func(c echo.Context, id uint) error {
					return models.ErrGeneric
				}},
			expectedErr: true,
		},
		{
			name:        "Fail on ViewUser",
			args:        args{id: 1},
			expectedErr: true,
			rbac: &mock.RBAC{
				EnforceUserFn: func(c echo.Context, id uint) error {
					return nil
				}},
			udb: &mockstore.User{
				ViewFn: func(db *gorm.DB, id uint) (*models.User, error) {
					if id != 1 {
						return nil, nil
					}
					return nil, models.ErrGeneric
				},
			},
		},
		{
			name: "Fail on PasswordMatch",
			args: args{id: 1, oldpass: "hunter123"},
			rbac: &mock.RBAC{
				EnforceUserFn: func(c echo.Context, id uint) error {
					return nil
				}},
			expectedErr: true,
			udb: &mockstore.User{
				ViewFn: func(db *gorm.DB, id uint) (*models.User, error) {
					return &models.User{
						Password: "HashedPassword",
					}, nil
				},
			},
			sec: &mock.Secure{
				HashMatchesPasswordFn: func(string, string) bool {
					return false
				},
			},
		},
		{
			name: "Fail on InsecurePassword",
			args: args{id: 1, oldpass: "hunter123"},
			rbac: &mock.RBAC{
				EnforceUserFn: func(c echo.Context, id uint) error {
					return nil
				}},
			expectedErr: true,
			udb: &mockstore.User{
				ViewFn: func(db *gorm.DB, id uint) (*models.User, error) {
					return &models.User{
						Password: "HashedPassword",
					}, nil
				},
			},
			sec: &mock.Secure{
				HashMatchesPasswordFn: func(string, string) bool {
					return true
				},
				PasswordFn: func(string, ...string) bool {
					return false
				},
			},
		},
		{
			name: "Success",
			args: args{id: 1, oldpass: "hunter123", newpass: "password"},
			rbac: &mock.RBAC{
				EnforceUserFn: func(c echo.Context, id uint) error {
					return nil
				}},
			udb: &mockstore.User{
				ViewFn: func(db *gorm.DB, id uint) (*models.User, error) {
					return &models.User{
						Password: "$2a$10$udRBroNGBeOYwSWCVzf6Lulg98uAoRCIi4t75VZg84xgw6EJbFNsG",
					}, nil
				},
				UpdateFn: func(*gorm.DB, *models.User) error {
					return nil
				},
			},
			sec: &mock.Secure{
				HashMatchesPasswordFn: func(string, string) bool {
					return true
				},
				PasswordFn: func(string, ...string) bool {
					return true
				},
				HashFn: func(string) string {
					return "hash3d"
				},
			},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			s := password.New(nil, tt.udb, tt.rbac, tt.sec)
			err := s.Change(nil, tt.args.id, tt.args.oldpass, tt.args.newpass)
			assert.Equal(t, tt.expectedErr, err != nil)
			// Check whether password was changed
		})
	}
}

func TestInitialize(t *testing.T) {
	p := password.Initialize(nil, nil, nil)
	if p == nil {
		t.Error("password service not initialized")
	}
}
