package mock

import (
	cerebrum "github.com/johncoleman83/cerebrum/pkg/utl/model"
)

// JWT mock
type JWT struct {
	GenerateTokenFn func(*cerebrum.User) (string, string, error)
}

// GenerateToken mock
func (j *JWT) GenerateToken(u *cerebrum.User) (string, string, error) {
	return j.GenerateTokenFn(u)
}
