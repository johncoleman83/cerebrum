package user

import (
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/johncoleman83/cerebrum/pkg/api/user/platform/mysqldb"
	"github.com/johncoleman83/cerebrum/pkg/utl/model"
)

// Service represents user application interface
type Service interface {
	Create(echo.Context, cerebrum.User) (*cerebrum.User, error)
	List(echo.Context, *cerebrum.Pagination) ([]cerebrum.User, error)
	View(echo.Context, uint) (*cerebrum.User, error)
	Delete(echo.Context, uint) error
	Update(echo.Context, *Update) (*cerebrum.User, error)
}

// New creates new user application service
func New(db *gorm.DB, udb UDB, rbac RBAC, sec Securer) *User {
	return &User{db: db, udb: udb, rbac: rbac, sec: sec}
}

// Initialize initalizes User application service with defaults
func Initialize(db *gorm.DB, rbac RBAC, sec Securer) *User {
	return New(db, mysqldb.NewUser(), rbac, sec)
}

// User represents user application service
type User struct {
	db   *gorm.DB
	udb  UDB
	rbac RBAC
	sec  Securer
}

// Securer represents security interface
type Securer interface {
	Hash(string) string
}

// UDB represents user repository interface
type UDB interface {
	Create(*gorm.DB, cerebrum.User) (*cerebrum.User, error)
	View(*gorm.DB, uint) (*cerebrum.User, error)
	List(*gorm.DB, *cerebrum.ListQuery, *cerebrum.Pagination) ([]cerebrum.User, error)
	Update(*gorm.DB, *cerebrum.User) error
	Delete(*gorm.DB, *cerebrum.User) error
}

// RBAC represents role-based-access-control interface
type RBAC interface {
	User(echo.Context) *cerebrum.AuthUser
	EnforceUser(echo.Context, uint) error
	AccountCreate(echo.Context, cerebrum.AccessRole, uint, uint) error
	IsLowerRole(echo.Context, cerebrum.AccessRole) error
}
