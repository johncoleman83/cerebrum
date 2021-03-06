package transport

import (
	"net/http"
	"strconv"

	"github.com/johncoleman83/cerebrum/pkg/api/password"

	"github.com/johncoleman83/cerebrum/pkg/utl/models"

	"github.com/labstack/echo"
)

// Custom errors
var (
	ErrPasswordsNotMaching = echo.NewHTTPError(http.StatusBadRequest, "passwords do not match")
)

// HTTP represents password http transport service
type HTTP struct {
	svc password.Service
}

// NewHTTP creates new password http service
func NewHTTP(svc password.Service, er *echo.Group) {
	h := HTTP{svc}
	pr := er.Group("/password")

	pr.PATCH("/:id", h.change)
}

// changeReq type for password change request
type changeReq struct {
	ID                 uint   `json:"-"`
	OldPassword        string `json:"old_password" validate:"required,min=8"`
	NewPassword        string `json:"new_password" validate:"required,min=8"`
	NewPasswordConfirm string `json:"new_password_confirm" validate:"required"`
}

// change Changes user's password;
// If user's old passowrd is correct, it will be replaced with new password
//
// usage: PATCH /v1/password/{id} password pwChange
//
// parameters:
// - name: id
//   in: path
//   description: id of user
//   type: integer
//   required: true
// - name: request
//   in: body
//   description: Request body
//   required: true
//   schema:
//     "$ref": "#/definitions/pwChange"
// responses:
//   "200":
//     "$ref": "#/responses/ok"
//   "400":
//     "$ref": "#/responses/errMsg"
//   "401":
//     "$ref": "#/responses/err"
//   "403":
//     "$ref": "#/responses/err"
//   "500":
//     "$ref": "#/responses/err"
func (h *HTTP) change(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return models.ErrBadRequest
	}

	p := new(changeReq)
	if err := c.Bind(p); err != nil {
		return err
	}

	if p.NewPassword != p.NewPasswordConfirm {
		return ErrPasswordsNotMaching
	}

	if err := h.svc.Change(c, uint(id), p.OldPassword, p.NewPassword); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}
