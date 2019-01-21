package mockstore

import (
	"github.com/jinzhu/gorm"

	"github.com/johncoleman83/cerebrum/pkg/utl/models"
)

// UserDBClient database mock
type UserDBClient struct {
	CreateFn         func(*gorm.DB, models.User) (*models.User, error)
	ViewFn           func(*gorm.DB, uint) (*models.User, error)
	FindByUsernameFn func(*gorm.DB, string) (*models.User, error)
	FindByTokenFn    func(*gorm.DB, string) (*models.User, error)
	ListFn           func(*gorm.DB, *models.ListQuery, *models.Pagination) ([]models.User, error)
	DeleteFn         func(*gorm.DB, *models.User) error
	UpdateFn         func(*gorm.DB, *models.User) error
}

// Create mock
func (u *UserDBClient) Create(db *gorm.DB, usr models.User) (*models.User, error) {
	return u.CreateFn(db, usr)
}

// View mock
func (u *UserDBClient) View(db *gorm.DB, id uint) (*models.User, error) {
	return u.ViewFn(db, id)
}

// FindByUsername mock
func (u *UserDBClient) FindByUsername(db *gorm.DB, uname string) (*models.User, error) {
	return u.FindByUsernameFn(db, uname)
}

// FindByToken mock
func (u *UserDBClient) FindByToken(db *gorm.DB, token string) (*models.User, error) {
	return u.FindByTokenFn(db, token)
}

// List mock
func (u *UserDBClient) List(db *gorm.DB, lq *models.ListQuery, p *models.Pagination) ([]models.User, error) {
	return u.ListFn(db, lq, p)
}

// Delete mock
func (u *UserDBClient) Delete(db *gorm.DB, usr *models.User) error {
	return u.DeleteFn(db, usr)
}

// Update mock
func (u *UserDBClient) Update(db *gorm.DB, usr *models.User) error {
	return u.UpdateFn(db, usr)
}
