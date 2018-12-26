package transport_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/johncoleman83/cerebrum/pkg/utl/model"

	"github.com/johncoleman83/cerebrum/pkg/api/user"
	"github.com/johncoleman83/cerebrum/pkg/api/user/transport"

	"github.com/johncoleman83/cerebrum/pkg/utl/mock"
	"github.com/johncoleman83/cerebrum/pkg/utl/mock/mockdb"
	"github.com/johncoleman83/cerebrum/pkg/utl/server"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	cases := []struct {
		name       string
		req        string
		wantStatus int
		wantResp   *cerebrum.User
		udb        *mockdb.User
		rbac       *mock.RBAC
		sec        *mock.Secure
	}{
		{
			name:       "Fail on validation",
			req:        `{"first_name":"John","last_name":"Doe","username":"ju","password":"hunter123","password_confirm":"hunter123","email":"johndoe@gmail.com","company_id":1,"location_id":2,"role_id":300}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "Fail on non-matching passwords",
			req:        `{"first_name":"John","last_name":"Doe","username":"juzernejm","password":"hunter123","password_confirm":"hunter1234","email":"johndoe@gmail.com","company_id":1,"location_id":2,"role_id":300}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Fail on invalid role",
			req:  `{"first_name":"John","last_name":"Doe","username":"juzernejm","password":"hunter123","password_confirm":"hunter123","email":"johndoe@gmail.com","company_id":1,"location_id":2,"role_id":50}`,
			rbac: &mock.RBAC{
				AccountCreateFn: func(c echo.Context, roleID cerebrum.AccessRole, companyID, locationID uint) error {
					return echo.ErrForbidden
				},
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Fail on RBAC",
			req:  `{"first_name":"John","last_name":"Doe","username":"juzernejm","password":"hunter123","password_confirm":"hunter123","email":"johndoe@gmail.com","company_id":1,"location_id":2,"role_id":200}`,
			rbac: &mock.RBAC{
				AccountCreateFn: func(c echo.Context, roleID cerebrum.AccessRole, companyID, locationID uint) error {
					return echo.ErrForbidden
				},
			},
			wantStatus: http.StatusForbidden,
		},

		{
			name: "Success",
			req:  `{"first_name":"John","last_name":"Doe","username":"juzernejm","password":"hunter123","password_confirm":"hunter123","email":"johndoe@gmail.com","company_id":1,"location_id":2,"role_id":200}`,
			rbac: &mock.RBAC{
				AccountCreateFn: func(c echo.Context, roleID cerebrum.AccessRole, companyID, locationID uint) error {
					return nil
				},
			},
			udb: &mockdb.User{
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
			wantResp: &cerebrum.User{
				Base: cerebrum.Base{Model: gorm.Model{
					ID:        1,
					CreatedAt: mock.TestTime(2018),
					UpdatedAt: mock.TestTime(2018),
				}},
				FirstName:  "John",
				LastName:   "Doe",
				Username:   "juzernejm",
				Email:      "johndoe@gmail.com",
				CompanyID:  1,
				LocationID: 2,
			},
			wantStatus: http.StatusOK,
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
			if tt.wantResp != nil {
				response := new(cerebrum.User)
				if err := json.NewDecoder(res.Body).Decode(response); err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, tt.wantResp, response)
			}
			assert.Equal(t, tt.wantStatus, res.StatusCode)
		})
	}
}

func TestList(t *testing.T) {
	type listResponse struct {
		Users []cerebrum.User `json:"users"`
		Page  int          `json:"page"`
	}
	cases := []struct {
		name       string
		req        string
		wantStatus int
		wantResp   *listResponse
		udb        *mockdb.User
		rbac       *mock.RBAC
		sec        *mock.Secure
	}{
		{
			name:       "Invalid request",
			req:        `?limit=2222&page=-1`,
			wantStatus: http.StatusBadRequest,
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
						Email:      "john@mail.com",
					}
				}},
			wantStatus: http.StatusForbidden,
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
						Email:      "john@mail.com",
					}
				}},
			udb: &mockdb.User{
				ListFn: func(db *gorm.DB, q *cerebrum.ListQuery, p *cerebrum.Pagination) ([]cerebrum.User, error) {
					if p.Limit == 100 && p.Offset == 100 {
						return []cerebrum.User{
							{
								Base: cerebrum.Base{Model: gorm.Model{
									ID:        10,
									CreatedAt: mock.TestTime(2001),
									UpdatedAt: mock.TestTime(2002),
								}},
								FirstName:  "John",
								LastName:   "Doe",
								Email:      "john@mail.com",
								CompanyID:  2,
								LocationID: 3,
								Role: &cerebrum.Role{
									ID:          1,
									AccessLevel: 1,
									Name:        "SUPER_ADMIN",
								},
							},
							{
								Base: cerebrum.Base{Model: gorm.Model{
									ID:        11,
									CreatedAt: mock.TestTime(2004),
									UpdatedAt: mock.TestTime(2005),
								}},
								FirstName:  "Joanna",
								LastName:   "Dye",
								Email:      "joanna@mail.com",
								CompanyID:  1,
								LocationID: 2,
								Role: &cerebrum.Role{
									ID:          2,
									AccessLevel: 2,
									Name:        "ADMIN",
								},
							},
						}, nil
					}
					return nil, cerebrum.ErrGeneric
				},
			},
			wantStatus: http.StatusOK,
			wantResp: &listResponse{
				Users: []cerebrum.User{
					{
						Base: cerebrum.Base{Model: gorm.Model{
							ID:        10,
							CreatedAt: mock.TestTime(2001),
							UpdatedAt: mock.TestTime(2002),
						}},
						FirstName:  "John",
						LastName:   "Doe",
						Email:      "john@mail.com",
						CompanyID:  2,
						LocationID: 3,
						Role: &cerebrum.Role{
							ID:          1,
							AccessLevel: 1,
							Name:        "SUPER_ADMIN",
						},
					},
					{
						Base: cerebrum.Base{Model: gorm.Model{
							ID:        11,
							CreatedAt: mock.TestTime(2004),
							UpdatedAt: mock.TestTime(2005),
						}},
						FirstName:  "Joanna",
						LastName:   "Dye",
						Email:      "joanna@mail.com",
						CompanyID:  1,
						LocationID: 2,
						Role: &cerebrum.Role{
							ID:          2,
							AccessLevel: 2,
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
			if tt.wantResp != nil {
				response := new(listResponse)
				if err := json.NewDecoder(res.Body).Decode(response); err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, tt.wantResp, response)
			}
			assert.Equal(t, tt.wantStatus, res.StatusCode)
		})
	}
}

