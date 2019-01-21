package user

import (
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"

	"github.com/johncoleman83/cerebrum/pkg/api/store"
	"github.com/johncoleman83/cerebrum/pkg/utl/models"
)

// Securer represents security interface
type Securer interface {
	Hash(string) string
	Password(string, ...string) bool
}

// UserDBClient represents user repository interface
type UserDBClient interface {
	Create(*gorm.DB, models.User) (*models.User, error)
	View(*gorm.DB, uint) (*models.User, error)
	List(*gorm.DB, *models.ListQuery, *models.Pagination) ([]models.User, error)
	Update(*gorm.DB, *models.User) error
	Delete(*gorm.DB, *models.User) error
}

// RBAC represents role-based-access-control interface
type RBAC interface {
	User(echo.Context) *models.AuthUser
	EnforceUser(echo.Context, uint) error
	AccountCreate(echo.Context, models.AccessRole, uint, uint) error
	IsLowerRole(echo.Context, models.AccessRole) error
}

// Service represents user application interface
type Service interface {
	Create(echo.Context, models.User) (*models.User, error)
	List(echo.Context, *models.Pagination) ([]models.User, error)
	View(echo.Context, uint) (*models.User, error)
	Delete(echo.Context, uint) error
	Update(echo.Context, *Update) (*models.User, error)
}

// User represents user application service
type User struct {
	db   *gorm.DB
	udb  UserDBClient
	rbac RBAC
	sec  Securer
}

// New creates new user application service
func New(db *gorm.DB, udb UserDBClient, rbac RBAC, sec Securer) *User {
	return &User{db: db, udb: udb, rbac: rbac, sec: sec}
}

// Initialize initalizes User application service with defaults
func Initialize(db *gorm.DB, rbac RBAC, sec Securer) *User {
	return New(db, store.NewUserDBClient(), rbac, sec)
}
