package password

import (
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"

	"github.com/johncoleman83/cerebrum/pkg/api/store"
	"github.com/johncoleman83/cerebrum/pkg/utl/models"
)

// Service represents password application interface
type Service interface {
	Change(echo.Context, uint, string, string) error
}

// UserDBClientInterface represents user repository interface
type UserDBClientInterface interface {
	View(*gorm.DB, uint) (*models.User, error)
	Update(*gorm.DB, *models.User) error
}

// Securer represents security interface
type Securer interface {
	Hash(string) string
	HashMatchesPassword(string, string) bool
	Password(string, ...string) bool
}

// RBAC represents role-based-access-control interface
type RBAC interface {
	EnforceUser(echo.Context, uint) error
}

// Password represents password application service
type Password struct {
	db   *gorm.DB
	udb  UserDBClientInterface
	rbac RBAC
	sec  Securer
}

// New creates new password application service
func New(db *gorm.DB, udb UserDBClientInterface, rbac RBAC, sec Securer) *Password {
	return &Password{
		db:   db,
		udb:  udb,
		rbac: rbac,
		sec:  sec,
	}
}

// Initialize initalizes password application service with defaults
func Initialize(db *gorm.DB, rbac RBAC, sec Securer) *Password {
	return New(db, store.NewUserDBClient(), rbac, sec)
}
