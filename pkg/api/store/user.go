// Package store contains the components necessary for api services
// to interact with the database
package store

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"

	cerebrum "github.com/johncoleman83/cerebrum/pkg/utl/model"
)

// Custom errors
var (
	ErrAlreadyExists  = echo.NewHTTPError(http.StatusBadRequest, "username or email already exists")
	ErrRecordNotFound = echo.NewHTTPError(http.StatusNotFound, "user not found")
)

// User represents the client for user table
type User struct{}

// NewUser returns a new user client for db interface
func NewUser() *User {
	return &User{}
}

// Create creates a new user on database
func (u *User) Create(db *gorm.DB, user cerebrum.User) (*cerebrum.User, error) {
	var checkUser = new(cerebrum.User)
	if err := db.Where(
		"lower(username) = ? or lower(email) = ?",
		strings.ToLower(user.Username),
		strings.ToLower(user.Email)).First(&checkUser).Error; err == nil {
		return nil, ErrAlreadyExists
	} else if !gorm.IsRecordNotFoundError(err) {
		return nil, err
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
		return user, ErrRecordNotFound
	} else if err != nil {
		log.Panicln(fmt.Sprintf("db connection error %v", err))
		return user, err
	}
	return user, nil
}

// FindByUsername queries for single user by username
func (u *User) FindByUsername(db *gorm.DB, uname string) (*cerebrum.User, error) {
	var user = new(cerebrum.User)
	if err := db.Set("gorm:auto_preload", true).Where("username = ?", uname).First(&user).Error; gorm.IsRecordNotFoundError(err) {
		return user, ErrRecordNotFound
	} else if err != nil {
		log.Panicln(fmt.Sprintf("db connection error %v", err))
		return user, err
	}
	return user, nil
}

// FindByToken queries for single user by token
func (u *User) FindByToken(db *gorm.DB, token string) (*cerebrum.User, error) {
	var user = new(cerebrum.User)
	if err := db.Set("gorm:auto_preload", true).Where("token = ?", token).First(&user).Error; gorm.IsRecordNotFoundError(err) {
		return user, ErrRecordNotFound
	} else if err != nil {
		log.Panicln(fmt.Sprintf("db connection error %v", err))
		return user, err
	}
	return user, nil
}

// List returns list of all users retrievable for the current user, depending on role
func (u *User) List(db *gorm.DB, qp *cerebrum.ListQuery, p *cerebrum.Pagination) ([]cerebrum.User, error) {
	var users []cerebrum.User
	// Inner Join users with Role
	if qp != nil {
		if err := db.Set("gorm:auto_preload", true).Offset(p.Offset).Limit(p.Limit).Find(&users).Where(qp.Query, qp.ID).Order("lastname asc").Error; gorm.IsRecordNotFoundError(err) {
			return users, ErrRecordNotFound
		} else if err != nil {
			return users, err
		}
	} else {
		if err := db.Set("gorm:auto_preload", true).Offset(p.Offset).Limit(p.Limit).Find(&users).Order("lastname asc").Error; gorm.IsRecordNotFoundError(err) {
			return users, ErrRecordNotFound
		} else if err != nil {
			return users, err
		}
	}
	return users, nil
}

// Update updates user's info
func (u *User) Update(db *gorm.DB, user *cerebrum.User) error {
	return db.Save(user).Error
}

// Delete sets deleted_at for a user
func (u *User) Delete(db *gorm.DB, user *cerebrum.User) error {
	return db.Delete(user).Error
}
