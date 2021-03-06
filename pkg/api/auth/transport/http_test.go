package transport_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"

	"github.com/johncoleman83/cerebrum/pkg/api/auth"
	"github.com/johncoleman83/cerebrum/pkg/api/auth/transport"
	jwtService "github.com/johncoleman83/cerebrum/pkg/utl/middleware/jsonwebtoken"
	"github.com/johncoleman83/cerebrum/pkg/utl/mock"
	"github.com/johncoleman83/cerebrum/pkg/utl/mock/mockstore"
	"github.com/johncoleman83/cerebrum/pkg/utl/models"
	"github.com/johncoleman83/cerebrum/pkg/utl/server"
)

func TestLogin(t *testing.T) {
	cases := []struct {
		name           string
		req            string
		expectedStatus int
		expectedResp   *models.AuthToken
		udb            *mockstore.UserDBClient
		jwt            *mock.JWT
		sec            *mock.Secure
	}{
		{
			name:           "Invalid request",
			req:            `{"username":"juzernejm"}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Fail on FindByUsername",
			req:            `{"username":"juzernejm","password":"hunter123"}`,
			expectedStatus: http.StatusInternalServerError,
			udb: &mockstore.UserDBClient{
				FindByUsernameFn: func(*gorm.DB, string) (*models.User, error) {
					return nil, models.ErrGeneric
				},
			},
		},
		{
			name:           "Success",
			req:            `{"username":"juzernejm","password":"hunter123"}`,
			expectedStatus: http.StatusOK,
			udb: &mockstore.UserDBClient{
				FindByUsernameFn: func(*gorm.DB, string) (*models.User, error) {
					return &models.User{
						Password: "hunter123",
					}, nil
				},
				UpdateFn: func(db *gorm.DB, u *models.User) error {
					return nil
				},
			},
			jwt: &mock.JWT{
				GenerateTokenFn: func(*models.User) (string, string, error) {
					return "jwttokenstring", mock.TestTime(2018).Format(time.RFC3339), nil
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
			expectedResp: &models.AuthToken{Token: "jwttokenstring", Expires: mock.TestTime(2018).Format(time.RFC3339), RefreshToken: "refreshtoken"},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			r := server.New()
			transport.NewHTTP(auth.New(nil, tt.udb, tt.jwt, tt.sec, nil), r, nil)
			ts := httptest.NewServer(r)
			defer ts.Close()
			path := ts.URL + "/login"
			res, err := http.Post(path, "application/json", bytes.NewBufferString(tt.req))
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()
			if tt.expectedResp != nil {
				response := new(models.AuthToken)
				if err := json.NewDecoder(res.Body).Decode(response); err != nil {
					t.Fatal(err)
				}
				tt.expectedResp.RefreshToken = response.RefreshToken
				assert.Equal(t, tt.expectedResp, response)
			}
			assert.Equal(t, tt.expectedStatus, res.StatusCode)
		})
	}
}

func TestRefresh(t *testing.T) {
	cases := []struct {
		name           string
		req            string
		expectedStatus int
		expectedResp   *models.RefreshToken
		udb            *mockstore.UserDBClient
		jwt            *mock.JWT
	}{
		{
			name:           "Fail on FindByToken",
			req:            "refreshtoken",
			expectedStatus: http.StatusInternalServerError,
			udb: &mockstore.UserDBClient{
				FindByTokenFn: func(*gorm.DB, string) (*models.User, error) {
					return nil, models.ErrGeneric
				},
			},
		},
		{
			name:           "Success",
			req:            "refreshtoken",
			expectedStatus: http.StatusOK,
			udb: &mockstore.UserDBClient{
				FindByTokenFn: func(*gorm.DB, string) (*models.User, error) {
					return &models.User{
						Username: "bugsbunny",
					}, nil
				},
			},
			jwt: &mock.JWT{
				GenerateTokenFn: func(*models.User) (string, string, error) {
					return "jwttokenstring", mock.TestTime(2018).Format(time.RFC3339), nil
				},
			},
			expectedResp: &models.RefreshToken{Token: "jwttokenstring", Expires: mock.TestTime(2018).Format(time.RFC3339)},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			r := server.New()
			transport.NewHTTP(auth.New(nil, tt.udb, tt.jwt, nil, nil), r, nil)
			ts := httptest.NewServer(r)
			defer ts.Close()
			path := ts.URL + "/refresh/" + tt.req
			res, err := http.Get(path)
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()
			if tt.expectedResp != nil {
				response := new(models.RefreshToken)
				if err := json.NewDecoder(res.Body).Decode(response); err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, tt.expectedResp, response)
			}
			assert.Equal(t, tt.expectedStatus, res.StatusCode)
		})
	}
}

func TestMe(t *testing.T) {
	cases := []struct {
		name           string
		expectedStatus int
		expectedResp   *models.User
		header         string
		udb            *mockstore.UserDBClient
		rbac           *mock.RBAC
	}{
		{
			name:           "Fail on user view",
			expectedStatus: http.StatusInternalServerError,
			udb: &mockstore.UserDBClient{
				ViewFn: func(*gorm.DB, uint) (*models.User, error) {
					return nil, models.ErrGeneric
				},
			},
			rbac: &mock.RBAC{
				UserFn: func(echo.Context) *models.AuthUser {
					return &models.AuthUser{ID: 1}
				},
			},
			header: mock.HeaderValid(),
		},
		{
			name:           "Success",
			expectedStatus: http.StatusOK,
			udb: &mockstore.UserDBClient{
				ViewFn: func(db *gorm.DB, id uint) (*models.User, error) {
					return &models.User{
						Base: models.Base{
							ID: id,
						},
						AccountID: 2,
						TeamID:    3,
						Email:     "bugs@mail.com",
						FirstName: "Bugs",
						LastName:  "Bunny",
					}, nil
				},
			},
			rbac: &mock.RBAC{
				UserFn: func(echo.Context) *models.AuthUser {
					return &models.AuthUser{ID: 1}
				},
			},
			header: mock.HeaderValid(),
			expectedResp: &models.User{
				Base: models.Base{
					ID: 1,
				},
				AccountID: 2,
				TeamID:    3,
				Email:     "bugs@mail.com",
				FirstName: "Bugs",
				LastName:  "Bunny",
			},
		},
	}

	client := &http.Client{}
	jwtMW := jwtService.New("jwtsecret", "HS256", 60)

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			r := server.New()
			transport.NewHTTP(auth.New(nil, tt.udb, nil, nil, tt.rbac), r, jwtMW.MWFunc())
			ts := httptest.NewServer(r)
			defer ts.Close()
			path := ts.URL + "/me"
			req, err := http.NewRequest("GET", path, nil)
			req.Header.Set("Authorization", tt.header)
			if err != nil {
				t.Fatal(err)
			}
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
