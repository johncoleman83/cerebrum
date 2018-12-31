package mysqldb

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"

	cerebrum "github.com/johncoleman83/cerebrum/pkg/utl/model"
)

// NewUser returns a new user database instance
func NewUser() *User {
	return &User{}
}

// User represents the client for user table
type User struct{}

// Custom errors
var (
	ErrAlreadyExists = echo.NewHTTPError(http.StatusInternalServerError, "Username or email already exists.")
)

// Create creates a new user on database
func (u *User) Create(db *gorm.DB, user cerebrum.User) (*cerebrum.User, error) {
	var checkUser = new(cerebrum.User)
	if err := db.Where(
		"lower(username) = ? or lower(email) = ?",
		strings.ToLower(user.Username),
		strings.ToLower(user.Email)).First(&checkUser).Error; !gorm.IsRecordNotFoundError(err) {
		return nil, ErrAlreadyExists
	}
	if err := db.Create(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// View returns single user by ID
func (u *User) View(db *gorm.DB, id uint) (*cerebrum.User, error) {
	var user = new(cerebrum.User)
	if err := db.Set("gorm:auto_preload", true).Where("id = ?", id).First(&user).Error; gorm.IsRecordNotFoundError(err) {
		return user, err
	} else if err != nil {
		log.Panicln(fmt.Sprintf("db connection error %v", err))
		return user, err
	}
	return user, nil
}

// Update updates user's contact info
func (u *User) Update(db *gorm.DB, user *cerebrum.User) error {
	return db.Save(user).Error
}

// List returns list of all users retrievable for the current user, depending on role
func (u *User) List(db *gorm.DB, qp *cerebrum.ListQuery, p *cerebrum.Pagination) ([]cerebrum.User, error) {
	var users []cerebrum.User
	// Inner Join users with Role
	if qp != nil {
		db.Set("gorm:auto_preload", true).Offset(p.Offset).Limit(p.Limit).Find(&users).Where(qp.Query, qp.ID).Order("lastname asc")
	} else {
		db.Set("gorm:auto_preload", true).Offset(p.Offset).Limit(p.Limit).Find(&users).Order("lastname asc")
	}
	return users, db.Error
}

// Delete sets deleted_at for a user
func (u *User) Delete(db *gorm.DB, user *cerebrum.User) error {
	return db.Delete(user).Error
}
