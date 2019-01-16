package user_test

import (
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"

	"github.com/johncoleman83/cerebrum/pkg/api/user"
	"github.com/johncoleman83/cerebrum/pkg/utl/mock"
	"github.com/johncoleman83/cerebrum/pkg/utl/mock/mockstore"
	cerebrum "github.com/johncoleman83/cerebrum/pkg/utl/model"
)

func TestCreate(t *testing.T) {
	type args struct {
		c   echo.Context
		req cerebrum.User
	}
	cases := []struct {
		name         string
		args         args
		expectedErr  bool
		expectedData *cerebrum.User
		udb          *mockstore.User
		rbac         *mock.RBAC
		sec          *mock.Secure
	}{
		{
			name: "Fail on is lower role",
			args: args{req: cerebrum.User{
				FirstName: "Braxton",
				LastName:  "Young",
				Username:  "BraxtonYoung",
				RoleID:    cerebrum.AccessRole(100),
				Password:  "Thranduil8822",
				Email:     "byoung@gmail.com",
			}},
			rbac: &mock.RBAC{
				AccountCreateFn: func(echo.Context, cerebrum.AccessRole, uint, uint) error {
					return cerebrum.ErrGeneric
				},
			},
			sec: &mock.Secure{
				HashFn: func(string) string {
					return "h4$h3d"
				},
				PasswordFn: func(string, ...string) bool {
					return true
				},
			},
			expectedErr: true,
		},
		{
			name: "Fail on is invalid password",
			args: args{req: cerebrum.User{
				FirstName: "Tina",
				LastName:  "Turner",
				Username:  "TinaTurner",
				RoleID:    cerebrum.AccessRole(200),
				Password:  "TinaTurnerMakesItRain",
				Email:     "tinaturner@gmail.com",
			}},
			rbac: &mock.RBAC{
				AccountCreateFn: func(echo.Context, cerebrum.AccessRole, uint, uint) error {
					return nil
				},
			},
			sec: &mock.Secure{
				HashFn: func(string) string {
					return "h4$h3d"
				},
				PasswordFn: func(string, ...string) bool {
					return false
				},
			},
			expectedErr: true,
		},
		{
			name: "Success",
			args: args{req: cerebrum.User{
				FirstName: "Oprah",
				LastName:  "Winfrey",
				Username:  "OprahWinfrey",
				RoleID:    cerebrum.AccessRole(100),
				Password:  "Thranduil8822",
				Email:     "owinfrey@gmail.com",
			}},
			udb: &mockstore.User{
				CreateFn: func(db *gorm.DB, u cerebrum.User) (*cerebrum.User, error) {
					u.CreatedAt = mock.TestTime(2000)
					u.UpdatedAt = mock.TestTime(2000)
					u.Base.ID = 1
					return &u, nil
				},
			},
			rbac: &mock.RBAC{
				AccountCreateFn: func(echo.Context, cerebrum.AccessRole, uint, uint) error {
					return nil
				}},
			sec: &mock.Secure{
				HashFn: func(string) string {
					return "h4$h3d"
				},
				PasswordFn: func(string, ...string) bool {
					return true
				},
			},
			expectedData: &cerebrum.User{
				Base: cerebrum.Base{
					Model: gorm.Model{
						ID:        1,
						CreatedAt: mock.TestTime(2000),
						UpdatedAt: mock.TestTime(2000),
					},
				},
				FirstName: "Oprah",
				LastName:  "Winfrey",
				Username:  "OprahWinfrey",
				RoleID:    cerebrum.AccessRole(100),
				Password:  "h4$h3d",
				Email:     "owinfrey@gmail.com",
			}}}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			s := user.New(nil, tt.udb, tt.rbac, tt.sec)
			usr, err := s.Create(tt.args.c, tt.args.req)
			assert.Equal(t, tt.expectedErr, err != nil)
			assert.Equal(t, tt.expectedData, usr)
		})
	}
}

