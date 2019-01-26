package transport_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/johncoleman83/cerebrum/pkg/utl/models"

	"github.com/johncoleman83/cerebrum/pkg/api/user"
	"github.com/johncoleman83/cerebrum/pkg/api/user/transport"

	"github.com/johncoleman83/cerebrum/pkg/utl/mock"
	"github.com/johncoleman83/cerebrum/pkg/utl/mock/mockstore"
	"github.com/johncoleman83/cerebrum/pkg/utl/server"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	cases := []struct {
		name           string
		req            string
		expectedStatus int
		expectedResp   *models.User
		udb            *mockstore.UserDBClient
		rbac           *mock.RBAC
		sec            *mock.Secure
	}{
		{
			name:           "Fail on bad params",
			req:            `{"firstname":"Vanessa","lastname":"Harris","username":"vanessaharris","password":"hunter123","password_confirm":"hunter123","email":"vanessaharris@gmail.com","account_id":1,"primary_team_id":2,"role_id":5}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Fail on validation with short username",
			req:            `{"first_name":"Frank","last_name":"Williams","username":"fw","password":"hunter123","password_confirm":"hunter123","email":"frankwilliams@gmail.com","account_id":1,"primary_team_id":2,"role_id":5}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Fail on validation of email",
			req:            `{"first_name":"Princton","last_name":"Thomas","username":"princetonthomas","password":"hunter123","password_confirm":"hunter123","email":"princetonthomas$gmail.com","account_id":1,"primary_team_id":2,"role_id":5}`,
			expectedStatus: http.StatusBadRequest,
		}, {
			name:           "Fail on non-matching passwords",
			req:            `{"first_name":"Blake","last_name":"Fields","username":"blakefields","password":"sampson","password_confirm":"sampson1","email":"blakefields@gmail.com","account_id":1,"primary_team_id":2,"role_id":5}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Fail on invalid role",
			req:  `{"first_name":"William","last_name":"Abbott","username":"williamabbot","password":"hunter123","password_confirm":"hunter123","email":"williamabbot@gmail.com","account_id":1,"primary_team_id":2,"role_id":5}`,
			rbac: &mock.RBAC{
				AccountCreateFn: func(c echo.Context, roleID models.AccessRole, accountID, teamID uint) error {
					return echo.ErrForbidden
				},
			},
			expectedStatus: http.StatusForbidden,
		},
		{
			name: "Fail on RBAC",
			req:  `{"first_name":"Sarah","last_name":"Smith","username":"sarahsmith","password":"hunter123","password_confirm":"hunter123","email":"sarahsmith@gmail.com","account_id":1,"primary_team_id":2,"role_id":5}`,
			rbac: &mock.RBAC{
				AccountCreateFn: func(c echo.Context, roleID models.AccessRole, accountID, teamID uint) error {
					return echo.ErrForbidden
				},
			},
			expectedStatus: http.StatusForbidden,
		},

		{
			name: "Success",
			req:  `{"first_name":"Edwin","last_name":"Abbott","username":"edwinabbott","password":"hunter123","password_confirm":"hunter123","email":"edwinabbott@gmail.com","account_id":1,"primary_team_id":2,"role_id":5}`,
			rbac: &mock.RBAC{
				AccountCreateFn: func(c echo.Context, roleID models.AccessRole, accountID, teamID uint) error {
					return nil
				},
			},
			udb: &mockstore.UserDBClient{
				CreateFn: func(db *gorm.DB, usr models.User) (*models.User, error) {
					usr.ID = 1
					usr.CreatedAt = mock.TestTime(2018)
					usr.UpdatedAt = mock.TestTime(2018)
					return &usr, nil
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
			expectedResp: &models.User{
				Base: models.Base{
					ID:        1,
					CreatedAt: mock.TestTime(2018),
					UpdatedAt: mock.TestTime(2018),
				},
				FirstName:     "Edwin",
				LastName:      "Abbott",
				Username:      "edwinabbott",
				Email:         "edwinabbott@gmail.com",
				AccountID:     1,
				PrimaryTeamID: 2,
				Role: models.Role{
					ID:          uint(5),
					AccessLevel: models.UserRole,
					Name:        "USER_ADMIN",
				},
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			r := server.New()
			rg := r.Group("")
			transport.NewHTTP(user.New(nil, tt.udb, tt.rbac, tt.sec), rg)
			ts := httptest.NewServer(r)
			defer ts.Close()
			path := ts.URL + "/users"
			res, err := http.Post(path, "application/json", bytes.NewBufferString(tt.req))
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()
			if tt.expectedResp != nil {
				response := new(models.User)
				if err := json.NewDecoder(res.Body).Decode(response); err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, tt.expectedResp, response)
			}
			assert.Equal(t, tt.expectedStatus, res.StatusCode)
		})
	}
}

func TestList(t *testing.T) {
	type listResponse struct {
		Users []models.User `json:"users"`
		Page  int           `json:"page"`
	}
	cases := []struct {
		name           string
		req            string
		expectedStatus int
		expectedResp   *listResponse
		udb            *mockstore.UserDBClient
		rbac           *mock.RBAC
		sec            *mock.Secure
	}{
		{
			name:           "Invalid request",
			req:            `?limit=2222&page=-1`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Fail on query list",
			req:  `?limit=100&page=1`,
			rbac: &mock.RBAC{
				UserFn: func(c echo.Context) *models.AuthUser {
					return &models.AuthUser{
						ID:            1,
						AccountID:     2,
						PrimaryTeamID: 3,
						AccessLevel:   models.UserRole,
						Email:         "barnabus@mail.com",
					}
				}},
			expectedStatus: http.StatusForbidden,
		},
		{
			name: "Success",
			req:  `?limit=100&page=1`,
			rbac: &mock.RBAC{
				UserFn: func(c echo.Context) *models.AuthUser {
					return &models.AuthUser{
						ID:            1,
						AccountID:     2,
						PrimaryTeamID: 3,
						AccessLevel:   models.SuperAdminRole,
						Email:         "pingpong@mail.com",
					}
				}},
			udb: &mockstore.UserDBClient{
				ListFn: func(db *gorm.DB, q *models.ListQuery, p *models.Pagination) ([]models.User, error) {
					if p.Limit == 100 && p.Offset == 100 {
						return []models.User{
							{
								Base: models.Base{
									ID:        10,
									CreatedAt: mock.TestTime(2001),
									UpdatedAt: mock.TestTime(2002),
								},
								FirstName:     "ilove",
								LastName:      "futbol",
								Email:         "futbol@mail.com",
								AccountID:     2,
								PrimaryTeamID: 3,
								Role: models.Role{
									ID:          1,
									AccessLevel: models.SuperAdminRole,
									Name:        "SUPER_ADMIN",
								},
							},
							{
								Base: models.Base{
									ID:        11,
									CreatedAt: mock.TestTime(2004),
									UpdatedAt: mock.TestTime(2005),
								},
								FirstName:     "Joanna",
								LastName:      "Dye",
								Email:         "joanna@mail.com",
								AccountID:     1,
								PrimaryTeamID: 2,
								Role: models.Role{
									ID:          1,
									AccessLevel: models.AdminRole,
									Name:        "ADMIN",
								},
							},
						}, nil
					}
					return nil, models.ErrGeneric
				},
			},
			expectedStatus: http.StatusOK,
			expectedResp: &listResponse{
				Users: []models.User{
					{
						Base: models.Base{
							ID:        10,
							CreatedAt: mock.TestTime(2001),
							UpdatedAt: mock.TestTime(2002),
						},
						FirstName:     "ilove",
						LastName:      "futbol",
						Email:         "futbol@mail.com",
						AccountID:     2,
						PrimaryTeamID: 3,
						Role: models.Role{
							ID:          1,
							AccessLevel: models.SuperAdminRole,
							Name:        "SUPER_ADMIN",
						},
					},
					{
						Base: models.Base{
							ID:        11,
							CreatedAt: mock.TestTime(2004),
							UpdatedAt: mock.TestTime(2005),
						},
						FirstName:     "Joanna",
						LastName:      "Dye",
						Email:         "joanna@mail.com",
						AccountID:     1,
						PrimaryTeamID: 2,
						Role: models.Role{
							ID:          1,
							AccessLevel: models.AdminRole,
							Name:        "ADMIN",
						},
					},
				}, Page: 1},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			r := server.New()
			rg := r.Group("")
			transport.NewHTTP(user.New(nil, tt.udb, tt.rbac, tt.sec), rg)
			ts := httptest.NewServer(r)
			defer ts.Close()
			path := ts.URL + "/users" + tt.req
			res, err := http.Get(path)
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()
			if tt.expectedResp != nil {
				response := new(listResponse)
				if err := json.NewDecoder(res.Body).Decode(response); err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, tt.expectedResp, response)
			}
			assert.Equal(t, tt.expectedStatus, res.StatusCode)
		})
	}
}

func TestView(t *testing.T) {
	cases := []struct {
		name           string
		req            string
		expectedStatus int
		expectedResp   *models.User
		udb            *mockstore.UserDBClient
		rbac           *mock.RBAC
		sec            *mock.Secure
	}{
		{
			name:           "Invalid request",
			req:            `a`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Fail on RBAC",
			req:  `1`,
			rbac: &mock.RBAC{
				EnforceUserFn: func(echo.Context, uint) error {
					return echo.ErrForbidden
				},
			},
			expectedStatus: http.StatusForbidden,
		},
		{
			name: "Success",
			req:  `1`,
			rbac: &mock.RBAC{
				EnforceUserFn: func(echo.Context, uint) error {
					return nil
				},
			},
			udb: &mockstore.UserDBClient{
				ViewFn: func(db *gorm.DB, id uint) (*models.User, error) {
					return &models.User{
						Base: models.Base{
							ID:        1,
							CreatedAt: mock.TestTime(2000),
							UpdatedAt: mock.TestTime(2000),
						},
						FirstName: "Rocinante",
						LastName:  "deLaMancha",
						Username:  "RocinantedeLaMancha",
					}, nil
				},
			},
			expectedStatus: http.StatusOK,
			expectedResp: &models.User{
				Base: models.Base{
					ID:        1,
					CreatedAt: mock.TestTime(2000),
					UpdatedAt: mock.TestTime(2000),
				},
				FirstName: "Rocinante",
				LastName:  "deLaMancha",
				Username:  "RocinantedeLaMancha",
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			r := server.New()
			rg := r.Group("")
			transport.NewHTTP(user.New(nil, tt.udb, tt.rbac, tt.sec), rg)
			ts := httptest.NewServer(r)
			defer ts.Close()
			path := ts.URL + "/users/" + tt.req
			res, err := http.Get(path)
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()
			if tt.expectedResp != nil {
				response := new(models.User)
				if err := json.NewDecoder(res.Body).Decode(response); err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, tt.expectedResp, response)
			}
			assert.Equal(t, tt.expectedStatus, res.StatusCode)
		})
	}
}

func TestUpdate(t *testing.T) {
	cases := []struct {
		name           string
		req            string
		id             string
		expectedStatus int
		expectedResp   *models.User
		udb            *mockstore.UserDBClient
		rbac           *mock.RBAC
		sec            *mock.Secure
	}{
		{
			name:           "Invalid request",
			id:             `a`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Fail on validation",
			id:             `1`,
			req:            `{"first_name":"j","last_name":"okocha","mobile":"123456","phone":"321321","address":"home"}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Fail on RBAC",
			id:   `1`,
			req:  `{"first_name":"jj","last_name":"okocha","mobile":"123456","phone":"321321","address":"home"}`,
			rbac: &mock.RBAC{
				EnforceUserFn: func(echo.Context, uint) error {
					return echo.ErrForbidden
				},
			},
			expectedStatus: http.StatusForbidden,
		},
		{
			name: "Success",
			id:   `1`,
			req:  `{"first_name":"jj","last_name":"okocha","phone":"321321","address":"home"}`,
			rbac: &mock.RBAC{
				EnforceUserFn: func(echo.Context, uint) error {
					return nil
				},
			},
			udb: &mockstore.UserDBClient{
				ViewFn: func(db *gorm.DB, id uint) (*models.User, error) {
					return &models.User{
						Base: models.Base{
							ID:        1,
							CreatedAt: mock.TestTime(2000),
							UpdatedAt: mock.TestTime(2000),
						},
						FirstName: "Nawj",
						LastName:  "Eode",
						Username:  "nawjeode",
						Address:   "Work",
						Phone:     "332223",
					}, nil
				},
				UpdateFn: func(db *gorm.DB, usr *models.User) error {
					usr.UpdatedAt = mock.TestTime(2010)
					usr.Mobile = "991991"
					return nil
				},
			},
			expectedStatus: http.StatusOK,
			expectedResp: &models.User{
				Base: models.Base{
					ID:        1,
					CreatedAt: mock.TestTime(2000),
					UpdatedAt: mock.TestTime(2010),
				},
				FirstName: "jj",
				LastName:  "okocha",
				Username:  "nawjeode",
				Phone:     "321321",
				Address:   "home",
				Mobile:    "991991",
			},
		},
	}

	client := http.Client{}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			r := server.New()
			rg := r.Group("")
			transport.NewHTTP(user.New(nil, tt.udb, tt.rbac, tt.sec), rg)
			ts := httptest.NewServer(r)
			defer ts.Close()
			path := ts.URL + "/users/" + tt.id
			req, _ := http.NewRequest("PATCH", path, bytes.NewBufferString(tt.req))
			req.Header.Set("Content-Type", "application/json")
			res, err := client.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()
			if tt.expectedResp != nil {
				response := new(models.User)
				if err := json.NewDecoder(res.Body).Decode(response); err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, tt.expectedResp, response)
			}
			assert.Equal(t, tt.expectedStatus, res.StatusCode)
		})
	}
}

