package auth

import (
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/labstack/echo"
	"github.com/johncoleman83/cerebrum/pkg/api/auth/platform/pgsql"
	"github.com/johncoleman83/cerebrum/pkg/utl/model"
)

// New creates new iam service
func New(db *pg.DB, udb UserDB, j TokenGenerator, sec Securer, rbac RBAC) *Auth {
	return &Auth{
		db:   db,
		udb:  udb,
		tg:   j,
		sec:  sec,
		rbac: rbac,
	}
}

// Initialize initializes auth application service
func Initialize(db *pg.DB, j TokenGenerator, sec Securer, rbac RBAC) *Auth {
	return New(db, pgsql.NewUser(), j, sec, rbac)
}

// Service represents auth service interface
type Service interface {
	Authenticate(echo.Context, string, string) (*cerebrum.AuthToken, error)
	Refresh(echo.Context, string) (*cerebrum.RefreshToken, error)
	Me(echo.Context) (*cerebrum.User, error)
}

// Auth represents auth application service
type Auth struct {
	db   *pg.DB
	udb  UserDB
	tg   TokenGenerator
	sec  Securer
	rbac RBAC
}

// UserDB represents user repository interface
type UserDB interface {
	View(orm.DB, int) (*cerebrum.User, error)
	FindByUsername(orm.DB, string) (*cerebrum.User, error)
	FindByToken(orm.DB, string) (*cerebrum.User, error)
	Update(orm.DB, *cerebrum.User) error
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
