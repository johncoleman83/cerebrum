package transport_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/johncoleman83/cerebrum/pkg/api/password"
	"github.com/johncoleman83/cerebrum/pkg/api/password/transport"

	"github.com/johncoleman83/cerebrum/pkg/utl/mock"
	"github.com/johncoleman83/cerebrum/pkg/utl/mock/mockstore"
	"github.com/johncoleman83/cerebrum/pkg/utl/models"
	"github.com/johncoleman83/cerebrum/pkg/utl/server"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

func TestChangePassword(t *testing.T) {
	cases := []struct {
		name           string
		req            string
		expectedStatus int
		id             string
		udb            *mockstore.User
		rbac           *mock.RBAC
		sec            *mock.Secure
	}{
		{
			name:           "NaN",
			expectedStatus: http.StatusBadRequest,
			id:             "abc",
		},
		{
			name:           "Fail on Bind",
			req:            `{"new_password":"new","old_password":"my_old_password", "new_password_confirm":"new"}`,
			expectedStatus: http.StatusBadRequest,
			id:             "1",
		},
		{
			name:           "Different passwords",
			req:            `{"new_password":"new_password","old_password":"my_old_password", "new_password_confirm":"new_password_cf"}`,
			expectedStatus: http.StatusBadRequest,
			id:             "1",
		},
		{
			name: "Fail on RBAC",
			req:  `{"new_password":"newpassw","old_password":"oldpassw", "new_password_confirm":"newpassw"}`,
			rbac: &mock.RBAC{
				EnforceUserFn: func(c echo.Context, id uint) error {
					return echo.ErrForbidden
				},
			},
			id:             "1",
			expectedStatus: http.StatusForbidden,
		},
		{
			name: "Success",
			req:  `{"new_password":"newpassw","old_password":"oldpassw", "new_password_confirm":"newpassw"}`,
			rbac: &mock.RBAC{
				EnforceUserFn: func(c echo.Context, id uint) error {
					return nil
				},
			},
			id: "1",
			udb: &mockstore.User{
				ViewFn: func(db *gorm.DB, id uint) (*models.User, error) {
					return &models.User{
						Password: "oldPassword",
					}, nil
				},
				UpdateFn: func(db *gorm.DB, usr *models.User) error {
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
					return "hashedPassword"
				},
			},
			expectedStatus: http.StatusOK,
		},
	}

	client := &http.Client{}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			r := server.New()
			rg := r.Group("")
			transport.NewHTTP(password.New(nil, tt.udb, tt.rbac, tt.sec), rg)
			ts := httptest.NewServer(r)
			defer ts.Close()
			path := ts.URL + "/password/" + tt.id
			req, err := http.NewRequest("PATCH", path, bytes.NewBufferString(tt.req))
			req.Header.Set("Content-Type", "application/json")
			if err != nil {
				t.Fatal(err)
			}
			res, err := client.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()
			assert.Equal(t, tt.expectedStatus, res.StatusCode)
		})
	}
}
