package mysqldb

import (
	"github.com/jinzhu/gorm"
	"github.com/johncoleman83/cerebrum/pkg/utl/model"
)

// NewUser returns a new user database instance
func NewUser() *User {
	return &User{}
}

// User represents the client for user table
type User struct{}

// View returns single user by ID
func (u *User) View(db *gorm.DB, id uint) (*cerebrum.User, error) {
	var user = new(cerebrum.User)
	sql := `SELECT "user".*, "role"."id" AS "role__id", "role"."access_level" AS "role__access_level", "role"."name" AS "role__name" 
	FROM "users" AS "user" LEFT JOIN "roles" AS "role" ON "role"."id" = "user"."role_id" 
	WHERE ("user"."id" = ? and deleted_at is null)
	LIMIT 1`
	db.Raw(sql, id).Scan(&user)
	return user, db.Error
}

// FindByUsername queries for single user by username
func (u *User) FindByUsername(db *gorm.DB, uname string) (*cerebrum.User, error) {
	var user = new(cerebrum.User)
	sql := `SELECT "user".*, "role"."id" AS "role__id", "role"."access_level" AS "role__access_level", "role"."name" AS "role__name" 
	FROM "users" AS "user" LEFT JOIN "roles" AS "role" ON "role"."id" = "user"."role_id" 
	WHERE ("user"."username" = ? and deleted_at is null)
	LIMIT 1`
	db.Raw(sql, uname).Scan(&user)
	return user, db.Error
}

// FindByToken queries for single user by token
func (u *User) FindByToken(db *gorm.DB, token string) (*cerebrum.User, error) {
	var user = new(cerebrum.User)
	sql := `SELECT "user".*, "role"."id" AS "role__id", "role"."access_level" AS "role__access_level", "role"."name" AS "role__name" 
	FROM "users" AS "user" LEFT JOIN "roles" AS "role" ON "role"."id" = "user"."role_id" 
	WHERE ("user"."token" = ? and deleted_at is null)
	LIMIT 1`
	db.Raw(sql, token).Scan(&user)
	return user, db.Error
}

// Update updates user's info
func (u *User) Update(db *gorm.DB, user *cerebrum.User) error {
	return db.Update(user).Error
}
