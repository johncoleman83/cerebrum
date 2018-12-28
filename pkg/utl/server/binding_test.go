package server_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator"
	"github.com/johncoleman83/cerebrum/pkg/utl/server"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

type Req struct {
	Name string `json:"name" validate:"required"`
}

func TestBind(t *testing.T) {
	cases := []struct {
		name         string
		req          string
		expectedErr  bool
		expectedData *Req
	}{
		{
			name:         "Fail on binding",
			expectedErr:  true,
			req:          `"bleja"`,
			expectedData: &Req{Name: ""},
		},
		{
			name:         "Fail on validation",
			expectedErr:  true,
			expectedData: &Req{Name: ""},
		},
		{
			name:         "Success",
			req:          `{"name":"John"}`,
			expectedData: &Req{Name: "John"},
		},
	}
	b := server.NewBinder()
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "", bytes.NewBufferString(tt.req))
			req.Header.Set("Content-Type", "application/json")
			e := echo.New()
			e.Validator = &server.CustomValidator{V: validator.New()}
			e.Binder = server.NewBinder()
			c := e.NewContext(req, w)
			r := new(Req)
			err := b.Bind(r, c)
			assert.Equal(t, tt.expectedData, r)
			assert.Equal(t, tt.expectedErr, err != nil)
		})
	}

}