func TestDelete(t *testing.T) {
	cases := []struct {
		name           string
		id             string
		expectedStatus int
		udb            *mockstore.UserDBClient
		rbac           *mock.RBAC
		sec            *mock.Secure
	}{
		{
			name:           "Invalid request",
			id:             `a`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Fail on RBAC",
			id:   `1`,
			udb: &mockstore.UserDBClient{
				ViewFn: func(db *gorm.DB, id uint) (*models.User, error) {
					return &models.User{
						Role: models.Role{
							AccessLevel: models.AccountAdminRole,
						},
					}, nil
				},
			},
			rbac: &mock.RBAC{
				IsLowerRoleFn: func(echo.Context, models.AccessRole) error {
					return echo.ErrForbidden
				},
			},
			expectedStatus: http.StatusForbidden,
		},
		{
			name: "Success",
			id:   `1`,
			udb: &mockstore.UserDBClient{
				ViewFn: func(db *gorm.DB, id uint) (*models.User, error) {
					return &models.User{
						Role: models.Role{
							AccessLevel: models.AccountAdminRole,
						},
					}, nil
				},
				DeleteFn: func(*gorm.DB, *models.User) error {
					return nil
				},
			},
			rbac: &mock.RBAC{
				IsLowerRoleFn: func(echo.Context, models.AccessRole) error {
					return nil
				},
			},
			expectedStatus: http.StatusOK,
		},
	}

	client := http.Client{}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			r := server.New()
			rg := r.Group("")
			transport.NewHTTP(user.New(nil, tt.udb, tt.rbac, tt.sec), rg)
			ts := httptest.NewServer(r)
			defer ts.Close()
			path := ts.URL + "/users/" + tt.id
			req, _ := http.NewRequest("DELETE", path, nil)
			res, err := client.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()
			assert.Equal(t, tt.expectedStatus, res.StatusCode)
		})
	}
}