func TestView(t *testing.T) {
	cases := []struct {
		name       string
		req        string
		wantStatus int
		wantResp   *cerebrum.User
		udb        *mockdb.User
		rbac       *mock.RBAC
		sec        *mock.Secure
	}{
		{
			name:       "Invalid request",
			req:        `a`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Fail on RBAC",
			req:  `1`,
			rbac: &mock.RBAC{
				EnforceUserFn: func(echo.Context, uint) error {
					return echo.ErrForbidden
				},
			},
			wantStatus: http.StatusForbidden,
		},
		{
			name: "Success",
			req:  `1`,
			rbac: &mock.RBAC{
				EnforceUserFn: func(echo.Context, uint) error {
					return nil
				},
			},
			udb: &mockdb.User{
				ViewFn: func(db *gorm.DB, id uint) (*cerebrum.User, error) {
					return &cerebrum.User{
						Base: cerebrum.Base{Model: gorm.Model{
							ID:        1,
							CreatedAt: mock.TestTime(2000),
							UpdatedAt: mock.TestTime(2000),
						}},
						FirstName: "John",
						LastName:  "Doe",
						Username:  "JohnDoe",
					}, nil
				},
			},
			wantStatus: http.StatusOK,
			wantResp: &cerebrum.User{
				Base: cerebrum.Base{Model: gorm.Model{
					ID:        1,
					CreatedAt: mock.TestTime(2000),
					UpdatedAt: mock.TestTime(2000),
				}},
				FirstName: "John",
				LastName:  "Doe",
				Username:  "JohnDoe",
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
			if tt.wantResp != nil {
				response := new(cerebrum.User)
				if err := json.NewDecoder(res.Body).Decode(response); err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, tt.wantResp, response)
			}
			assert.Equal(t, tt.wantStatus, res.StatusCode)
		})
	}
}

func TestUpdate(t *testing.T) {
	cases := []struct {
		name       string
		req        string
		id         string
		wantStatus int
		wantResp   *cerebrum.User
		udb        *mockdb.User
		rbac       *mock.RBAC
		sec        *mock.Secure
	}{
		{
			name:       "Invalid request",
			id:         `a`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "Fail on validation",
			id:         `1`,
			req:        `{"first_name":"j","last_name":"okocha","mobile":"123456","phone":"321321","address":"home"}`,
			wantStatus: http.StatusBadRequest,
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
			wantStatus: http.StatusForbidden,
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
			udb: &mockdb.User{
				ViewFn: func(db *gorm.DB, id uint) (*cerebrum.User, error) {
					return &cerebrum.User{
						Base: cerebrum.Base{Model: gorm.Model{
							ID:        1,
							CreatedAt: mock.TestTime(2000),
							UpdatedAt: mock.TestTime(2000),
						}},
						FirstName: "John",
						LastName:  "Doe",
						Username:  "JohnDoe",
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
			wantStatus: http.StatusOK,
			wantResp: &cerebrum.User{
				Base: cerebrum.Base{Model: gorm.Model{
					ID:        1,
					CreatedAt: mock.TestTime(2000),
					UpdatedAt: mock.TestTime(2010),
				}},
				FirstName: "jj",
				LastName:  "okocha",
				Username:  "JohnDoe",
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
			if tt.wantResp != nil {
				response := new(cerebrum.User)
				if err := json.NewDecoder(res.Body).Decode(response); err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, tt.wantResp, response)
			}
			assert.Equal(t, tt.wantStatus, res.StatusCode)
		})
	}
}

func TestDelete(t *testing.T) {
	cases := []struct {
		name       string
		id         string
		wantStatus int
		udb        *mockdb.User
		rbac       *mock.RBAC
		sec        *mock.Secure
	}{
		{
			name:       "Invalid request",
			id:         `a`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Fail on RBAC",
			id:   `1`,
			udb: &mockdb.User{
				ViewFn: func(db *gorm.DB, id uint) (*cerebrum.User, error) {
					return &cerebrum.User{
						Role: &cerebrum.Role{
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
			wantStatus: http.StatusForbidden,
		},
		{
			name: "Success",
			id:   `1`,
			udb: &mockdb.User{
				ViewFn: func(db *gorm.DB, id uint) (*cerebrum.User, error) {
					return &cerebrum.User{
						Role: &cerebrum.Role{
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
			wantStatus: http.StatusOK,
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
			assert.Equal(t, tt.wantStatus, res.StatusCode)
		})
	}
}
