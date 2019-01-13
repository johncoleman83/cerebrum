// Package transport contians HTTP service for authentication
package transport

import (
	"net/http"

	"github.com/johncoleman83/cerebrum/pkg/api/auth"

	"github.com/labstack/echo"
)

// HTTP represents auth http service
type HTTP struct {
	svc auth.Service
}

// NewHTTP creates new auth http service
func NewHTTP(svc auth.Service, e *echo.Echo, mw echo.MiddlewareFunc) {
	h := HTTP{svc}

	e.POST("/login", h.login)
	e.GET("/refresh/:token", h.refresh)
	e.GET("/me", h.me, mw)
}

// credentials contains a username and password
type credentials struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// login Logs in user by username and password
//
// usage: POST /login auth login
//
// responses:
//  200: loginResp
//  400: errMsg
//  401: errMsg
// 	403: err
//  404: errMsg
//  500: err
func (h *HTTP) login(c echo.Context) error {
	cred := new(credentials)
	if err := c.Bind(cred); err != nil {
		return err
	}
	r, err := h.svc.Authenticate(c, cred.Username, cred.Password)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, r)
}

// refresh Refreshes jwt token by checking if refresh token exists in db
//
// usage: GET /refresh/{token} auth refresh
//
// parameters:
// - name: token
//   in: path
//   description: refresh token
//   type: string
//   required: true
//
// responses:
//   "200":
//     "$ref": "#/responses/refreshResp"
//   "400":
//     "$ref": "#/responses/errMsg"
//   "401":
//     "$ref": "#/responses/err"
//   "500":
//     "$ref": "#/responses/err"
func (h *HTTP) refresh(c echo.Context) error {
	r, err := h.svc.Refresh(c, c.Param("token"))
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, r)
}

// me Gets user's info from session.
//
// usage: GET /me auth meReq
//
// responses:
//  200: userResp
//  500: err
func (h *HTTP) me(c echo.Context) error {
	user, err := h.svc.Me(c)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, user)
}
