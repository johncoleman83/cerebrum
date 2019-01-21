package user_test

import (
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"

	"github.com/johncoleman83/cerebrum/pkg/api/user"
	"github.com/johncoleman83/cerebrum/pkg/utl/mock"
	"github.com/johncoleman83/cerebrum/pkg/utl/mock/mockstore"
	"github.com/johncoleman83/cerebrum/pkg/utl/models"
)

func TestCreate(t *testing.T) {
	type args struct {
		c   echo.Context
		req models.User
	}
	cases := []struct {
		name         string
		args         args
		expectedErr  bool
		expectedData *models.User
		udb          *mockstore.UserDBClient
		rbac         *mock.RBAC
		sec          *mock.Secure
	}{
		{
			name: "Fail on is lower role",
			args: args{req: models.User{
				FirstName: "Braxton",
				LastName:  "Young",
				Username:  "BraxtonYoung",
				RoleID:    1,
				Password:  "Thranduil8822",
				Email:     "byoung@gmail.com",
			}},
			rbac: &mock.RBAC{
				AccountCreateFn: func(echo.Context, models.AccessRole, uint, uint) error {
					return models.ErrGeneric
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
			args: args{req: models.User{
				FirstName: "Tina",
				LastName:  "Turner",
				Username:  "TinaTurner",
				RoleID:    5,
				Password:  "TinaTurnerMakesItRain",
				Email:     "tinaturner@gmail.com",
			}},
			rbac: &mock.RBAC{
				AccountCreateFn: func(echo.Context, models.AccessRole, uint, uint) error {
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
			args: args{req: models.User{
				FirstName: "Oprah",
				LastName:  "Winfrey",
				Username:  "OprahWinfrey",
				RoleID:    1,
				Password:  "Thranduil8822",
				Email:     "owinfrey@gmail.com",
			}},
			udb: &mockstore.UserDBClient{
				CreateFn: func(db *gorm.DB, u models.User) (*models.User, error) {
					u.CreatedAt = mock.TestTime(2000)
					u.UpdatedAt = mock.TestTime(2000)
					u.Base.ID = 1
					return &u, nil
				},
			},
			rbac: &mock.RBAC{
				AccountCreateFn: func(echo.Context, models.AccessRole, uint, uint) error {
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
			expectedData: &models.User{
				Base: models.Base{
					Model: gorm.Model{
						ID:        1,
						CreatedAt: mock.TestTime(2000),
						UpdatedAt: mock.TestTime(2000),
					},
				},
				FirstName: "Oprah",
				LastName:  "Winfrey",
				Username:  "OprahWinfrey",
				RoleID:    1,
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
		expectedData *models.User
		expectedErr  error
		udb          *mockstore.UserDBClient
		rbac         *mock.RBAC
	}{
		{
			name: "Fail on RBAC",
			args: args{id: 5},
			rbac: &mock.RBAC{
				EnforceUserFn: func(c echo.Context, id uint) error {
					return models.ErrGeneric
				}},
			expectedErr: models.ErrGeneric,
		},
		{
			name: "Success",
			args: args{id: 1},
			expectedData: &models.User{
				Base: models.Base{
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
			udb: &mockstore.UserDBClient{
				ViewFn: func(db *gorm.DB, id uint) (*models.User, error) {
					if id == 1 {
						return &models.User{
							Base: models.Base{
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
		pgn *models.Pagination
	}
	cases := []struct {
		name         string
		args         args
		expectedData []models.User
		expectedErr  bool
		udb          *mockstore.UserDBClient
		rbac         *mock.RBAC
	}{
		{
			name: "Fail on query List",
			args: args{c: nil, pgn: &models.Pagination{
				Limit:  100,
				Offset: 200,
			}},
			expectedErr: true,
			rbac: &mock.RBAC{
				UserFn: func(c echo.Context) *models.AuthUser {
					return &models.AuthUser{
						ID:            1,
						AccountID:     2,
						PrimaryTeamID: 3,
						AccessLevel:   models.UserRole,
					}
				}}},
		{
			name: "Success",
			args: args{c: nil, pgn: &models.Pagination{
				Limit:  100,
				Offset: 200,
			}},
			rbac: &mock.RBAC{
				UserFn: func(c echo.Context) *models.AuthUser {
					return &models.AuthUser{
						ID:            1,
						AccountID:     2,
						PrimaryTeamID: 3,
						AccessLevel:   models.AdminRole,
					}
				}},
			udb: &mockstore.UserDBClient{
				ListFn: func(*gorm.DB, *models.ListQuery, *models.Pagination) ([]models.User, error) {
					return []models.User{
						{
							Base: models.Base{
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
							Base: models.Base{
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
			expectedData: []models.User{
				{
					Base: models.Base{
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
					Base: models.Base{
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
		udb         *mockstore.UserDBClient
		rbac        *mock.RBAC
	}{
		{
			name:        "Fail on ViewUser",
			args:        args{id: 1},
			expectedErr: models.ErrGeneric,
			udb: &mockstore.UserDBClient{
				ViewFn: func(db *gorm.DB, id uint) (*models.User, error) {
					if id != 1 {
						return nil, nil
					}
					return nil, models.ErrGeneric
				},
			},
		},
		{
			name: "Fail on RBAC",
			args: args{id: 1},
			udb: &mockstore.UserDBClient{
				ViewFn: func(db *gorm.DB, id uint) (*models.User, error) {
					return &models.User{
						Base: models.Base{
							Model: gorm.Model{
								ID:        id,
								CreatedAt: mock.TestTime(1999),
								UpdatedAt: mock.TestTime(2000),
							},
						},
						FirstName: "Abigail",
						LastName:  "Gunnings",
						Role: models.Role{
							AccessLevel: models.UserRole,
						},
					}, nil
				},
			},
			rbac: &mock.RBAC{
				IsLowerRoleFn: func(echo.Context, models.AccessRole) error {
					return models.ErrGeneric
				}},
			expectedErr: models.ErrGeneric,
		},
		{
			name: "Success",
			args: args{id: 1},
			udb: &mockstore.UserDBClient{
				ViewFn: func(db *gorm.DB, id uint) (*models.User, error) {
					return &models.User{
						Base: models.Base{
							Model: gorm.Model{
								ID:        id,
								CreatedAt: mock.TestTime(1999),
								UpdatedAt: mock.TestTime(2000),
							},
						},
						FirstName: "Yignacio",
						LastName:  "Valley",
						Role: models.Role{
							AccessLevel: models.AdminRole,
							ID:          2,
							Name:        "Admin",
						},
					}, nil
				},
				DeleteFn: func(db *gorm.DB, usr *models.User) error {
					return nil
				},
			},
			rbac: &mock.RBAC{
				IsLowerRoleFn: func(echo.Context, models.AccessRole) error {
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
		expectedData *models.User
		expectedErr  error
		udb          *mockstore.UserDBClient
		rbac         *mock.RBAC
	}{
		{
			name: "Fail on RBAC",
			args: args{upd: &user.Update{
				ID: 1,
			}},
			rbac: &mock.RBAC{
				EnforceUserFn: func(c echo.Context, id uint) error {
					return models.ErrGeneric
				}},
			expectedErr: models.ErrGeneric,
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
			expectedErr: models.ErrGeneric,
			udb: &mockstore.UserDBClient{
				ViewFn: func(db *gorm.DB, id uint) (*models.User, error) {
					if id != 1 {
						return nil, nil
					}
					return nil, models.ErrGeneric
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
			expectedErr: models.ErrGeneric,
			udb: &mockstore.UserDBClient{
				ViewFn: func(db *gorm.DB, id uint) (*models.User, error) {
					return &models.User{
						Base: models.Base{
							Model: gorm.Model{
								ID:        1,
								CreatedAt: mock.TestTime(1990),
								UpdatedAt: mock.TestTime(1991),
							},
						},
						AccountID:     1,
						PrimaryTeamID: 2,
						RoleID:        5,
						FirstName:     "Joanna",
						LastName:      "Dimsley",
						Mobile:        "334455",
						Phone:         "444555",
						Address:       "Work Address",
						Email:         "golang@go.org",
					}, nil
				},
				UpdateFn: func(db *gorm.DB, usr *models.User) error {
					return models.ErrGeneric
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
			expectedData: &models.User{
				Base: models.Base{
					Model: gorm.Model{
						ID:        1,
						CreatedAt: mock.TestTime(1990),
						UpdatedAt: mock.TestTime(2000),
					},
				},
				AccountID:     1,
				PrimaryTeamID: 2,
				RoleID:        5,
				FirstName:     "Bethany",
				LastName:      "Christian",
				Mobile:        "123456",
				Phone:         "234567",
				Address:       "Work Address",
				Email:         "golang@go.org",
			},
			udb: &mockstore.UserDBClient{
				ViewFn: func(db *gorm.DB, id uint) (*models.User, error) {
					return &models.User{
						Base: models.Base{
							Model: gorm.Model{
								ID:        1,
								CreatedAt: mock.TestTime(1990),
								UpdatedAt: mock.TestTime(1991),
							},
						},
						AccountID:     1,
						PrimaryTeamID: 2,
						RoleID:        5,
						FirstName:     "AnewName",
						LastName:      "ALastName",
						Mobile:        "334455",
						Phone:         "444555",
						Address:       "Work Address",
						Email:         "golang@go.org",
					}, nil
				},
				UpdateFn: func(db *gorm.DB, usr *models.User) error {
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
