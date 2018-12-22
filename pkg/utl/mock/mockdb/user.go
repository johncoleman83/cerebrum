package mockdb

import (
	"github.com/go-pg/pg/orm"
	"github.com/johncoleman83/cerebrum/pkg/utl/model"
)

// User database mock
type User struct {
	CreateFn         func(orm.DB, cerebrum.User) (*cerebrum.User, error)
	ViewFn           func(orm.DB, int) (*cerebrum.User, error)
	FindByUsernameFn func(orm.DB, string) (*cerebrum.User, error)
	FindByTokenFn    func(orm.DB, string) (*cerebrum.User, error)
	ListFn           func(orm.DB, *cerebrum.ListQuery, *cerebrum.Pagination) ([]cerebrum.User, error)
	DeleteFn         func(orm.DB, *cerebrum.User) error
	UpdateFn         func(orm.DB, *cerebrum.User) error
}

// Create mock
func (u *User) Create(db orm.DB, usr cerebrum.User) (*cerebrum.User, error) {
	return u.CreateFn(db, usr)
}

// View mock
func (u *User) View(db orm.DB, id int) (*cerebrum.User, error) {
	return u.ViewFn(db, id)
}

// FindByUsername mock
func (u *User) FindByUsername(db orm.DB, uname string) (*cerebrum.User, error) {
	return u.FindByUsernameFn(db, uname)
}

// FindByToken mock
func (u *User) FindByToken(db orm.DB, token string) (*cerebrum.User, error) {
	return u.FindByTokenFn(db, token)
}

// List mock
func (u *User) List(db orm.DB, lq *cerebrum.ListQuery, p *cerebrum.Pagination) ([]cerebrum.User, error) {
	return u.ListFn(db, lq, p)
}

// Delete mock
func (u *User) Delete(db orm.DB, usr *cerebrum.User) error {
	return u.DeleteFn(db, usr)
}

// Update mock
func (u *User) Update(db orm.DB, usr *cerebrum.User) error {
	return u.UpdateFn(db, usr)
}
