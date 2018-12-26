package auth

import (
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/johncoleman83/cerebrum/pkg/api/auth/platform/mysqldb"
	"github.com/johncoleman83/cerebrum/pkg/utl/model"
)

// New creates new iam service
func New(db *gorm.DB, udb UserDB, j TokenGenerator, sec Securer, rbac RBAC) *Auth {
	return &Auth{
		db:   db,
		udb:  udb,
		tg:   j,
		sec:  sec,
		rbac: rbac,
	}
}

// Initialize initializes auth application service
func Initialize(db *gorm.DB, j TokenGenerator, sec Securer, rbac RBAC) *Auth {
	return New(db, mysqldb.NewUser(), j, sec, rbac)
}

// Service represents auth service interface
type Service interface {
	Authenticate(echo.Context, string, string) (*cerebrum.AuthToken, error)
	Refresh(echo.Context, string) (*cerebrum.RefreshToken, error)
	Me(echo.Context) (*cerebrum.User, error)
}

// Auth represents auth application service
type Auth struct {
	db   *gorm.DB
	udb  UserDB
	tg   TokenGenerator
	sec  Securer
	rbac RBAC
}

// UserDB represents user repository interface
type UserDB interface {
	View(*gorm.DB, uint) (*cerebrum.User, error)
	FindByUsername(*gorm.DB, string) (*cerebrum.User, error)
	FindByToken(*gorm.DB, string) (*cerebrum.User, error)
	Update(*gorm.DB, *cerebrum.User) error
}

// TokenGenerator represents token generator (jwt) interface
type TokenGenerator interface {
	GenerateToken(*cerebrum.User) (string, string, error)
}

// Securer represents security interface
type Securer interface {
	HashMatchesPassword(string, string) bool
	Token(string) string
}

// RBAC represents role-based-access-control interface
type RBAC interface {
	User(echo.Context) *cerebrum.AuthUser
}
