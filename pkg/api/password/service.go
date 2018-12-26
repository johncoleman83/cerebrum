package password

import (
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/johncoleman83/cerebrum/pkg/api/password/platform/mysqldb"
	"github.com/johncoleman83/cerebrum/pkg/utl/model"
)

// Service represents password application interface
type Service interface {
	Change(echo.Context, uint, string, string) error
}

// New creates new password application service
func New(db *gorm.DB, udb UserDB, rbac RBAC, sec Securer) *Password {
	return &Password{
		db:   db,
		udb:  udb,
		rbac: rbac,
		sec:  sec,
	}
}

// Initialize initalizes password application service with defaults
func Initialize(db *gorm.DB, rbac RBAC, sec Securer) *Password {
	return New(db, mysqldb.NewUser(), rbac, sec)
}

// Password represents password application service
type Password struct {
	db   *gorm.DB
	udb  UserDB
	rbac RBAC
	sec  Securer
}

// UserDB represents user repository interface
type UserDB interface {
	View(*gorm.DB, uint) (*cerebrum.User, error)
	Update(*gorm.DB, *cerebrum.User) error
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
