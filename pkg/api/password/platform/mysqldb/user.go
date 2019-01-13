package mysqldb

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"

	cerebrum "github.com/johncoleman83/cerebrum/pkg/utl/model"
)

// User represents the client for user table
type User struct{}

// NewUser returns a new user client for db interface
func NewUser() *User {
	return &User{}
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

// Update updates user's info
func (u *User) Update(db *gorm.DB, user *cerebrum.User) error {
	return db.Save(user).Error
}
