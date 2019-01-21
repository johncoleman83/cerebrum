package auth

import (
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"

	"github.com/johncoleman83/cerebrum/pkg/api/store"
	"github.com/johncoleman83/cerebrum/pkg/utl/models"
)

// Service represents auth service interface
type Service interface {
	Authenticate(echo.Context, string, string) (*models.AuthToken, error)
	Refresh(echo.Context, string) (*models.RefreshToken, error)
	Me(echo.Context) (*models.User, error)
}

// UserDBClientInterface represents user repository interface
type UserDBClientInterface interface {
	View(*gorm.DB, uint) (*models.User, error)
	FindByUsername(*gorm.DB, string) (*models.User, error)
	FindByToken(*gorm.DB, string) (*models.User, error)
	Update(*gorm.DB, *models.User) error
}

// TokenGenerator represents token generator (jwt) interface
type TokenGenerator interface {
	GenerateToken(*models.User) (string, string, error)
}

// Securer represents security interface
type Securer interface {
	HashMatchesPassword(string, string) bool
	Token(string) string
}

// RBAC represents role-based-access-control interface
type RBAC interface {
	User(echo.Context) *models.AuthUser
}

// Auth represents auth application service
type Auth struct {
	db   *gorm.DB
	udb  UserDBClientInterface
	tg   TokenGenerator
	sec  Securer
	rbac RBAC
}

// New creates new iam service
func New(db *gorm.DB, udb UserDBClientInterface, j TokenGenerator, sec Securer, rbac RBAC) *Auth {
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
	return New(db, store.NewUserDBClient(), j, sec, rbac)
}
