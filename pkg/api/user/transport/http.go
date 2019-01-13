// Package transport contains the HTTP service for user interactions
package transport

import (
	"net/http"
	"strconv"

	"github.com/johncoleman83/cerebrum/pkg/api/user"

	cerebrum "github.com/johncoleman83/cerebrum/pkg/utl/model"

	"github.com/labstack/echo"
)

// Custom errors
var (
	ErrPasswordsNotMaching = echo.NewHTTPError(http.StatusBadRequest, "passwords do not match")
)

// HTTP represents user http service
type HTTP struct {
	svc user.Service
}

// NewHTTP creates new user http service
func NewHTTP(svc user.Service, er *echo.Group) {
	h := HTTP{svc}
	ur := er.Group("/users")

	ur.POST("", h.create)
	ur.GET("", h.list)
	ur.GET("/:id", h.view)
	ur.PATCH("/:id", h.update)
	ur.DELETE("/:id", h.delete)
}

// createReq is a used to serialize the request payload to a struct
type createReq struct {
	FirstName       string `json:"first_name" validate:"required"`
	LastName        string `json:"last_name" validate:"required"`
	Username        string `json:"username" validate:"required,min=3,alphanum"`
	Password        string `json:"password" validate:"required,min=8"`
	PasswordConfirm string `json:"password_confirm" validate:"required"`
	Email           string `json:"email" validate:"required,email"`

	CompanyID  uint                `json:"company_id" validate:"required"`
	LocationID uint                `json:"location_id" validate:"required"`
	RoleID     cerebrum.AccessRole `json:"role_id" validate:"required"`
}

// create Creates new user account
//
// usage: POST /v1/users users userCreate
//
// responses:
//  200: userResp
//  400: errMsg
//  401: err
//  403: errMsg
//  500: err
func (h *HTTP) create(c echo.Context) error {
	r := new(createReq)

	if err := c.Bind(r); err != nil {

		return err
	}

	if r.Password != r.PasswordConfirm {
		return ErrPasswordsNotMaching
	}

	if r.RoleID < cerebrum.SuperAdminRole || r.RoleID > cerebrum.UserRole {
		return cerebrum.ErrBadRequest
	}

	usr, err := h.svc.Create(c, cerebrum.User{
		Username:   r.Username,
		Password:   r.Password,
		Email:      r.Email,
		FirstName:  r.FirstName,
		LastName:   r.LastName,
		CompanyID:  r.CompanyID,
		LocationID: r.LocationID,
		RoleID:     r.RoleID,
	})

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, usr)
}

// listResponse contains the users list and page for the list response
type listResponse struct {
	Users []cerebrum.User `json:"users"`
	Page  int             `json:"page"`
}

// list Returns list of users. Depending on the user role requesting it:
// it may return all users for SuperAdmin/Admin users,
// all company/location users for Company/Location admins
// and an error for non-admin users.
//
// usage: GET /v1/users users listUsers
//
// parameters:
// - name: limit
//   in: query
//   description: number of results
//   type: integer
//   required: false
// - name: page
//   in: query
//   description: page number
//   type: integer
//   required: false
//
// responses:
//   "200":
//     "$ref": "#/responses/userListResp"
//   "400":
//     "$ref": "#/responses/errMsg"
//   "401":
//     "$ref": "#/responses/err"
//   "403":
//     "$ref": "#/responses/err"
//   "500":
//     "$ref": "#/responses/err"
func (h *HTTP) list(c echo.Context) error {
	p := new(cerebrum.PaginationReq)
	if err := c.Bind(p); err != nil {
		return err
	}

	result, err := h.svc.List(c, p.NewPagination())

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, listResponse{result, p.Page})
}

// view returns a single user with same id as request id
//
// usage: GET /v1/users/{id} users getUser
//
// parameters:
// - name: id
//   in: path
//   description: id of user
//   type: integer
//   required: true
//
// responses:
//   "200":
//     "$ref": "#/responses/userResp"
//   "400":
//     "$ref": "#/responses/err"
//   "401":
//     "$ref": "#/responses/err"
//   "403":
//     "$ref": "#/responses/err"
//   "404":
//     "$ref": "#/responses/err"
//   "500":
//     "$ref": "#/responses/err"
func (h *HTTP) view(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return cerebrum.ErrBadRequest
	}

	result, err := h.svc.View(c, uint(id))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, result)
}

// updateReq is used to serialize the request payload to a struct
type updateReq struct {
	ID        uint    `json:"-"`
	FirstName *string `json:"first_name,omitempty" validate:"omitempty,min=2"`
	LastName  *string `json:"last_name,omitempty" validate:"omitempty,min=2"`
	Mobile    *string `json:"mobile,omitempty"`
	Phone     *string `json:"phone,omitempty"`
	Address   *string `json:"address,omitempty"`
}

// update updates user's contact information -> first name, last name, mobile, phone, address
//
// usage: PATCH /v1/users/{id} users userUpdate
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
//     "$ref": "#/definitions/userUpdate"
//
// responses:
//   "200":
//     "$ref": "#/responses/userResp"
//   "400":
//     "$ref": "#/responses/errMsg"
//   "401":
//     "$ref": "#/responses/err"
//   "403":
//     "$ref": "#/responses/err"
//   "500":
//     "$ref": "#/responses/err"
func (h *HTTP) update(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return cerebrum.ErrBadRequest
	}

	req := new(updateReq)
	if err := c.Bind(req); err != nil {
		return err
	}

	usr, err := h.svc.Update(c, &user.Update{
		ID:        uint(id),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Mobile:    req.Mobile,
		Phone:     req.Phone,
		Address:   req.Address,
	})

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, usr)
}

// delete deletes a user with requested ID.
//
// usage: DELETE /v1/users/{id} users userDelete
//
// parameters:
// - name: id
//   in: path
//   description: id of user
//   type: integer
//   required: true
//
// responses:
//   "200":
//     "$ref": "#/responses/ok"
//   "400":
//     "$ref": "#/responses/err"
//   "401":
//     "$ref": "#/responses/err"
//   "403":
//     "$ref": "#/responses/err"
//   "500":
//     "$ref": "#/responses/err"
func (h *HTTP) delete(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return cerebrum.ErrBadRequest
	}

	if err := h.svc.Delete(c, uint(id)); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}
