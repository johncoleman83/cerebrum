package mock

import (
	"github.com/johncoleman83/cerebrum/pkg/utl/models"
)

// JWT mock
type JWT struct {
	GenerateTokenFn func(*models.User) (string, string, error)
}

// GenerateToken mock
func (j *JWT) GenerateToken(u *models.User) (string, string, error) {
	return j.GenerateTokenFn(u)
}
