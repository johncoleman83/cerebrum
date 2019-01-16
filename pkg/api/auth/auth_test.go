package auth_test

import (
	"testing"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"

	"github.com/johncoleman83/cerebrum/pkg/api/auth"
	"github.com/johncoleman83/cerebrum/pkg/utl/mock"
	"github.com/johncoleman83/cerebrum/pkg/utl/mock/mockstore"
	cerebrum "github.com/johncoleman83/cerebrum/pkg/utl/model"
)

func TestAuthenticate(t *testing.T) {
	type args struct {
		user string
		pass string
	}
	cases := []struct {
		name         string
		args         args
		expectedData *cerebrum.AuthToken
		expectedErr  bool
		udb          *mockstore.User
		jwt          *mock.JWT
		sec          *mock.Secure
	}{
		{
			name:        "Fail on finding user",
			args:        args{user: "juzernejm"},
			expectedErr: true,
			udb: &mockstore.User{
				FindByUsernameFn: func(db *gorm.DB, user string) (*cerebrum.User, error) {
					return nil, cerebrum.ErrGeneric
				},
			},
		},
		{
			name:        "Fail on wrong password",
			args:        args{user: "juzernejm", pass: "notHashedPassword"},
			expectedErr: true,
			udb: &mockstore.User{
				FindByUsernameFn: func(db *gorm.DB, user string) (*cerebrum.User, error) {
					return &cerebrum.User{
						Username: user,
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
			name:        "Fail on token generation",
			args:        args{user: "juzernejm", pass: "pass"},
			expectedErr: true,
			udb: &mockstore.User{
				FindByUsernameFn: func(db *gorm.DB, user string) (*cerebrum.User, error) {
					return &cerebrum.User{
						Username: user,
						Password: "pass",
					}, nil
				},
			},
			sec: &mock.Secure{
				HashMatchesPasswordFn: func(string, string) bool {
					return true
				},
			},
			jwt: &mock.JWT{
				GenerateTokenFn: func(u *cerebrum.User) (string, string, error) {
					return "", "", cerebrum.ErrGeneric
				},
			},
		},
		{
			name:        "Fail on updating last login",
			args:        args{user: "juzernejm", pass: "pass"},
			expectedErr: true,
			udb: &mockstore.User{
				FindByUsernameFn: func(db *gorm.DB, user string) (*cerebrum.User, error) {
					return &cerebrum.User{
						Username: user,
						Password: "pass",
					}, nil
				},
				UpdateFn: func(db *gorm.DB, u *cerebrum.User) error {
					return cerebrum.ErrGeneric
				},
			},
			sec: &mock.Secure{
				HashMatchesPasswordFn: func(string, string) bool {
					return true
				},
				TokenFn: func(string) string {
					return "refreshtoken"
				},
			},
			jwt: &mock.JWT{
				GenerateTokenFn: func(u *cerebrum.User) (string, string, error) {
					return "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9", mock.TestTime(2000).Format(time.RFC3339), nil
				},
			},
		},
		{
			name: "Success",
			args: args{user: "juzernejm", pass: "pass"},
			udb: &mockstore.User{
				FindByUsernameFn: func(db *gorm.DB, user string) (*cerebrum.User, error) {
					return &cerebrum.User{
						Username: user,
						Password: "password",
					}, nil
				},
				UpdateFn: func(db *gorm.DB, u *cerebrum.User) error {
					return nil
				},
			},
			jwt: &mock.JWT{
				GenerateTokenFn: func(u *cerebrum.User) (string, string, error) {
					return "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9", mock.TestTime(2000).Format(time.RFC3339), nil
				},
			},
			sec: &mock.Secure{
				HashMatchesPasswordFn: func(string, string) bool {
					return true
				},
				TokenFn: func(string) string {
					return "refreshtoken"
				},
			},
			expectedData: &cerebrum.AuthToken{
				Token:        "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
				Expires:      mock.TestTime(2000).Format(time.RFC3339),
				RefreshToken: "refreshtoken",
			},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			s := auth.New(nil, tt.udb, tt.jwt, tt.sec, nil)
			token, err := s.Authenticate(nil, tt.args.user, tt.args.pass)
			if tt.expectedData != nil {
				tt.expectedData.RefreshToken = token.RefreshToken
				assert.Equal(t, tt.expectedData, token)
			}
			assert.Equal(t, tt.expectedErr, err != nil)
		})
	}
}
func TestRefresh(t *testing.T) {
	type args struct {
		c     echo.Context
		token string
	}
	cases := []struct {
		name         string
		args         args
		expectedData *cerebrum.RefreshToken
		expectedErr  bool
		udb          *mockstore.User
		jwt          *mock.JWT
	}{
		{
			name:        "Fail on finding token",
			args:        args{token: "refreshtoken"},
			expectedErr: true,
			udb: &mockstore.User{
				FindByTokenFn: func(db *gorm.DB, token string) (*cerebrum.User, error) {
					return nil, cerebrum.ErrGeneric
				},
			},
		},
		{
			name:        "Fail on token generation",
			args:        args{token: "refreshtoken"},
			expectedErr: true,
			udb: &mockstore.User{
				FindByTokenFn: func(db *gorm.DB, token string) (*cerebrum.User, error) {
					return &cerebrum.User{
						Username: "username",
						Password: "password",
						Token:    token,
					}, nil
				},
			},
			jwt: &mock.JWT{
				GenerateTokenFn: func(u *cerebrum.User) (string, string, error) {
					return "", "", cerebrum.ErrGeneric
				},
			},
		},
		{
			name: "Success",
			args: args{token: "refreshtoken"},
			udb: &mockstore.User{
				FindByTokenFn: func(db *gorm.DB, token string) (*cerebrum.User, error) {
					return &cerebrum.User{
						Username: "username",
						Password: "password",
						Token:    token,
					}, nil
				},
			},
			jwt: &mock.JWT{
				GenerateTokenFn: func(u *cerebrum.User) (string, string, error) {
					return "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9", mock.TestTime(2000).Format(time.RFC3339), nil
				},
			},
			expectedData: &cerebrum.RefreshToken{
				Token:   "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
				Expires: mock.TestTime(2000).Format(time.RFC3339),
			},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			s := auth.New(nil, tt.udb, tt.jwt, nil, nil)
			token, err := s.Refresh(tt.args.c, tt.args.token)
			assert.Equal(t, tt.expectedData, token)
			assert.Equal(t, tt.expectedErr, err != nil)
		})
	}
}

func TestMe(t *testing.T) {
	cases := []struct {
		name         string
		expectedData *cerebrum.User
		udb          *mockstore.User
		rbac         *mock.RBAC
		expectedErr  bool
	}{
		{
			name: "Success",
			rbac: &mock.RBAC{
				UserFn: func(echo.Context) *cerebrum.AuthUser {
					return &cerebrum.AuthUser{ID: 9}
				},
			},
			udb: &mockstore.User{
				ViewFn: func(db *gorm.DB, id uint) (*cerebrum.User, error) {
					return &cerebrum.User{
						Base: cerebrum.Base{
							Model: gorm.Model{
								ID:        id,
								CreatedAt: mock.TestTime(1999),
								UpdatedAt: mock.TestTime(2000),
							},
						},
						FirstName: "Blazing",
						LastName:  "Saddles",
						Role: cerebrum.Role{
							AccessLevel: cerebrum.UserRole,
						},
					}, nil
				},
			},
			expectedData: &cerebrum.User{
				Base: cerebrum.Base{
					Model: gorm.Model{
						ID:        9,
						CreatedAt: mock.TestTime(1999),
						UpdatedAt: mock.TestTime(2000),
					},
				},
				FirstName: "Blazing",
				LastName:  "Saddles",
				Role: cerebrum.Role{
					AccessLevel: cerebrum.UserRole,
				},
			},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			s := auth.New(nil, tt.udb, nil, nil, tt.rbac)
			user, err := s.Me(nil)
			assert.Equal(t, tt.expectedData, user)
			assert.Equal(t, tt.expectedErr, err != nil)
		})
	}
}

func TestInitialize(t *testing.T) {
	a := auth.Initialize(nil, nil, nil, nil)
	if a == nil {
		t.Error("auth service not initialized")
	}
}
