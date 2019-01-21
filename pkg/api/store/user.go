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

	"github.com/johncoleman83/cerebrum/pkg/utl/models"
)

// Custom errors
var (
	ErrAlreadyExists  = echo.NewHTTPError(http.StatusBadRequest, "username or email already exists")
	ErrRecordNotFound = echo.NewHTTPError(http.StatusNotFound, "user not found")
)

// UserDBClient represents the client for user table
type UserDBClient struct{}

// NewUserDBClient returns a new user client for db interface
func NewUserDBClient() *UserDBClient {
	return &UserDBClient{}
}

// Create creates a new user on database
func (u *UserDBClient) Create(db *gorm.DB, user models.User) (*models.User, error) {
	var checkUser = new(models.User)
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
func (u *UserDBClient) View(db *gorm.DB, id uint) (*models.User, error) {
	var user = new(models.User)
	if err := db.Set("gorm:auto_preload", true).Where("id = ?", id).First(&user).Error; gorm.IsRecordNotFoundError(err) {
		return user, ErrRecordNotFound
	} else if err != nil {
		log.Panicln(fmt.Sprintf("db connection error %v", err))
		return user, err
	}
	return user, nil
}

// FindByUsername queries for single user by username
func (u *UserDBClient) FindByUsername(db *gorm.DB, uname string) (*models.User, error) {
	var user = new(models.User)
	if err := db.Set("gorm:auto_preload", true).Where("username = ?", uname).First(&user).Error; gorm.IsRecordNotFoundError(err) {
		return user, ErrRecordNotFound
	} else if err != nil {
		log.Panicln(fmt.Sprintf("db connection error %v", err))
		return user, err
	}
	return user, nil
}

// FindByToken queries for single user by token
func (u *UserDBClient) FindByToken(db *gorm.DB, token string) (*models.User, error) {
	var user = new(models.User)
	if err := db.Set("gorm:auto_preload", true).Where("token = ?", token).First(&user).Error; gorm.IsRecordNotFoundError(err) {
		return user, ErrRecordNotFound
	} else if err != nil {
		log.Panicln(fmt.Sprintf("db connection error %v", err))
		return user, err
	}
	return user, nil
}

// List returns list of all users retrievable for the current user, depending on role
func (u *UserDBClient) List(db *gorm.DB, qp *models.ListQuery, p *models.Pagination) ([]models.User, error) {
	var users []models.User
	// Inner Join users with Role
	if qp != nil {
		if err := db.Set("gorm:auto_preload", true).Offset(p.Offset).Limit(p.Limit).Where(qp.Query, qp.ID).Find(&users).Order("lastname asc").Error; err != nil {
			log.Panicln(fmt.Sprintf("db connection error %v", err))
			return users, err
		}
	} else {
		if err := db.Set("gorm:auto_preload", true).Offset(p.Offset).Limit(p.Limit).Find(&users).Order("lastname asc").Error; err != nil {
			log.Panicln(fmt.Sprintf("db connection error %v", err))
			return users, err
		}
	}
	return users, nil
}

// Update updates user's info
func (u *UserDBClient) Update(db *gorm.DB, user *models.User) error {
	return db.Save(user).Error
}

// Delete sets deleted_at for a user
func (u *UserDBClient) Delete(db *gorm.DB, user *models.User) error {
	return db.Delete(user).Error
}
