package mockdb

import (
	"github.com/jinzhu/gorm"

	cerebrum "github.com/johncoleman83/cerebrum/pkg/utl/model"
)

// User database mock
type User struct {
	CreateFn         func(*gorm.DB, cerebrum.User) (*cerebrum.User, error)
	ViewFn           func(*gorm.DB, uint) (*cerebrum.User, error)
	FindByUsernameFn func(*gorm.DB, string) (*cerebrum.User, error)
	FindByTokenFn    func(*gorm.DB, string) (*cerebrum.User, error)
	ListFn           func(*gorm.DB, *cerebrum.ListQuery, *cerebrum.Pagination) ([]cerebrum.User, error)
	DeleteFn         func(*gorm.DB, *cerebrum.User) error
	UpdateFn         func(*gorm.DB, *cerebrum.User) error
}

// Create mock
func (u *User) Create(db *gorm.DB, usr cerebrum.User) (*cerebrum.User, error) {
	return u.CreateFn(db, usr)
}

// View mock
func (u *User) View(db *gorm.DB, id uint) (*cerebrum.User, error) {
	return u.ViewFn(db, id)
}

// FindByUsername mock
func (u *User) FindByUsername(db *gorm.DB, uname string) (*cerebrum.User, error) {
	return u.FindByUsernameFn(db, uname)
}

// FindByToken mock
func (u *User) FindByToken(db *gorm.DB, token string) (*cerebrum.User, error) {
	return u.FindByTokenFn(db, token)
}

// List mock
func (u *User) List(db *gorm.DB, lq *cerebrum.ListQuery, p *cerebrum.Pagination) ([]cerebrum.User, error) {
	return u.ListFn(db, lq, p)
}

// Delete mock
func (u *User) Delete(db *gorm.DB, usr *cerebrum.User) error {
	return u.DeleteFn(db, usr)
}

// Update mock
func (u *User) Update(db *gorm.DB, usr *cerebrum.User) error {
	return u.UpdateFn(db, usr)
}
