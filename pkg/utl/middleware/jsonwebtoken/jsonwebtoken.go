// Package jsonwebtoken contains logic for using JSON web tokens
package jsonwebtoken

import (
	"net/http"
	"strings"
	"time"

	jwtGo "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"

	"github.com/johncoleman83/cerebrum/pkg/utl/models"
)

// Service provides a Json-Web-Token authentication implementation
type Service struct {
	// Secret key used for signing.
	key []byte

	// Duration for which the jwt token is valid.
	duration time.Duration

	// JWT signing algorithm
	algo jwtGo.SigningMethod
}

// New generates new JWT service necessery for auth middleware
func New(secret, algo string, d int) *Service {
	signingMethod := jwtGo.GetSigningMethod(algo)
	if signingMethod == nil {
		panic("invalid jwt signing method")
	}
	return &Service{
		key:      []byte(secret),
		algo:     signingMethod,
		duration: time.Duration(d) * time.Minute,
	}
}

// MWFunc makes JWT implement the Middleware interface.
func (j *Service) MWFunc() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token, err := j.ParseToken(c)
			if err != nil || !token.Valid {
				return c.NoContent(http.StatusUnauthorized)
			}

			claims := token.Claims.(jwtGo.MapClaims)

			id := uint(claims["id"].(float64))
			accountID := uint(claims["c"].(float64))
			teamID := uint(claims["l"].(float64))
			username := claims["u"].(string)
			email := claims["e"].(string)
			role := models.AccessRole(claims["r"].(float64))

			c.Set("id", id)
			c.Set("account_id", accountID)
			c.Set("team_id", teamID)
			c.Set("username", username)
			c.Set("email", email)
			c.Set("role", role)

			return next(c)
		}
	}
}

// ParseToken parses token from Authorization header
func (j *Service) ParseToken(c echo.Context) (*jwtGo.Token, error) {

	token := c.Request().Header.Get("Authorization")
	if token == "" {
		return nil, models.ErrGeneric
	}
	parts := strings.SplitN(token, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return nil, models.ErrGeneric
	}

	return jwtGo.Parse(parts[1], func(token *jwtGo.Token) (interface{}, error) {
		if j.algo != token.Method {
			return nil, models.ErrGeneric
		}
		return j.key, nil
	})

}

// GenerateToken generates new JWT token and populates it with user data
func (j *Service) GenerateToken(u *models.User) (string, string, error) {
	expire := time.Now().Add(j.duration)

	token := jwtGo.NewWithClaims((j.algo), jwtGo.MapClaims{
		"id":  u.ID,
		"u":   u.Username,
		"e":   u.Email,
		"r":   u.Role.AccessLevel,
		"c":   u.AccountID,
		"l":   u.TeamID,
		"exp": expire.Unix(),
	})

	tokenString, err := token.SignedString(j.key)

	return tokenString, expire.Format(time.RFC3339), err
}
