package user

import (
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/labstack/echo"
	"github.com/johncoleman83/cerebrum/pkg/api/user/platform/pgsql"
	"github.com/johncoleman83/cerebrum/pkg/utl/model"
)

// Service represents user application interface
type Service interface {
	Create(echo.Context, cerebrum.User) (*cerebrum.User, error)
	List(echo.Context, *cerebrum.Pagination) ([]cerebrum.User, error)
	View(echo.Context, int) (*cerebrum.User, error)
	Delete(echo.Context, int) error
	Update(echo.Context, *Update) (*cerebrum.User, error)
}

// New creates new user application service
func New(db *pg.DB, udb UDB, rbac RBAC, sec Securer) *User {
	return &User{db: db, udb: udb, rbac: rbac, sec: sec}
}

// Initialize initalizes User application service with defaults
func Initialize(db *pg.DB, rbac RBAC, sec Securer) *User {
	return New(db, pgsql.NewUser(), rbac, sec)
}

// User represents user application service
type User struct {
	db   *pg.DB
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
	Create(orm.DB, cerebrum.User) (*cerebrum.User, error)
	View(orm.DB, int) (*cerebrum.User, error)
	List(orm.DB, *cerebrum.ListQuery, *cerebrum.Pagination) ([]cerebrum.User, error)
	Update(orm.DB, *cerebrum.User) error
	Delete(orm.DB, *cerebrum.User) error
}

// RBAC represents role-based-access-control interface
type RBAC interface {
	User(echo.Context) *cerebrum.AuthUser
	EnforceUser(echo.Context, int) error
	AccountCreate(echo.Context, cerebrum.AccessRole, int, int) error
	IsLowerRole(echo.Context, cerebrum.AccessRole) error
}
