package mysqldb

import (
	"net/http"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"

	"github.com/johncoleman83/cerebrum/pkg/utl/model"
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
func (u *User) Create(db *gorm.DB, usr cerebrum.User) (*cerebrum.User, error) {
	var user = new(cerebrum.User)
	found := db.Where(
		"lower(username) = ? or lower(email) = ? and deleted_at is null",
		strings.ToLower(usr.Username), strings.ToLower(usr.Email)).First(&user).RecordNotFound()

	if found || db.Error != nil {
		return nil, ErrAlreadyExists

	}

	if err := db.Create(&usr).Error; err != nil {
		return nil, err
	}
	return &usr, nil
}

// View returns single user by ID
func (u *User) View(db *gorm.DB, id uint) (*cerebrum.User, error) {
	var user = new(cerebrum.User)
	sql := `SELECT "user".*, "role"."id" AS "role__id", "role"."access_level" AS "role__access_level", "role"."name" AS "role__name" 
	FROM "users" AS "user" LEFT JOIN "roles" AS "role" ON "role"."id" = "user"."role_id" 
	WHERE ("user"."id" = ? and deleted_at is null)`
	db.Raw(sql, id).Scan(&user)

	return user, nil
}

// Update updates user's contact info
func (u *User) Update(db *gorm.DB, user *cerebrum.User) error {
	return db.Update(user).Error
}

// List returns list of all users retrievable for the current user, depending on role
func (u *User) List(db *gorm.DB, qp *cerebrum.ListQuery, p *cerebrum.Pagination) ([]cerebrum.User, error) {
	var users []cerebrum.User
	// Inner Join users with Role
	if qp != nil {
		db.Find(&users).Related("role").Where(qp.Query, qp.ID).Where("deleted_at is null").Limit(p.Limit).Offset(p.Offset).Order("user.id desc")
	} else {
		db.Find(&users).Related("role").Where("deleted_at is null").Limit(p.Limit).Offset(p.Offset).Order("user.id desc")
	}
	return users, db.Error
}

// Delete sets deleted_at for a user
func (u *User) Delete(db *gorm.DB, user *cerebrum.User) error {
	return db.Delete(user).Error
}
