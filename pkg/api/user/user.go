package user

import (
	"net/http"

	"github.com/labstack/echo"

	"github.com/johncoleman83/cerebrum/pkg/utl/models"
	"github.com/johncoleman83/cerebrum/pkg/utl/query"
	"github.com/johncoleman83/cerebrum/pkg/utl/structs"
)

// Custom errors
var (
	ErrInsecurePassword = echo.NewHTTPError(http.StatusBadRequest, "insecure password")
)

// Create creates a new user account
func (u *RequestHandler) Create(c echo.Context, req models.User) (*models.User, error) {
	if err := u.rbac.AccountCreate(c, req.Role.AccessLevel, req.AccountID, req.TeamID); err != nil {
		return nil, err
	}
	if ok := u.sec.Password(req.Password, req.FirstName, req.LastName, req.Username, req.Email); !ok {
		return nil, ErrInsecurePassword
	}
	req.Password = u.sec.Hash(req.Password)
	return u.udb.Create(u.db, req)
}

// List returns list of users
func (u *RequestHandler) List(c echo.Context, p *models.Pagination) ([]models.User, error) {
	au := u.rbac.User(c)
	q, err := query.List(au)
	if err != nil {
		return nil, err
	}
	return u.udb.List(u.db, q, p)
}

// View returns single user
func (u *RequestHandler) View(c echo.Context, id uint) (*models.User, error) {
	if err := u.rbac.EnforceUser(c, id); err != nil {
		return nil, err
	}
	return u.udb.View(u.db, id)
}

// Delete deletes a user
func (u *RequestHandler) Delete(c echo.Context, id uint) error {
	user, err := u.udb.View(u.db, id)
	if err != nil {
		return err
	}
	if err := u.rbac.IsLowerRole(c, user.Role.AccessLevel); err != nil {
		return err
	}
	return u.udb.Delete(u.db, user)
}

// Update contains user's information used for updating
type Update struct {
	ID        uint
	FirstName *string
	LastName  *string
	Mobile    *string
	Phone     *string
	Address   *string
}

// Update updates user's contact information
func (u *RequestHandler) Update(c echo.Context, req *Update) (*models.User, error) {
	if err := u.rbac.EnforceUser(c, req.ID); err != nil {
		return nil, err
	}

	user, err := u.udb.View(u.db, req.ID)
	if err != nil {
		return nil, err
	}

	structs.Merge(user, req)
	if err := u.udb.Update(u.db, user); err != nil {
		return nil, err
	}

	return user, nil
}