func TestView(t *testing.T) {
	type args struct {
		c  echo.Context
		id uint
	}
	cases := []struct {
		name         string
		args         args
		expectedData *cerebrum.User
		expectedErr  error
		udb          *mockstore.User
		rbac         *mock.RBAC
	}{
		{
			name: "Fail on RBAC",
			args: args{id: 5},
			rbac: &mock.RBAC{
				EnforceUserFn: func(c echo.Context, id uint) error {
					return cerebrum.ErrGeneric
				}},
			expectedErr: cerebrum.ErrGeneric,
		},
		{
			name: "Success",
			args: args{id: 1},
			expectedData: &cerebrum.User{
				Base: cerebrum.Base{
					Model: gorm.Model{
						ID:        1,
						CreatedAt: mock.TestTime(2000),
						UpdatedAt: mock.TestTime(2000),
					},
				},
				FirstName: "William",
				LastName:  "Faukner",
				Username:  "WilliamFaukner",
			},
			rbac: &mock.RBAC{
				EnforceUserFn: func(c echo.Context, id uint) error {
					return nil
				}},
			udb: &mockstore.User{
				ViewFn: func(db *gorm.DB, id uint) (*cerebrum.User, error) {
					if id == 1 {
						return &cerebrum.User{
							Base: cerebrum.Base{
								Model: gorm.Model{
									ID:        1,
									CreatedAt: mock.TestTime(2000),
									UpdatedAt: mock.TestTime(2000),
								},
							},
							FirstName: "William",
							LastName:  "Faukner",
							Username:  "WilliamFaukner",
						}, nil
					}
					return nil, nil
				}},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			s := user.New(nil, tt.udb, tt.rbac, nil)
			usr, err := s.View(tt.args.c, tt.args.id)
			assert.Equal(t, tt.expectedData, usr)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestList(t *testing.T) {
	type args struct {
		c   echo.Context
		pgn *cerebrum.Pagination
	}
	cases := []struct {
		name         string
		args         args
		expectedData []cerebrum.User
		expectedErr  bool
		udb          *mockstore.User
		rbac         *mock.RBAC
	}{
		{
			name: "Fail on query List",
			args: args{c: nil, pgn: &cerebrum.Pagination{
				Limit:  100,
				Offset: 200,
			}},
			expectedErr: true,
			rbac: &mock.RBAC{
				UserFn: func(c echo.Context) *cerebrum.AuthUser {
					return &cerebrum.AuthUser{
						ID:         1,
						CompanyID:  2,
						LocationID: 3,
						Role:       cerebrum.UserRole,
					}
				}}},
		{
			name: "Success",
			args: args{c: nil, pgn: &cerebrum.Pagination{
				Limit:  100,
				Offset: 200,
			}},
			rbac: &mock.RBAC{
				UserFn: func(c echo.Context) *cerebrum.AuthUser {
					return &cerebrum.AuthUser{
						ID:         1,
						CompanyID:  2,
						LocationID: 3,
						Role:       cerebrum.AdminRole,
					}
				}},
			udb: &mockstore.User{
				ListFn: func(*gorm.DB, *cerebrum.ListQuery, *cerebrum.Pagination) ([]cerebrum.User, error) {
					return []cerebrum.User{
						{
							Base: cerebrum.Base{
								Model: gorm.Model{
									ID:        1,
									CreatedAt: mock.TestTime(1999),
									UpdatedAt: mock.TestTime(2000),
								},
							},
							FirstName: "Samantha",
							LastName:  "Mills",
							Email:     "Samanthamills@gmail.com",
							Username:  "Samanthamills",
						},
						{
							Base: cerebrum.Base{
								Model: gorm.Model{
									ID:        2,
									CreatedAt: mock.TestTime(2001),
									UpdatedAt: mock.TestTime(2002),
								},
							},
							FirstName: "Preston",
							LastName:  "Phelps",
							Email:     "Prestonphelps@aol.com",
							Username:  "Prestonphelps",
						},
					}, nil
				}},
			expectedData: []cerebrum.User{
				{
					Base: cerebrum.Base{
						Model: gorm.Model{
							ID:        1,
							CreatedAt: mock.TestTime(1999),
							UpdatedAt: mock.TestTime(2000),
						},
					},
					FirstName: "Samantha",
					LastName:  "Mills",
					Email:     "Samanthamills@gmail.com",
					Username:  "Samanthamills",
				},
				{
					Base: cerebrum.Base{
						Model: gorm.Model{
							ID:        2,
							CreatedAt: mock.TestTime(2001),
							UpdatedAt: mock.TestTime(2002),
						},
					},
					FirstName: "Preston",
					LastName:  "Phelps",
					Email:     "Prestonphelps@aol.com",
					Username:  "Prestonphelps",
				},
			},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			s := user.New(nil, tt.udb, tt.rbac, nil)
			usrs, err := s.List(tt.args.c, tt.args.pgn)
			assert.Equal(t, tt.expectedData, usrs)
			assert.Equal(t, tt.expectedErr, err != nil)
		})
	}

}

func TestDelete(t *testing.T) {
	type args struct {
		c  echo.Context
		id uint
	}
	cases := []struct {
		name        string
		args        args
		expectedErr error
		udb         *mockstore.User
		rbac        *mock.RBAC
	}{
		{
			name:        "Fail on ViewUser",
			args:        args{id: 1},
			expectedErr: cerebrum.ErrGeneric,
			udb: &mockstore.User{
				ViewFn: func(db *gorm.DB, id uint) (*cerebrum.User, error) {
					if id != 1 {
						return nil, nil
					}
					return nil, cerebrum.ErrGeneric
				},
			},
		},
		{
			name: "Fail on RBAC",
			args: args{id: 1},
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
						FirstName: "Abigail",
						LastName:  "Gunnings",
						Role: cerebrum.Role{
							AccessLevel: cerebrum.UserRole,
						},
					}, nil
				},
			},
			rbac: &mock.RBAC{
				IsLowerRoleFn: func(echo.Context, cerebrum.AccessRole) error {
					return cerebrum.ErrGeneric
				}},
			expectedErr: cerebrum.ErrGeneric,
		},
		{
			name: "Success",
			args: args{id: 1},
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
						FirstName: "Yignacio",
						LastName:  "Valley",
						Role: cerebrum.Role{
							AccessLevel: cerebrum.AdminRole,
							ID:          cerebrum.AdminRole,
							Name:        "Admin",
						},
					}, nil
				},
				DeleteFn: func(db *gorm.DB, usr *cerebrum.User) error {
					return nil
				},
			},
			rbac: &mock.RBAC{
				IsLowerRoleFn: func(echo.Context, cerebrum.AccessRole) error {
					return nil
				}},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			s := user.New(nil, tt.udb, tt.rbac, nil)
			err := s.Delete(tt.args.c, tt.args.id)
			if err != tt.expectedErr {
				t.Errorf("Expected error %v, received %v", tt.expectedErr, err)
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	type args struct {
		c   echo.Context
		upd *user.Update
	}
	cases := []struct {
		name         string
		args         args
		expectedData *cerebrum.User
		expectedErr  error
		udb          *mockstore.User
		rbac         *mock.RBAC
	}{
		{
			name: "Fail on RBAC",
			args: args{upd: &user.Update{
				ID: 1,
			}},
			rbac: &mock.RBAC{
				EnforceUserFn: func(c echo.Context, id uint) error {
					return cerebrum.ErrGeneric
				}},
			expectedErr: cerebrum.ErrGeneric,
		},
		{
			name: "Fail on ViewUser",
			args: args{upd: &user.Update{
				ID: 1,
			}},
			rbac: &mock.RBAC{
				EnforceUserFn: func(c echo.Context, id uint) error {
					return nil
				}},
			expectedErr: cerebrum.ErrGeneric,
			udb: &mockstore.User{
				ViewFn: func(db *gorm.DB, id uint) (*cerebrum.User, error) {
					if id != 1 {
						return nil, nil
					}
					return nil, cerebrum.ErrGeneric
				},
			},
		},
		{
			name: "Fail on Update",
			args: args{upd: &user.Update{
				ID: 1,
			}},
			rbac: &mock.RBAC{
				EnforceUserFn: func(c echo.Context, id uint) error {
					return nil
				}},
			expectedErr: cerebrum.ErrGeneric,
			udb: &mockstore.User{
				ViewFn: func(db *gorm.DB, id uint) (*cerebrum.User, error) {
					return &cerebrum.User{
						Base: cerebrum.Base{
							Model: gorm.Model{
								ID:        1,
								CreatedAt: mock.TestTime(1990),
								UpdatedAt: mock.TestTime(1991),
							},
						},
						CompanyID:  1,
						LocationID: 2,
						RoleID:     cerebrum.AccessRole(200),
						FirstName:  "Joanna",
						LastName:   "Dimsley",
						Mobile:     "334455",
						Phone:      "444555",
						Address:    "Work Address",
						Email:      "golang@go.org",
					}, nil
				},
				UpdateFn: func(db *gorm.DB, usr *cerebrum.User) error {
					return cerebrum.ErrGeneric
				},
			},
		},
		{
			name: "Success",
			args: args{upd: &user.Update{
				ID:        1,
				FirstName: mock.Str2Ptr("Bethany"),
				LastName:  mock.Str2Ptr("Christian"),
				Mobile:    mock.Str2Ptr("123456"),
				Phone:     mock.Str2Ptr("234567"),
			}},
			rbac: &mock.RBAC{
				EnforceUserFn: func(c echo.Context, id uint) error {
					return nil
				}},
			expectedData: &cerebrum.User{
				Base: cerebrum.Base{
					Model: gorm.Model{
						ID:        1,
						CreatedAt: mock.TestTime(1990),
						UpdatedAt: mock.TestTime(2000),
					},
				},
				CompanyID:  1,
				LocationID: 2,
				RoleID:     cerebrum.AccessRole(200),
				FirstName:  "Bethany",
				LastName:   "Christian",
				Mobile:     "123456",
				Phone:      "234567",
				Address:    "Work Address",
				Email:      "golang@go.org",
			},
			udb: &mockstore.User{
				ViewFn: func(db *gorm.DB, id uint) (*cerebrum.User, error) {
					return &cerebrum.User{
						Base: cerebrum.Base{
							Model: gorm.Model{
								ID:        1,
								CreatedAt: mock.TestTime(1990),
								UpdatedAt: mock.TestTime(1991),
							},
						},
						CompanyID:  1,
						LocationID: 2,
						RoleID:     cerebrum.AccessRole(200),
						FirstName:  "AnewName",
						LastName:   "ALastName",
						Mobile:     "334455",
						Phone:      "444555",
						Address:    "Work Address",
						Email:      "golang@go.org",
					}, nil
				},
				UpdateFn: func(db *gorm.DB, usr *cerebrum.User) error {
					usr.UpdatedAt = mock.TestTime(2000)
					return nil
				},
			},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			s := user.New(nil, tt.udb, tt.rbac, nil)
			usr, err := s.Update(tt.args.c, tt.args.upd)
			assert.Equal(t, tt.expectedData, usr)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestInitialize(t *testing.T) {
	u := user.Initialize(nil, nil, nil)
	if u == nil {
		t.Error("User service not initialized")
	}
}
