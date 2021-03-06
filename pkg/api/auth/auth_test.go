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
	"github.com/johncoleman83/cerebrum/pkg/utl/models"
)

func TestAuthenticate(t *testing.T) {
	type args struct {
		user string
		pass string
	}
	cases := []struct {
		name         string
		args         args
		expectedData *models.AuthToken
		expectedErr  bool
		udb          *mockstore.UserDBClient
		jwt          *mock.JWT
		sec          *mock.Secure
	}{
		{
			name:        "Fail on finding user",
			args:        args{user: "juzernejm"},
			expectedErr: true,
			udb: &mockstore.UserDBClient{
				FindByUsernameFn: func(db *gorm.DB, user string) (*models.User, error) {
					return nil, models.ErrGeneric
				},
			},
		},
		{
			name:        "Fail on wrong password",
			args:        args{user: "juzernejm", pass: "notHashedPassword"},
			expectedErr: true,
			udb: &mockstore.UserDBClient{
				FindByUsernameFn: func(db *gorm.DB, user string) (*models.User, error) {
					return &models.User{
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
			udb: &mockstore.UserDBClient{
				FindByUsernameFn: func(db *gorm.DB, user string) (*models.User, error) {
					return &models.User{
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
				GenerateTokenFn: func(u *models.User) (string, string, error) {
					return "", "", models.ErrGeneric
				},
			},
		},
		{
			name:        "Fail on updating last login",
			args:        args{user: "juzernejm", pass: "pass"},
			expectedErr: true,
			udb: &mockstore.UserDBClient{
				FindByUsernameFn: func(db *gorm.DB, user string) (*models.User, error) {
					return &models.User{
						Username: user,
						Password: "pass",
					}, nil
				},
				UpdateFn: func(db *gorm.DB, u *models.User) error {
					return models.ErrGeneric
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
				GenerateTokenFn: func(u *models.User) (string, string, error) {
					return "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9", mock.TestTime(2000).Format(time.RFC3339), nil
				},
			},
		},
		{
			name: "Success",
			args: args{user: "juzernejm", pass: "pass"},
			udb: &mockstore.UserDBClient{
				FindByUsernameFn: func(db *gorm.DB, user string) (*models.User, error) {
					return &models.User{
						Username: user,
						Password: "password",
					}, nil
				},
				UpdateFn: func(db *gorm.DB, u *models.User) error {
					return nil
				},
			},
			jwt: &mock.JWT{
				GenerateTokenFn: func(u *models.User) (string, string, error) {
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
			expectedData: &models.AuthToken{
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
		expectedData *models.RefreshToken
		expectedErr  bool
		udb          *mockstore.UserDBClient
		jwt          *mock.JWT
	}{
		{
			name:        "Fail on finding token",
			args:        args{token: "refreshtoken"},
			expectedErr: true,
			udb: &mockstore.UserDBClient{
				FindByTokenFn: func(db *gorm.DB, token string) (*models.User, error) {
					return nil, models.ErrGeneric
				},
			},
		},
		{
			name:        "Fail on token generation",
			args:        args{token: "refreshtoken"},
			expectedErr: true,
			udb: &mockstore.UserDBClient{
				FindByTokenFn: func(db *gorm.DB, token string) (*models.User, error) {
					return &models.User{
						Username: "username",
						Password: "password",
						Token:    token,
					}, nil
				},
			},
			jwt: &mock.JWT{
				GenerateTokenFn: func(u *models.User) (string, string, error) {
					return "", "", models.ErrGeneric
				},
			},
		},
		{
			name: "Success",
			args: args{token: "refreshtoken"},
			udb: &mockstore.UserDBClient{
				FindByTokenFn: func(db *gorm.DB, token string) (*models.User, error) {
					return &models.User{
						Username: "username",
						Password: "password",
						Token:    token,
					}, nil
				},
			},
			jwt: &mock.JWT{
				GenerateTokenFn: func(u *models.User) (string, string, error) {
					return "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9", mock.TestTime(2000).Format(time.RFC3339), nil
				},
			},
			expectedData: &models.RefreshToken{
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
		expectedData *models.User
		udb          *mockstore.UserDBClient
		rbac         *mock.RBAC
		expectedErr  bool
	}{
		{
			name: "Success",
			rbac: &mock.RBAC{
				UserFn: func(echo.Context) *models.AuthUser {
					return &models.AuthUser{ID: 9}
				},
			},
			udb: &mockstore.UserDBClient{
				ViewFn: func(db *gorm.DB, id uint) (*models.User, error) {
					return &models.User{
						Base: models.Base{
							ID:        id,
							CreatedAt: mock.TestTime(1999),
							UpdatedAt: mock.TestTime(2000),
						},
						FirstName: "Blazing",
						LastName:  "Saddles",
						Role: models.Role{
							AccessLevel: models.UserRole,
						},
					}, nil
				},
			},
			expectedData: &models.User{
				Base: models.Base{
					ID:        9,
					CreatedAt: mock.TestTime(1999),
					UpdatedAt: mock.TestTime(2000),
				},
				FirstName: "Blazing",
				LastName:  "Saddles",
				Role: models.Role{
					AccessLevel: models.UserRole,
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
