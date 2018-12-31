package jwt_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"

	"github.com/johncoleman83/cerebrum/pkg/utl/middleware/jwt"
	"github.com/johncoleman83/cerebrum/pkg/utl/mock"
	cerebrum "github.com/johncoleman83/cerebrum/pkg/utl/model"
)

func echoHandler(mw ...echo.MiddlewareFunc) *echo.Echo {
	e := echo.New()
	for _, v := range mw {
		e.Use(v)
	}
	e.GET("/hello", hwHandler)
	return e
}

func hwHandler(c echo.Context) error {
	return c.String(200, "Hello World")
}

func TestMWFunc(t *testing.T) {
	cases := []struct {
		name           string
		expectedStatus int
		header         string
		signMethod     string
	}{
		{
			name:           "Empty header",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Header not containing Bearer",
			header:         "notBearer",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Invalid header",
			header:         mock.HeaderInvalid(),
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Success",
			header:         mock.HeaderValid(),
			expectedStatus: http.StatusOK,
		},
	}
	jwtMW := jwt.New("jwtsecret", "HS256", 60)
	ts := httptest.NewServer(echoHandler(jwtMW.MWFunc()))
	defer ts.Close()
	path := ts.URL + "/hello"
	client := &http.Client{}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", path, nil)
			req.Header.Set("Authorization", tt.header)
			res, err := client.Do(req)
			if err != nil {
				t.Fatal("Cannot create http request")
			}
			assert.Equal(t, tt.expectedStatus, res.StatusCode)
		})
	}
}

func TestGenerateToken(t *testing.T) {
	cases := []struct {
		name          string
		expectedToken string
		algo          string
		req           *cerebrum.User
	}{
		{
			name: "Invalid algo",
			algo: "invalid",
		},
		{
			name: "Success",
			algo: "HS256",
			req: &cerebrum.User{
				Base: cerebrum.Base{Model: gorm.Model{
					ID: 1,
				}},
				Username: "johndoe",
				Email:    "johndoe@mail.com",
				Role: cerebrum.Role{
					AccessLevel: cerebrum.SuperAdminRole,
				},
				CompanyID:  1,
				LocationID: 1,
			},
			expectedToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			if tt.algo != "HS256" {
				assert.Panics(t, func() {
					jwt.New("jwtsecret", tt.algo, 60)
				}, "The code did not panic")
				return
			}
			jwt := jwt.New("jwtsecret", tt.algo, 60)
			str, _, err := jwt.GenerateToken(tt.req)
			assert.Nil(t, err)
			assert.Equal(t, tt.expectedToken, strings.Split(str, ".")[0])
		})
	}
}
