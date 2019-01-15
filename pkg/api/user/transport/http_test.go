package transport_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	cerebrum "github.com/johncoleman83/cerebrum/pkg/utl/model"

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
		expectedResp   *cerebrum.User
		udb            *mockstore.User
		rbac           *mock.RBAC
		sec            *mock.Secure
	}{
		{
			name:           "Fail on bad params",
			req:            `{"firstname":"Vanessa","lastname":"Harris","username":"vanessaharris","password":"hunter123","password_confirm":"hunter123","email":"vanessaharris@gmail.com","company_id":1,"location_id":2,"role_id":200}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Fail on validation with short username",
			req:            `{"first_name":"Frank","last_name":"Williams","username":"fw","password":"hunter123","password_confirm":"hunter123","email":"frankwilliams@gmail.com","company_id":1,"location_id":2,"role_id":200}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Fail on validation of email",
			req:            `{"first_name":"Princton","last_name":"Thomas","username":"princetonthomas","password":"hunter123","password_confirm":"hunter123","email":"princetonthomas$gmail.com","company_id":1,"location_id":2,"role_id":200}`,
			expectedStatus: http.StatusBadRequest,
		}, {
			name:           "Fail on non-matching passwords",
			req:            `{"first_name":"Blake","last_name":"Fields","username":"blakefields","password":"sampson","password_confirm":"sampson1","email":"blakefields@gmail.com","company_id":1,"location_id":2,"role_id":200}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Fail on invalid role",
			req:  `{"first_name":"William","last_name":"Abbott","username":"williamabbot","password":"hunter123","password_confirm":"hunter123","email":"williamabbot@gmail.com","company_id":1,"location_id":2,"role_id":199}`,
			rbac: &mock.RBAC{
				AccountCreateFn: func(c echo.Context, roleID cerebrum.AccessRole, companyID, locationID uint) error {
					return echo.ErrForbidden
				},
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Fail on RBAC",
			req:  `{"first_name":"Sarah","last_name":"Smith","username":"sarahsmith","password":"hunter123","password_confirm":"hunter123","email":"sarahsmith@gmail.com","company_id":1,"location_id":2,"role_id":200}`,
			rbac: &mock.RBAC{
				AccountCreateFn: func(c echo.Context, roleID cerebrum.AccessRole, companyID, locationID uint) error {
					return echo.ErrForbidden
				},
			},
			expectedStatus: http.StatusForbidden,
		},

		{
			name: "Success",
			req:  `{"first_name":"Edwin","last_name":"Abbott","username":"edwinabbott","password":"hunter123","password_confirm":"hunter123","email":"edwinabbott@gmail.com","company_id":1,"location_id":2,"role_id":200}`,
			rbac: &mock.RBAC{
				AccountCreateFn: func(c echo.Context, roleID cerebrum.AccessRole, companyID, locationID uint) error {
					return nil
				},
			},
			udb: &mockstore.User{
				CreateFn: func(db *gorm.DB, usr cerebrum.User) (*cerebrum.User, error) {
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
			},
			expectedResp: &cerebrum.User{
				Base: cerebrum.Base{
					Model: gorm.Model{
						ID:        1,
						CreatedAt: mock.TestTime(2018),
						UpdatedAt: mock.TestTime(2018),
					},
				},
				FirstName:  "Edwin",
				LastName:   "Abbott",
				Username:   "edwinabbott",
				Email:      "edwinabbott@gmail.com",
				CompanyID:  1,
				LocationID: 2,
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
				response := new(cerebrum.User)
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
		Users []cerebrum.User `json:"users"`
		Page  int             `json:"page"`
	}
	cases := []struct {
		name           string
		req            string
		expectedStatus int
		expectedResp   *listResponse
		udb            *mockstore.User
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
				UserFn: func(c echo.Context) *cerebrum.AuthUser {
					return &cerebrum.AuthUser{
						ID:         1,
						CompanyID:  2,
						LocationID: 3,
						Role:       cerebrum.UserRole,
						Email:      "barnabus@mail.com",
					}
				}},
			expectedStatus: http.StatusForbidden,
		},
		{
			name: "Success",
			req:  `?limit=100&page=1`,
			rbac: &mock.RBAC{
				UserFn: func(c echo.Context) *cerebrum.AuthUser {
					return &cerebrum.AuthUser{
						ID:         1,
						CompanyID:  2,
						LocationID: 3,
						Role:       cerebrum.SuperAdminRole,
						Email:      "pingpong@mail.com",
					}
				}},
			udb: &mockstore.User{
				ListFn: func(db *gorm.DB, q *cerebrum.ListQuery, p *cerebrum.Pagination) ([]cerebrum.User, error) {
					if p.Limit == 100 && p.Offset == 100 {
						return []cerebrum.User{
							{
								Base: cerebrum.Base{
									Model: gorm.Model{
										ID:        10,
										CreatedAt: mock.TestTime(2001),
										UpdatedAt: mock.TestTime(2002),
									},
								},
								FirstName:  "ilove",
								LastName:   "futbol",
								Email:      "futbol@mail.com",
								CompanyID:  2,
								LocationID: 3,
								Role: cerebrum.Role{
									ID:          cerebrum.SuperAdminRole,
									AccessLevel: cerebrum.SuperAdminRole,
									Name:        "SUPER_ADMIN",
								},
							},
							{
								Base: cerebrum.Base{
									Model: gorm.Model{
										ID:        11,
										CreatedAt: mock.TestTime(2004),
										UpdatedAt: mock.TestTime(2005),
									},
								},
								FirstName:  "Joanna",
								LastName:   "Dye",
								Email:      "joanna@mail.com",
								CompanyID:  1,
								LocationID: 2,
								Role: cerebrum.Role{
									ID:          cerebrum.AdminRole,
									AccessLevel: cerebrum.AdminRole,
									Name:        "ADMIN",
								},
							},
						}, nil
					}
					return nil, cerebrum.ErrGeneric
				},
			},
			expectedStatus: http.StatusOK,
			expectedResp: &listResponse{
				Users: []cerebrum.User{
					{
						Base: cerebrum.Base{
							Model: gorm.Model{
								ID:        10,
								CreatedAt: mock.TestTime(2001),
								UpdatedAt: mock.TestTime(2002),
							},
						},
						FirstName:  "ilove",
						LastName:   "futbol",
						Email:      "futbol@mail.com",
						CompanyID:  2,
						LocationID: 3,
						Role: cerebrum.Role{
							ID:          cerebrum.SuperAdminRole,
							AccessLevel: cerebrum.SuperAdminRole,
							Name:        "SUPER_ADMIN",
						},
					},
					{
						Base: cerebrum.Base{
							Model: gorm.Model{
								ID:        11,
								CreatedAt: mock.TestTime(2004),
								UpdatedAt: mock.TestTime(2005),
							},
						},
						FirstName:  "Joanna",
						LastName:   "Dye",
						Email:      "joanna@mail.com",
						CompanyID:  1,
						LocationID: 2,
						Role: cerebrum.Role{
							ID:          cerebrum.AdminRole,
							AccessLevel: cerebrum.AdminRole,
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
		expectedResp   *cerebrum.User
		udb            *mockstore.User
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
			udb: &mockstore.User{
				ViewFn: func(db *gorm.DB, id uint) (*cerebrum.User, error) {
					return &cerebrum.User{
						Base: cerebrum.Base{
							Model: gorm.Model{
								ID:        1,
								CreatedAt: mock.TestTime(2000),
								UpdatedAt: mock.TestTime(2000),
							},
						},
						FirstName: "Rocinante",
						LastName:  "deLaMancha",
						Username:  "RocinantedeLaMancha",
					}, nil
				},
			},
			expectedStatus: http.StatusOK,
			expectedResp: &cerebrum.User{
				Base: cerebrum.Base{
					Model: gorm.Model{
						ID:        1,
						CreatedAt: mock.TestTime(2000),
						UpdatedAt: mock.TestTime(2000),
					},
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
				response := new(cerebrum.User)
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
		expectedResp   *cerebrum.User
		udb            *mockstore.User
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
			udb: &mockstore.User{
				ViewFn: func(db *gorm.DB, id uint) (*cerebrum.User, error) {
					return &cerebrum.User{
						Base: cerebrum.Base{
							Model: gorm.Model{
								ID:        1,
								CreatedAt: mock.TestTime(2000),
								UpdatedAt: mock.TestTime(2000),
							},
						},
						FirstName: "Nawj",
						LastName:  "Eode",
						Username:  "nawjeode",
						Address:   "Work",
						Phone:     "332223",
					}, nil
				},
				UpdateFn: func(db *gorm.DB, usr *cerebrum.User) error {
					usr.UpdatedAt = mock.TestTime(2010)
					usr.Mobile = "991991"
					return nil
				},
			},
			expectedStatus: http.StatusOK,
			expectedResp: &cerebrum.User{
				Base: cerebrum.Base{
					Model: gorm.Model{
						ID:        1,
						CreatedAt: mock.TestTime(2000),
						UpdatedAt: mock.TestTime(2010),
					},
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
				response := new(cerebrum.User)
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
		udb            *mockstore.User
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
			udb: &mockstore.User{
				ViewFn: func(db *gorm.DB, id uint) (*cerebrum.User, error) {
					return &cerebrum.User{
						Role: cerebrum.Role{
							AccessLevel: cerebrum.CompanyAdminRole,
						},
					}, nil
				},
			},
			rbac: &mock.RBAC{
				IsLowerRoleFn: func(echo.Context, cerebrum.AccessRole) error {
					return echo.ErrForbidden
				},
			},
			expectedStatus: http.StatusForbidden,
		},
		{
			name: "Success",
			id:   `1`,
			udb: &mockstore.User{
				ViewFn: func(db *gorm.DB, id uint) (*cerebrum.User, error) {
					return &cerebrum.User{
						Role: cerebrum.Role{
							AccessLevel: cerebrum.CompanyAdminRole,
						},
					}, nil
				},
				DeleteFn: func(*gorm.DB, *cerebrum.User) error {
					return nil
				},
			},
			rbac: &mock.RBAC{
				IsLowerRoleFn: func(echo.Context, cerebrum.AccessRole) error {
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
